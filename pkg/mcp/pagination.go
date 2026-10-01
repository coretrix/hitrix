package mcp

import "fmt"

// Pagination supplies bounded defaults for read-only commands.
func Pagination(page, pageSize int) (int, int, error) {
	if page == 0 {
		page = 1
	}

	if pageSize == 0 {
		pageSize = 50
	}

	if page < 1 || page > 1000000 || pageSize < 1 || pageSize > 100 {
		return 0, 0, fmt.Errorf("page must be between 1 and 1000000 and page_size between 1 and 100")
	}

	return page, pageSize, nil
}
