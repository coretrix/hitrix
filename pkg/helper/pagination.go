package helper

type URLQueryPager struct {
	// example = ?current_page=1&page_size=25
	CurrentPage int `binding:"min=1" form:"current_page"`
	PageSize    int `binding:"min=1" form:"page_size"`
}
