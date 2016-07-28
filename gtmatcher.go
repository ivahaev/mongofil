package mongofil

import "github.com/buger/jsonparser"

type GtMatcher struct {
	propName  string
	condition value
	eq        bool
}

func NewGtMatcher(propName string, in interface{}, eq bool) (Matcher, error) {
	val, err := createValue(in)
	if err != nil {
		return nil, err
	}
	m := GtMatcher{propName: propName, condition: val, eq: eq}
	return &m, nil
}

func (m *GtMatcher) Match(doc []byte) bool {
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
		if m.eq {
			return v >= conVal
		}
		return v > conVal
	case jsonparser.Number:
		if m.condition.typ != Number {
			return false
		}
		conVal, _ := m.condition.getFloat()
		v, _ := jsonparser.ParseFloat(val)
		if m.eq {
			return v >= conVal
		}
		return v > conVal
	}
	return false
}