package cli

// Flag defines a flagValue.
type Flag interface {
	// FlagValuer defines methods to get and set flagValue value.
	FlagValuer

	// LongName if non-empty, defines the long name of the flagValue like --long.
	LongName() string
	// ShortName if non-empty, defines the short name of the flagValue like -s.
	ShortName() string
	// Description describes the flagValue for display purpose.
	Description() string
}

// NewCustomFlag creates a flagValue based on any underlying destination type.
//
//	longName is the long flagValue name, like --longname ; cannot be empty.
//	shortName is the short flagValue name ; usually 1 character, like -s ; can be empty.
//	description is a short text explaining the flagValue ; can be empty.
//	valuer provide the way to set value to the destination.
func NewCustomFlag(longName, shortName string, valuer FlagValuer, description string) Flag {
	if longName == "" && shortName == "" {
		panic("longName and/or shortName must be non-empty")
	}

	if shortName != "" && len(shortName) > 1 {
		panic("shortName must be one character long")
	}

	if valuer == nil {
		panic("a non-nil valuer is required")
	}

	return &flagValue{
		FlagValuer: valuer,

		longName:    longName,
		shortName:   shortName,
		description: description,
	}
}

type flagValue struct {
	FlagValuer

	longName    string
	shortName   string
	description string
}

func (f flagValue) LongName() string    { return f.longName }
func (f flagValue) ShortName() string   { return f.shortName }
func (f flagValue) Description() string { return f.description }
