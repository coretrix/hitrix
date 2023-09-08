package metrics

type ResponseDTORMetrics struct {
	Rows []*Row
}

type Row struct {
	ID        uint64
	CreatedAt int64
}
