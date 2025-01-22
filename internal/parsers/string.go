package parsers

type String struct{}

func (String) Name() string { return "string" }

func (String) Parse(s string) (string, error) { return s, nil }
