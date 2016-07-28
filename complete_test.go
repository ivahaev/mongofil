package mongofil

import (
	"encoding/json"
	"testing"

	"github.com/franela/goblin"
)

type testPair struct {
	q             []byte
	doc           []byte
	shouldMatched bool
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
				g.Assert(err == nil).IsTrue()
				g.Assert(matched).Equal(p.shouldMatched)
			}
		})
	})
}
