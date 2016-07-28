package mongofil

import "errors"

// Errors
var (
	ErrInMustBeArray  = errors.New("$in must points to array of interface{}")
	ErrAndMustBeArray = errors.New("$and must points to array of map[string]interface{}")
	ErrOrMustBeArray  = errors.New("$or must points to array of map[string]interface{}")
	ErrNotMustBeMap   = errors.New("$not must points to map[string]interface{}")
)

type Query struct {
	and []Matcher
	or  []Matcher
	nor []Matcher
	not []Matcher
}

func NewQuery(query map[string]interface{}) (*Query, error) {
	q := Query{and: []Matcher{}, or: []Matcher{}, nor: []Matcher{}, not: []Matcher{}}
	for field, v := range query {
		switch field {
		case "$and":
			val, ok := v.([]interface{})
			if !ok {
				return nil, ErrAndMustBeArray
			}
			for _, subv := range val {
				subval, ok := subv.(map[string]interface{})
				if !ok {
					return nil, ErrAndMustBeArray
				}
				m, err := NewQuery(subval)
				if err != nil {
					return nil, err
				}
				q.and = append(q.and, m)
			}
		case "$or":
			val, ok := v.([]interface{})
			if !ok {
				return nil, ErrAndMustBeArray
			}
			for _, subv := range val {
				subval, ok := subv.(map[string]interface{})
				if !ok {
					return nil, ErrAndMustBeArray
				}
				m, err := NewQuery(subval)
				if err != nil {
					return nil, err
				}
				q.or = append(q.or, m)
			}
		default:
			switch v.(type) {
			case string, float64, bool:
				em, err := NewEqMatcher(field, v, false)
				if err != nil {
					return nil, err
				}
				q.and = append(q.and, em)
			case map[string]interface{}:
				val := v.(map[string]interface{})
				if val["$ne"] != nil {
					em, err := NewEqMatcher(field, val["$ne"], true)
					if err != nil {
						return nil, err
					}
					q.and = append(q.and, em)
				}
				if val["$not"] != nil {
					notMap, ok := val["$not"].(map[string]interface{})
					if !ok {
						return nil, ErrNotMustBeMap
					}
					nm, err := NewQuery(map[string]interface{}{field: notMap})
					if err != nil {
						return nil, err
					}
					q.not = append(q.not, nm)
				}
				if val["$in"] != nil {
					arr, ok := val["$in"].([]interface{})
					if !ok {
						return nil, ErrInMustBeArray
					}
					inm, err := NewInMatcher(field, arr, false)
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
					ninm, err := NewInMatcher(field, arr, true)
					if err != nil {
						return nil, err
					}
					q.and = append(q.and, ninm)
				}
				if val["$exists"] != nil {
					em, err := NewExistsMatcher(field, val["$exists"])
					if err != nil {
						return nil, err
					}
					q.and = append(q.and, em)
				}
				if val["$gt"] != nil {
					em, err := NewGtMatcher(field, val["$gt"], false)
					if err != nil {
						return nil, err
					}
					q.and = append(q.and, em)
				}
				if val["$gte"] != nil {
					em, err := NewGtMatcher(field, val["$gte"], true)
					if err != nil {
						return nil, err
					}
					q.and = append(q.and, em)
				}
				if val["$lt"] != nil {
					em, err := NewLtMatcher(field, val["$lt"], false)
					if err != nil {
						return nil, err
					}
					q.and = append(q.and, em)
				}
				if val["$lte"] != nil {
					em, err := NewLtMatcher(field, val["$lte"], true)
					if err != nil {
						return nil, err
					}
					q.and = append(q.and, em)
				}
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
	if len(q.not) != 0 {
		var matched bool
		for i := 0; i < len(q.not) && !matched; i++ {
			if q.not[i].Match(doc) {
				return false
			}
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
