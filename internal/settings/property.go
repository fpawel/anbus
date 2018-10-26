package settings

type Property struct {
	Hint, Name,
	Value,
	DefaultValue string
	Min, Max  *ValueEx
	ValueType ValueType
	List      []string
}

type ValueEx struct {
	Value float64
}
