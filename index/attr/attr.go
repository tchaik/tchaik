// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package attr defines types and helpers for accessing typed attributes.
package attr

// Getter is an interface which is implemented by types which want to export typed
// values.
type Getter interface {
	GetInt(string) int
	GetString(string) string
	GetStrings(string) []string
}

// Interface is a type which defines behaviour necessary to implement a typed attribute.
type Interface interface {
	// Name returns the name of the attribute.
	Name() string

	// IsEmpty returns true iff `x` is a representation of the empty value of this attribute.
	IsEmpty(x interface{}) bool

	// Value returns the value of this attribute in `g`.
	Value(g Getter) interface{}

	// Intersect returns the intersection of `x` and `y`, assumed to be the type of the
	// attribute.
	Intersect(x, y interface{}) interface{}
}

type valueType struct {
	name  string
	empty interface{}
	get   func(Getter) interface{}
}

// Name implements Interface.
func (v *valueType) Name() string {
	return v.name
}

// IsEmpty implements Interface.
func (v *valueType) IsEmpty(x interface{}) bool {
	return v.empty == x
}

// Value implements Interface.
func (v *valueType) Value(g Getter) interface{} {
	return v.get(g)
}

// Intersect implements Interface.
func (v *valueType) Intersect(x, y interface{}) interface{} {
	if x == y {
		return x
	}
	return v.empty
}

// String constructs an implementation of Interface with the field name `f` to access a string
// attribute of an implementation of Getter.
func String(f string) Interface {
	return &valueType{
		name:  f,
		empty: "",
		get: func(g Getter) interface{} {
			return g.GetString(f)
		},
	}
}

// Int constructs an implementation of Interface with the field name `f` to access an int
// attribute of an implementation of Getter.
func Int(f string) Interface {
	return &valueType{
		name:  f,
		empty: 0,
		get: func(g Getter) interface{} {
			return g.GetInt(f)
		},
	}
}

type stringsType struct {
	valueType
}

// IsEmpty implements Interface.
func (p *stringsType) IsEmpty(x interface{}) bool {
	if x == nil {
		return true
	}
	xs := x.([]string)
	return len(xs) == 0
}

// Intersect implements Interface.
func (p *stringsType) Intersect(x, y interface{}) interface{} {
	if x == nil || y == nil {
		return nil
	}
	xs := x.([]string)
	ys := y.([]string)
	return stringSliceIntersect(xs, ys)
}

// Strings constructs an implementation of Interface with the field name `f` to access a string slice
// attribute of an implementation of Getter.
func Strings(f string) Interface {
	return &stringsType{
		valueType{
			name:  f,
			empty: nil,
			get: func(g Getter) interface{} {
				return g.GetStrings(f)
			},
		},
	}
}

// stringSliceIntersect computes the intersection of two string slices (ignoring ordering).
func stringSliceIntersect(s, t []string) []string {
	var res []string
	m := make(map[string]bool, len(s))
	for _, x := range s {
		m[x] = true
	}
	for _, y := range t {
		if m[y] {
			res = append(res, y)
		}
	}
	return res
}
