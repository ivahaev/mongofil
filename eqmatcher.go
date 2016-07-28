package mongofil

import "github.com/buger/jsonparser"

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
	}
	return m.invert
}
