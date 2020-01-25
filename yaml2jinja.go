package yaml2jinja

// New creates Yaml2Jinja object
func New() Yaml2Jinja {
	return Yaml2Jinja{}
}

type line struct {
	structName string
	line       string
}

// Yaml2Jinja to store converted result
type Yaml2Jinja struct {
	Visited   map[line]bool
	StructMap map[string]string
}
