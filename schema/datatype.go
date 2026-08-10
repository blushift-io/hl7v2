package schema

// DataTypeType represents the classification of a data type or component.
type DataTypeType int

const (
	// DataTypeTypeDataType indicates a standalone HL7 data type.
	DataTypeTypeDataType DataTypeType = iota //DataType
	// DataTypeTypeComponent indicates a composite component field.
	DataTypeTypeComponent //Component
)

// DataTypes represents a slice of DataType pointers.
type DataTypes []*DataType

// Len returns the number of data types in the collection.
func (ds DataTypes) Len() int {
	return len(ds)
}

// Less reports whether the data type at index i sorts before the data type at index j by ID.
func (ds DataTypes) Less(i, j int) bool {
	return ds[i].ID < ds[j].ID
}

// Swap exchanges the elements at indices i and j.
func (ds DataTypes) Swap(i, j int) {
	ds[i], ds[j] = ds[j], ds[i]
}

// DataType represents an HL7 data type definition.
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

// Required reports whether the data type must be present.
func (d *DataType) Required() bool {
	return d.MinRepetitions > 0
}

// Repeatable reports whether the data type can occur multiple times.
func (d *DataType) Repeatable() bool {
	return d.MaxRepetitions != 1
}

// Table returns the associated table for the data type, if any.
func (d *DataType) Table() *Table {
	if d.s == nil {
		return nil
	}

	return d.s.Table(d.TableID)
}

// IsPrimitive reports whether the data type has no sub-fields.
func (d *DataType) IsPrimitive() bool {
	return len(d.Fields) == 0
}

// PrimitiveType returns the underlying Go or primitive type name corresponding to the HL7 data type.
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

// GetDataType retrieves the full DataType schema object from the parent schema.
func (d *DataType) GetDataType() *DataType {
	if d.s == nil {
		return nil
	}

	return d.s.DataType(d.DataType)
}

// GetFields retrieves the child field DataType definitions for composite data types.
func (d *DataType) GetFields() []*DataType {
	var fs []*DataType
	for _, f := range d.Fields {
		fs = append(fs, d.s.DataType(f.ID))
	}

	return fs
}
