package socket

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

const testNamespace = "test"

type testGoroutineManager struct{}

func (manager *testGoroutineManager) Goroutine(fn func()) {
	go fn()
}

func (manager *testGoroutineManager) GoroutineWithRestart(fn func()) {
	go fn()
}

func TestServeHTTPManagesSocketLifecycle(t *testing.T) {
	registered := make(chan *Socket, 1)
	unregistered := make(chan *Socket, 1)
	handlerErrors := make(chan error, 2)
	registry := newTestRegistry(registered, unregistered)

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		err := registry.ServeHTTP(writer, request, ConnectionOptions{
			ID:        "socket-1",
			Namespace: testNamespace,
			Context:   request.Context(),
			ReadMessageHandler: func(socketHolder *Socket, rawData []byte) {
				if err := socketHolder.Emit(string(rawData)); err != nil {
					handlerErrors <- err
				}
			},
		})
		handlerErrors <- err
	}))
	defer server.Close()

	connection, _, err := websocket.DefaultDialer.Dial(
		"ws"+strings.TrimPrefix(server.URL, "http"),
		nil,
	)
	require.NoError(t, err)

	select {
	case socketHolder := <-registered:
		require.Equal(t, "socket-1", socketHolder.ID)
	case <-time.After(time.Second):
		t.Fatal("socket was not registered")
	}

	require.NoError(t, registry.Emit("socket-1", map[string]string{"source": "server"}))
	_, serverMessage, err := connection.ReadMessage()
	require.NoError(t, err)
	require.JSONEq(t, `{"source":"server"}`, string(serverMessage))

	require.NoError(t, connection.WriteMessage(websocket.TextMessage, []byte("client")))
	_, echoedMessage, err := connection.ReadMessage()
	require.NoError(t, err)
	require.JSONEq(t, `"client"`, string(echoedMessage))

	require.NoError(t, connection.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
	))
	require.NoError(t, connection.Close())

	select {
	case socketHolder := <-unregistered:
		require.Equal(t, "socket-1", socketHolder.ID)
	case <-time.After(time.Second):
		t.Fatal("socket was not unregistered")
	}

	_, ok := registry.Sockets.Load("socket-1")
	require.False(t, ok)

	select {
	case handlerErr := <-handlerErrors:
		require.NoError(t, handlerErr)
	case <-time.After(time.Second):
		t.Fatal("socket handler did not finish")
	}
}

func TestServeHTTPValidatesConnectionBeforeUpgrade(t *testing.T) {
	registry := newTestRegistry(nil, nil)

	testCases := []struct {
		name     string
		options  ConnectionOptions
		expected error
	}{
		{
			name:     "missing socket ID",
			options:  ConnectionOptions{Namespace: testNamespace},
			expected: ErrMissingSocketID,
		},
		{
			name:     "missing namespace",
			options:  ConnectionOptions{ID: "socket-1"},
			expected: ErrMissingNamespace,
		},
		{
			name:     "unknown namespace",
			options:  ConnectionOptions{ID: "socket-1", Namespace: "unknown"},
			expected: ErrUnknownNamespace,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/socket", nil)

			err := registry.ServeHTTP(response, request, testCase.options)

			require.Error(t, err)
			require.True(t, errors.Is(err, testCase.expected))
			require.Equal(t, http.StatusOK, response.Code)
			require.Empty(t, response.Header().Get("Upgrade"))
		})
	}
}

func TestServeHTTPReplacesDuplicateConnectionID(t *testing.T) {
	registered := make(chan *Socket, 2)
	unregistered := make(chan *Socket, 2)
	registry := newTestRegistry(registered, unregistered)

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_ = registry.ServeHTTP(writer, request, ConnectionOptions{
			ID:        "same-socket",
			Namespace: testNamespace,
			Context:   request.Context(),
		})
	}))
	defer server.Close()

	websocketURL := "ws" + strings.TrimPrefix(server.URL, "http")
	firstConnection, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
	require.NoError(t, err)
	firstSocket := receiveSocket(t, registered, "first socket was not registered")

	secondConnection, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
	require.NoError(t, err)
	defer secondConnection.Close()

	unregisteredFirst := receiveSocket(t, unregistered, "first socket was not unregistered")
	require.Same(t, firstSocket, unregisteredFirst)
	secondSocket := receiveSocket(t, registered, "second socket was not registered")
	require.NotSame(t, firstSocket, secondSocket)

	current, ok := registry.Sockets.Load("same-socket")
	require.True(t, ok)
	require.Same(t, secondSocket, current)

	_, _, err = firstConnection.ReadMessage()
	require.Error(t, err)
	require.NoError(t, firstConnection.Close())

	require.NoError(t, secondConnection.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
	))
	require.Same(t, secondSocket, receiveSocket(t, unregistered, "second socket was not unregistered"))
}

func TestDefaultOriginPolicy(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "https://api.example.com/socket", nil)

	require.True(t, checkSameOrigin(request))

	request.Header.Set("Origin", "https://api.example.com")
	require.True(t, checkSameOrigin(request))

	request.Header.Set("Origin", "https://attacker.example.com")
	require.False(t, checkSameOrigin(request))
}

func TestEmitDoesNotBlockSlowOrClosedConnections(t *testing.T) {
	socketHolder := &Socket{
		Ctx: context.Background(),
		Connection: &Connection{
			Send: make(chan []byte, 1),
		},
		done: make(chan struct{}),
	}

	require.NoError(t, socketHolder.Emit("first"))
	require.ErrorIs(t, socketHolder.Emit("second"), ErrSendBufferFull)

	close(socketHolder.done)
	require.ErrorIs(t, socketHolder.Emit("third"), ErrSocketClosed)
}

func newTestRegistry(registered chan *Socket, unregistered chan *Socket) *Registry {
	return NewSocketRegistry(
		NamespaceEventHandlerMap{
			testNamespace: {
				RegisterHandler: func(socketHolder *Socket) {
					if registered != nil {
						registered <- socketHolder
					}
				},
				UnregisterHandler: func(socketHolder *Socket) {
					if unregistered != nil {
						unregistered <- socketHolder
					}
				},
			},
		},
		&testGoroutineManager{},
	)
}

func receiveSocket(t *testing.T, channel <-chan *Socket, timeoutMessage string) *Socket {
	t.Helper()

	select {
	case socketHolder := <-channel:
		return socketHolder
	case <-time.After(time.Second):
		t.Fatal(timeoutMessage)
	}

	return nil
}
