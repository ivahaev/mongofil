package mongofil

import (
	"strings"

	"github.com/buger/jsonparser"
)

type eqMatcher struct {
	propName  string
	condition value
	invert    bool
}

func newEqMatcher(propName string, in interface{}, invert bool) (Matcher, error) {
	val, err := createValue(in)
	if err != nil {
		return nil, err
	}
	m := eqMatcher{propName: propName, condition: val, invert: invert}
	return &m, nil
}

func (m *eqMatcher) Match(doc []byte) bool {
	val, typ, _, err := jsonparser.Get(doc, m.propName)
	if err != nil {
		return m.invert
	}
	return m.matchValue(val, typ, m.propName)
}

func (m *eqMatcher) matchValue(val []byte, typ jsonparser.ValueType, propName string) bool {
	switch typ {
	case jsonparser.String:
		if m.condition.typ != String {
			return m.invert
		}
		conVal, _ := m.condition.getString()
		v, _ := jsonparser.ParseString(val)
		if conVal == v {
			return !m.invert
		}
	case jsonparser.Number:
		if m.condition.typ != Number {
			return m.invert
		}
		conVal, _ := m.condition.getFloat()
		v, _ := jsonparser.ParseFloat(val)
		if conVal == v {
			return !m.invert
		}
	case jsonparser.Array:
		keys := strings.Split(propName, ".")
		if len(keys) == 0 {
			keys = nil
		} else {
			keys = keys[1:]
		}
		var matched bool
		jsonparser.ArrayEach(val, func(value []byte, dataType jsonparser.ValueType, _ int, _ error) {
			if matched {
				return
			}
			matched = m.matchValue(value, dataType, strings.Join(keys, "."))
		}, keys...)
		return matched && !m.invert
	}
	return m.invert
}
