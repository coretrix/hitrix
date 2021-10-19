# Pagination
You can use:
```go
package helper

type URLQueryPager struct {
	// example = ?current_page=1&page_size=25
	CurrentPage int `binding:"min=1" form:"current_page"`
	PageSize    int `binding:"min=1" form:"page_size"`
}
```
in your code that needs pagination like:

```go
package mypackage

import "github.com/coretrix/hitrix/pkg/helper"

type SomeURLQuery struct {
	helper.URLQueryPager
	OtherField1 string `form:"other_field_1"`
	OtherField2 int `form:"other_field_2"`
}
```