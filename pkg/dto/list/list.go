package list

type RequestDTOList struct {
	Page     *int `binding:"required"`
	PageSize *int `binding:"required"`
	Filter   map[string]interface{}
	Search   map[string]interface{}
	SearchOR map[string]interface{}
	Sort     map[string]interface{}
}
