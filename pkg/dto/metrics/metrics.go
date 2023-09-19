package delayedqueue

type Row struct {
	Queue string
	Total int64
}

type List struct {
	Rows []Row
}
