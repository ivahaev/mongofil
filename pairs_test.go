package mongofil

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/franela/goblin"
)

type testPair struct {
	q             []byte
	doc           []byte
	shouldMatched bool
	err           error
}

var pairs = []testPair{
	{
		q:             []byte(`{}`),
		doc:           []byte(`{}`),
		shouldMatched: true,
	},
	{
		q:             []byte(`{"name": "Batman"}`),
		doc:           []byte(`{"name": "Joker"}`),
		shouldMatched: false,
	},
	{
		q:             []byte(`{"name": "Batman"}`),
		doc:           []byte(`{"name": "Batman"}`),
		shouldMatched: true,
	},
	{
		q:             []byte(`{"name": {"$ne": "Joker"}}`),
		doc:           []byte(`{"name": "Batman"}`),
		shouldMatched: true,
	},
	{
		q:             []byte(`{"$or": [{"name": "Joker"}, {"name": "Batman"}]}`),
		doc:           []byte(`{"name": "Batman"}`),
		shouldMatched: true,
	},
	{
		q:             []byte(`{"$or": [{"name": "Joker"}, {"name": "Batman"}]}`),
		doc:           []byte(`{"name": "Superman"}`),
		shouldMatched: false,
	},
	{
		q:             []byte(`{"$or": [{"name": "Joker"}, {"name": "Batman"}], "$and": [{"lastName": "Petrov"}, {"age": 230}]}`),
		doc:           []byte(`{"name": "Batman", "lastName": "Petrov", "age": 230}`),
		shouldMatched: true,
	},
	{
		q:             []byte(`{"$or": [{"name": "Joker"}, {"name": "Batman"}], "$and": [{"lastName": "Ivanov"}, {"age": 230}]}`),
		doc:           []byte(`{"name": "Batman", "lastName": "Petrov", "age": 230}`),
		shouldMatched: false,
	},
	{
		q:             []byte(`{"$or": [{"name": "Joker"}, {"name": "Batman"}], "$and": [{"lastName": {"$ne": "Petrov"}}, {"age": 230}]}`),
		doc:           []byte(`{"name": "Batman", "lastName": "Petrov", "age": 230}`),
		shouldMatched: false,
	},
	{
		q:             []byte(`{"$or": [{"name": "Joker"}, {"name": "Batman"}], "$and": [{"lastName": {"$ne": "Ivanov"}}, {"age": 230}]}`),
		doc:           []byte(`{"name": "Batman", "lastName": "Petrov", "age": 230}`),
		shouldMatched: true,
	},
	{
		q:             []byte(`{"$or": [{"name": "Joker"}, {"name": "Batman"}], "$and": [{"lastName": {"$ne": "Ivanov"}}, {"$or": [{"age": 200}, {"age": {"$gt": 220}  }]  }]}`),
		doc:           []byte(`{"name": "Batman", "lastName": "Petrov", "age": 230}`),
		shouldMatched: true,
	},
	{
		q:             []byte(`{"$or": [{"name": "Joker"}, {"name": "Batman"}], "$and": [{"lastName": {"$ne": "Ivanov"}}, {"$or": [{"age": 200}, {"age": {"$gt": 230}  }]  }]}`),
		doc:           []byte(`{"name": "Batman", "lastName": "Petrov", "age": 230}`),
		shouldMatched: false,
	},
	{
		q:             []byte(`{"$or": [{"name": "Joker"}, {"name": "Batman"}], "$and": [{"lastName": {"$ne": "Ivanov"}}, {"$or": [{"age": 200}, {"age": {"$gte": 230}  }]  }]}`),
		doc:           []byte(`{"name": "Batman", "lastName": "Petrov", "age": 230}`),
		shouldMatched: true,
	},
	{
		q:             []byte(`{"$or": [{"name": "Joker"}, {"name": "Batman"}], "$and": [{"lastName": {"$ne": "Ivanov"}}, {"$or": [{"age": 200}, {"age": {"$lte": 230}  }]  }]}`),
		doc:           []byte(`{"name": "Batman", "lastName": "Petrov", "age": 230}`),
		shouldMatched: true,
	},
	{
		q:             []byte(`{"$or": [{"name": "Joker"}, {"name": "Batman"}], "$and": [{"lastName": {"$ne": "Ivanov"}}, {"$or": [{"age": 200}, {"age": {"$lt": 230}  }]  }]}`),
		doc:           []byte(`{"name": "Batman", "lastName": "Petrov", "age": 230}`),
		shouldMatched: false,
	},
	{
		q:             []byte(`{"$and": [{"name": "Vasya"}, {"age": {"$not": {"$gte": 40}}}]}`),
		doc:           []byte(`{"name": "Vasya", "age": "30"}`),
		shouldMatched: true,
	},
	{
		q:             []byte(`{"$and": [{"name": "Vasya"}, {"age": {"$not": {"$gte": 40}}}]}`),
		doc:           []byte(`{"name": "Vasya", "age": 50}`),
		shouldMatched: false,
	},
	{
		q:             []byte(`{"$and": [{"name": "Vasya"}, {"lastName": {"$regex": "^.*ov$", "$options": "i"}}]}`),
		doc:           []byte(`{"name": "Vasya", "lastName": "Ivanov", "age": 50}`),
		shouldMatched: true,
	},
	{
		q:             []byte(`{"$and": [{"name": "Vasya"}, {"lastName": {"$regex": "^.*ov$", "$options": "i"}}]}`),
		doc:           []byte(`{"name": "Vasya", "lastName": "IVANOV", "age": 50}`),
		shouldMatched: true,
	},
	{
		q:             []byte(`{"$nor": [{"name": "Vasya"}, {"lastName": {"$regex": "^.*ovitch$", "$options": "i"}}]}`),
		doc:           []byte(`{"name": "Vasya", "lastName": "IVANOV", "age": 50}`),
		shouldMatched: false,
	},
	{
		q:             []byte(`{"$nor": [{"name": "vasya"}, {"lastName": {"$regex": "^.*ovitch$", "$options": "i"}}]}`),
		doc:           []byte(`{"name": "Vasya", "lastName": "IVANOV", "age": 50}`),
		shouldMatched: true,
	},
	{
		q:             []byte(`{"$and": [{"name": "Vasya"}, {"lastName": {"$regex": "^.*ov$", "$options": "i"}}]}`),
		doc:           []byte(`{"name": "Vasya", "lastName": "Rabinovitch", "age": 50}`),
		shouldMatched: false,
	},
	{
		q:             []byte(`{"$and": [{"name": "Vasya"}, 42]}`),
		doc:           []byte(`{"name": "Vasya", "lastName": "Rabinovitch", "age": 50}`),
		shouldMatched: false,
		err:           ErrAndMustBeArray,
	},
	{
		q:             []byte(`{"$and": 42}`),
		doc:           []byte(`{"name": "Vasya", "lastName": "Rabinovitch", "age": 50}`),
		shouldMatched: false,
		err:           ErrAndMustBeArray,
	},
	{
		q:             []byte(`{"$or": [{"name": "Vasya"}, 42]}`),
		doc:           []byte(`{"name": "Vasya", "lastName": "Rabinovitch", "age": 50}`),
		shouldMatched: false,
		err:           ErrOrMustBeArray,
	},
	{
		q:             []byte(`{"$or": 42}`),
		doc:           []byte(`{"name": "Vasya", "lastName": "Rabinovitch", "age": 50}`),
		shouldMatched: false,
		err:           ErrOrMustBeArray,
	},
}

func TestMatchPairs(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Pairs testing", func() {
		g.It("should pass test", func() {
			for _, p := range pairs {
				var q map[string]interface{}
				err := json.Unmarshal(p.q, &q)
				g.Assert(err == nil).IsTrue()
				matched, err := Match(q, p.doc)
				g.Assert(err == p.err).IsTrue(func() string {
					if err != p.err {
						fmt.Println(fmt.Sprintf(`Failed query: %s, with document: %s. Error expected: %v, but received: %v`, string(p.q), string(p.doc), p.err, err))
					}
					if p.err != nil {
						return p.err.Error()
					}
					return ""
				}())
				if matched != p.shouldMatched {
					fmt.Println(fmt.Sprintf(`Failed query: %s, with document: %s. Expected: %v`, string(p.q), string(p.doc), p.shouldMatched))
				}
				g.Assert(matched).Equal(p.shouldMatched)
			}
		})
	})
}
