package mongofil

import (
	"testing"

	"github.com/franela/goblin"
)

func TestLtMatch(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("#newLtMatcher", func() {
		g.It("should return nil error when valid value passed", func() {
			_, err := newLtMatcher("prop", 42, false)
			g.Assert(err != nil).IsFalse("error should be nil")

			_, err = newLtMatcher("prop", "42", false)
			g.Assert(err != nil).IsFalse("error should be nil")
		})
		g.It("should return error when invalid value passed", func() {
			_, err := newLtMatcher("prop", []int{}, false)
			g.Assert(err != nil).IsTrue("error should not be nil")
		})
	})

	g.Describe("#Match", func() {
		doc := []byte(`{"title": "Batman & Robin", "duration": 112}`)
		g.It("should return true if document matched query", func() {
			m, err := newLtMatcher("duration", 200, false)
			g.Assert(err != nil).IsFalse("error should be nil")
			g.Assert(m.Match(doc)).IsTrue()

			m, err = newLtMatcher("duration", 112, true)
			g.Assert(err != nil).IsFalse("error should be nil")
			g.Assert(m.Match(doc)).IsTrue()

			m, err = newLtMatcher("title", "Cars", false)
			g.Assert(err != nil).IsFalse("error should be nil")
			g.Assert(m.Match(doc)).IsTrue()

			m, err = newLtMatcher("title", "Batman & Robin", true)
			g.Assert(err != nil).IsFalse("error should be nil")
			g.Assert(m.Match(doc)).IsTrue()
		})
		g.It("should return false if document not matched query", func() {
			m, err := newLtMatcher("duration", 112, false)
			g.Assert(err != nil).IsFalse("error should be nil")
			g.Assert(m.Match(doc)).IsFalse()

			m, err = newLtMatcher("duration", 100, true)
			g.Assert(err != nil).IsFalse("error should be nil")
			g.Assert(m.Match(doc)).IsFalse()

			m, err = newLtMatcher("title", "Aurum", false)
			g.Assert(err != nil).IsFalse("error should be nil")
			g.Assert(m.Match(doc)).IsFalse()

			m, err = newLtMatcher("title", "Batman & Robin", false)
			g.Assert(err != nil).IsFalse("error should be nil")
			g.Assert(m.Match(doc)).IsFalse()
		})
		g.It("should return false if invalid JSON passed", func() {
			m, err := newLtMatcher("duration", 112, false)
			g.Assert(err != nil).IsFalse("error should be nil")
			g.Assert(m.Match([]byte(`{`))).IsFalse()
		})
		g.It("should return false if type of condition and value mismatched", func() {
			m, err := newLtMatcher("duration", "112", false)
			g.Assert(err != nil).IsFalse("error should be nil")
			g.Assert(m.Match(doc)).IsFalse()

			m, err = newLtMatcher("title", 42, false)
			g.Assert(err != nil).IsFalse("error should be nil")
			g.Assert(m.Match(doc)).IsFalse()

			m, err = newLtMatcher("title", "Batman", false)
			g.Assert(err != nil).IsFalse("error should be nil")
			g.Assert(m.Match([]byte(`{"title": true}`))).IsFalse()
		})
	})
}
