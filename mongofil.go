package mongofil

import (
	"errors"
	"regexp"

	"github.com/buger/jsonparser"
)

var (
	ErrInMustBeArray = errors.New("$in must points to array of interface{}")
)

type Q struct {
	In        []interface{}
	Nin       []interface{}
	Exists    bool
	Gte       float64
	Gt        float64
	Lte       float64
	Lt        float64
	Eq        interface{}
	Ne        interface{}
	Mod       []interface{}
	All       []interface{}
	And       []interface{}
	Or        []interface{}
	Nor       []interface{}
	Size      int
	Regex     *regexp.Regexp
	Where     interface{}
	ElemMatch interface{}
	Not       interface{}
}

type Query struct {
	and []Matcher
	or  []Matcher
}

func NewQuery(query map[string]interface{}) (*Query, error) {
	q := Query{and: []Matcher{}, or: []Matcher{}}
	for k, v := range query {
		switch v.(type) {
		case string, float64, bool:
			em, err := NewEqMatcher(k, v)
			if err != nil {
				return nil, err
			}
			q.and = append(q.and, em)
		case map[string]interface{}:
			val := v.(map[string]interface{})
			if val["$in"] != nil {
				arr, ok := val["$in"].([]interface{})
				if !ok {
					return nil, ErrInMustBeArray
				}
				inm, err := NewInMatcher(k, arr)
				if err != nil {
					return nil, err
				}
				q.and = append(q.and, inm)
			}
		}
	}
	return &q, nil
}

func (q *Query) Match(doc []byte) bool {
	for _, _q := range q.and {
		if !_q.Match(doc) {
			return false
		}
	}
	if len(q.or) != 0 {
		for _, _q := range q.or {
			if _q.Match(doc) {
				return true
			}
		}
		return false
	}
	return true
}

// Match returns true if JSON matched query
func Match(query map[string]interface{}, doc []byte) (bool, error) {
	q, err := NewQuery(query)
	if err != nil {
		return false, err
	}
	return q.Match(doc), nil
	for k, _v := range query {
		val, typ, _, err := jsonparser.Get(doc, k)
		if err != nil {
			continue
		}
		switch typ {
		case jsonparser.String:
			v, ok := _v.(string)
			if !ok {
				continue
			}
			return string(val) == v, nil
		}
	}
	return false, nil
}
