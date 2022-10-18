package list

type RequestDTOList struct {
	Page     *int `binding:"required"`
	PageSize *int `binding:"required"`
	Search   map[string]interface{}
	SearchOR map[string]interface{}
	Sort     map[string]interface{}
}
