package mongofil

import (
	"testing"

	"github.com/franela/goblin"
)

func TestExistMatch(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("#newExistsMatcher", func() {
		g.It("should return nil error when condition is boolean", func() {
			_, err := newExistsMatcher("prop", true)
			g.Assert(err != nil).IsFalse("error should be nil")
		})
		g.It("should return error when condition is not boolean", func() {
			_, err := newExistsMatcher("prop", "true")
			g.Assert(err != nil).IsTrue("error should not be nil")
		})
	})
}
