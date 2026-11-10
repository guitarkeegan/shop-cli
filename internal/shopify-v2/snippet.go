package shopifyv2

type SnippetParamType int

const (
	StringParam SnippetParamType = iota
	NumberParam
	BooleanParam
	ObjectParam
)

var snippetParamTypeNames = map[SnippetParamType]string{
	StringParam:  "string",
	NumberParam:  "number",
	BooleanParam: "boolean",
	ObjectParam:  "object",
}

var stringToSnippetParamType = map[string]SnippetParamType{
	"string":  StringParam,
	"number":  NumberParam,
	"boolean": BooleanParam,
	"object":  ObjectParam,
}

func (s SnippetParamType) String() string {
	if name, ok := snippetParamTypeNames[s]; ok {
		return name
	}
	return "string"
}

// ParamTypeFromString converts a string to a SnippetParamType.
// Returns StringParam for unrecognized values.
func ParamTypeFromString(s string) SnippetParamType {
	if pt, ok := stringToSnippetParamType[s]; ok {
		return pt
	}
	return StringParam
}

// ParamTypes is the list of valid snippet parameter type strings.
var ParamTypes = []string{"string", "number", "boolean", "object"}

// ParamTypeToString maps SnippetParamType int values to their string names.
var ParamTypeToString = map[int]string{
	int(StringParam):  "string",
	int(NumberParam):  "number",
	int(BooleanParam): "boolean",
	int(ObjectParam):  "object",
}

type Param struct {
	Type        SnippetParamType
	Name        string
	Required    bool
	Description string
}

type Snippet struct {
	DocText     string
	Params      []Param
	Examples    []string
	Description string
}
