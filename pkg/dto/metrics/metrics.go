package metrics

//{
//"Memory" :  [{"Name": "admin-api", "Data": [{date, val}]}]
//}

type AppRMetrics struct {
	AppName uint64
	Rows    []*Row
}

type Row struct {
	Value     interface{}
	CreatedAt int64
}

type Series struct {
	Data       map[string][]Row
	XAxisTitle string
}
