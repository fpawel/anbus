package types

type ReadVar struct {
	Place,
	VarIndex int
	Value float64
	Error string
}

type EmptyRecord struct{}
