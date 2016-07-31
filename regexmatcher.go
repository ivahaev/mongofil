package mongofil

import (
	"errors"
	"regexp"

	"github.com/buger/jsonparser"
)

var (
	errInvalidRegex = errors.New("invalid regex")
)

type regexMatcher struct {
	propName string
	regex    *regexp.Regexp
}

func newRegexMatcher(propName string, in interface{}) (Matcher, error) {
	condition, ok := in.(map[string]interface{})
	if !ok {
		return nil, errInvalidRegex
	}
	_expr, ok := condition["$regex"]
	if !ok {
		return nil, errInvalidRegex
	}
	expr, ok := _expr.(string)
	if !ok {
		return nil, errInvalidRegex
	}
	var options string
	if _opt, ok := condition["$options"]; ok {
		options, ok = _opt.(string)
		if !ok {
			return nil, errInvalidRegex
		}
	}
	var rgx *regexp.Regexp
	var err error
	if options == "" {
		rgx, err = regexp.Compile(expr)
	} else {
		rgx, err = regexp.Compile("(?" + options + ")" + expr)
	}
	if err != nil {
		return nil, err
	}
	m := regexMatcher{propName, rgx}
	return &m, nil
}

func (m *regexMatcher) Match(doc []byte) bool {
	val, typ, _, err := jsonparser.Get(doc, m.propName)
	if err != nil || typ != jsonparser.String {
		return false
	}
	return m.regex.Match(val)
}
