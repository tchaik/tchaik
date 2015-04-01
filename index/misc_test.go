// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package index

import (
	"reflect"
	"testing"
)

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
			in: col{
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
			in: col{
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
					"Group-Three": col{
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
		// One group with one track, empty common string field
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
			out:    []interface{}{""},
		},

		// One group with one track, empty common int field
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
			out:    []interface{}{0},
		},

		// One group with one track, empty common string and int fields
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
			out:    []interface{}{"", 0},
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
			out:    []interface{}{"", 0},
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
			out:    []interface{}{0},
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
			out:    []interface{}{"", "Composer One"},
		},

		// One collection, one group, one track
		{
			in: col{
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
			out:    []interface{}{"Artist One", ""},
		},

		// One collection, three groups, many tracks!
		{
			in: col{
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
					"Group-Three": col{
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
			out:    []interface{}{"Artist One", 0},
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
			out:   "",
		},

		// One collection, one group, one track
		{
			in: col{
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
			out:   0,
		},

		// One collection, one group, one track
		{
			in: col{
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

		if got != tt.out {
			t.Errorf("[%d] got %#v, expected %#v", ii, got, tt.out)
		}
	}
}
