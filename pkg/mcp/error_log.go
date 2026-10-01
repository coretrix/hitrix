package mcp

import (
	"context"
	"fmt"
	"sort"
	"strings"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/coretrix/hitrix/service"
	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
)

type ErrorLogInput struct {
	Group          string `json:"group,omitempty" jsonschema:"Log group: error (default), service, warning, or missing-translation"`
	ErrorID        string `json:"error_id,omitempty" jsonschema:"Exact stored entry ID for reading one error in detail"`
	AppName        string `json:"app_name,omitempty" jsonschema:"Case-insensitive app name substring filter"`
	Query          string `json:"query,omitempty" jsonschema:"Search ID, message, file, request, stack, app, service, time, and ticket link"`
	Trace          string `json:"trace,omitempty" jsonschema:"Case-insensitive stack trace substring filter"`
	SortBy         string `json:"sort_by,omitempty" jsonschema:"Sort by time (default) for recent errors or counter for most frequent errors"`
	IncludeRequest bool   `json:"include_request,omitempty" jsonschema:"Include captured requests with common credentials masked; defaults to false"`
	IncludeStack   bool   `json:"include_stack,omitempty" jsonschema:"Include stack traces with common credentials masked; defaults to false"`
	Page           int    `json:"page,omitempty" jsonschema:"Page number starting at 1; defaults to 1"`
	PageSize       int    `json:"page_size,omitempty" jsonschema:"Entries per page, from 1 to 100; defaults to 50"`
}

type ErrorLogEntry struct {
	ID string
	*errorlogger.Entry
}

type ErrorLogOutput struct {
	Group    string
	Rows     []ErrorLogEntry
	Total    int
	Page     int
	PageSize int
}

// RegisterErrorLog adds the built-in read-only log tool.
func RegisterErrorLog(server *sdk.Server) {
	NewErrorLogCommand(func(ctx context.Context, group string) map[string]string {
		ormService := service.DI().OrmEngineForContext(ctx)
		appService := service.DI().App()

		if appService.RedisPools != nil && appService.RedisPools.Persistent != "" {
			return ormService.GetRedis(appService.RedisPools.Persistent).HGetAll(group)
		}

		return ormService.GetRedis().HGetAll(group)
	})(server)
}

// NewErrorLogCommand permits isolated tests or an alternative read-only log store.
func NewErrorLogCommand(read func(context.Context, string) map[string]string) Command {
	return func(server *sdk.Server) {
		sdk.AddTool(server, &sdk.Tool{
			Name: "get_error_logs",
			Description: "Read and search stored Hitrix error logs, service errors, warnings, or missing translations. " +
				"Returns a page of entries with IDs, messages, source files and lines, stack traces, " +
				"captured requests, app/service names, occurrence counters, last occurrence times, and Jira links. " +
				"Filter by error_id for one entry, app_name, query, or trace. Sort by time for recent errors or " +
				"counter for top recurring errors. Request and stack content is returned only when requested. " +
				"Reads this endpoint's environment directly; select the demo or production URL in your MCP client. " +
				"This tool never removes logs or creates tickets.",
			Annotations: &sdk.ToolAnnotations{ReadOnlyHint: true},
		}, func(ctx context.Context, _ *sdk.CallToolRequest, input ErrorLogInput) (*sdk.CallToolResult, ErrorLogOutput, error) {
			return readErrorLogs(ctx, input, read)
		})
	}
}

func readErrorLogs(ctx context.Context, input ErrorLogInput,
	read func(context.Context, string) map[string]string) (*sdk.CallToolResult, ErrorLogOutput, error) {
	if input.Group == "" {
		input.Group = errorlogger.GroupError
	}

	switch input.Group {
	case errorlogger.GroupError, errorlogger.GroupService, errorlogger.GroupWarning, errorlogger.GroupMissingTranslation:
	default:
		return nil, ErrorLogOutput{}, fmt.Errorf("invalid log group")
	}

	page, pageSize, err := Pagination(input.Page, input.PageSize)
	if err != nil {
		return nil, ErrorLogOutput{}, err
	}

	if input.SortBy != "" && input.SortBy != "time" && input.SortBy != "counter" {
		return nil, ErrorLogOutput{}, fmt.Errorf("sort_by must be time or counter")
	}

	entries, err := errorlogger.ParseEntries(read(ctx, input.Group))
	if err != nil {
		return nil, ErrorLogOutput{}, fmt.Errorf("could not decode stored error logs")
	}

	rows := make([]ErrorLogEntry, 0, len(entries))

	for entryID, entry := range entries {
		if !matchesError(entryID, entry, input) {
			continue
		}

		rows = append(rows, ErrorLogEntry{ID: entryID, Entry: entry})
	}

	sort.Slice(rows, func(first, second int) bool {
		if input.SortBy == "counter" && rows[first].Counter != rows[second].Counter {
			return rows[first].Counter > rows[second].Counter
		}

		if rows[first].Timestamp != rows[second].Timestamp {
			return rows[first].Timestamp > rows[second].Timestamp
		}

		return rows[first].ID < rows[second].ID
	})

	output := ErrorLogOutput{Group: input.Group, Total: len(rows), Page: page, PageSize: pageSize, Rows: []ErrorLogEntry{}}
	start := (page - 1) * pageSize

	if start < len(rows) {
		end := start + pageSize
		if end > len(rows) {
			end = len(rows)
		}

		output.Rows = rows[start:end]
	}

	for _, row := range output.Rows {
		row.Message = maskSensitiveText(row.Message)
		if input.IncludeRequest {
			row.Request = maskSensitiveText(row.Request)
		} else {
			row.Request = ""
		}

		if input.IncludeStack {
			row.Stack = maskSensitiveText(row.Stack)
		} else {
			row.Stack = ""
		}
	}

	return nil, output, nil
}

func matchesError(entryID string, entry *errorlogger.Entry, input ErrorLogInput) bool {
	contains := func(value, query string) bool {
		return strings.Contains(strings.ToLower(value), strings.ToLower(query))
	}

	if input.ErrorID != "" && input.ErrorID != entryID {
		return false
	}

	if !contains(entry.AppName, input.AppName) || !contains(entry.Stack, input.Trace) {
		return false
	}

	for _, value := range []string{entryID, entry.Message, entry.File, entry.Request, entry.Stack,
		entry.AppName, entry.Service, entry.Time, entry.TicketLink} {
		if contains(value, input.Query) {
			return true
		}
	}

	return false
}
