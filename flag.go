package cli

import (
	"time"
)

// Flag defines how to configure a flag.
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

// FlagCustom creates a flag.
func FlagCustom(name string, shorthand string, dest interface{}, description string) Flag {
	return &flag{name: name, shorthand: shorthand, description: description, destination: dest}
}

// FlagBool defines a boolean flag.
func FlagBool(name string, shorthand string, dest *bool, description string) Flag {
	return FlagCustom(name, shorthand, dest, description)
}

// FlagBoolSlice defines a boolean slice flag.
func FlagBoolSlice(name string, shorthand string, dest *[]bool, description string) Flag {
	return FlagCustom(name, shorthand, dest, description)
}

// FlagString defines a string flag.
func FlagString(name string, shorthand string, dest *string, description string) Flag {
	return FlagCustom(name, shorthand, dest, description)
}

// FlagStringSlice defines a string slice flag.
func FlagStringSlice(name string, shorthand string, dest *[]string, description string) Flag {
	return FlagCustom(name, shorthand, dest, description)
}

// FlagInt defines a int flag.
func FlagInt(name string, shorthand string, dest *int, description string) Flag {
	return FlagCustom(name, shorthand, dest, description)
}

// FlagIntSlice defines a int slice flag.
func FlagIntSlice(name string, shorthand string, dest *[]int, description string) Flag {
	return FlagCustom(name, shorthand, dest, description)
}

// FlagUint defines a uint flag.
func FlagUint(name string, shorthand string, dest *uint, description string) Flag {
	return FlagCustom(name, shorthand, dest, description)
}

// FlagUIntSlice defines a uint slice flag.
func FlagUIntSlice(name string, shorthand string, dest *[]uint, description string) Flag {
	return FlagCustom(name, shorthand, dest, description)
}

// FlagFloat32 defines a float32 flag.
func FlagFloat32(name string, shorthand string, dest *float32, description string) Flag {
	return FlagCustom(name, shorthand, dest, description)
}

// FlagFloat32Slice defines a float32 slice flag.
func FlagFloat32Slice(name string, shorthand string, dest *[]float32, description string) Flag {
	return FlagCustom(name, shorthand, dest, description)
}

// FlagFloat64 defines a float64 flag.
func FlagFloat64(name string, shorthand string, dest *float64, description string) Flag {
	return FlagCustom(name, shorthand, dest, description)
}

// FlagFloat64Slice defines a float64 slice flag.
func FlagFloat64Slice(name string, shorthand string, dest *[]float64, description string) Flag {
	return FlagCustom(name, shorthand, dest, description)
}

// FlagDuration defines a time.Duration flag.
func FlagDuration(name string, shorthand string, dest *time.Duration, description string) Flag {
	return FlagCustom(name, shorthand, dest, description)
}

// FlagDurationSlice defines a duration slice flag.
func FlagDurationSlice(name string, shorthand string, dest *[]time.Duration, description string) Flag {
	return FlagCustom(name, shorthand, dest, description)
}
