package mongofil

import "github.com/buger/jsonparser"

type InMatcher struct {
	propName   string
	conditions []value
	strVal     *string
	numVal     *float64
	boolVal    *bool
}

func NewInMatcher(propName string, arr []interface{}) (Matcher, error) {
	m := InMatcher{propName: propName, conditions: make([]value, len(arr))}
	for i, v := range arr {
		val, err := createValue(v)
		if err != nil {
			return nil, err
		}
		m.conditions[i] = val
	}
	return &m, nil
}

func (m *InMatcher) Match(doc []byte) bool {
	m.strVal = nil
	m.numVal = nil
	m.boolVal = nil
	val, typ, _, err := jsonparser.Get(doc, m.propName)
	if err != nil {
		return false
	}
	for _, _v := range m.conditions {
		switch typ {
		case jsonparser.String:
			v, err := _v.getString()
			if err != nil {
				continue
			}
			if m.strVal == nil {
				s, _ := jsonparser.ParseString(val)
				m.strVal = &s
			}
			if *m.strVal == v {
				return true
			}
		case jsonparser.Number:
			v, err := _v.getFloat()
			if err != nil {
				continue
			}
			if m.numVal == nil {
				f, _ := jsonparser.ParseFloat(val)
				m.numVal = &f
			}
			if *m.numVal == v {
				return true
			}
		case jsonparser.Boolean:
			v, err := _v.getBool()
			if err != nil {
				continue
			}
			if m.numVal == nil {
				b, _ := jsonparser.ParseBoolean(val)
				m.boolVal = &b
			}
			if *m.boolVal == v {
				return true
			}
		}
	}
	return false
}
