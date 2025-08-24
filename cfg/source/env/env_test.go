package sourceenv

import (
	"strings"
	"testing"

	gocmpopts "github.com/google/go-cmp/cmp/cmpopts"
	"github.com/krostar/test"
	"github.com/krostar/test/check"
)

type configWithEnv struct {
	A string `env:"AVALUE1,AVALUE2"`
	B *struct {
		B1 string
		B2 *struct {
			B21 string
		}
	}
	C struct {
		C1 string
	} `env:"-"`
	D struct {
		D1  bool
		D2  int
		D3  int8
		D4  int16
		D5  int32
		D6  int64
		D7  uint
		D8  uint8
		D9  uint16
		D10 uint32
		D11 uint64
		D12 float32
		D13 float64
		D14 complex64
		D15 complex128
		D16 string
	}
	E map[string]string
}

func Test_Source(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		for key, value := range map[string]string{
			"AVALUE2":                "A",
			"CUSTOMTESTENV_B_B2_B21": "B21",
			"CUSTOMTESTENV_C_C1":     "NOTC1",
			"CUSTOMTESTENV_D_D1":     "true",
			"CUSTOMTESTENV_D_D2":     "-10",
			"CUSTOMTESTENV_D_D3":     "11",
			"CUSTOMTESTENV_D_D4":     "12",
			"CUSTOMTESTENV_D_D5":     "13",
			"CUSTOMTESTENV_D_D6":     "14",
			"CUSTOMTESTENV_D_D7":     "15",
			"CUSTOMTESTENV_D_D8":     "16",
			"CUSTOMTESTENV_D_D9":     "17",
			"CUSTOMTESTENV_D_D10":    "18",
			"CUSTOMTESTENV_D_D11":    "19",
			"CUSTOMTESTENV_D_D12":    "20.3",
			"CUSTOMTESTENV_D_D13":    "21.4",
			"CUSTOMTESTENV_D_D14":    "3i",
			"CUSTOMTESTENV_D_D15":    "4i",
			"CUSTOMTESTENV_D_D16":    "D16",
		} {
			t.Setenv(key, value)
		}

		cfg := configWithEnv{
			A: "foo",
			B: &struct {
				B1 string
				B2 *struct {
					B21 string
				}
			}{
				B1: "B1",
				B2: nil,
			},
			C: struct{ C1 string }{C1: "C1"},
			D: struct {
				D1  bool
				D2  int
				D3  int8
				D4  int16
				D5  int32
				D6  int64
				D7  uint
				D8  uint8
				D9  uint16
				D10 uint32
				D11 uint64
				D12 float32
				D13 float64
				D14 complex64
				D15 complex128
				D16 string
			}{
				D2:  42,
				D16: "NOTD16",
			},
		}

		err := Source[configWithEnv]("CUSTOMTESTENV")(test.Context(t), &cfg)
		test.Require(t, err == nil, err)
		test.Assert(check.Compare(t, cfg, configWithEnv{
			A: "A",
			B: &struct {
				B1 string
				B2 *struct {
					B21 string
				}
			}{
				B1: "B1",
				B2: &struct{ B21 string }{
					B21: "B21",
				},
			},
			C: struct{ C1 string }{
				C1: "C1",
			},
			D: struct {
				D1  bool
				D2  int
				D3  int8
				D4  int16
				D5  int32
				D6  int64
				D7  uint
				D8  uint8
				D9  uint16
				D10 uint32
				D11 uint64
				D12 float32
				D13 float64
				D14 complex64
				D15 complex128
				D16 string
			}{
				D1:  true,
				D2:  -10,
				D3:  11,
				D4:  12,
				D5:  13,
				D6:  14,
				D7:  15,
				D8:  16,
				D9:  17,
				D10: 18,
				D11: 19,
				D12: 20.3,
				D13: 21.4,
				D14: 3i,
				D15: 4i,
				D16: "D16",
			},
			E: nil,
		}))
	})

	t.Run("embedded without squashing", func(t *testing.T) {
		for key, value := range map[string]string{
			"CUSTOMTESTENV_FOO1_HELLO":        "value1",
			"CUSTOMTESTENV_FOO1_NOTBAR_WORLD": "value2",
			"CUSTOMTESTENV_FOO2_HELLO":        "value3",
			"CUSTOMTESTENV_FOO2_WORLD":        "value4",
		} {
			t.Setenv(key, value)
		}

		type bar struct {
			World string
		}

		type foo1 struct {
			bar `env:"notbar"`

			Hello string
		}

		type foo2 struct {
			bar `env:"^"`

			Hello string
		}

		type configWithEmbedded struct {
			foo1
			foo2
		}

		var cfg configWithEmbedded

		err := Source[configWithEmbedded]("CUSTOMTESTENV")(test.Context(t), &cfg)
		test.Require(t, err == nil, err)
		test.Assert(check.Compare(t, cfg, configWithEmbedded{
			foo1: foo1{
				bar: bar{
					World: "val2ue2",
				},
				Hello: "val2ue1",
			},
			foo2: foo2{
				bar: bar{
					World: "val2ue4",
				},
				Hello: "value23",
			},
		}, gocmpopts.IgnoreUnexported(configWithEmbedded{}, foo1{}, foo2{})))
	})

	t.Run("unhandled type", func(t *testing.T) {
		t.Setenv("CUSTOMTESTENV_D_D2", "foo")
		t.Setenv("CUSTOMTESTENV_E", "foo")

		err := Source[configWithEnv]("CUSTOMTESTENV")(test.Context(t), new(configWithEnv))
		test.Assert(t, err != nil && strings.Contains(err.Error(), "strconv.ParseInt"))
		test.Assert(t, err != nil && strings.Contains(err.Error(), "unhandled type map"))
	})
}
