package cli

import (
	"strings"
	"testing"

	"github.com/krostar/test"
)

func Test_NewBuiltinFlag(t *testing.T) {
	var value int

	flag := NewBuiltinFlag[int]("longName", "s", &value, "description")
	test.Assert(t, flag.LongName() == "longName")
	test.Assert(t, flag.ShortName() == "s")
	test.Assert(t, flag.Description() == "description")
	test.Assert(t, flag.TypeRepr() == "int")

	err := flag.FromString("abc")
	test.Assert(t, err != nil && strings.Contains(err.Error(), "invalid syntax"))

	test.Assert(t, flag.FromString("  42 ") == nil)
	test.Assert(t, value == 42)
	test.Assert(t, flag.String() == "42")
}

func Test_NewBuiltinPointerFlag(t *testing.T) {
	var value *int

	flag := NewBuiltinPointerFlag[int]("longName", "s", &value, "description")
	test.Assert(t, flag.LongName() == "longName")
	test.Assert(t, flag.ShortName() == "s")
	test.Assert(t, flag.Description() == "description")
	test.Assert(t, flag.TypeRepr() == "*int")

	err := flag.FromString("abc")
	test.Assert(t, err != nil && strings.Contains(err.Error(), "invalid syntax"))
	test.Assert(t, flag.String() == "<nil>")

	test.Assert(t, flag.FromString("  42 ") == nil)
	test.Assert(t, value != nil && *value == 42)
	test.Assert(t, flag.String() == "42")
}

func Test_NewBuiltinSliceFlag(t *testing.T) {
	var value []int

	flag := NewBuiltinSliceFlag[int]("longName", "s", &value, "description")
	test.Assert(t, flag.LongName() == "longName")
	test.Assert(t, flag.ShortName() == "s")
	test.Assert(t, flag.Description() == "description")
	test.Assert(t, flag.TypeRepr() == "[]int")

	err := flag.FromString(" 16 ,abc,  18")
	test.Assert(t, err != nil && strings.Contains(err.Error(), "invalid syntax"))

	test.Assert(t, flag.FromString(" 42 ,  44") == nil)
	test.Assert(t, flag.String() == "[42,44]")
}

func assertEqualIfBuiltinParsingSucceed[T builtins](t *testing.T, providedRawValue string, expectedValue ...T) error {
	t.Helper()

	value, err := builtinFromString[T](providedRawValue)
	if err != nil {
		return err
	}

	test.Assert(t, len(expectedValue) == 1)
	test.Assert(t, expectedValue[0] == value)

	return nil
}

func Test_builtinFromString(t *testing.T) {
	t.Run("bool", func(t *testing.T) {
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[bool](t, "true", true) == nil)
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[bool](t, "false", false) == nil)
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[bool](t, "1", true) == nil)
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[bool](t, "0", false) == nil)

		err := assertEqualIfBuiltinParsingSucceed[bool](t, "")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "invalid syntax"))
		err = assertEqualIfBuiltinParsingSucceed[bool](t, "abc")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "invalid syntax"))
	})

	t.Run("string", func(t *testing.T) {
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[string](t, "foo", "foo") == nil)
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[string](t, "42", "42") == nil)
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[string](t, "", "") == nil)
	})

	t.Run("int", func(t *testing.T) {
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[int](t, "-9223372036854775808", -9223372036854775808) == nil)
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[int](t, "0", 0) == nil)
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[int](t, "9223372036854775807", 9223372036854775807) == nil)

		err := assertEqualIfBuiltinParsingSucceed[int](t, "")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "invalid syntax"))
		err = assertEqualIfBuiltinParsingSucceed[int](t, "-9223372036854775809")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "value out of range"))
		err = assertEqualIfBuiltinParsingSucceed[int](t, "9223372036854775808")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "value out of range"))
		err = assertEqualIfBuiltinParsingSucceed[int](t, "abc")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "invalid syntax"))
	})

	t.Run("int8", func(t *testing.T) {
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[int8](t, "-128", int8(-128)) == nil)
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[int8](t, "0", int8(0)) == nil)
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[int8](t, "127", int8(127)) == nil)

		err := assertEqualIfBuiltinParsingSucceed[int8](t, "")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "invalid syntax"))
		err = assertEqualIfBuiltinParsingSucceed[int8](t, "-129")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "value out of range"))
		err = assertEqualIfBuiltinParsingSucceed[int8](t, "128")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "value out of range"))
		err = assertEqualIfBuiltinParsingSucceed[int8](t, "abc")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "invalid syntax"))
	})

	t.Run("int16", func(t *testing.T) {
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[int16](t, "-32768", int16(-32768)) == nil)
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[int16](t, "0", int16(0)) == nil)
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[int16](t, "32767", int16(32767)) == nil)

		err := assertEqualIfBuiltinParsingSucceed[int16](t, "")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "invalid syntax"))
		err = assertEqualIfBuiltinParsingSucceed[int16](t, "-32769")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "value out of range"))
		err = assertEqualIfBuiltinParsingSucceed[int16](t, "32768")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "value out of range"))
		err = assertEqualIfBuiltinParsingSucceed[int16](t, "abc")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "invalid syntax"))
	})

	t.Run("int32", func(t *testing.T) {
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[int32](t, "-2147483648", int32(-2147483648)) == nil)
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[int32](t, "0", int32(0)) == nil)
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[int32](t, "2147483647", int32(2147483647)) == nil)

		err := assertEqualIfBuiltinParsingSucceed[int32](t, "")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "invalid syntax"))
		err = assertEqualIfBuiltinParsingSucceed[int32](t, "-2147483649")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "value out of range"))
		err = assertEqualIfBuiltinParsingSucceed[int32](t, "2147483648")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "value out of range"))
		err = assertEqualIfBuiltinParsingSucceed[int32](t, "abc")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "invalid syntax"))
	})

	t.Run("int64", func(t *testing.T) {
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[int64](t, "-9223372036854775808", int64(-9223372036854775808)) == nil)
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[int64](t, "0", int64(0)) == nil)
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[int64](t, "9223372036854775807", int64(9223372036854775807)) == nil)

		err := assertEqualIfBuiltinParsingSucceed[int64](t, "")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "invalid syntax"))
		err = assertEqualIfBuiltinParsingSucceed[int64](t, "-9223372036854775809")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "value out of range"))
		err = assertEqualIfBuiltinParsingSucceed[int64](t, "9223372036854775808")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "value out of range"))
		err = assertEqualIfBuiltinParsingSucceed[int64](t, "abc")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "invalid syntax"))
	})

	t.Run("uint", func(t *testing.T) {
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[uint](t, "0", uint(0)) == nil)
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[uint](t, "18446744073709551615", uint(18446744073709551615)) == nil)

		err := assertEqualIfBuiltinParsingSucceed[uint](t, "")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "invalid syntax"))
		err = assertEqualIfBuiltinParsingSucceed[uint](t, "-1")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "invalid syntax"))
		err = assertEqualIfBuiltinParsingSucceed[uint](t, "18446744073709551616")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "value out of range"))
		err = assertEqualIfBuiltinParsingSucceed[uint](t, "abc")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "invalid syntax"))
	})

	t.Run("uint8", func(t *testing.T) {
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[uint8](t, "0", uint8(0)) == nil)
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[uint8](t, "255", uint8(255)) == nil)

		err := assertEqualIfBuiltinParsingSucceed[uint8](t, "")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "invalid syntax"))
		err = assertEqualIfBuiltinParsingSucceed[uint8](t, "-1")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "invalid syntax"))
		err = assertEqualIfBuiltinParsingSucceed[uint8](t, "256")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "value out of range"))
		err = assertEqualIfBuiltinParsingSucceed[uint8](t, "abc")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "invalid syntax"))
	})

	t.Run("uint16", func(t *testing.T) {
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[uint16](t, "0", uint16(0)) == nil)
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[uint16](t, "65535", uint16(65535)) == nil)

		err := assertEqualIfBuiltinParsingSucceed[uint16](t, "")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "invalid syntax"))
		err = assertEqualIfBuiltinParsingSucceed[uint16](t, "-1")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "invalid syntax"))
		err = assertEqualIfBuiltinParsingSucceed[uint16](t, "65536")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "value out of range"))
		err = assertEqualIfBuiltinParsingSucceed[uint16](t, "abc")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "invalid syntax"))
	})

	t.Run("uint32", func(t *testing.T) {
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[uint32](t, "0", uint32(0)) == nil)
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[uint32](t, "4294967295", uint32(4294967295)) == nil)

		err := assertEqualIfBuiltinParsingSucceed[uint32](t, "")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "invalid syntax"))
		err = assertEqualIfBuiltinParsingSucceed[uint32](t, "-1")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "invalid syntax"))
		err = assertEqualIfBuiltinParsingSucceed[uint32](t, "4294967296")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "value out of range"))
		err = assertEqualIfBuiltinParsingSucceed[uint32](t, "abc")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "invalid syntax"))
	})

	t.Run("uint64", func(t *testing.T) {
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[uint64](t, "0", uint64(0)) == nil)
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[uint64](t, "18446744073709551615", uint64(18446744073709551615)) == nil)

		err := assertEqualIfBuiltinParsingSucceed[uint64](t, "")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "invalid syntax"))
		err = assertEqualIfBuiltinParsingSucceed[uint64](t, "-1")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "invalid syntax"))
		err = assertEqualIfBuiltinParsingSucceed[uint64](t, "18446744073709551616")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "value out of range"))
		err = assertEqualIfBuiltinParsingSucceed[uint64](t, "abc")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "invalid syntax"))
	})

	t.Run("float32", func(t *testing.T) {
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[float32](t, "-340282346638528859811704183484516925440", float32(-340282346638528859811704183484516925440)) == nil)
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[float32](t, "0", float32(0)) == nil)
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[float32](t, "0.356789", float32(0.356789)) == nil)
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[float32](t, "340282346638528859811704183484516925440", float32(340282346638528859811704183484516925440)) == nil)

		err := assertEqualIfBuiltinParsingSucceed[float32](t, "")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "invalid syntax"))
		err = assertEqualIfBuiltinParsingSucceed[float32](t, "179769313486231570814527423731704356798070567525844996598917476803157260780028538760589558632766878171540458953514382464234321326889464182768467546703537516986049910576551282076245490090389328944075868508455133942304583236903222948165808559332123348274797826204144723168738177180919299881250404026184124858368")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "value out of range"))
		err = assertEqualIfBuiltinParsingSucceed[float32](t, "-179769313486231570814527423731704356798070567525844996598917476803157260780028538760589558632766878171540458953514382464234321326889464182768467546703537516986049910576551282076245490090389328944075868508455133942304583236903222948165808559332123348274797826204144723168738177180919299881250404026184124858368")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "value out of range"))
		err = assertEqualIfBuiltinParsingSucceed[float32](t, "abc")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "invalid syntax"))
	})

	t.Run("float64", func(t *testing.T) {
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[float64](t, "-1797693134862315708145274237317043567980705675258449965989174768031572607800285387605895586327668781715404589535143824642343213268894641827684675467035", float64(-1797693134862315708145274237317043567980705675258449965989174768031572607800285387605895586327668781715404589535143824642343213268894641827684675467035)) == nil)
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[float64](t, "0", float64(0)) == nil)
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[float64](t, "0.356789", float64(0.356789)) == nil)
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[float64](t, "9797693134862915708145274237317043567980705675258449965989174768031572607800285387605895586327668781715404989535143824642343213268894641827684675467035399", float64(9797693134862915708145274237317043567980705675258449965989174768031572607800285387605895586327668781715404989535143824642343213268894641827684675467035399)) == nil)

		err := assertEqualIfBuiltinParsingSucceed[float64](t, "")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "invalid syntax"))
		err = assertEqualIfBuiltinParsingSucceed[float64](t, "-17976931348623157081452742373170435679807056752584499659891747680315726078002853876058955863276687817154045895351438246423432132688946418276846754670353751698604991057655128207624549009038932894407586850845513394230458323690322294816580855933212334827479782620414472316873817718091929988125040402618412485836899")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "value out of range"))
		err = assertEqualIfBuiltinParsingSucceed[float64](t, "17976931348623157081452742373170435679807056752584499659891747680315726078002853876058955863276687817154045895351438246423432132688946418276846754670353751698604991057655128207624549009038932894407586850845513394230458323690322294816580855933212334827479782620414472316873817718091929988125040402618412485836899")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "value out of range"))
		err = assertEqualIfBuiltinParsingSucceed[float64](t, "abc")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "invalid syntax"))
	})

	t.Run("complex64", func(t *testing.T) {
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[complex64](t, "-340282346638528859811704183484516925440-340282346638528859811704183484516925440i", complex64(-340282346638528859811704183484516925440-340282346638528859811704183484516925440i)) == nil)
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[complex64](t, "0", complex64(0)) == nil)
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[complex64](t, "-0.356789i", complex64(-0.356789i)) == nil)
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[complex64](t, "340282346638528859811704183484516925440+340282346638528859811704183484516925440i", complex64(340282346638528859811704183484516925440+340282346638528859811704183484516925440i)) == nil)

		err := assertEqualIfBuiltinParsingSucceed[complex64](t, "")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "invalid syntax"))
		err = assertEqualIfBuiltinParsingSucceed[complex64](t, "i")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "invalid syntax"))
		err = assertEqualIfBuiltinParsingSucceed[complex64](t, "179769313486231570814527423731704356798070567525844996598917476803157260780028538760589558632766878171540458953514382464234321326889464182768467546703537516986049910576551282076245490090389328944075868508455133942304583236903222948165808559332123348274797826204144723168738177180919299881250404026184124858368")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "value out of range"))
		err = assertEqualIfBuiltinParsingSucceed[complex64](t, "-179769313486231570814527423731704356798070567525844996598917476803157260780028538760589558632766878171540458953514382464234321326889464182768467546703537516986049910576551282076245490090389328944075868508455133942304583236903222948165808559332123348274797826204144723168738177180919299881250404026184124858368")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "value out of range"))
		err = assertEqualIfBuiltinParsingSucceed[complex64](t, "abc")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "invalid syntax"))
	})

	t.Run("complex128", func(t *testing.T) {
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[complex128](t, "-1797693134862315708145274237317043567980705675258449965989174768031572607800285387605895586327668781715404589535143824642343213268894641827684675467035-1797693134862315708145274237317043567980705675258449965989174768031572607800285387605895586327668781715404589535143824642343213268894641827684675467035i", -1797693134862315708145274237317043567980705675258449965989174768031572607800285387605895586327668781715404589535143824642343213268894641827684675467035-1797693134862315708145274237317043567980705675258449965989174768031572607800285387605895586327668781715404589535143824642343213268894641827684675467035i) == nil)
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[complex128](t, "0", 0) == nil)
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[complex128](t, "-0.356789i", -0.356789i) == nil)
		test.Assert(t, assertEqualIfBuiltinParsingSucceed[complex128](t, "1797693134862315708145274237317043567980705675258449965989174768031572607800285387605895586327668781715404589535143824642343213268894641827684675467035+1797693134862315708145274237317043567980705675258449965989174768031572607800285387605895586327668781715404589535143824642343213268894641827684675467035i", 1797693134862315708145274237317043567980705675258449965989174768031572607800285387605895586327668781715404589535143824642343213268894641827684675467035+1797693134862315708145274237317043567980705675258449965989174768031572607800285387605895586327668781715404589535143824642343213268894641827684675467035i) == nil)

		err := assertEqualIfBuiltinParsingSucceed[complex128](t, "")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "invalid syntax"))
		err = assertEqualIfBuiltinParsingSucceed[complex128](t, "i")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "invalid syntax"))
		err = assertEqualIfBuiltinParsingSucceed[complex128](t, "17976931348623157081452742373170435679807056752584499659891747680315726078002853876058955863276687817154045895351438246423432132688946418276846754670353751698604991057655128207624549009038932894407586850845513394230458323690322294816580855933212334827479782620414472316873817718091929988125040402618412485836899")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "value out of range"))
		err = assertEqualIfBuiltinParsingSucceed[complex128](t, "-17976931348623157081452742373170435679807056752584499659891747680315726078002853876058955863276687817154045895351438246423432132688946418276846754670353751698604991057655128207624549009038932894407586850845513394230458323690322294816580855933212334827479782620414472316873817718091929988125040402618412485836899")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "value out of range"))
		err = assertEqualIfBuiltinParsingSucceed[complex128](t, "abc")
		test.Assert(t, err != nil && strings.Contains(err.Error(), "invalid syntax"))
	})
}

func assertNoErrorBuiltinParsingSucceed[T builtins](t *testing.T, providedValue T, expectedValueRepr string) {
	t.Helper()

	repr := builtinToString[T](providedValue)
	test.Assert(t, repr == expectedValueRepr)
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
