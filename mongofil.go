package mongofil

import (
	"errors"
	"regexp"
)

// Errors
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
	nor []Matcher
}

func NewQuery(query map[string]interface{}) (*Query, error) {
	q := Query{and: []Matcher{}, or: []Matcher{}, nor: []Matcher{}}
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
				inm, err := NewInMatcher(k, arr, false)
				if err != nil {
					return nil, err
				}
				q.and = append(q.and, inm)
			}
			if val["$nin"] != nil {
				arr, ok := val["$nin"].([]interface{})
				if !ok {
					return nil, ErrInMustBeArray
				}
				inm, err := NewInMatcher(k, arr, true)
				if err != nil {
					return nil, err
				}
				q.and = append(q.and, inm)
			}
			if val["$exists"] != nil {
				em, err := NewExistsMatcher(k, val["$exists"])
				if err != nil {
					return nil, err
				}
				q.and = append(q.and, em)
			}
			if val["$gt"] != nil {
				em, err := NewGtMatcher(k, val["$gt"], false)
				if err != nil {
					return nil, err
				}
				q.and = append(q.and, em)
			}
			if val["$gte"] != nil {
				em, err := NewGtMatcher(k, val["$gte"], true)
				if err != nil {
					return nil, err
				}
				q.and = append(q.and, em)
			}
			if val["$lt"] != nil {
				em, err := NewLtMatcher(k, val["$lt"], false)
				if err != nil {
					return nil, err
				}
				q.and = append(q.and, em)
			}
			if val["$lte"] != nil {
				em, err := NewLtMatcher(k, val["$lte"], true)
				if err != nil {
					return nil, err
				}
				q.and = append(q.and, em)
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
		var matched bool
		for i := 0; i < len(q.or) && !matched; i++ {
			if q.or[i].Match(doc) {
				matched = true
			}
		}
		if !matched {
			return false
		}
	}
	if len(q.nor) != 0 {
		var matched bool
		for i := 0; i < len(q.nor) && !matched; i++ {
			if q.nor[i].Match(doc) {
				matched = true
			}
		}
		if matched {
			return false
		}
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
}
