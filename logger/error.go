package logger

// RowError is an error that is logged when reading an erroneous row from the
// provide CSV file.
type RowError struct {
	Row      int
	ID       string
	SKU      string
	Handle   string
	ImageURL string
	Reason   error
}
