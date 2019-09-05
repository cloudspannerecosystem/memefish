package analyzer

type Catalog struct {
	Tables map[string]*TableSchema
}

type TableSchema struct {
	Columns []*ColumnSchema
}

type ColumnSchema struct {
	Name string
	Type Type
}
