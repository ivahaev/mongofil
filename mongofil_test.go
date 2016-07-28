package mongofil

import (
	"testing"

	"github.com/franela/goblin"
)

func TestMatch(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("#Match", func() {
		g.Describe("simple query", func() {
			g.It("should return true if JSON is matched simple query with string value", func() {
				query := map[string]interface{}{"name": "Vasya"}
				json := []byte(`{"name": "Vasya", "lastName": "Ivanov"}`)
				matched, err := Match(query, json)
				g.Assert(err == nil).IsTrue()
				g.Assert(matched).IsTrue()
			})
			g.It("should return false if JSON is not matched simple query with string value", func() {
				query := map[string]interface{}{"name": "vasya"}
				json := []byte(`{"name": "Vasya", "lastName": "Ivanov"}`)
				matched, err := Match(query, json)
				g.Assert(err == nil).IsTrue()
				g.Assert(matched).IsFalse()
			})
		})

		g.Describe("$in statement", func() {
			g.It("should return true if JSON is matched query with $in statement with string value", func() {
				query := map[string]interface{}{"name": map[string]interface{}{"$in": []interface{}{"Petya", 1, "Vasya"}}}
				json := []byte(`{"name": "Vasya", "lastName": "Ivanov"}`)
				matched, err := Match(query, json)
				g.Assert(err == nil).IsTrue()
				g.Assert(matched).IsTrue()
			})
			g.It("should return false if JSON is not matched query with $in statement with string value", func() {
				query := map[string]interface{}{"name": map[string]interface{}{"$in": []interface{}{"Petya", 1, "vasya"}}}
				json := []byte(`{"name": "Vasya", "lastName": "Ivanov"}`)
				matched, err := Match(query, json)
				g.Assert(err == nil).IsTrue()
				g.Assert(matched).IsFalse()
			})

			g.It("should return true if JSON is matched query with $in statement with number value", func() {
				query := map[string]interface{}{"name": map[string]interface{}{"$in": []interface{}{"Petya", 1, "Vasya"}}}
				json := []byte(`{"name": 1, "lastName": "Ivanov"}`)
				matched, err := Match(query, json)
				g.Assert(err == nil).IsTrue()
				g.Assert(matched).IsTrue()
			})
			g.It("should return false if JSON is not matched query with $in statement with number value", func() {
				query := map[string]interface{}{"name": map[string]interface{}{"$in": []interface{}{"Petya", 1, "vasya"}}}
				json := []byte(`{"name": 2, "lastName": "Ivanov"}`)
				matched, err := Match(query, json)
				g.Assert(err == nil).IsTrue()
				g.Assert(matched).IsFalse()
			})
		})

		g.Describe("$nin statement", func() {
			g.It("should return true if JSON is matched query with $nin statement with string value", func() {
				query := map[string]interface{}{"name": map[string]interface{}{"$nin": []interface{}{"Petya", 1, "vasya"}}}
				json := []byte(`{"name": "Vasya", "lastName": "Ivanov"}`)
				matched, err := Match(query, json)
				g.Assert(err == nil).IsTrue()
				g.Assert(matched).IsTrue()
			})
			g.It("should return false if JSON is not matched query with $in statement with string value", func() {
				query := map[string]interface{}{"name": map[string]interface{}{"$nin": []interface{}{"Petya", 1, "Vasya"}}}
				json := []byte(`{"name": "Vasya", "lastName": "Ivanov"}`)
				matched, err := Match(query, json)
				g.Assert(err == nil).IsTrue()
				g.Assert(matched).IsFalse()
			})

			g.It("should return true if JSON is matched query with $in statement with number value", func() {
				query := map[string]interface{}{"name": map[string]interface{}{"$nin": []interface{}{"Petya", 2, "Vasya"}}}
				json := []byte(`{"name": 1, "lastName": "Ivanov"}`)
				matched, err := Match(query, json)
				g.Assert(err == nil).IsTrue()
				g.Assert(matched).IsTrue()
			})
			g.It("should return false if JSON is not matched query with $in statement with number value", func() {
				query := map[string]interface{}{"name": map[string]interface{}{"$nin": []interface{}{"Petya", 2, "vasya"}}}
				json := []byte(`{"name": 2, "lastName": "Ivanov"}`)
				matched, err := Match(query, json)
				g.Assert(err == nil).IsTrue()
				g.Assert(matched).IsFalse()
			})
		})

		g.Describe("$exists statement", func() {
			g.It("should return true if property value is exists", func() {
				query := map[string]interface{}{"name": map[string]interface{}{"$exists": true}}
				json := []byte(`{"name": 2, "lastName": "Ivanov"}`)
				matched, err := Match(query, json)
				g.Assert(err == nil).IsTrue()
				g.Assert(matched).IsTrue()
			})
			g.It("should return true if property value is not exists", func() {
				query := map[string]interface{}{"middleName": map[string]interface{}{"$exists": false}}
				json := []byte(`{"name": 2, "lastName": "Ivanov"}`)
				matched, err := Match(query, json)
				g.Assert(err == nil).IsTrue()
				g.Assert(matched).IsTrue()
			})

			g.It("should return false if property value is exists", func() {
				query := map[string]interface{}{"name": map[string]interface{}{"$exists": false}}
				json := []byte(`{"name": 2, "lastName": "Ivanov"}`)
				matched, err := Match(query, json)
				g.Assert(err == nil).IsTrue()
				g.Assert(matched).IsFalse()
			})
			g.It("should return false if property value is not exists", func() {
				query := map[string]interface{}{"middleName": map[string]interface{}{"$exists": true}}
				json := []byte(`{"name": 2, "lastName": "Ivanov"}`)
				matched, err := Match(query, json)
				g.Assert(err == nil).IsTrue()
				g.Assert(matched).IsFalse()
			})
		})
	})
}
