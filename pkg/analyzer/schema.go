package analyzer

type Catalog struct {
	Tables map[string]*TableSchema
}

type TableSchema struct {
	Name    string
	Columns []*ColumnSchema
}

func (table *TableSchema) ToType() *StructType {
	fields := make([]*StructField, len(table.Columns))
	for i, c := range table.Columns {
		fields[i] = &StructField{
			Name: c.Name,
			Type: c.Type,
		}
	}
	return &StructType{Fields: fields}
}

type ColumnSchema struct {
	Name string
	Type Type
}
