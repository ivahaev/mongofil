package mongofil

import "github.com/buger/jsonparser"

type existsMatcher struct {
	propName  string
	condition bool
}

func newExistsMatcher(propName string, in interface{}) (Matcher, error) {
	condition, ok := in.(bool)
	if !ok {
		return nil, errIsNotBool
	}
	m := existsMatcher{propName: propName, condition: condition}
	return &m, nil
}

func (m *existsMatcher) Match(doc []byte) bool {
	_, _, _, err := jsonparser.Get(doc, m.propName)

	if m.condition {
		return err == nil
	}
	return err == jsonparser.KeyPathNotFoundError
}
