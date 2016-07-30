[![Build Status](https://travis-ci.org/ivahaev/mongofil.svg?branch=master)](https://travis-ci.org/ivahaev/mongofil)
[![Go Report Card](https://goreportcard.com/badge/github.com/ivahaev/mongofil)](https://goreportcard.com/report/github.com/ivahaev/mongofil)
[![License](https://img.shields.io/badge/license-MIT%20v3-blue.svg)](https://github.com/github.com/ivahaev/blob/master/LICENSE)

# mongofil
Checking if JSON matches mongodb query and finding matched documents in JSON array.
Project is on very early stage. Contributors welcome!


## Roadmap:

- [ ] Support for nested documents and arrays
- [ ] Filter JSON array to only matched documents
- [ ] Caching extracted values to use them in multiple operators

### Implement query operators
#### Comparison
- [x] $eq	    Matches values that are equal to a specified value.
- [x] $gt	    Matches values that are greater than a specified value.
- [x] $gte	Matches values that are greater than or equal to a specified value.
- [x] $lt	    Matches values that are less than a specified value.
- [x] $lte	Matches values that are less than or equal to a specified value.
- [x] $ne	    Matches all values that are not equal to a specified value.
- [x] $in	    Matches any of the values specified in an array.
- [x] $nin	Matches none of the values specified in an array.

#### Logical
- [x] $or	Joins query clauses with a logical OR returns all documents that match the conditions of either clause.
- [x] $and	Joins query clauses with a logical AND returns all documents that match the conditions of both clauses.
- [x] $not	Inverts the effect of a query expression and returns documents that do not match the query expression.
- [x] $nor	Joins query clauses with a logical NOR returns all documents that fail to match both clauses.

#### Element
- [x] $exists	Matches documents that have the specified field.
- [ ] $type	    Selects documents if a field is of the specified type.

#### Evaluation
- [ ] $mod	Performs a modulo operation on the value of a field and selects documents with a specified result.
- [ ] $regex	Selects documents where values match a specified regular expression.
- [ ] $text	Performs text search.
- [ ] $where	Matches documents that satisfy a JavaScript expression.

#### Geospatial
- [ ] $geoWithin	Selects geometries within a bounding GeoJSON geometry. The 2dsphere and 2d indexes support $geoWithin.
- [ ] $geoIntersects	Selects geometries that intersect with a GeoJSON geometry. The 2dsphere index supports $geoIntersects.
- [ ] $near	Returns geospatial objects in proximity to a point. Requires a geospatial index. The 2dsphere and 2d indexes support $near.
- [ ] $nearSphere	Returns geospatial objects in proximity to a point on a sphere. Requires a geospatial index. The 2dsphere and 2d indexes support $nearSphere.

#### Array
- [ ] $all	Matches arrays that contain all elements specified in the query.
- [ ] $elemMatch	Selects documents if element in the array field matches all the specified $elemMatch conditions.
- [ ] $size	Selects documents if the array field is a specified size.