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

func FlagCustom(name string, shorthand string, dest interface{}, description string) Flag {
	return &flag{name: name, shorthand: shorthand, description: description, destination: dest}
}

func FlagBool(name string, shorthand string, dest *bool, description string) Flag {
	return FlagCustom(name, shorthand, dest, description)
}

func FlagBoolSlice(name string, shorthand string, dest *[]bool, description string) Flag {
	return FlagCustom(name, shorthand, dest, description)
}

func FlagString(name string, shorthand string, dest *string, description string) Flag {
	return FlagCustom(name, shorthand, dest, description)
}

func FlagStringSlice(name string, shorthand string, dest *[]string, description string) Flag {
	return FlagCustom(name, shorthand, dest, description)
}

func FlagInt(name string, shorthand string, dest *int, description string) Flag {
	return FlagCustom(name, shorthand, dest, description)
}

func FlagIntSlice(name string, shorthand string, dest *[]int, description string) Flag {
	return FlagCustom(name, shorthand, dest, description)
}

func FlagUint(name string, shorthand string, dest *uint, description string) Flag {
	return FlagCustom(name, shorthand, dest, description)
}

func FlagUIntSlice(name string, shorthand string, dest *[]uint, description string) Flag {
	return FlagCustom(name, shorthand, dest, description)
}

func FlagFloat32(name string, shorthand string, dest *float32, description string) Flag {
	return FlagCustom(name, shorthand, dest, description)
}

func FlagFloat32Slice(name string, shorthand string, dest *[]float32, description string) Flag {
	return FlagCustom(name, shorthand, dest, description)
}

func FlagFloat64(name string, shorthand string, dest *float64, description string) Flag {
	return FlagCustom(name, shorthand, dest, description)
}

func FlagFloat64Slice(name string, shorthand string, dest *[]float64, description string) Flag {
	return FlagCustom(name, shorthand, dest, description)
}

func FlagDuration(name string, shorthand string, dest *time.Duration, description string) Flag {
	return FlagCustom(name, shorthand, dest, description)
}

func FlagDurationSlice(name string, shorthand string, dest *[]time.Duration, description string) Flag {
	return FlagCustom(name, shorthand, dest, description)
}
