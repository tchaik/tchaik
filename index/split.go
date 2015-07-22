package index

import "strings"

type subGrpTrks struct {
	Group
	tracks []Track
}

func (g *subGrpTrks) Tracks() []Track {
	if g.tracks != nil {
		return g.tracks
	}
	return g.Group.Tracks()
}

// splitMultiple applies strings.Split to `x` with each strings in `s`
// recursively.
func splitMultiple(x string, s []string) []string {
	if len(x) == 0 {
		return nil
	}
	if len(s) == 0 {
		return []string{x}
	}
	r := strings.Split(x, s[0])
	for _, sx := range s[1:] {
		t := make([]string, 0, len(r))
		for _, ry := range r {
			t = append(t, strings.Split(ry, sx)...)
		}
		r = t
	}
	res := make([]string, 0, len(r))
	for _, x := range r {
		y := strings.TrimSpace(x)
		if y != "" {
			res = append(res, y)
		}
	}
	return res
}

// SplitNameList is a transform function which goes through the tracks
// in the group and
func SplitNameList(fields ...string) func(Group) Group {
	return func(g Group) Group {
		tracks := g.Tracks()
		nt := splitNameList(fields, tracks)
		return &subGrpTrks{
			Group:  g,
			tracks: nt,
		}
	}
}

type stringsTrack struct {
	Track
	m map[string][]string
}

func (s *stringsTrack) GetStrings(k string) []string {
	v, ok := s.m[k]
	if !ok {
		return s.Track.GetStrings(k)
	}
	return v
}

var splitNameSeparator = []string{"/", ",", ";", ":", "&", "and"}

func splitNameList(fields []string, tracks []Track) []Track {
	result := make([]Track, len(tracks))
	for i, t := range tracks {
		m := make(map[string][]string)
		for _, f := range fields {
			m[f] = splitMultiple(t.GetString(f), splitNameSeparator)
		}
		result[i] = &stringsTrack{
			Track: t,
			m:     m,
		}
	}
	return result
}
