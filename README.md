# Regex to struct

## General idea
This is a very simple library that allows you to convert a regex into a struct. It's intended to be used for simple text parsing around 
dummy bots.

The struct will have a field for each capture group in the regex.

## Usage

```go
import(
    	r "github.com/fclairamb/restruct"
)

type Human struct {
    Name   string `restruct:"name"` // Specifying the field
    Age    int  // No tag, "age" will be used
    Height *int // A pointer will be set to nil if the capture group is empty
}

rs := &r.Restruct{
    RegexToStructs: []*r.RegexToStruct{
        {
            ID:     "age",
            Regex:  `^(?P<name>\w+) is ((?P<age>\d+)( years old)?$`,
            Struct: &Human{},
        },
        {
            ID:     "height",
            Regex:  `^(?P<name>\w+) is (?P<height>\d+) cm tall$`,
            Struct: &Human{},
        },
    },
}

m := rs.Match("John is 42 years old")
if m != nil {
    h := m.Struct.(*Human)
    fmt.Printf("name = %s, age = %d", h.Name, h.Age)
}
```