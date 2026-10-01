package mcp_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/coretrix/hitrix"
	"github.com/coretrix/hitrix/pkg/mcp"
	"github.com/coretrix/hitrix/service"
	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
	hitrixRegistry "github.com/coretrix/hitrix/service/registry"
)

const testToken = "isolated-mcp-test-token" //nolint:gosec // Credential used only in isolated httptest requests.

func newRouter(testingContext *testing.T, read func(context.Context, string) map[string]string) http.Handler {
	testingContext.Helper()
	appFolder := testingContext.TempDir()
	configFolder := filepath.Join(appFolder, "config")
	require.NoError(testingContext, os.MkdirAll(filepath.Join(configFolder, "mcp-test"), 0700))
	require.NoError(testingContext, os.WriteFile(filepath.Join(configFolder, "mcp-test", "config.yaml"),
		[]byte("cors: [https://allowed.example]\n"), 0600))
	require.NoError(testingContext, os.WriteFile(filepath.Join(configFolder, "hitrix.yaml"),
		[]byte("mcp:\n  auth_token: isolated-mcp-test-token\n"), 0600))
	testingContext.Setenv("APP_FOLDER", appFolder)
	testingContext.Setenv("APP_MODE", "test")

	_, closeServer := hitrix.New("mcp-test", "test").RegisterDIGlobalService(
		hitrixRegistry.ServiceProviderConfigDirectory(configFolder),
	).RegisterMCPCommands([]mcp.Command{mcp.NewErrorLogCommand(read)}).Build()

	testingContext.Cleanup(func() {
		closeServer()
		service.SetContainer(nil)
	})

	return hitrix.InitGin(nil)
}

func rpc(testingContext *testing.T, router http.Handler, method string, params interface{}) map[string]json.RawMessage {
	testingContext.Helper()

	encoded, err := json.Marshal(map[string]interface{}{"jsonrpc": "2.0", "id": 1, "method": method, "params": params})
	require.NoError(testingContext, err)

	request := httptest.NewRequest(http.MethodPost, mcp.Endpoint, bytes.NewReader(encoded))
	request.Header.Set("Authorization", "Bearer "+testToken)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json, text/event-stream")
	request.Header.Set("Mcp-Protocol-Version", "2025-03-26")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	require.Equal(testingContext, http.StatusOK, recorder.Code, recorder.Body.String())
	var result map[string]json.RawMessage
	require.NoError(testingContext, json.Unmarshal(recorder.Body.Bytes(), &result))
	require.NotContains(testingContext, result, "error", recorder.Body.String())

	return result
}

func callLogs(testingContext *testing.T, router http.Handler, arguments map[string]interface{}) *sdk.CallToolResult {
	testingContext.Helper()
	result := rpc(testingContext, router, "tools/call", map[string]interface{}{"name": "get_error_logs", "arguments": arguments})
	var toolResult sdk.CallToolResult
	require.NoError(testingContext, json.Unmarshal(result["result"], &toolResult))

	return &toolResult
}

func TestMCPInitializeDiscoverAndReadLogs(testingContext *testing.T) {
	message, err := json.Marshal(errorlogger.ErrorMessage{
		File: "worker.go", Line: 42, AppName: "worker", Service: "payments", Message: "provider failed",
		Request: []byte("POST /purchase\nAuthorization: Bearer private-header\nCookie: session=private-cookie\n" +
			`{"password":"private-password","token":"private-token","license":"CA1234AA"}`),
		Stack: []byte("worker.process:42"),
	})
	require.NoError(testingContext, err)

	reads := 0
	router := newRouter(testingContext, func(ctx context.Context, group string) map[string]string {
		reads++

		assert.Equal(testingContext, "service", group)
		assert.Equal(testingContext, mcp.Endpoint, service.GinFromContext(ctx).Request.URL.Path)

		return map[string]string{"first": string(message), "first:counter": "3", "first:time": "1700000000",
			"first:ticket_link": "https://tickets.example/ONE-1", "ignored": ""}
	})

	for _, protocolVersion := range []string{"2025-03-26", "2025-06-18", "2025-11-25"} {
		response := rpc(testingContext, router, "initialize", map[string]interface{}{
			"protocolVersion": protocolVersion, "capabilities": map[string]interface{}{},
			"clientInfo": map[string]string{"name": "test-client", "version": "1"},
		})
		var initialized sdk.InitializeResult
		require.NoError(testingContext, json.Unmarshal(response["result"], &initialized))
		require.Equal(testingContext, protocolVersion, initialized.ProtocolVersion)
	}

	listed := rpc(testingContext, router, "tools/list", map[string]interface{}{})
	var tools sdk.ListToolsResult
	require.NoError(testingContext, json.Unmarshal(listed["result"], &tools))
	require.Len(testingContext, tools.Tools, 1)
	require.Equal(testingContext, "get_error_logs", tools.Tools[0].Name)
	require.True(testingContext, tools.Tools[0].Annotations.ReadOnlyHint)
	require.Zero(testingContext, reads)

	result := callLogs(testingContext, router, map[string]interface{}{
		"group": "service", "query": "PROVIDER", "trace": "PROCESS", "app_name": "WORK",
		"include_request": true, "include_stack": true,
	})
	require.False(testingContext, result.IsError)
	encoded, err := json.Marshal(result.StructuredContent)
	require.NoError(testingContext, err)
	var output mcp.ErrorLogOutput
	require.NoError(testingContext, json.Unmarshal(encoded, &output))
	require.Equal(testingContext, 1, output.Total)
	require.Equal(testingContext, 3, output.Rows[0].Counter)
	require.Equal(testingContext, "payments", output.Rows[0].Service)
	require.Equal(testingContext, "worker.process:42", output.Rows[0].Stack)

	for _, sensitive := range []string{"private-header", "private-cookie", "private-password", "private-token"} {
		require.NotContains(testingContext, string(encoded), sensitive)
	}

	require.Contains(testingContext, string(encoded), "CA1234AA")
	require.Equal(testingContext, 1, reads)
}

func TestMCPAuthenticationAndOrigin(testingContext *testing.T) {
	router := newRouter(testingContext, func(context.Context, string) map[string]string {
		testingContext.Error("rejected request accessed the log store")

		return nil
	})

	for _, method := range []string{http.MethodPost, http.MethodGet, http.MethodDelete} {
		for _, authorization := range []string{"", "Bearer invalid", testToken, "Basic " + testToken} {
			request := httptest.NewRequest(method, mcp.Endpoint, nil)
			request.Header.Set("Authorization", authorization)

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)
			require.Equal(testingContext, http.StatusUnauthorized, recorder.Code)
			require.Equal(testingContext, `Bearer realm="mcp"`, recorder.Header().Get("WWW-Authenticate"))
		}
	}

	request := httptest.NewRequest(http.MethodPost, mcp.Endpoint, nil)
	request.Header.Set("Authorization", "Bearer "+testToken)
	request.Header.Set("Origin", "https://untrusted.example")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	require.Equal(testingContext, http.StatusForbidden, recorder.Code)
	require.Panics(testingContext, func() { mcp.NewHandler(mcp.NewServer("test", nil), "", nil) })

	request = httptest.NewRequest(http.MethodGet, mcp.Endpoint, nil)
	request.Header.Set("Authorization", "Bearer "+testToken)
	request.Header.Set("Accept", "text/event-stream")
	request.Header.Set("Origin", "https://allowed.example")

	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	require.Equal(testingContext, http.StatusMethodNotAllowed, recorder.Code)
}

func TestMCPLogSortingPaginationDefaultsAndValidation(testingContext *testing.T) {
	router := newRouter(testingContext, func(context.Context, string) map[string]string {
		return map[string]string{
			"older": `{"Message":"older","Request":"c2VjcmV0","Stack":"dHJhY2U="}`, "older:time": "100", "older:counter": "20",
			"newer": `{"Message":"newer"}`, "newer:time": "200", "newer:counter": "2",
		}
	})
	for _, item := range []struct {
		arguments map[string]interface{}
		firstID   string
	}{
		{map[string]interface{}{"page_size": 1}, "newer"},
		{map[string]interface{}{"page_size": 1, "sort_by": "counter"}, "older"},
		{map[string]interface{}{"page_size": 1, "page": 2}, "older"},
		{map[string]interface{}{"error_id": "older"}, "older"},
	} {
		result := callLogs(testingContext, router, item.arguments)
		require.False(testingContext, result.IsError)
		encoded, err := json.Marshal(result.StructuredContent)
		require.NoError(testingContext, err)
		var output mcp.ErrorLogOutput
		require.NoError(testingContext, json.Unmarshal(encoded, &output))
		require.Equal(testingContext, item.firstID, output.Rows[0].ID)
		require.Empty(testingContext, output.Rows[0].Request)
		require.Empty(testingContext, output.Rows[0].Stack)
	}

	for _, arguments := range []map[string]interface{}{
		{"group": "anything"}, {"page": -1}, {"page_size": 101}, {"sort_by": "invalid"},
	} {
		require.True(testingContext, callLogs(testingContext, router, arguments).IsError)
	}

	require.False(testingContext, callLogs(testingContext, router, map[string]interface{}{"page": 5}).IsError)
}

func TestMCPStorePanicReturnsToolError(testingContext *testing.T) {
	router := newRouter(testingContext, func(context.Context, string) map[string]string {
		panic("store unavailable; private credentials")
	})
	result := callLogs(testingContext, router, map[string]interface{}{})
	require.True(testingContext, result.IsError)
	require.NotContains(testingContext, result.Content[0].(*sdk.TextContent).Text, "private credentials")
}
