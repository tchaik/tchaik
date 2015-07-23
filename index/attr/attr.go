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
	Field() string
	IsEmpty(x interface{}) bool
	Value(Getter) interface{}
	Intersect(x, y interface{}) interface{}
}

type valueType struct {
	field string
	empty interface{}
	get   func(Getter) interface{}
}

func (v *valueType) Field() string {
	return v.field
}

func (v *valueType) IsEmpty(x interface{}) bool {
	return v.empty == x
}

func (v *valueType) Value(g Getter) interface{} {
	return v.get(g)
}

func (v *valueType) Intersect(x, y interface{}) interface{} {
	if x == y {
		return x
	}
	return v.empty
}

func String(f string) Interface {
	return &valueType{
		field: f,
		empty: "",
		get: func(g Getter) interface{} {
			return g.GetString(f)
		},
	}
}

func Int(f string) Interface {
	return &valueType{
		field: f,
		empty: 0,
		get: func(g Getter) interface{} {
			return g.GetInt(f)
		},
	}
}

type stringsType struct {
	valueType
}

func (p *stringsType) IsEmpty(x interface{}) bool {
	if x == nil {
		return true
	}
	xs := x.([]string)
	return len(xs) == 0
}

func (p *stringsType) Intersect(x, y interface{}) interface{} {
	if x == nil || y == nil {
		return nil
	}
	xs := x.([]string)
	ys := y.([]string)
	return stringSliceIntersect(xs, ys)
}

func Strings(f string) Interface {
	return &stringsType{
		valueType{
			field: f,
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
