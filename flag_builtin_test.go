package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewBuiltinFlag(t *testing.T) {
	var value int

	flag := NewBuiltinFlag[int]("longName", "s", &value, "description")
	assert.Equal(t, "longName", flag.LongName())
	assert.Equal(t, "s", flag.ShortName())
	assert.Equal(t, "description", flag.Description())
	assert.Equal(t, "int", flag.TypeRepr())

	require.Error(t, flag.FromString("abc"))

	require.NoError(t, flag.FromString("  42 "))
	repr := flag.String()
	assert.Equal(t, "42", repr)
}

func Test_NewBuiltinSliceFlag(t *testing.T) {
	var value []int

	flag := NewBuiltinSliceFlag[int]("longName", "s", &value, "description")
	assert.Equal(t, "longName", flag.LongName())
	assert.Equal(t, "s", flag.ShortName())
	assert.Equal(t, "description", flag.Description())
	assert.Equal(t, "[]int", flag.TypeRepr())

	require.Error(t, flag.FromString(" 16 ,abc,  18"))

	require.NoError(t, flag.FromString(" 42 ,  44"))
	repr := flag.String()
	assert.Equal(t, "[42,44]", repr)
}

func assertEqualIfBuiltinParsingSucceed[T builtins](t *testing.T, providedRawValue string, expectedValue ...T) error {
	t.Helper()

	value, err := builtinFromString[T](providedRawValue)
	if err != nil {
		return err
	}

	if len(expectedValue) != 1 {
		t.Fatal("expected to have 1 value in case of success")
	}

	assert.Equal(t, expectedValue[0], value)
	return nil
}

func Test_builtinFromString(t *testing.T) {
	t.Run("bool", func(t *testing.T) {
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[bool](t, "true", true))
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[bool](t, "false", false))
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[bool](t, "1", true))
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[bool](t, "0", false))

		require.Error(t, assertEqualIfBuiltinParsingSucceed[bool](t, ""))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[bool](t, "abc"))
	})

	t.Run("string", func(t *testing.T) {
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[string](t, "foo", "foo"))
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[string](t, "42", "42"))
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[string](t, "", ""))
	})

	t.Run("int", func(t *testing.T) {
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[int](t, "-9223372036854775808", -9223372036854775808))
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[int](t, "0", 0))
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[int](t, "9223372036854775807", 9223372036854775807))

		require.Error(t, assertEqualIfBuiltinParsingSucceed[int](t, ""))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[int](t, "-9223372036854775809"))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[int](t, "9223372036854775808"))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[int](t, "abc"))
	})

	t.Run("int8", func(t *testing.T) {
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[int8](t, "-128", int8(-128)))
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[int8](t, "0", int8(0)))
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[int8](t, "127", int8(127)))

		require.Error(t, assertEqualIfBuiltinParsingSucceed[int8](t, ""))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[int8](t, "-129"))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[int8](t, "128"))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[int8](t, "abc"))
	})

	t.Run("int16", func(t *testing.T) {
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[int16](t, "-32768", int16(-32768)))
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[int16](t, "0", int16(0)))
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[int16](t, "32767", int16(32767)))

		require.Error(t, assertEqualIfBuiltinParsingSucceed[int16](t, ""))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[int16](t, "-32769"))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[int16](t, "32768"))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[int16](t, "abc"))
	})

	t.Run("int32", func(t *testing.T) {
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[int32](t, "-2147483648", int32(-2147483648)))
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[int32](t, "0", int32(0)))
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[int32](t, "2147483647", int32(2147483647)))

		require.Error(t, assertEqualIfBuiltinParsingSucceed[int32](t, ""))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[int32](t, "-2147483649"))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[int32](t, "2147483648"))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[int32](t, "abc"))
	})

	t.Run("int64", func(t *testing.T) {
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[int64](t, "-9223372036854775808", int64(-9223372036854775808)))
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[int64](t, "0", int64(0)))
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[int64](t, "9223372036854775807", int64(9223372036854775807)))

		require.Error(t, assertEqualIfBuiltinParsingSucceed[int64](t, ""))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[int64](t, "-9223372036854775809"))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[int64](t, "9223372036854775808"))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[int64](t, "abc"))
	})

	t.Run("uint", func(t *testing.T) {
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[uint](t, "0", uint(0)))
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[uint](t, "18446744073709551615", uint(18446744073709551615)))

		require.Error(t, assertEqualIfBuiltinParsingSucceed[uint](t, ""))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[uint](t, "-1"))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[uint](t, "18446744073709551616"))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[uint](t, "abc"))
	})

	t.Run("uint8", func(t *testing.T) {
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[uint8](t, "0", uint8(0)))
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[uint8](t, "255", uint8(255)))

		require.Error(t, assertEqualIfBuiltinParsingSucceed[uint8](t, ""))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[uint8](t, "-1"))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[uint8](t, "256"))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[uint8](t, "abc"))
	})

	t.Run("uint16", func(t *testing.T) {
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[uint16](t, "0", uint16(0)))
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[uint16](t, "65535", uint16(65535)))

		require.Error(t, assertEqualIfBuiltinParsingSucceed[uint16](t, ""))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[uint16](t, "-1"))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[uint16](t, "65536"))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[uint16](t, "abc"))
	})

	t.Run("uint32", func(t *testing.T) {
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[uint32](t, "0", uint32(0)))
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[uint32](t, "4294967295", uint32(4294967295)))

		require.Error(t, assertEqualIfBuiltinParsingSucceed[uint32](t, ""))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[uint32](t, "-1"))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[uint32](t, "4294967296"))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[uint32](t, "abc"))
	})

	t.Run("uint64", func(t *testing.T) {
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[uint64](t, "0", uint64(0)))
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[uint64](t, "18446744073709551615", uint64(18446744073709551615)))

		require.Error(t, assertEqualIfBuiltinParsingSucceed[uint64](t, ""))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[uint64](t, "-1"))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[uint64](t, "18446744073709551616"))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[uint64](t, "abc"))
	})

	t.Run("float32", func(t *testing.T) {
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[float32](t, "-340282346638528859811704183484516925440", float32(-340282346638528859811704183484516925440)))
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[float32](t, "0", float32(0)))
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[float32](t, "0.356789", float32(0.356789)))
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[float32](t, "340282346638528859811704183484516925440", float32(340282346638528859811704183484516925440)))

		require.Error(t, assertEqualIfBuiltinParsingSucceed[float32](t, ""))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[float32](t, "179769313486231570814527423731704356798070567525844996598917476803157260780028538760589558632766878171540458953514382464234321326889464182768467546703537516986049910576551282076245490090389328944075868508455133942304583236903222948165808559332123348274797826204144723168738177180919299881250404026184124858368"))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[float32](t, "-179769313486231570814527423731704356798070567525844996598917476803157260780028538760589558632766878171540458953514382464234321326889464182768467546703537516986049910576551282076245490090389328944075868508455133942304583236903222948165808559332123348274797826204144723168738177180919299881250404026184124858368"))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[float32](t, "abc"))
	})

	t.Run("float64", func(t *testing.T) {
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[float64](t, "-1797693134862315708145274237317043567980705675258449965989174768031572607800285387605895586327668781715404589535143824642343213268894641827684675467035", float64(-1797693134862315708145274237317043567980705675258449965989174768031572607800285387605895586327668781715404589535143824642343213268894641827684675467035)))
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[float64](t, "0", float64(0)))
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[float64](t, "0.356789", float64(0.356789)))
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[float64](t, "9797693134862915708145274237317043567980705675258449965989174768031572607800285387605895586327668781715404989535143824642343213268894641827684675467035399", float64(9797693134862915708145274237317043567980705675258449965989174768031572607800285387605895586327668781715404989535143824642343213268894641827684675467035399)))

		require.Error(t, assertEqualIfBuiltinParsingSucceed[float64](t, ""))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[float64](t, "-17976931348623157081452742373170435679807056752584499659891747680315726078002853876058955863276687817154045895351438246423432132688946418276846754670353751698604991057655128207624549009038932894407586850845513394230458323690322294816580855933212334827479782620414472316873817718091929988125040402618412485836899"))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[float64](t, "17976931348623157081452742373170435679807056752584499659891747680315726078002853876058955863276687817154045895351438246423432132688946418276846754670353751698604991057655128207624549009038932894407586850845513394230458323690322294816580855933212334827479782620414472316873817718091929988125040402618412485836899"))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[float64](t, "abc"))
	})

	t.Run("complex64", func(t *testing.T) {
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[complex64](t, "-340282346638528859811704183484516925440-340282346638528859811704183484516925440i", complex64(-340282346638528859811704183484516925440-340282346638528859811704183484516925440i)))
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[complex64](t, "0", complex64(0)))
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[complex64](t, "-0.356789i", complex64(-0.356789i)))
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[complex64](t, "340282346638528859811704183484516925440+340282346638528859811704183484516925440i", complex64(340282346638528859811704183484516925440+340282346638528859811704183484516925440i)))

		require.Error(t, assertEqualIfBuiltinParsingSucceed[complex64](t, ""))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[complex64](t, "i"))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[complex64](t, "179769313486231570814527423731704356798070567525844996598917476803157260780028538760589558632766878171540458953514382464234321326889464182768467546703537516986049910576551282076245490090389328944075868508455133942304583236903222948165808559332123348274797826204144723168738177180919299881250404026184124858368"))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[complex64](t, "-179769313486231570814527423731704356798070567525844996598917476803157260780028538760589558632766878171540458953514382464234321326889464182768467546703537516986049910576551282076245490090389328944075868508455133942304583236903222948165808559332123348274797826204144723168738177180919299881250404026184124858368"))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[complex64](t, "abc"))
	})

	t.Run("complex128", func(t *testing.T) {
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[complex128](t, "-1797693134862315708145274237317043567980705675258449965989174768031572607800285387605895586327668781715404589535143824642343213268894641827684675467035-1797693134862315708145274237317043567980705675258449965989174768031572607800285387605895586327668781715404589535143824642343213268894641827684675467035i", -1797693134862315708145274237317043567980705675258449965989174768031572607800285387605895586327668781715404589535143824642343213268894641827684675467035-1797693134862315708145274237317043567980705675258449965989174768031572607800285387605895586327668781715404589535143824642343213268894641827684675467035i))
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[complex128](t, "0", 0))
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[complex128](t, "-0.356789i", -0.356789i))
		require.NoError(t, assertEqualIfBuiltinParsingSucceed[complex128](t, "1797693134862315708145274237317043567980705675258449965989174768031572607800285387605895586327668781715404589535143824642343213268894641827684675467035+1797693134862315708145274237317043567980705675258449965989174768031572607800285387605895586327668781715404589535143824642343213268894641827684675467035i", 1797693134862315708145274237317043567980705675258449965989174768031572607800285387605895586327668781715404589535143824642343213268894641827684675467035+1797693134862315708145274237317043567980705675258449965989174768031572607800285387605895586327668781715404589535143824642343213268894641827684675467035i))

		require.Error(t, assertEqualIfBuiltinParsingSucceed[complex128](t, ""))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[complex128](t, "i"))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[complex128](t, "17976931348623157081452742373170435679807056752584499659891747680315726078002853876058955863276687817154045895351438246423432132688946418276846754670353751698604991057655128207624549009038932894407586850845513394230458323690322294816580855933212334827479782620414472316873817718091929988125040402618412485836899"))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[complex128](t, "-17976931348623157081452742373170435679807056752584499659891747680315726078002853876058955863276687817154045895351438246423432132688946418276846754670353751698604991057655128207624549009038932894407586850845513394230458323690322294816580855933212334827479782620414472316873817718091929988125040402618412485836899"))
		require.Error(t, assertEqualIfBuiltinParsingSucceed[complex128](t, "abc"))
	})
}

func assertNoErrorBuiltinParsingSucceed[T builtins](t *testing.T, providedValue T, expectedValueRepr string) {
	t.Helper()

	repr := builtinToString[T](providedValue)
	assert.Equal(t, expectedValueRepr, repr)
}

func Test_builtinToString(t *testing.T) {
	t.Run("bool", func(t *testing.T) {
		assertNoErrorBuiltinParsingSucceed[bool](t, true, "true")
		assertNoErrorBuiltinParsingSucceed[bool](t, false, "false")
	})

	t.Run("string", func(t *testing.T) {
		assertNoErrorBuiltinParsingSucceed[string](t, "foo", "foo")
	})

	t.Run("int", func(t *testing.T) {
		assertNoErrorBuiltinParsingSucceed[int](t, 42, "42")
	})

	t.Run("int8", func(t *testing.T) {
		assertNoErrorBuiltinParsingSucceed[int8](t, int8(42), "42")
	})

	t.Run("int16", func(t *testing.T) {
		assertNoErrorBuiltinParsingSucceed[int16](t, int16(42), "42")
	})

	t.Run("int32", func(t *testing.T) {
		assertNoErrorBuiltinParsingSucceed[int32](t, int32(42), "42")
	})

	t.Run("int64", func(t *testing.T) {
		assertNoErrorBuiltinParsingSucceed[int64](t, int64(42), "42")
	})

	t.Run("uint", func(t *testing.T) {
		assertNoErrorBuiltinParsingSucceed[uint](t, uint(42), "42")
	})

	t.Run("uint8", func(t *testing.T) {
		assertNoErrorBuiltinParsingSucceed[uint8](t, uint8(42), "42")
	})

	t.Run("uint16", func(t *testing.T) {
		assertNoErrorBuiltinParsingSucceed[uint16](t, uint16(42), "42")
	})

	t.Run("uint32", func(t *testing.T) {
		assertNoErrorBuiltinParsingSucceed[uint32](t, uint32(42), "42")
	})

	t.Run("uint64", func(t *testing.T) {
		assertNoErrorBuiltinParsingSucceed[uint64](t, uint64(42), "42")
	})

	t.Run("float32", func(t *testing.T) {
		assertNoErrorBuiltinParsingSucceed[float32](t, float32(42), "42.0000")
		assertNoErrorBuiltinParsingSucceed[float32](t, float32(42.21), "42.2100")
		assertNoErrorBuiltinParsingSucceed[float32](t, float32(42.123456), "42.1235")
	})

	t.Run("float64", func(t *testing.T) {
		assertNoErrorBuiltinParsingSucceed[float64](t, 42, "42.0000")
		assertNoErrorBuiltinParsingSucceed[float64](t, 42.21, "42.2100")
		assertNoErrorBuiltinParsingSucceed[float64](t, 42.123456, "42.1235")
	})

	t.Run("complex64", func(t *testing.T) {
		assertNoErrorBuiltinParsingSucceed[complex64](t, complex64(42+21i), "(42.0000+21.0000i)")
	})

	t.Run("complex128", func(t *testing.T) {
		assertNoErrorBuiltinParsingSucceed[complex128](t, 42+21i, "(42.0000+21.0000i)")
	})
}
