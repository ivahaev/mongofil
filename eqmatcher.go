package mongofil

import "github.com/buger/jsonparser"

type EqMatcher struct {
	propName  string
	condition value
}

func NewEqMatcher(propName string, in interface{}) (Matcher, error) {
	val, err := createValue(in)
	if err != nil {
		return nil, err
	}
	m := EqMatcher{propName: propName, condition: val}
	return &m, nil
}

func (m *EqMatcher) Match(doc []byte) bool {
	val, typ, _, err := jsonparser.Get(doc, m.propName)
	if err != nil {
		return false
	}
	switch typ {
	case jsonparser.String:
		if m.condition.typ != String {
			return false
		}
		conVal, _ := m.condition.getString()
		v, _ := jsonparser.ParseString(val)
		return conVal == v
	}
	return false
}
