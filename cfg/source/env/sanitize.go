package sourceenv

import (
	"strings"
)

// SanitizeName normalizes s into a POSIX-portable environment variable name.
//
// POSIX portability rules this function enforces:
//   - First rune: ASCII letter or '_' (never a digit or other character)
//   - Subsequent runes: only ASCII letters, digits, or '_' (no spaces, hyphens,
//     dots, Unicode, or '=' anywhere in the name)
//
// The returned string is safe for use with os.Setenv and in exec environments
// on POSIX systems.
func SanitizeName(s string) string {
	if s == "" {
		return "_"
	}

	s = strings.ToUpper(s)

	{ // replace unhandled runes with _
		var b strings.Builder
		b.Grow(len(s))

		for i, r := range s {
			if (r >= 'A' && r <= 'Z') || (i != 0 && r >= '0' && r <= '9') {
				b.WriteRune(r)
			} else {
				b.WriteRune('_')
			}
		}

		s = b.String()
	}

	for strings.Contains(s, "__") {
		s = strings.ReplaceAll(s, "__", "_")
	}

	s = strings.Trim(s, "_")
	if s == "" {
		return "_"
	}

	if s[0] >= '0' && s[0] <= '9' {
		s = "_" + s
	}

	return s
}
