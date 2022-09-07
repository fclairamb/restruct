// Package restruct provides a way to match a string with a regex and fill a struct with the result
package restruct

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	float32Bits = 32
	float64Bits = 64
)

// RegexToStruct defines the link between the regex and the struct
type RegexToStruct struct {
	ID            string         // ID of the match (optional)
	Regex         string         // Regex to match (before compilation)
	Struct        interface{}    // Struct instance to fill (can be shared between multiple RegexToStruct)
	compiledRegex *regexp.Regexp // Compiled regex
	subexToField  map[int]int    // Map of subexpression index to field index
}

func (r *RegexToStruct) compile() error {
	compiledRegex, err := regexp.Compile(r.Regex)
	if err != nil {
		return &CompilationError{Err: err, ID: r.ID}
	}

	r.compiledRegex = compiledRegex

	r.subexToField = make(map[int]int)

	reNameToID := make(map[string]int)

	for i, name := range r.compiledRegex.SubexpNames() {
		if name == "" {
			continue
		}

		reNameToID[name] = i
	}

	typeT := reflect.TypeOf(r.Struct)

	if typeT.Kind() == reflect.Ptr {
		typeT = typeT.Elem()
	}

	nbFields := typeT.NumField()
	for i := 0; i < nbFields; i++ {
		var tagValue string
		{
			field := typeT.Field(i)
			tagValue = field.Tag.Get("restruct")

			if tagValue == "" {
				tagValue = strings.ToLower(field.Name)
			} else if tagValue == "-" {
				continue
			}
		}

		if reID, ok := reNameToID[tagValue]; ok {
			r.subexToField[reID] = i
		}
	}

	return nil
}

// Restruct is the core type
type Restruct struct {
	RegexToStructs []*RegexToStruct
	compiled       bool
}

// FieldFillingError is an error that occurs when filling a field
type FieldFillingError struct {
	FieldName string
	Err       error
}

func (e *FieldFillingError) Error() string {
	return fmt.Sprintf("could not fill field %s: %s", e.FieldName, e.Err)
}

// CompilationError is an error that occurs when compiling rules
type CompilationError struct {
	Err error
	ID  string
}

func (e *CompilationError) Error() string {
	return fmt.Sprintf("could not compile rule: %s", e.Err)
}

// compile compiles the regexes
func (r *Restruct) compile() error {
	for _, regexToStruct := range r.RegexToStructs {
		if err := regexToStruct.compile(); err != nil {
			return err
		}
	}

	r.compiled = true

	return nil
}

// MatchString will return a possible match with a filled struct
func (r *Restruct) MatchString(s string) (*RegexToStruct, error) {
	if !r.compiled {
		if err := r.compile(); err != nil {
			return nil, err
		}
	}

	for _, regexToStruct := range r.RegexToStructs {
		in, err := r.matchRegexString(regexToStruct, s)
		if err != nil {
			return nil, err
		}

		if in != nil {
			return regexToStruct, err
		}
	}

	return nil, nil
}

func (r *Restruct) matchRegexString(rs *RegexToStruct, str string) (interface{}, error) {
	match := rs.compiledRegex.FindStringSubmatch(str)
	if match == nil {
		return nil, nil
	}

	return rs.fillStruct(rs.Struct, match)
}

func (r *RegexToStruct) fillStruct(s interface{}, match []string) (interface{}, error) {
	for i, reValue := range match {
		fieldIndex, ok := r.subexToField[i]
		if !ok {
			continue
		}

		stValue := reflect.ValueOf(s).Elem().Field(fieldIndex)

		if reValue == "" {
			stValue.Set(reflect.Zero(stValue.Type()))

			continue
		}

		if err := fillField(stValue, reValue); err != nil {
			return nil, &FieldFillingError{FieldName: stValue.Type().Name(), Err: err}
		}
	}

	return s, nil
}

func fillField(stValue reflect.Value, reValue string) error {
	switch stValue.Kind() { //nolint:exhaustive
	case reflect.String:
		stValue.SetString(reValue)

	case reflect.Int:
		intValue, err := strconv.Atoi(reValue)
		if err != nil {
			return err //nolint:wrapcheck // no need to wrap as it's done later
		}

		stValue.SetInt(int64(intValue))

	case reflect.Float64:
		floatValue, err := strconv.ParseFloat(reValue, float64Bits)
		if err != nil {
			return err //nolint:wrapcheck // no need to wrap as it's done later
		}

		stValue.SetFloat(floatValue)

	case reflect.Float32:
		floatValue, err := strconv.ParseFloat(reValue, float32Bits)
		if err != nil {
			return err //nolint:wrapcheck // no need to wrap as it's done later
		}

		stValue.SetFloat(floatValue)

	case reflect.Bool:
		boolValue, err := strconv.ParseBool(reValue)
		if err != nil {
			return err //nolint:wrapcheck // no need to wrap as it's done later
		}

		stValue.SetBool(boolValue)

	case reflect.Pointer:
		if stValue.IsNil() {
			stValue.Set(reflect.New(stValue.Type().Elem()))
		}

		return fillField(stValue.Elem(), reValue)
	}

	return nil
}
