package catalog

type Table struct {
	Name    string   `json:"name"`
	Columns []Column `json:"columns"`
}

type Column struct {
	Name string `json:"name"`
	//Type ColumnType `json:"type"`
	Type      ColumnType `json:"type"`
	Modifiers []Modifier `json:"modifiers"`
	//	Modifiers todo
}

type ColumnType string

const INTEGER ColumnType = "INTEGER"
const VARCHAR ColumnType = "VARCHAR"

var ValidColumnTypes = []ColumnType{INTEGER, VARCHAR}

type Modifier string

const PK Modifier = "PK"

var ValidModifiers = []Modifier{PK}
