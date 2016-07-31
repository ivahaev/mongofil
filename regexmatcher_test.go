package mongofil

import (
	"testing"

	"github.com/franela/goblin"
)

func TestRegexMatch(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("#newRegexMatcher", func() {
		g.It("should return nil error when condition is right map", func() {
			_, err := newRegexMatcher("prop", map[string]interface{}{"$regex": "aa"})
			g.Assert(err != nil).IsFalse("error should be nil")
		})
		g.It("should return error when condition is not right map", func() {
			_, err := newRegexMatcher("prop", "//")
			g.Assert(err != nil).IsTrue("error should not be nil")

			_, err = newRegexMatcher("prop", map[string]interface{}{"$regexp": "aa"})
			g.Assert(err != nil).IsTrue("error should not be nil")

			_, err = newRegexMatcher("prop", map[string]interface{}{"$regex": 42})
			g.Assert(err != nil).IsTrue("error should not be nil")

			_, err = newRegexMatcher("prop", map[string]interface{}{"$regex": "aa", "$options": 42})
			g.Assert(err != nil).IsTrue("error should not be nil")
		})

		g.It("should return error when invalid regex passed", func() {
			_, err := newRegexMatcher("prop", map[string]interface{}{"$regex": "(aa"})
			g.Assert(err != nil).IsTrue("error should not be nil")
		})
	})

	g.Describe("#Match", func() {
		doc := []byte(`{"title": " Batman & Robin ", "duration": 112}`)
		g.It("should return true if document matched query", func() {
			m, err := newRegexMatcher("title", map[string]interface{}{"$regex": "Batman"})
			g.Assert(err != nil).IsFalse("error should be nil")
			g.Assert(m.Match(doc)).IsTrue()

			m, err = newRegexMatcher("title", map[string]interface{}{"$regex": "batman", "$options": "i"})
			g.Assert(err != nil).IsFalse("error should be nil")
			g.Assert(m.Match(doc)).IsTrue()
		})
		g.It("should return false if document not matched query", func() {
			m, err := newRegexMatcher("title", map[string]interface{}{"$regex": "batman"})
			g.Assert(err != nil).IsFalse("error should be nil")
			g.Assert(m.Match(doc)).IsFalse()

			m, err = newRegexMatcher("title", map[string]interface{}{"$regex": "Robinovich", "$options": "i"})
			g.Assert(err != nil).IsFalse("error should be nil")
			g.Assert(m.Match(doc)).IsFalse()

			g.Assert(m.Match([]byte(`"title": 42`))).IsFalse()
		})
	})
}
