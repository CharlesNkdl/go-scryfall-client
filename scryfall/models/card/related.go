package card

type Related struct {
	Id        string    `json:"id"`
	Object    string    `json:"object"`
	Component Component `json:"component"`
	Name      string    `json:"name"`
	TypeLine  string    `json:"type_line"`
	Uri       string    `json:"uri"`
}
