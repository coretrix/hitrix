package delayedqueue

type Row struct {
	Queue      string
	Total      int64
	LatestItem *uint64
}

type List struct {
	Rows []Row
}
