// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package index

import "github.com/tchaik/tchaik/index/attr"

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
func CommonGroupAttr(a []attr.Interface, g Group) Group {
	if c, ok := g.(Collection); ok {
		return commonCollectionTrackAttr(a, c)
	}
	return commonGroupTrackAttr(a, g)
}

func commonCollectionTrackAttr(l []attr.Interface, c Collection) Collection {
	grps := make(map[Key]Group, len(c.Keys()))
	flds := make(map[string]interface{}, len(l))

	keys := c.Keys()
	if len(keys) > 0 {
		k0 := keys[0]
		g0 := c.Get(k0)
		g0 = CommonGroupAttr(l, g0)
		grps[k0] = g0

		for _, a := range l {
			flds[a.Name()] = g0.Field(a.Name())
		}

		if len(keys) > 1 {
			for _, k := range keys[1:] {
				g1 := c.Get(k)
				g1 = CommonGroupAttr(l, g1)
				grps[k] = g1

				for _, a := range l {
					f := a.Name()
					flds[f] = a.Intersect(flds[f], g1.Field(f))
				}
			}
		}
	}

	for _, a := range l {
		f := a.Name()
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

func commonGroupTrackAttr(l []attr.Interface, g Group) Group {
	flds := make(map[string]interface{}, len(l))
	tracks := g.Tracks()

	if len(tracks) > 0 {
		t0 := tracks[0]
		for _, a := range l {
			f := a.Name()
			flds[f] = a.Value(t0)
		}

		if len(tracks) > 1 {
			for _, t := range tracks[1:] {
				for _, a := range l {
					f := a.Name()
					flds[f] = a.Intersect(flds[f], a.Value(t))
				}
			}
		}
	}

	for _, a := range l {
		f := a.Name()
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

// FirstTrackAttr wraps the given Group adding a string field `field` with the value taken
// from the first track.
func FirstTrackAttr(a attr.Interface, g Group) Group {
	t := firstTrack(g)
	if t == nil {
		return g
	}

	v := a.Value(t)
	if a.IsEmpty(v) {
		return g
	}
	m := map[string]interface{}{
		a.Name(): v,
	}
	return fieldsGroup(m, g)
}
