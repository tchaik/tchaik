// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package index

// subColFlds is a Collection wrapper which overrides the Field
// method.
type subColFlds struct {
	Collection
	flds map[string]interface{}
}

func (scf subColFlds) Field(f string) interface{} {
	if x, ok := scf.flds[f]; ok {
		return x
	}
	return scf.Collection.Field(f)
}

// subGrpFlds is a Group wrapper which overrids the Field method
type subGrpFlds struct {
	Group
	tracks []Track
	flds   map[string]interface{}
}

func (sgf subGrpFlds) Tracks() []Track {
	if sgf.tracks != nil {
		return sgf.tracks
	}
	return sgf.Group.Tracks()
}

// Field implements Group.
func (sgf subGrpFlds) Field(field string) interface{} {
	if x, ok := sgf.flds[field]; ok {
		return x
	}
	return sgf.Group.Field(field)
}

// SumGroupIntAttr recurses through the Group and assigns the field with the sum
// of fields from children (Groups or Tracks).
func SumGroupIntAttr(field string, g Group) Group {
	if c, ok := g.(Collection); ok {
		return sumCollectionIntAttr(field, c)
	}
	return sumGroupIntAttr(field, g)
}

func sumCollectionIntAttr(field string, c Collection) Collection {
	nc := subCol{
		Collection: c,
		grps:       make(map[Key]Group, len(c.Keys())),
		flds:       make(map[string]interface{}),
	}
	var total int
	for _, k := range c.Keys() {
		g := c.Get(k)
		g = SumGroupIntAttr(field, g)
		total += g.Field(field).(int)
		nc.grps[k] = g
	}
	nc.flds = map[string]interface{}{
		field: total,
	}
	return nc
}

func sumGroupIntAttr(field string, g Group) Group {
	ng := subGrpFlds{
		Group: g,
		flds:  map[string]interface{}{},
	}
	var total int
	for _, t := range g.Tracks() {
		total += t.GetInt(field)
	}
	ng.flds[field] = total
	return ng
}

// CommonGroupAttr recurses through the Group and assigns fields on all sub groups
// which are common amoungst their children (Groups or Tracks).  If there is no common
// field, then the associated Field value is not set.
func CommonGroupAttr(attrs []Attr, g Group) Group {
	if c, ok := g.(Collection); ok {
		return commonCollectionTrackAttr(attrs, c)
	}
	return commonGroupTrackAttr(attrs, g)
}

func commonCollectionTrackAttr(attrs []Attr, c Collection) Collection {
	grps := make(map[Key]Group, len(c.Keys()))
	flds := make(map[string]interface{}, len(attrs))

	keys := c.Keys()
	if len(keys) > 0 {
		k0 := keys[0]
		g0 := c.Get(k0)
		g0 = CommonGroupAttr(attrs, g0)
		grps[k0] = g0

		for _, a := range attrs {
			flds[a.Field()] = g0.Field(a.Field())
		}

		if len(keys) > 1 {
			for _, k := range keys[1:] {
				g1 := c.Get(k)
				g1 = CommonGroupAttr(attrs, g1)
				grps[k] = g1

				for _, a := range attrs {
					f := a.Field()
					flds[f] = a.Intersect(flds[f], g1.Field(f))
				}
			}
		}
	}

	for _, a := range attrs {
		f := a.Field()
		if v, ok := flds[f]; ok && a.IsEmpty(v) {
			delete(flds, f)
		}
	}

	return subCol{
		Collection: c,
		grps:       grps,
		flds:       flds,
	}
}

func commonGroupTrackAttr(attrs []Attr, g Group) Group {
	flds := make(map[string]interface{}, len(attrs))
	tracks := g.Tracks()

	if len(tracks) > 0 {
		t0 := tracks[0]
		for _, a := range attrs {
			f := a.Field()
			flds[f] = a.Value(t0)
		}

		if len(tracks) > 1 {
			for _, t := range tracks[1:] {
				for _, a := range attrs {
					f := a.Field()
					flds[f] = a.Intersect(flds[f], a.Value(t))
				}
			}
		}
	}

	for _, a := range attrs {
		f := a.Field()
		if v, ok := flds[f]; ok && a.IsEmpty(v) {
			delete(flds, f)
		}
	}

	return subGrpFlds{
		Group: g,
		flds:  flds,
	}
}

// subGrpName is a Group wrapper which overrides Name.
type subGrpName struct {
	Group
	name string
}

// Name implements Group.
func (s subGrpName) Name() string {
	return s.name
}

// RemoveEmptyCollections recursively goes through each sub Collection contained
// in the Group and removes any which don't have any tracks/groups in them.
func RemoveEmptyCollections(g Group) Group {
	gc, ok := g.(Collection)
	if ok {
		keys := gc.Keys()
		if len(keys) == 1 {
			gc0 := gc.Get(keys[0])
			_, col := gc0.(Collection)
			if !col && gc0.Name() == "" {
				return subGrpName{
					name:  gc.Name(),
					Group: gc0,
				}
			}
		}
		ngc := subCol{
			Collection: gc,
			grps:       make(map[Key]Group, len(gc.Keys())),
		}
		for _, k := range keys {
			ngc.grps[k] = RemoveEmptyCollections(gc.Get(k))
		}
		return ngc
	}
	return g
}

func firstTrack(g Group) Track {
	c, ok := g.(Collection)
	if ok {
		keys := c.Keys()
		if len(keys) > 0 {
			return firstTrack(c.Get(keys[0]))
		}
		return nil
	}

	ts := g.Tracks()
	if len(ts) > 0 {
		return ts[0]
	}
	return nil
}

func fieldsGroup(m map[string]interface{}, g Group) Group {
	if c, ok := g.(Collection); ok {
		return subColFlds{
			Collection: c,
			flds:       m,
		}
	}
	return subGrpFlds{
		Group: g,
		flds:  m,
	}
}

// Attr is a type which wraps a closure to get an attribute from an implementation of the
// Attr interface.
type Attr struct {
	field     string
	empty     interface{}
	isEmpty   func(x interface{}) bool
	fn        func(t Track) interface{}
	intersect func(x, y interface{}) interface{}
}

// Field returns the underlying field name.
func (g Attr) Field() string {
	return g.field
}

// IsEmpty returns true iff the given value represents the empty value of the underlying attribute
// type.
func (g Attr) IsEmpty(x interface{}) bool {
	return g.isEmpty(x)
}

// Value returns the value of the attribute for the given track.
func (g Attr) Value(t Track) interface{} {
	return g.fn(t)
}

// Intersect returns the intersection of the two attribute values.
func (g Attr) Intersect(x, y interface{}) interface{} {
	return g.intersect(x, y)
}

// StringAttr constructs an Attr which will retrieve the string field from an implementation
// of Track.
func StringAttr(field string) Attr {
	return Attr{
		field: field,
		empty: "",
		isEmpty: func(x interface{}) bool {
			return x == ""
		},
		intersect: func(x, y interface{}) interface{} {
			if x == y {
				return x
			}
			return ""
		},
		fn: func(t Track) interface{} {
			return t.GetString(field)
		},
	}
}

// StringSliceEqual is a function used to compare two interface{} types which are assumed
// to be of type []string (or interface{}(nil)).
func StringSliceEqual(x, y interface{}) bool {
	// Annoyingly we have to cater for zero values from map[string]interface{}
	// which don't have the correct type wrapping the nil.
	if x == nil || y == nil {
		return x == nil && y == nil
	}
	xs := x.([]string) // NB: panics here are acceptable: should not be called on a non-'Strings' field.
	ys := y.([]string)
	if len(xs) != len(ys) {
		return false
	}
	for i, xss := range xs {
		if ys[i] != xss {
			return false
		}
	}
	return true
}

// StringSliceIntersect computes the intersection of two string slices (ignoring ordering).
func StringSliceIntersect(s, t []string) []string {
	var res []string
	m := make(map[string]bool)
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

// StringsIntersect computes the intersection of the two interface values assumed to be
// of type []string.
func StringsIntersect(x, y interface{}) interface{} {
	if x == nil || y == nil {
		return nil
	}
	xs := x.([]string)
	ys := y.([]string)
	return StringSliceIntersect(xs, ys)
}

// StringsAttr returns an Attr which will retrieve the strings field from an implementation of Track.
func StringsAttr(field string) Attr {
	return Attr{
		field:     field,
		empty:     nil,
		intersect: StringsIntersect,
		isEmpty: func(x interface{}) bool {
			if x == nil {
				return true
			}
			xs := x.([]string)
			return len(xs) == 0
		},
		fn: func(t Track) interface{} {
			return t.GetStrings(field)
		},
	}
}

// IntAttr constructs an Attr which will retrieve the int field from an implementation of Track.
func IntAttr(field string) Attr {
	return Attr{
		field:   field,
		empty:   0,
		isEmpty: func(x interface{}) bool { return x == 0 },
		intersect: func(x, y interface{}) interface{} {
			if x == y {
				return x
			}
			return 0
		},
		fn: func(t Track) interface{} { return t.GetInt(field) },
	}
}

// FirstTrackAttr wraps the given Group adding a string field `field` with the value taken
// from the first track.
func FirstTrackAttr(attr Attr, g Group) Group {
	t := firstTrack(g)
	if t == nil {
		return g
	}

	v := attr.Value(t)
	if attr.IsEmpty(v) {
		return g
	}
	m := map[string]interface{}{
		attr.field: v,
	}
	return fieldsGroup(m, g)
}
