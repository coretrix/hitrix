// Package mcp exposes opt-in, read-only tools over a single Streamable HTTP endpoint.
package mcp

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/coretrix/hitrix/service"
)

const Endpoint = "/dev/mcp"

// Command registers tools on the shared server. Commands must only read data.
type Command = func(*sdk.Server)

// NewServer lets projects combine Hitrix commands with their own commands.
func NewServer(name string, commands []Command) *sdk.Server {
	server := sdk.NewServer(&sdk.Implementation{Name: name, Version: "1.0.0"}, &sdk.ServerOptions{
		Instructions: "All tools are read-only. Read application data and logs without modifying records, " +
			"deleting logs, reindexing, deploying, or running arbitrary SQL.",
	})
	server.AddReceivingMiddleware(recoverToolPanic)

	for _, command := range commands {
		command(server)
	}

	return server
}

// NewHandler requires a bearer token even in local mode. Empty tokens fail closed.
func NewHandler(server *sdk.Server, token string, allowedOrigins []string) http.Handler {
	if strings.TrimSpace(token) == "" {
		panic("mcp.auth_token must not be empty")
	}

	tokenHash := sha256.Sum256([]byte(token))
	transport := sdk.NewStreamableHTTPHandler(func(_ *http.Request) *sdk.Server {
		return server
	}, &sdk.StreamableHTTPOptions{Stateless: true, JSONResponse: true})

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		authorization := strings.Fields(request.Header.Get("Authorization"))
		if len(authorization) != 2 || !strings.EqualFold(authorization[0], "Bearer") {
			unauthorized(writer)

			return
		}

		providedHash := sha256.Sum256([]byte(authorization[1]))
		if subtle.ConstantTimeCompare(providedHash[:], tokenHash[:]) != 1 {
			unauthorized(writer)

			return
		}

		if origin := request.Header.Get("Origin"); origin != "" && !originAllowed(origin, allowedOrigins) {
			http.Error(writer, "origin is not allowed", http.StatusForbidden)

			return
		}

		transport.ServeHTTP(writer, request)
	})
}

func unauthorized(writer http.ResponseWriter) {
	writer.Header().Set("WWW-Authenticate", `Bearer realm="mcp"`)
	http.Error(writer, "unauthorized", http.StatusUnauthorized)
}

func originAllowed(origin string, allowedOrigins []string) bool {
	for _, allowedOrigin := range allowedOrigins {
		if origin == allowedOrigin {
			return true
		}
	}

	return false
}

// RegisterRoutes keeps MCP separate from dev-panel JWT/cookie authentication.
func RegisterRoutes(ginEngine *gin.Engine) {
	appService := service.DI().App()
	if len(appService.MCPCommands) == 0 {
		return
	}

	configService := service.DI().Config()
	allowedOrigins, _ := configService.Strings("cors")
	handler := NewHandler(NewServer(appService.Name, appService.MCPCommands),
		configService.MustString("mcp.auth_token"), allowedOrigins)
	ginEngine.Any(Endpoint, gin.WrapH(handler))
}

// ORM operations panic on failures; SDK handlers execute outside Gin's recovery.
func recoverToolPanic(next sdk.MethodHandler) sdk.MethodHandler {
	return func(ctx context.Context, method string, request sdk.Request) (result sdk.Result, err error) {
		defer func() {
			if recovered := recover(); recovered != nil {
				// The data store may be down; avoid logging back into the same Redis store.
				slog.ErrorContext(ctx, "MCP request failed", "method", method)

				if method == "tools/call" {
					toolResult := &sdk.CallToolResult{}
					toolResult.SetError(fmt.Errorf("tool could not read data"))
					result, err = toolResult, nil
				} else {
					result, err = nil, fmt.Errorf("MCP request failed")
				}
			}
		}()

		return next(ctx, method, request)
	}
}
