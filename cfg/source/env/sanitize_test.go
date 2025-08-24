package sourceenv

import (
	"strings"
	"testing"

	"github.com/krostar/test"
	"github.com/krostar/test/check"
)

func Test_SanitizeName(t *testing.T) {
	for name, tc := range map[string]struct {
		input    string
		expected string
	}{
		"empty string": {
			input:    "",
			expected: "_",
		},
		"single character letter": {
			input:    "A",
			expected: "A",
		},
		"single character digit": {
			input:    "1",
			expected: "_",
		},
		"single underscore": {
			input:    "_",
			expected: "_",
		},
		"simple valid name": {
			input:    "FOO",
			expected: "FOO",
		},
		"lowercase converted to uppercase": {
			input:    "foo",
			expected: "FOO",
		},
		"mixed case converted to uppercase": {
			input:    "FoO",
			expected: "FOO",
		},
		"name with digits": {
			input:    "FOO123",
			expected: "FOO123",
		},
		"name starting with digit": {
			input:    "123FOO",
			expected: "_23FOO",
		},
		"name with underscores": {
			input:    "FOO_BAR",
			expected: "FOO_BAR",
		},
		"name with spaces": {
			input:    "FOO BAR",
			expected: "FOO_BAR",
		},
		"name with hyphens": {
			input:    "FOO-BAR",
			expected: "FOO_BAR",
		},
		"name with dots": {
			input:    "FOO.BAR",
			expected: "FOO_BAR",
		},
		"name with equals": {
			input:    "FOO=BAR",
			expected: "FOO_BAR",
		},
		"name with multiple consecutive special chars": {
			input:    "FOO---BAR",
			expected: "FOO_BAR",
		},
		"name with leading underscores": {
			input:    "___FOO",
			expected: "FOO",
		},
		"name with trailing underscores": {
			input:    "FOO___",
			expected: "FOO",
		},
		"name with leading and trailing underscores": {
			input:    "___FOO___",
			expected: "FOO",
		},
		"name with only underscores and special chars": {
			input:    "___---___",
			expected: "_",
		},
		"unicode characters": {
			input:    "FOO_ÄÖÜ_BAR",
			expected: "FOO_BAR",
		},
		"complex mixed input": {
			input:    "123foo-bar.baz=qux___",
			expected: "_23FOO_BAR_BAZ_QUX",
		},
		"all special characters": {
			input:    "!@#$%^&*()",
			expected: "_",
		},
		"name with numbers in middle": {
			input:    "FOO123BAR456",
			expected: "FOO123BAR456",
		},
		"single special character": {
			input:    "-",
			expected: "_",
		},
		"whitespace only": {
			input:    "   ",
			expected: "_",
		},
		"digits only": {
			input:    "123",
			expected: "_23",
		},
	} {
		t.Run(name, func(t *testing.T) {
			result := SanitizeName(tc.input)
			test.Assert(check.Compare(t, result, tc.expected))
		})
	}
}

func Fuzz_SanitizeName(f *testing.F) {
	f.Add("")                                                     // empty
	f.Add("_")                                                    // bare underscore
	f.Add("a")                                                    // single letter
	f.Add("A")                                                    // single uppercase
	f.Add("0")                                                    // single digit (invalid first)
	f.Add("__")                                                   // only underscores
	f.Add("___a___")                                              // leading/trailing underscores
	f.Add("FOO")                                                  // already valid
	f.Add("FOO_BAR")                                              // already valid
	f.Add("foo_bar")                                              // needs uppercasing
	f.Add("1abc")                                                 // starts with digit
	f.Add("-abc")                                                 // invalid first rune
	f.Add(".abc")                                                 // invalid first rune
	f.Add("=abc")                                                 // invalid first rune
	f.Add("foo-bar")                                              // hyphen
	f.Add("foo.bar")                                              // dot
	f.Add("foo=bar")                                              // equals
	f.Add("foo bar")                                              // space
	f.Add("foo\tbar")                                             // tab
	f.Add("foo\nbar")                                             // newline
	f.Add("foo/bar")                                              // slash
	f.Add("foo\\bar")                                             // backslash
	f.Add(`a*b?c[d]e{f}g$h(i)j'k"l\m|n;o&p<q>r~s#t!:,u@v+ w%x^y`) // --- Shell metachar soup ---
	f.Add("FOO_ÄÖÜ_BAR")                                          // Latin-1 letters
	f.Add("привет")                                               // Cyrillic
	f.Add("γειάσου")                                              // Greek
	f.Add("مرحبا")                                                // Arabic (RTL)
	f.Add("こんにちは")                                                // Japanese
	f.Add("ℌ𝔢𝔩𝔩𝔬")                                                // math Fraktur letters (non-ASCII)
	f.Add("e\u0301")                                              // "e" + combining acute
	f.Add("\u200Djoin")                                           // ZWJ
	f.Add("\u200Bhidden")                                         // ZWSP
	f.Add("foo\u00A0bar")                                         // NBSP
	f.Add("\u202Ertl")                                            // RLO (directionality)
	f.Add("foo😀bar")                                              // emoji in middle
	f.Add("✅ok")                                                  // emoji first
	f.Add("ＦＯＯ＿ＢＡＲ")                                              // fullwidth ASCII lookalikes
	f.Add(strings.Repeat("_", 256))                               // long underscores
	f.Add(strings.Repeat("A", 1024))                              // long valid
	f.Add(strings.Repeat("-", 1024))                              // long invalid
	f.Add(strings.Repeat("9", 512) + "A")                         // long digit prefix then letter
	f.Add("a" + strings.Repeat("_", 512) + "b")                   // collapse many underscores
	f.Add(strings.Repeat("é", 512))                               // long non-ASCII
	f.Add("__9a")                                                 // trim exposes leading digit
	f.Add("__=foo__")                                             // trims around invalids
	f.Add("9__")                                                  // digit then underscores
	f.Add("_9")                                                   // underscore then digit
	f.Add(string([]byte{0xff}))                                   // single invalid byte
	f.Add("A" + string([]byte{0xc3, 0x28}) + "Z")                 // malformed UTF-8 sequence

	f.Fuzz(func(t *testing.T, input string) {
		result := SanitizeName(input)
		test.Assert(t, result != "")

		first := rune(result[0])
		test.Assert(t, (first >= 'A' && first <= 'Z') || first == '_', "first character must be letter or underscore: %q -> %q", input, result)

		for _, r := range result[1:] {
			test.Assert(t, (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_', "subsequent character must be letter, digit, or underscore: %q -> %q", input, result)
		}

		for i := range len(result) - 1 {
			test.Assert(t, result[i] != '_' || result[i+1] != '_', "regardless of wrong characters, there should never be two consecutive underscore: %q -> %q", input, result)
		}
	})
}
