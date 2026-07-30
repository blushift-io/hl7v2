package schema

type DataTypeType int

const (
	DataTypeTypeDataType  DataTypeType = iota //DataType
	DataTypeTypeComponent                     //Component
)

type DataTypes []*DataType

func (ds DataTypes) Len() int {
	return len(ds)
}

func (ds DataTypes) Less(i, j int) bool {
	return ds[i].ID < ds[j].ID
}

func (ds DataTypes) Swap(i, j int) {
	ds[i], ds[j] = ds[j], ds[i]
}

type DataType struct {
	s *Schema

	ID             string       `json:"id"`
	Type           DataTypeType `json:"type"`
	Name           string       `json:"name"`
	Description    string       `json:"description"`
	DataType       string       `json:"data_type"`
	DataTypeName   string       `json:"data_type_name"`
	Length         int          `json:"length"`
	MaxRepetitions int          `json:"max_repetitions"`
	MinRepetitions int          `json:"min_repetitions"`
	TableID        string       `json:"table_id"`
	TableName      string       `json:"table_name"`
	Sample         string       `json:"sample"`
	Fields         []DataType   `json:"fields"`
	Position       string       `json:"position"`
}

func (d *DataType) Required() bool {
	return d.MinRepetitions > 0
}

func (d *DataType) Repeatable() bool {
	return d.MaxRepetitions != 1
}

func (d *DataType) Table() *Table {
	if d.s == nil {
		return nil
	}

	return d.s.Table(d.TableID)
}

func (d *DataType) IsPrimitive() bool {
	return len(d.Fields) == 0
}

func (d *DataType) PrimitiveType() string {
	switch d.DataType {
	case "ST", "TX", "FT", "ID", "IS", "TN":
		return "string"
	case "TM", "DT", "TS":
		return "time.Time"
	case "NM":
		return "float64"
	case "SI":
		return "int64"
	default:
		return ""
	}
}

func (d *DataType) GetDataType() *DataType {
	if d.s == nil {
		return nil
	}

	return d.s.DataType(d.DataType)
}

func (d *DataType) GetFields() []*DataType {
	var fs []*DataType
	for _, f := range d.Fields {
		fs = append(fs, d.s.DataType(f.ID))
	}

	return fs
}
