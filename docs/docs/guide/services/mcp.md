# Read-only MCP endpoint

Register commands in the application's entrypoint before `Build()`:

```go
import (
    hitrixMCP "github.com/coretrix/hitrix/pkg/mcp"
    sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// Add to the existing Hitrix registry chain:
registry.RegisterMCPCommands([]hitrixMCP.Command{
    hitrixMCP.RegisterErrorLog,
    RegisterProjectLookup,
})
```

Each command is a function accepting `*sdk.Server`. Use `sdk.AddTool` to bind
typed inputs to a handler. Only register handlers that read data, and set
`Annotations: &sdk.ToolAnnotations{ReadOnlyHint: true}`. A project's handlers
can use `service.DI().OrmEngineForContext(ctx)` for the request's ORM engine.
The built-in command is optional; each application chooses its own command list.

Hitrix registers a single `/dev/mcp` endpoint during `InitGin`. It uses the
official MCP Go SDK, with stateless Streamable HTTP and JSON responses. It supports
MCP initialization, tool discovery, and tool calls without needing session storage.
GET streaming is unnecessary for these request/response tools and returns 405.

Set `mcp.auth_token` to a generated secret in the application's Hitrix config.
Every MCP request requires `Authorization: Bearer <configured token>`, including
requests in local mode. Missing/empty configuration prevents server startup;
missing or incorrect credentials return 401. This token is independent of
dev-panel JWTs and login cookies. No auth credentials belong in the open-source
Hitrix repository. If a client sends an Origin, it must exactly match an entry
in the application's `cors` configuration.

## Error-log tool

`get_error_logs` reads the persistent Redis pool registered with `RegisterRedisPools`,
falling back to the default Redis pool. It uses the same log decoding as the
dev-panel error-log API. It does not delete logs or create Jira tickets.

Arguments:

| Argument | Behavior |
| --- | --- |
| `group` | `error` (default), `service`, `warning`, or `missing-translation` |
| `error_id` | Optional exact entry ID |
| `app_name` | Optional case-insensitive app-name substring |
| `query` | Optional case-insensitive search across entry fields |
| `trace` | Optional case-insensitive stack substring |
| `sort_by` | `time` (default) or `counter`, descending |
| `include_request` / `include_stack` | Default false; opt into captured content |
| `page` / `page_size` | Default 1 / 50; page size is capped at 100 |

The result includes `Group`, `Rows`, filtered `Total`, `Page`, and `PageSize`.
Rows retain the dev-panel's field names and add the entry `ID`. Common credentials
are masked in messages and included request/stack content. The current Redis hash
is read and filtered in memory, so pagination bounds the response rather than the
Redis read. The selected endpoint URL determines the environment.

For example, tool arguments to inspect recurring service errors:

```json
{"group":"service","sort_by":"counter","page_size":10,"query":"timeout"}
```

Tests in `tests/mcp` exercise the endpoint through `InitGin`, using in-memory
readers and temporary configuration. They require no database or cache service.
