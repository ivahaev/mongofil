package mongofil

import "errors"

// Errors
var (
	ErrEmptyFieldName = errors.New("empty field name")
	ErrInMustBeArray  = errors.New("$in must points to array of interface{}")
	ErrAndMustBeArray = errors.New("$and must points to array of map[string]interface{}")
	ErrOrMustBeArray  = errors.New("$or must points to array of map[string]interface{}")
	ErrNorMustBeArray = errors.New("$nor must points to array of map[string]interface{}")
	ErrNotMustBeMap   = errors.New("$not must points to map[string]interface{}")
)

// Query is a "compiled" query
type Query struct {
	and []Matcher
	or  []Matcher
	nor []Matcher
	not []Matcher
}

// NewQuery returns new compiled query
func NewQuery(query map[string]interface{}) (*Query, error) {
	q := &Query{and: []Matcher{}, or: []Matcher{}, nor: []Matcher{}, not: []Matcher{}}
	for field, v := range query {
		if field == "" {
			return nil, ErrEmptyFieldName
		}
		if field[0] == '$' {
			err := q.appendRootCondition(field, v)
			if err != nil {
				return nil, err
			}
			continue
		}

		switch v.(type) {
		case string, float64, bool:
			em, err := newEqMatcher(field, v, false)
			if err != nil {
				return nil, err
			}
			q.and = append(q.and, em)
		case int:
			v = float64(v.(int))
			em, err := newEqMatcher(field, v, false)
			if err != nil {
				return nil, err
			}
			q.and = append(q.and, em)
		case map[string]interface{}:
			val := v.(map[string]interface{})
			err := q.appendFieldCondition(field, val)
			if err != nil {
				return nil, err
			}
		}
	}

	return q, nil
}

// Match returns true if document matched compiled query
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
// returns error if query contains errors
func Match(query map[string]interface{}, doc []byte) (bool, error) {
	q, err := NewQuery(query)
	if err != nil {
		return false, err
	}
	return q.Match(doc), nil
}

func (q *Query) appendRootCondition(typ string, v interface{}) error {
	switch typ {
	case "$and":
		val, ok := v.([]interface{})
		if !ok {
			return ErrAndMustBeArray
		}
		for _, subv := range val {
			subval, ok := subv.(map[string]interface{})
			if !ok {
				return ErrAndMustBeArray
			}
			m, err := NewQuery(subval)
			if err != nil {
				return err
			}
			q.and = append(q.and, m)
		}
	case "$or":
		val, ok := v.([]interface{})
		if !ok {
			return ErrOrMustBeArray
		}
		for _, subv := range val {
			subval, ok := subv.(map[string]interface{})
			if !ok {
				return ErrOrMustBeArray
			}
			m, err := NewQuery(subval)
			if err != nil {
				return err
			}
			q.or = append(q.or, m)
		}
	case "$nor":
		val, ok := v.([]interface{})
		if !ok {
			return ErrNorMustBeArray
		}
		for _, subv := range val {
			subval, ok := subv.(map[string]interface{})
			if !ok {
				return ErrNorMustBeArray
			}
			m, err := NewQuery(subval)
			if err != nil {
				return err
			}
			q.nor = append(q.nor, m)
		}
	default:
		return errors.New("unknown operator " + typ)
	}
	return nil
}

func (q *Query) appendFieldCondition(field string, val map[string]interface{}) error {
	if val["$eq"] != nil {
		em, err := newEqMatcher(field, val["$eq"], false)
		if err != nil {
			return err
		}
		q.and = append(q.and, em)
	}
	if val["$ne"] != nil {
		em, err := newEqMatcher(field, val["$ne"], true)
		if err != nil {
			return err
		}
		q.and = append(q.and, em)
	}
	if val["$not"] != nil {
		notMap, ok := val["$not"].(map[string]interface{})
		if !ok {
			return ErrNotMustBeMap
		}
		nm, err := NewQuery(map[string]interface{}{field: notMap})
		if err != nil {
			return err
		}
		q.not = append(q.not, nm)
	}
	if val["$in"] != nil {
		arr, ok := val["$in"].([]interface{})
		if !ok {
			return ErrInMustBeArray
		}
		inm, err := newInMatcher(field, arr, false)
		if err != nil {
			return err
		}
		q.and = append(q.and, inm)
	}
	if val["$nin"] != nil {
		arr, ok := val["$nin"].([]interface{})
		if !ok {
			return ErrInMustBeArray
		}
		ninm, err := newInMatcher(field, arr, true)
		if err != nil {
			return err
		}
		q.and = append(q.and, ninm)
	}
	if val["$exists"] != nil {
		em, err := newExistsMatcher(field, val["$exists"])
		if err != nil {
			return err
		}
		q.and = append(q.and, em)
	}
	if val["$gt"] != nil {
		em, err := newGtMatcher(field, val["$gt"], false)
		if err != nil {
			return err
		}
		q.and = append(q.and, em)
	}
	if val["$gte"] != nil {
		em, err := newGtMatcher(field, val["$gte"], true)
		if err != nil {
			return err
		}
		q.and = append(q.and, em)
	}
	if val["$lt"] != nil {
		em, err := newLtMatcher(field, val["$lt"], false)
		if err != nil {
			return err
		}
		q.and = append(q.and, em)
	}
	if val["$lte"] != nil {
		em, err := newLtMatcher(field, val["$lte"], true)
		if err != nil {
			return err
		}
		q.and = append(q.and, em)
	}
	if val["$regex"] != nil {
		rm, err := newRegexMatcher(field, val)
		if err != nil {
			return err
		}
		q.and = append(q.and, rm)
	}
	return nil
}
