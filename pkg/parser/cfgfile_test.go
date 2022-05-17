package parser

import (
	"io/ioutil"
	"testing"

	"github.com/matryer/is"
)

func TestGetDivisionFilename(t *testing.T) {
	is := is.New(t)
	contents, err := ioutil.ReadFile("./testdata/config.tsh")
	is.NoErr(err)
	fn, err := GetDivisionFilename(contents, "A")
	is.NoErr(err)
	is.Equal(fn, "a.t")

	fn, err = GetDivisionFilename(contents, "B")
	is.NoErr(err)
	is.Equal(fn, "b.t")
	_, err = GetDivisionFilename(contents, "Dogs")
	is.Equal(err.Error(), "division Dogs not found")
}
