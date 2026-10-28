package indexes

// Keep IDs as decimal strings to preserve uint64 precision in JavaScript.
type RequestDTOReindexEntity struct {
	ID string `binding:"required"`
}
