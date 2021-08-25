// Package app exposes app-related information that are injected at compile time.
package app

import (
	"errors"
	"fmt"
	"time"
	"unicode"
)

// build injection requires these global variables
// nolint: gochecknoglobals
var (
	name       = "app"                  // app name
	version    = "0.0.0-0-gmaster"      // see https://semver.org/ for a description of the format
	buildAtRaw = "1970-01-01T00:00:00Z" // build date in RFC3339 format

	app *App // initialized through init() function
)

// this is required as it depends on build-time variable injection
// nolint: gochecknoinits
func init() { Init(name, version, buildAtRaw) }

// App stores application build data.
type App struct {
	Name             string
	Version          string
	AlphaNumericName string
	BuiltAt          time.Time
}

// New creates a new version of app.
// It is usually not manually created (see init() function).
func New(name, version, buildAtRaw string) (*App, error) {
	var (
		app = App{
			Name:    name,
			Version: version,
		}
		err error
	)

	if app.Name == "" {
		return nil, errors.New("app name cannot be empty")
	}

	if app.Version == "" {
		return nil, errors.New("app version cannot be empty")
	}

	app.BuiltAt, err = time.Parse(time.RFC3339, buildAtRaw)
	if err != nil {
		return nil, fmt.Errorf("buildAtRaw must respect RFC3339 format")
	}

	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			app.AlphaNumericName += string(r)
		}
	}

	return &app, nil
}

// Init resets the global app with the provided values.
// It should only be set once, during main initialization.
func Init(name, version, buildAtRaw string) {
	var err error

	app, err = New(name, version, buildAtRaw)
	if err != nil {
		panic(fmt.Errorf("unable to initialize app package: %w", err))
	}
}

// Copy returns a copy of globally initialized app.
func Copy() *App {
	app := *app
	return &app
}

// Name returns the app name as set during app Init.
func Name() string { return app.Name }

// AlphaNumericName returns the app name with only alphanumeric characters.
func AlphaNumericName() string { return app.AlphaNumericName }

// Version returns the app version as set during app Init.
func Version() string { return app.Version }

// BuiltAt returns the app built time as set during app Init.
func BuiltAt() time.Time { return app.BuiltAt }
