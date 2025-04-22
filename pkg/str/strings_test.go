package str_test

import (
	"testing"

	"github.com/hiendaovinh/toolkit/v2/pkg/str"
	"github.com/stretchr/testify/assert"
)

func TestTrimStrings(t *testing.T) {
	in := []string{" 1", "2 ", " 3 "}
	out := str.StrTrim(in...)
	assert.Equal(t, out, []string{"1", "2", "3"})
}

func TestIsNumber(t *testing.T) {
	x := "1234000000000"
	out := str.StrIsNumber(x)
	assert.Equal(t, out, true)
}

func TestIsNumber2(t *testing.T) {
	x := "123400a000000"
	out := str.StrIsNumber(x)
	assert.Equal(t, out, false)
}

func TestIsNumber3(t *testing.T) {
	x := "123400.000000"
	out := str.StrIsNumber(x)
	assert.Equal(t, out, true)
}

func TestIsNumber4(t *testing.T) {
	x := "123.400.000000"
	out := str.StrIsNumber(x)
	assert.Equal(t, out, false)
}

func TestIsNumber5(t *testing.T) {
	x := "123,400000000"
	out := str.StrIsNumber(x)
	assert.Equal(t, out, false)
}

func TestTrimSpecialCharacter(t *testing.T) {
	in := "&&||!(){}[]^\"~*?:"
	out := str.StrTrimSpecialCharacter(in)
	assert.Equal(t, "", out)
}
