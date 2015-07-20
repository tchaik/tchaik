// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package index

import (
	"reflect"
	"testing"
)

type testCol struct {
	name string
	keys []Key
	grps map[Key]Group
	flds map[string]interface{}
}

func (c testCol) Keys() []Key                    { return c.keys }
func (c testCol) Name() string                   { return c.name }
func (c testCol) Get(k Key) Group                { return c.grps[k] }
func (c testCol) Field(field string) interface{} { return c.flds[field] }

func (c testCol) Tracks() []Track {
	return collectionTracks(c)
}

func TestSumGroupIntAttr(t *testing.T) {
	table := []struct {
		in    Group
		field string
		total int
	}{
		// One group with one track
		{
			in: group{
				name: "Group One",
				tracks: []Track{
					testTrack{
						Name:     "Track One",
						Duration: 1,
					},
				},
			},
			field: "Duration",
			total: 1,
		},

		// One group, two tracks
		{
			in: group{
				name: "Group Two",
				tracks: []Track{
					testTrack{
						Name:     "Track One",
						Duration: 1,
					},
					testTrack{
						Name:     "Track Two",
						Duration: 2,
					},
				},
			},
			field: "Duration",
			total: 3,
		},

		// One collection, one group, three tracks
		{
			in: testCol{
				name: "Group Three (collection)",
				keys: []Key{"Group-Three-One"},
				grps: map[Key]Group{
					"Group-Three-One": group{
						name: "Group One (Three)",
						tracks: []Track{
							testTrack{
								Name:     "Track One",
								Duration: 1,
							},
							testTrack{
								Name:     "Track Two",
								Duration: 2,
							},
							testTrack{
								Name:     "Track One",
								Duration: 4,
							},
						},
					},
				},
			},
			field: "Duration",
			total: 7,
		},

		// One collection, three groups, many tracks!
		{
			in: testCol{
				name: "Root",
				keys: []Key{"Group-One", "Group-Two", "Group-Three"},
				grps: map[Key]Group{
					"Group-One": group{
						name: "Group One",
						tracks: []Track{
							testTrack{
								Name:     "Track One",
								Duration: 1,
							},
						},
					},
					"Group-Two": group{
						name: "Group Two",
						tracks: []Track{
							testTrack{
								Name:     "Track One",
								Duration: 10,
							},
							testTrack{
								Name:     "Track Two",
								Duration: 20,
							},
						},
					},
					"Group-Three": testCol{
						name: "Group Three (collection)",
						keys: []Key{"Group-Three-One"},
						grps: map[Key]Group{
							"Group-Three-One": group{
								name: "Group One (Three)",
								tracks: []Track{
									testTrack{
										Name:     "Track One",
										Duration: 100,
									},
									testTrack{
										Name:     "Track Two",
										Duration: 200,
									},
									testTrack{
										Name:     "Track One",
										Duration: 400,
									},
								},
							},
						},
					},
				},
			},
			field: "Duration",
			total: 731,
		},
	}

	for ii, tt := range table {
		g := SumGroupIntAttr(tt.field, tt.in)
		got := g.Field(tt.field)
		if got != tt.total {
			t.Errorf("[%d] Fields()['%v'] = %d, expected %d", ii, tt.field, got, tt.total)
		}
	}
}

func TestCommonGroupAttr(t *testing.T) {
	table := []struct {
		in     Group
		fields []Attr
		out    []interface{}
	}{
		// One group with one track, unset common field
		{
			in: group{
				name: "Group One",
				tracks: []Track{
					testTrack{
						Name: "Track One",
					},
				},
			},
			fields: []Attr{StringAttr("Artist")},
			out:    []interface{}{nil},
		},

		// One group with one track, unset common int field
		{
			in: group{
				name: "Group One",
				tracks: []Track{
					testTrack{
						Name: "Empty common field",
					},
				},
			},
			fields: []Attr{IntAttr("Year")},
			out:    []interface{}{nil},
		},

		// One group with one track, unset common (string) and (int) fields
		{
			in: group{
				name: "Group One",
				tracks: []Track{
					testTrack{
						Name: "Track One",
					},
				},
			},
			fields: []Attr{StringAttr("Artist"), IntAttr("Year")},
			out:    []interface{}{nil, nil},
		},

		// One group with two tracks, empty first string & int fields
		{
			in: group{
				name: "Group One",
				tracks: []Track{
					testTrack{
						Name: "Track One",
					},
					testTrack{
						Name:   "Track Two",
						Artist: "Artist One",
						Year:   1984,
					},
				},
			},
			fields: []Attr{StringAttr("Artist"), IntAttr("Year")},
			out:    []interface{}{nil, nil},
		},

		// One group with two tracks, empty first field
		{
			in: group{
				name: "Group One",
				tracks: []Track{
					testTrack{
						Name: "Track One",
					},
					testTrack{
						Name: "Track Two",
						Year: 1985,
					},
				},
			},
			fields: []Attr{IntAttr("Year")},
			out:    []interface{}{nil},
		},

		// One group with one track, common string/int fields
		{
			in: group{
				name: "Group One",
				tracks: []Track{
					testTrack{
						Name:     "Track One",
						Artist:   "Artist One",
						Composer: "Composer One",
						Year:     1984,
					},
				},
			},
			fields: []Attr{StringAttr("Artist"), IntAttr("Year")},
			out:    []interface{}{"Artist One", 1984},
		},

		// One group with one track, common string fields for Artist and Composer
		{
			in: group{
				name: "Group One",
				tracks: []Track{
					testTrack{
						Name:     "Track One",
						Artist:   "Artist One",
						Composer: "Composer One",
					},
				},
			},
			fields: []Attr{StringAttr("Artist"), StringAttr("Composer")},
			out:    []interface{}{"Artist One", "Composer One"},
		},

		// One group with two tracks, common fields across pairs
		{
			in: group{
				name: "Group One",
				tracks: []Track{
					testTrack{
						Name:     "Track One",
						Artist:   "Artist One",
						Composer: "Composer One",
					},
					testTrack{
						Name:     "Track Two",
						Artist:   "Artist Two",
						Composer: "Composer One",
					},
				},
			},
			fields: []Attr{StringAttr("Artist"), StringAttr("Composer")},
			out:    []interface{}{nil, "Composer One"},
		},

		// One collection, one group, one track
		{
			in: testCol{
				name: "Group Three (collection)",
				keys: []Key{"Group-Three-One"},
				grps: map[Key]Group{
					"Group-Three-One": group{
						name: "Group One (Three)",
						tracks: []Track{
							testTrack{
								Name:   "Track One",
								Artist: "Artist One",
							},
						},
					},
				},
			},
			fields: []Attr{StringAttr("Artist"), StringAttr("Composer")},
			out:    []interface{}{"Artist One", nil},
		},

		// One collection, three groups, many tracks!
		{
			in: testCol{
				name: "Root",
				keys: []Key{"Group-One", "Group-Two", "Group-Three"},
				grps: map[Key]Group{
					"Group-One": group{
						name: "Group One",
						tracks: []Track{
							testTrack{
								Name:   "Track One",
								Artist: "Artist One",
								Year:   1984,
							},
						},
					},
					"Group-Two": group{
						name: "Group Two",
						tracks: []Track{
							testTrack{
								Name:   "Track One",
								Artist: "Artist One",
								Year:   1984,
							},
							testTrack{
								Name:   "Track Two",
								Artist: "Artist One",
								Year:   1984,
							},
						},
					},
					"Group-Three": testCol{
						name: "Group Three (collection)",
						keys: []Key{"Group-Three-One"},
						grps: map[Key]Group{
							"Group-Three-One": group{
								name: "Group One (Three)",
								tracks: []Track{
									testTrack{
										Name:   "Track One",
										Artist: "Artist One",
										Year:   1984,
									},
									testTrack{
										Name:   "Track Two",
										Artist: "Artist One",
										Year:   2000,
									},
									testTrack{
										Name:   "Track One",
										Artist: "Artist One",
										Year:   1984,
									},
								},
							},
						},
					},
				},
			},
			fields: []Attr{StringAttr("Artist"), IntAttr("Year")},
			out:    []interface{}{"Artist One", nil},
		},
	}

	for ii, tt := range table {
		g := CommonGroupAttr(tt.fields, tt.in)
		got := make([]interface{}, len(tt.out))
		for i, f := range tt.fields {
			got[i] = g.Field(f.Field())
		}

		if !reflect.DeepEqual(got, tt.out) {
			t.Errorf("[%d] got %#v, expected %#v", ii, got, tt.out)
		}
	}
}

func TestFirstTrackAttr(t *testing.T) {
	table := []struct {
		in    Group
		field Attr
		out   interface{}
	}{
		// One group with no tracks
		{
			in: group{
				name: "Group One",
			},
			field: StringAttr("Name"),
			out:   nil,
		},

		{
			in: group{
				name: "Group One",
			},
			field: StringsAttr("Artist"),
			out:   nil,
		},

		// One group with one track
		{
			in: group{
				name: "Group One",
				tracks: []Track{
					testTrack{
						Name: "Track One",
					},
				},
			},
			field: StringAttr("Name"),
			out:   "Track One",
		},

		{
			in: group{
				name: "Group One",
				tracks: []Track{
					testTrack{
						Artist: "Track One",
					},
				},
			},
			field: StringsAttr("Artist"),
			out:   []string{"Track One"},
		},

		// One group with two tracks, empty first field
		{
			in: group{
				name: "Group One",
				tracks: []Track{
					testTrack{},
					testTrack{
						Name:   "Track Two",
						Artist: "Artist One",
					},
				},
			},
			field: StringAttr("Artist"),
			out:   nil,
		},

		{
			in: group{
				name: "Group One",
				tracks: []Track{
					testTrack{},
					testTrack{
						Name:   "Track Two",
						Artist: "Artist One",
					},
				},
			},
			field: StringsAttr("Artist"),
			out:   []string(nil),
		},

		// One collection, one group, one track
		{
			in: testCol{
				name: "Group One (collection)",
				keys: []Key{"Group-One-One"},
				grps: map[Key]Group{
					"Group-One-One": group{
						name: "Group One (One)",
						tracks: []Track{
							testTrack{
								Name: "Track One",
							},
						},
					},
				},
			},
			field: StringAttr("Name"),
			out:   "Track One",
		},

		// One group with one track
		{
			in: group{
				name: "Group One",
				tracks: []Track{
					testTrack{
						Duration: 1234,
					},
				},
			},
			field: IntAttr("Duration"),
			out:   1234,
		},

		// One group with two tracks, empty first field
		{
			in: group{
				name: "Group One",
				tracks: []Track{
					testTrack{},
					testTrack{
						Duration: 1234,
						Artist:   "Artist One",
					},
				},
			},
			field: IntAttr("Duration"),
			out:   nil,
		},

		// One collection, one group, one track
		{
			in: testCol{
				name: "Group One (collection)",
				keys: []Key{"Group-One-One"},
				grps: map[Key]Group{
					"Group-One-One": group{
						name: "Group One (One)",
						tracks: []Track{
							testTrack{
								Duration: 1234,
							},
						},
					},
				},
			},
			field: IntAttr("Duration"),
			out:   1234,
		},
	}

	for ii, tt := range table {
		g := FirstTrackAttr(tt.field, tt.in)
		got := g.Field(tt.field.Field())

		if !reflect.DeepEqual(got, tt.out) {
			t.Errorf("[%d] got %#v, expected %#v", ii, got, tt.out)
		}
	}
}
