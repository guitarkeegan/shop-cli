package fulltui

import "fmt"

type param struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Required    bool   `json:"-"`
}

func (p param) String() string {
	return fmt.Sprintf("\nName: %s\nType: %s\nDescription: %s\nRequired: %t\n", p.Name, p.Type, p.Description, p.Required)
}

type snippetDoc struct {
	Description string  `json:"description"`
	Example     string  `json:"example"`
	Param       []param `json:"params"`
}

func (sd *snippetDoc) String() string {
	return fmt.Sprintf("Description: %s, Example: %s, Param: %+v", sd.Description, sd.Example, sd.Param)
}
