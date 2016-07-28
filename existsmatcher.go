package mongofil

import "github.com/buger/jsonparser"

type ExistsMatcher struct {
	propName  string
	condition bool
}

func NewExistsMatcher(propName string, in interface{}) (Matcher, error) {
	condition, ok := in.(bool)
	if !ok {
		return nil, errIsNotBool
	}
	m := ExistsMatcher{propName: propName, condition: condition}
	return &m, nil
}

func (m *ExistsMatcher) Match(doc []byte) bool {
	_, _, _, err := jsonparser.Get(doc, m.propName)

	if m.condition {
		return err == nil
	}
	return err == jsonparser.KeyPathNotFoundError
}
