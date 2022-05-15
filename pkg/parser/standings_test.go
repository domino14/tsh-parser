package parser

import (
	"bytes"
	"flag"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"testing"

	"github.com/matryer/is"
)

var goldenFileUpdate bool

func init() {
	flag.BoolVar(&goldenFileUpdate, "update", false, "update golden files")
}

func slurp(filename string) string {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return string(contents)
}

func updateGolden(filename string, bts []byte) {
	// write the bts to filename
	os.WriteFile(filename, bts, 0600)
}

func compareGolden(t *testing.T, goldenFile string, actualRepr []byte) {
	is := is.New(t)
	if goldenFileUpdate {
		updateGolden(goldenFile, actualRepr)
	} else {
		expected := slurp(goldenFile)
		is.Equal(expected, string(actualRepr))
	}
}

func serializePtMap(m map[string]int) []byte {
	keys := []string{}
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b bytes.Buffer
	for _, k := range keys {
		b.WriteString(k)
		b.WriteString(": ")
		b.WriteString(strconv.Itoa(m[k]))
		b.WriteString("\n")
	}
	return b.Bytes()
}

func TestCreatePtMap(t *testing.T) {
	is := is.New(t)
	ptMap, err := createPtMap("../../cfg/pts_mgi.csv")
	is.NoErr(err)
	compareGolden(t, "./testdata/pts_mgi.golden", serializePtMap(ptMap))
}
