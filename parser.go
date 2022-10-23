package parser

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/google/shlex"
	"github.com/jessevdk/go-flags"
)

// ParseQuery is a convenience function that parses a string into flags and
// a single arg. Flags are set using reflection (flags.ParseArgs) on an input
// struct. Any such flags are then parsed out of the original input before
// it is returned.
func ParseQuery(arg string, opts interface{}) string {
	argv, _ := shlex.Split(arg)
	flags.ParseArgs(opts, argv)
	return query(opts, arg)
}

// query removes string and bool flag types from an input string.
func query(opts interface{}, arg string) string {
	optsType := reflect.TypeOf(opts).Elem()
	optsValue := reflect.ValueOf(opts).Elem()

	for i := 0; i < optsType.NumField(); i++ {
		short, long := tags(optsType.Field(i).Tag)
		value := optsValue.Field(i)

		switch value.Kind() {
		case reflect.Bool:
			flagFormat := "-%s"
			if short != "" {
				arg = strings.ReplaceAll(arg, fmt.Sprintf(flagFormat, short), "")
			}
			if long != "" {
				arg = strings.ReplaceAll(arg, fmt.Sprintf(flagFormat, long), "")
			}
		case reflect.String:
			s := value.String()
			flagFormat := stringFlagFormat(s)

			if short != "" {
				arg = strings.ReplaceAll(arg, fmt.Sprintf(flagFormat, short, s), "")
			}
			if long != "" {
				arg = strings.ReplaceAll(arg, fmt.Sprintf(flagFormat, long, s), "")
			}
		default:
			// TODO: Add support for other flag types
		}
	}

	return strings.TrimSpace(arg)
}

// tags returns string values for `short` and `long` struct tags.
func tags(st reflect.StructTag) (string, string) {
	short := st.Get("short")
	long := st.Get("long")
	return short, long
}

// stringFlagFormat returns a format string for an input string.
func stringFlagFormat(s string) string {
	flagFormat := ""
	if strings.Contains(s, " ") {
		flagFormat = "-%s=\"%s\""
	} else {
		flagFormat = "-%s=%s"
	}
	return flagFormat
}
