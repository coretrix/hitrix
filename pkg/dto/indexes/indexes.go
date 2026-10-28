package indexes

type ResponseDTOList struct {
	Indexes []Index
}

type Index struct {
	Name             string
	Pool             string
	CanReindexEntity bool
	TotalDocs        uint64
	TotalSize        uint64
}
