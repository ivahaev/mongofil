package mongofil

import "github.com/buger/jsonparser"

type inMatcher struct {
	propName   string
	conditions []value
	strVal     *string
	numVal     *float64
	boolVal    *bool
	invert     bool
}

func newInMatcher(propName string, arr []interface{}, invert bool) (Matcher, error) {
	m := inMatcher{propName: propName, conditions: make([]value, len(arr)), invert: invert}
	for i, v := range arr {
		val, err := createValue(v)
		if err != nil {
			return nil, err
		}
		m.conditions[i] = val
	}
	return &m, nil
}

func (m *inMatcher) Match(doc []byte) bool {
	defer func() {
		m.strVal = nil
		m.numVal = nil
		m.boolVal = nil
	}()
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
				return !m.invert
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
				return !m.invert
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
				return !m.invert
			}
		}
	}
	return m.invert
}
