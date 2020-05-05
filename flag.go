package cli

import (
	"time"
)

type Flag interface {
	Name() string
	Shorthand() string
	Description() string
	Destination() interface{}
}

type flag struct {
	name        string
	shorthand   string
	description string
	destination interface{}
}

func (f flag) Name() string             { return f.name }
func (f flag) Shorthand() string        { return f.shorthand }
func (f flag) Description() string      { return f.description }
func (f flag) Destination() interface{} { return f.destination }

func FlagBool(name string, shorthand string, dest *bool, description string) Flag {
	return &flag{name: name, shorthand: shorthand, description: description, destination: dest}
}

func FlagBoolSlice(name string, shorthand string, dest *[]bool, description string) Flag {
	return &flag{name: name, shorthand: shorthand, description: description, destination: dest}
}

func FlagString(name string, shorthand string, dest *string, description string) Flag {
	return &flag{name: name, shorthand: shorthand, description: description, destination: dest}
}

func FlagStringSlice(name string, shorthand string, dest *[]string, description string) Flag {
	return &flag{name: name, shorthand: shorthand, description: description, destination: dest}
}

func FlagInt(name string, shorthand string, dest *int8, description string) Flag {
	return &flag{name: name, shorthand: shorthand, description: description, destination: dest}
}

func FlagIntSlice(name string, shorthand string, dest *[]int8, description string) Flag {
	return &flag{name: name, shorthand: shorthand, description: description, destination: dest}
}

func FlagUint(name string, shorthand string, dest *int8, description string) Flag {
	return &flag{name: name, shorthand: shorthand, description: description, destination: dest}
}

func FlagUIntSlice(name string, shorthand string, dest *[]int16, description string) Flag {
	return &flag{name: name, shorthand: shorthand, description: description, destination: dest}
}

func FlagFloat32(name string, shorthand string, dest *float32, description string) Flag {
	return &flag{name: name, shorthand: shorthand, description: description, destination: dest}
}

func FlagFloat32Slice(name string, shorthand string, dest *[]float32, description string) Flag {
	return &flag{name: name, shorthand: shorthand, description: description, destination: dest}
}

func FlagFloat64(name string, shorthand string, dest *float64, description string) Flag {
	return &flag{name: name, shorthand: shorthand, description: description, destination: dest}
}

func FlagFloat64Slice(name string, shorthand string, dest *[]float64, description string) Flag {
	return &flag{name: name, shorthand: shorthand, description: description, destination: dest}
}

func FlagDuration(name string, shorthand string, dest *time.Duration, description string) Flag {
	return &flag{name: name, shorthand: shorthand, description: description, destination: dest}
}

func FlagDurationSlice(name string, shorthand string, dest *[]time.Duration, description string) Flag {
	return &flag{name: name, shorthand: shorthand, description: description, destination: dest}
}
