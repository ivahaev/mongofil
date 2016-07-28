package mongofil

import "errors"

var (
	errIsNotString      = errors.New("value is not a string")
	errIsNotNumber      = errors.New("value is not a number")
	errIsNotBool        = errors.New("value is not a boolean")
	errValueIsNotExists = errors.New("value is not exists")
	errUnknownValueType = errors.New("unknown value type")
)

// ValueType is representation of data types available in valid JSON data.
type ValueType int

const (
	NotExist = ValueType(iota)
	String
	Number
	Object
	Array
	Boolean
	Null
	Unknown

	RegExp = 10
)

type value struct {
	typ     ValueType
	strVal  *string
	numVal  *float64
	boolVal *bool
}

func createValue(in interface{}) (value, error) {
	res := value{}
	if v, ok := in.(string); ok {
		res.typ = String
		res.strVal = &v
	} else if v, ok := in.(float64); ok {
		res.typ = Number
		res.numVal = &v
	} else if v, ok := in.(int); ok {
		res.typ = Number
		f := float64(v)
		res.numVal = &f
	} else if v, ok := in.(bool); ok {
		res.typ = Boolean
		res.boolVal = &v
	}
	if res.typ == 0 {
		return res, errUnknownValueType
	}
	return res, nil
}

func (v value) getString() (string, error) {
	if v.typ == 0 {
		return "", errValueIsNotExists
	}
	if v.typ != String {
		return "", errIsNotString
	}
	return *v.strVal, nil
}

func (v value) getFloat() (float64, error) {
	if v.typ == 0 {
		return 0, errValueIsNotExists
	}
	if v.typ != Number {
		return 0, errIsNotNumber
	}
	return *v.numVal, nil
}

func (v value) getBool() (bool, error) {
	if v.typ == 0 {
		return false, errValueIsNotExists
	}
	if v.typ != Number {
		return false, errIsNotBool
	}
	return *v.boolVal, nil
}
