package file

type Image struct {
	ID        uint64 `json:",omitempty"`
	URL       string
	Namespace string
	Primary   bool
	Hidden    bool
}
