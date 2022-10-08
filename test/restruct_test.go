package restruct_test

import (
	"testing"

	r "github.com/fclairamb/restruct"
	"github.com/stretchr/testify/assert"
)

func TestStandard(t *testing.T) {
	a := assert.New(t)

	type Human struct {
		Name   string
		Age    int
		Height *int
	}

	rs := &r.Restruct{
		RegexToStructs: []*r.RegexToStruct{
			{
				ID:     "age",
				Regex:  `^(?P<name>\w+) is ((?P<age>\d+)( years old)?|old|great)$`,
				Struct: &Human{},
			},
			{
				ID:     "height",
				Regex:  `^(?P<name>\w+) is (?P<height>\d+) cm tall$`,
				Struct: &Human{},
			},
		},
	}

	t.Run("pointer", func(t *testing.T) {
		m, err := rs.MatchString("John is 178 cm tall")
		a.NoError(err)
		a.NotNil(m)
		a.Equal("height", m.ID)

		th, ok := m.Struct.(*Human)
		a.True(ok)
		a.Equal("John", th.Name)
		a.NotNil(th.Height)
		a.Equal(178, *th.Height)
	})

	// We have a match
	t.Run("direct", func(t *testing.T) {
		for _, str := range []string{"John is 42 years old", "John is 42"} {
			m, err := rs.MatchString(str)
			a.NoError(err)
			a.NotNil(m)
			a.Equal("age", m.ID)

			th, ok := m.Struct.(*Human)
			a.True(ok)
			a.Equal("John", th.Name)
			a.Equal(42, th.Age)
		}
	})

	t.Run("optional", func(t *testing.T) {
		for _, str := range []string{"John is old", "John is great"} {
			m, err := rs.MatchString(str)
			a.NoError(err)
			a.NotNil(m)

			th, ok := m.Struct.(*Human)
			a.True(ok)
			a.Equal("John", th.Name)
			a.Equal(0, th.Age)
		}
	})

	// We don't have a match
	m, err := rs.MatchString("John was 42 years old")
	a.NoError(err)
	a.Nil(m)
}

func TestFillingError(t *testing.T) {
	a := assert.New(t)

	type Something struct {
		A int
	}

	rs := &r.Restruct{
		RegexToStructs: []*r.RegexToStruct{
			{
				Regex:  `(?P<a>\w+)`,
				Struct: &Something{},
			},
		},
	}

	m, err := rs.MatchString("abc")
	a.Error(err)
	a.Nil(m)
}

func TestTypeIgnore(t *testing.T) {
	a := assert.New(t)

	type StA struct {
		A int `restruct:"-"`
		B int `restruct:"a"`
	}

	type StB struct {
		B int `restruct:"a"`
		A int `restruct:"-"`
	}

	rs := &r.Restruct{
		RegexToStructs: []*r.RegexToStruct{
			{
				Regex:  `a:(?P<a>\d+)`,
				Struct: &StA{},
			},
			{
				Regex:  `b:(?P<a>\d+)`,
				Struct: &StB{},
			},
		},
	}

	m, err := rs.MatchString("a:123")
	a.NoError(err)
	a.NotNil(m)
	a.Equal(123, m.Struct.(*StA).B)
	a.Equal(0, m.Struct.(*StA).A)

	m, err = rs.MatchString("b:123")
	a.NoError(err)
	a.NotNil(m)
	a.Equal(123, m.Struct.(*StB).B)
	a.Equal(0, m.Struct.(*StB).A)
}

func TestTypeFloat(t *testing.T) {
	a := assert.New(t)

	type Something struct {
		A float64
		B float32
		C *float64
		D *float32
	}

	rs := &r.Restruct{
		RegexToStructs: []*r.RegexToStruct{
			{
				Regex:  `((A:(?P<a>[\w\.]+)|B:(?P<b>[\w\.]+)|C:(?P<c>[\w\.]+)|D:(?P<d>[\w\.]+))\s*)+`,
				Struct: &Something{},
			},
		}}

	m, err := rs.MatchString("A:abc") // bad float64
	a.Error(err)
	a.Nil(m)

	m, err = rs.MatchString("B:abc") // bad float32
	a.Error(err)
	a.Nil(m)

	m, err = rs.MatchString("A:1.2 B:3.4 C:5.6 D:7.8")
	a.NoError(err)
	a.NotNil(m)
	s := m.Struct.(*Something)
	a.Equal(1.2, s.A)
	a.Equal(float32(3.4), s.B)
	a.Equal(5.6, *s.C)
	a.Equal(float32(7.8), *s.D)

	m, err = rs.MatchString("A:9")
	a.NoError(err)
	a.NotNil(m)
	s = m.Struct.(*Something)
	a.Equal(9.0, s.A)
	a.Equal(float32(0), s.B)
	a.Nil(s.C)
	a.Nil(s.D)
}

func TestTypeInt(t *testing.T) {
	a := assert.New(t)

	type Something struct {
		A int
		B *int
	}

	rs := &r.Restruct{
		RegexToStructs: []*r.RegexToStruct{
			{
				Regex:  `((A:(?P<a>[\w\.]+)|B:(?P<b>[\w\.]+))\s*)+`,
				Struct: &Something{},
			},
		}}

	m, err := rs.MatchString("A:abc")
	a.Error(err)
	a.Nil(m)

	m, err = rs.MatchString("A:1 B:2")
	a.NoError(err)
	a.NotNil(m)
	s := m.Struct.(*Something)
	a.Equal(1, s.A)
	a.Equal(2, *s.B)
}

func TestTypeBool(t *testing.T) {
	a := assert.New(t)

	type Something struct {
		A bool
		B *bool
	}

	rs := &r.Restruct{
		RegexToStructs: []*r.RegexToStruct{
			{
				Regex:  `((A:(?P<a>[\w\.]+)|B:(?P<b>[\w\.]+))\s*)+`,
				Struct: &Something{},
			},
		}}

	m, err := rs.MatchString("A:abc")
	a.Error(err)
	a.Nil(m)

	m, err = rs.MatchString("A:true B:true")
	a.NoError(err)
	a.NotNil(m)
	s := m.Struct.(*Something)
	a.True(s.A)
	a.True(*s.B)
}

func TestBadRegex(t *testing.T) {
	a := assert.New(t)

	rs := &r.Restruct{
		RegexToStructs: []*r.RegexToStruct{
			{
				Regex: `(?P<\w+) is old`,
			},
		},
	}

	{
		errType := &r.CompilationError{}
		a.ErrorAs(rs.Compile(), &errType)
	}

	m, err := rs.MatchString("anything")
	a.Error(err)
	a.ErrorContains(err, "could not compile rule: error parsing regexp:")
	a.Nil(m)
}

func TestNotAPointer(t *testing.T) {
	a := assert.New(t)

	type Something struct {
		A string
	}

	rs := &r.Restruct{
		RegexToStructs: []*r.RegexToStruct{
			{
				Regex:  `(?P<a>\w+)`,
				Struct: Something{},
			},
		},
	}

	_, err := rs.MatchString("anything")
	a.ErrorIs(err, r.ErrNotAStructPointer)
}

func TestNotAStructPointer(t *testing.T) {
	a := assert.New(t)

	rs := &r.Restruct{
		RegexToStructs: []*r.RegexToStruct{
			{
				Regex:  `(?P<a>\w+)`,
				Struct: []string{},
			},
		},
	}

	_, err := rs.MatchString("anything")
	a.ErrorIs(err, r.ErrNotAStructPointer)
}

func TestBadField(t *testing.T) {
	a := assert.New(t)

	type Something struct {
		A int
	}

	rs := &r.Restruct{
		RegexToStructs: []*r.RegexToStruct{
			{
				Regex:  `(?P<a>\w+)`,
				Struct: &Something{},
			},
		},
	}

	m, err := rs.MatchString("anything")
	a.Error(err)
	a.ErrorContains(err, "could not fill field A: strconv.Atoi")
	a.Nil(m)
}

func BenchmarkSmallStruct(b *testing.B) {
	type Human struct {
		Name string
		Age  int
	}

	rs := &r.Restruct{
		RegexToStructs: []*r.RegexToStruct{
			{
				ID:     "age",
				Regex:  `(?P<name>\w+) is (?P<age>\d+) years old`,
				Struct: &Human{},
			},
		},
	}

	for i := 0; i < b.N; i++ {
		_, _ = rs.MatchString("Johh is 42 years old")
	}
}

func BenchmarkLoadAndExec(b *testing.B) {
	type Human struct {
		Name string
		Age  int
	}

	rules := []*r.RegexToStruct{
		{
			ID:     "age",
			Regex:  `(?P<name>\w+) is (?P<age>\d+) years old`,
			Struct: &Human{},
		},
	}

	for i := 0; i < b.N; i++ {
		rs := &r.Restruct{
			RegexToStructs: rules,
		}
		_, _ = rs.MatchString("Johh is 42 years old")
	}
}

func BenchmarkThreeRules(b *testing.B) {
	type Human struct {
		Name string
		Age  int
	}

	rs := &r.Restruct{
		RegexToStructs: []*r.RegexToStruct{
			{
				ID:     "height",
				Regex:  `^(?P<name>\w+) is (?P<height>\d+) cm tall$`,
				Struct: &Human{},
			},
			{
				ID:     "weight",
				Regex:  `^(?P<name>\w+) weighs (?P<weight>\d+) kg$`,
				Struct: &Human{},
			},
			{
				ID:     "age",
				Regex:  `(?P<name>\w+) is (?P<age>\d+) years old`,
				Struct: &Human{},
			},
		},
	}

	for i := 0; i < b.N; i++ {
		_, _ = rs.MatchString("Johh is 42 years old")
	}
}

func BenchmarkBiggerStruct(b *testing.B) {
	type Human struct {
		Name   string
		Age    int
		Height *int
		Male   bool
	}

	rs := &r.Restruct{
		RegexToStructs: []*r.RegexToStruct{
			{
				ID:     "age",
				Regex:  `(?P<name>\w+) is (?P<age>\d+) years old`,
				Struct: &Human{},
			},
		},
	}

	for i := 0; i < b.N; i++ {
		_, _ = rs.MatchString("Johh is 42 years old")
	}
}
