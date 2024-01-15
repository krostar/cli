package cli

import (
	"testing"

	"gotest.tools/v3/assert"
)

func Test_NewBuiltinFlag(t *testing.T) {
	var value int

	flag := NewBuiltinFlag[int]("longName", "s", &value, "description")
	assert.Check(t, flag.LongName() == "longName")
	assert.Check(t, flag.ShortName() == "s")
	assert.Check(t, flag.Description() == "description")
	assert.Check(t, flag.TypeRepr() == "int")

	assert.ErrorContains(t, flag.FromString("abc"), "invalid syntax")

	assert.NilError(t, flag.FromString("  42 "))
	repr := flag.String()
	assert.Check(t, repr == "42")
}

func Test_NewBuiltinSliceFlag(t *testing.T) {
	var value []int

	flag := NewBuiltinSliceFlag[int]("longName", "s", &value, "description")
	assert.Check(t, flag.LongName() == "longName")
	assert.Check(t, flag.ShortName() == "s")
	assert.Check(t, flag.Description() == "description")
	assert.Check(t, flag.TypeRepr() == "[]int")

	assert.ErrorContains(t, flag.FromString(" 16 ,abc,  18"), "invalid syntax")

	assert.NilError(t, flag.FromString(" 42 ,  44"))
	repr := flag.String()
	assert.Equal(t, "[42,44]", repr)
}

func assertEqualIfBuiltinParsingSucceed[T builtins](t *testing.T, providedRawValue string, expectedValue ...T) error {
	t.Helper()

	value, err := builtinFromString[T](providedRawValue)
	if err != nil {
		return err
	}

	assert.Assert(t, len(expectedValue) == 1)

	assert.Check(t, expectedValue[0] == value)
	return nil
}

func Test_builtinFromString(t *testing.T) {
	t.Run("bool", func(t *testing.T) {
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[bool](t, "true", true))
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[bool](t, "false", false))
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[bool](t, "1", true))
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[bool](t, "0", false))

		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[bool](t, ""), "invalid syntax")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[bool](t, "abc"), "invalid syntax")
	})

	t.Run("string", func(t *testing.T) {
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[string](t, "foo", "foo"))
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[string](t, "42", "42"))
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[string](t, "", ""))
	})

	t.Run("int", func(t *testing.T) {
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[int](t, "-9223372036854775808", -9223372036854775808))
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[int](t, "0", 0))
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[int](t, "9223372036854775807", 9223372036854775807))

		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[int](t, ""), "invalid syntax")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[int](t, "-9223372036854775809"), "value out of range")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[int](t, "9223372036854775808"), "value out of range")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[int](t, "abc"), "invalid syntax")
	})

	t.Run("int8", func(t *testing.T) {
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[int8](t, "-128", int8(-128)))
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[int8](t, "0", int8(0)))
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[int8](t, "127", int8(127)))

		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[int8](t, ""), "invalid syntax")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[int8](t, "-129"), "value out of range")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[int8](t, "128"), "value out of range")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[int8](t, "abc"), "invalid syntax")
	})

	t.Run("int16", func(t *testing.T) {
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[int16](t, "-32768", int16(-32768)))
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[int16](t, "0", int16(0)))
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[int16](t, "32767", int16(32767)))

		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[int16](t, ""), "invalid syntax")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[int16](t, "-32769"), "value out of range")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[int16](t, "32768"), "value out of range")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[int16](t, "abc"), "invalid syntax")
	})

	t.Run("int32", func(t *testing.T) {
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[int32](t, "-2147483648", int32(-2147483648)))
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[int32](t, "0", int32(0)))
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[int32](t, "2147483647", int32(2147483647)))

		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[int32](t, ""), "invalid syntax")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[int32](t, "-2147483649"), "value out of range")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[int32](t, "2147483648"), "value out of range")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[int32](t, "abc"), "invalid syntax")
	})

	t.Run("int64", func(t *testing.T) {
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[int64](t, "-9223372036854775808", int64(-9223372036854775808)))
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[int64](t, "0", int64(0)))
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[int64](t, "9223372036854775807", int64(9223372036854775807)))

		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[int64](t, ""), "invalid syntax")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[int64](t, "-9223372036854775809"), "value out of range")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[int64](t, "9223372036854775808"), "value out of range")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[int64](t, "abc"), "invalid syntax")
	})

	t.Run("uint", func(t *testing.T) {
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[uint](t, "0", uint(0)))
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[uint](t, "18446744073709551615", uint(18446744073709551615)))

		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[uint](t, ""), "invalid syntax")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[uint](t, "-1"), "invalid syntax")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[uint](t, "18446744073709551616"), "value out of range")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[uint](t, "abc"), "invalid syntax")
	})

	t.Run("uint8", func(t *testing.T) {
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[uint8](t, "0", uint8(0)))
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[uint8](t, "255", uint8(255)))

		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[uint8](t, ""), "invalid syntax")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[uint8](t, "-1"), "invalid syntax")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[uint8](t, "256"), "value out of range")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[uint8](t, "abc"), "invalid syntax")
	})

	t.Run("uint16", func(t *testing.T) {
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[uint16](t, "0", uint16(0)))
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[uint16](t, "65535", uint16(65535)))

		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[uint16](t, ""), "invalid syntax")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[uint16](t, "-1"), "invalid syntax")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[uint16](t, "65536"), "value out of range")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[uint16](t, "abc"), "invalid syntax")
	})

	t.Run("uint32", func(t *testing.T) {
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[uint32](t, "0", uint32(0)))
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[uint32](t, "4294967295", uint32(4294967295)))

		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[uint32](t, ""), "invalid syntax")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[uint32](t, "-1"), "invalid syntax")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[uint32](t, "4294967296"), "value out of range")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[uint32](t, "abc"), "invalid syntax")
	})

	t.Run("uint64", func(t *testing.T) {
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[uint64](t, "0", uint64(0)))
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[uint64](t, "18446744073709551615", uint64(18446744073709551615)))

		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[uint64](t, ""), "invalid syntax")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[uint64](t, "-1"), "invalid syntax")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[uint64](t, "18446744073709551616"), "value out of range")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[uint64](t, "abc"), "invalid syntax")
	})

	t.Run("float32", func(t *testing.T) {
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[float32](t, "-340282346638528859811704183484516925440", float32(-340282346638528859811704183484516925440)))
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[float32](t, "0", float32(0)))
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[float32](t, "0.356789", float32(0.356789)))
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[float32](t, "340282346638528859811704183484516925440", float32(340282346638528859811704183484516925440)))

		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[float32](t, ""), "invalid syntax")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[float32](t, "179769313486231570814527423731704356798070567525844996598917476803157260780028538760589558632766878171540458953514382464234321326889464182768467546703537516986049910576551282076245490090389328944075868508455133942304583236903222948165808559332123348274797826204144723168738177180919299881250404026184124858368"), "value out of range")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[float32](t, "-179769313486231570814527423731704356798070567525844996598917476803157260780028538760589558632766878171540458953514382464234321326889464182768467546703537516986049910576551282076245490090389328944075868508455133942304583236903222948165808559332123348274797826204144723168738177180919299881250404026184124858368"), "value out of range")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[float32](t, "abc"), "invalid syntax")
	})

	t.Run("float64", func(t *testing.T) {
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[float64](t, "-1797693134862315708145274237317043567980705675258449965989174768031572607800285387605895586327668781715404589535143824642343213268894641827684675467035", float64(-1797693134862315708145274237317043567980705675258449965989174768031572607800285387605895586327668781715404589535143824642343213268894641827684675467035)))
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[float64](t, "0", float64(0)))
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[float64](t, "0.356789", float64(0.356789)))
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[float64](t, "9797693134862915708145274237317043567980705675258449965989174768031572607800285387605895586327668781715404989535143824642343213268894641827684675467035399", float64(9797693134862915708145274237317043567980705675258449965989174768031572607800285387605895586327668781715404989535143824642343213268894641827684675467035399)))

		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[float64](t, ""), "invalid syntax")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[float64](t, "-17976931348623157081452742373170435679807056752584499659891747680315726078002853876058955863276687817154045895351438246423432132688946418276846754670353751698604991057655128207624549009038932894407586850845513394230458323690322294816580855933212334827479782620414472316873817718091929988125040402618412485836899"), "value out of range")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[float64](t, "17976931348623157081452742373170435679807056752584499659891747680315726078002853876058955863276687817154045895351438246423432132688946418276846754670353751698604991057655128207624549009038932894407586850845513394230458323690322294816580855933212334827479782620414472316873817718091929988125040402618412485836899"), "value out of range")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[float64](t, "abc"), "invalid syntax")
	})

	t.Run("complex64", func(t *testing.T) {
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[complex64](t, "-340282346638528859811704183484516925440-340282346638528859811704183484516925440i", complex64(-340282346638528859811704183484516925440-340282346638528859811704183484516925440i)))
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[complex64](t, "0", complex64(0)))
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[complex64](t, "-0.356789i", complex64(-0.356789i)))
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[complex64](t, "340282346638528859811704183484516925440+340282346638528859811704183484516925440i", complex64(340282346638528859811704183484516925440+340282346638528859811704183484516925440i)))

		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[complex64](t, ""), "invalid syntax")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[complex64](t, "i"), "invalid syntax")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[complex64](t, "179769313486231570814527423731704356798070567525844996598917476803157260780028538760589558632766878171540458953514382464234321326889464182768467546703537516986049910576551282076245490090389328944075868508455133942304583236903222948165808559332123348274797826204144723168738177180919299881250404026184124858368"), "value out of range")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[complex64](t, "-179769313486231570814527423731704356798070567525844996598917476803157260780028538760589558632766878171540458953514382464234321326889464182768467546703537516986049910576551282076245490090389328944075868508455133942304583236903222948165808559332123348274797826204144723168738177180919299881250404026184124858368"), "value out of range")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[complex64](t, "abc"), "invalid syntax")
	})

	t.Run("complex128", func(t *testing.T) {
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[complex128](t, "-1797693134862315708145274237317043567980705675258449965989174768031572607800285387605895586327668781715404589535143824642343213268894641827684675467035-1797693134862315708145274237317043567980705675258449965989174768031572607800285387605895586327668781715404589535143824642343213268894641827684675467035i", -1797693134862315708145274237317043567980705675258449965989174768031572607800285387605895586327668781715404589535143824642343213268894641827684675467035-1797693134862315708145274237317043567980705675258449965989174768031572607800285387605895586327668781715404589535143824642343213268894641827684675467035i))
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[complex128](t, "0", 0))
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[complex128](t, "-0.356789i", -0.356789i))
		assert.NilError(t, assertEqualIfBuiltinParsingSucceed[complex128](t, "1797693134862315708145274237317043567980705675258449965989174768031572607800285387605895586327668781715404589535143824642343213268894641827684675467035+1797693134862315708145274237317043567980705675258449965989174768031572607800285387605895586327668781715404589535143824642343213268894641827684675467035i", 1797693134862315708145274237317043567980705675258449965989174768031572607800285387605895586327668781715404589535143824642343213268894641827684675467035+1797693134862315708145274237317043567980705675258449965989174768031572607800285387605895586327668781715404589535143824642343213268894641827684675467035i))

		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[complex128](t, ""), "invalid syntax")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[complex128](t, "i"), "invalid syntax")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[complex128](t, "17976931348623157081452742373170435679807056752584499659891747680315726078002853876058955863276687817154045895351438246423432132688946418276846754670353751698604991057655128207624549009038932894407586850845513394230458323690322294816580855933212334827479782620414472316873817718091929988125040402618412485836899"), "value out of range")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[complex128](t, "-17976931348623157081452742373170435679807056752584499659891747680315726078002853876058955863276687817154045895351438246423432132688946418276846754670353751698604991057655128207624549009038932894407586850845513394230458323690322294816580855933212334827479782620414472316873817718091929988125040402618412485836899"), "value out of range")
		assert.ErrorContains(t, assertEqualIfBuiltinParsingSucceed[complex128](t, "abc"), "invalid syntax")
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
