package cli

// Flag represents a command-line flag. It combines the FlagValuer interface
// with methods to access flag metadata (long name, short name, description).
type Flag interface {
	// FlagValuer defines methods to get and set flag's value.
	FlagValuer

	// LongName returns the long name of the flag (e.g., "--verbose").
	LongName() string
	// ShortName returns the short name of the flag (e.g., "-v"). May be empty.
	ShortName() string
	// Description returns a description of the flag's purpose.
	Description() string
}

// NewFlag creates a new Flag instance.
//
//	longName is the long flag name, like --longname ; cannot be empty.
//	shortName is the short flag name ; usually 1 character, like -s ; can be empty.
//	valuer provide the way to set value to the destination.
//	description is a short text explaining the flag ; can be empty.
//
// It panics if invalid inputs are provided.
func NewFlag(longName, shortName string, valuer FlagValuer, description string) Flag {
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
