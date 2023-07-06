package indexes

type ResponseDTOList struct {
	Indexes []Index
}

type Index struct {
	Name      string
	TotalDocs uint64
	TotalSize uint64
}
