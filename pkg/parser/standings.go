package parser

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
)

var (
	singlePosRegex = regexp.MustCompile(`(\d+)`)
	rangePosRegex  = regexp.MustCompile(`(\d+)-(\d+)`)
)

func createPtMap(schemaFile string) (map[string]int, error) {
	f, err := os.Open(schemaFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	// header format of the schema file should be as follows:
	// posrange,ttype1,ttype2,ttype3,....,ttypeN
	idx := 0

	// create a map of position and tournament type to points.
	ptmap := make(map[string]int)
	var ttypes []string
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		// rec is a list of records.
		if idx == 0 {
			ttypes = rec[1:]
			idx += 1
			continue
		}

		posrange := rec[0]
		positions := []int{}
		res := rangePosRegex.FindStringSubmatch(posrange)
		if len(res) == 0 {
			// didn't match, try number
			res := singlePosRegex.FindStringSubmatch(posrange)
			if len(res) == 0 {
				return nil, errors.New("malformed first column " + posrange)
			}
			pl, err := strconv.Atoi(res[1])
			if err != nil {
				return nil, err
			}
			positions = append(positions, pl)
		} else {
			beg, err := strconv.Atoi(res[1])
			if err != nil {
				return nil, err
			}
			end, err := strconv.Atoi(res[2])
			if err != nil {
				return nil, err
			}
			for i := beg; i <= end; i++ {
				positions = append(positions, i)
			}
		}
		for idx, ptstr := range rec[1:] {
			for _, p := range positions {
				pts, err := strconv.Atoi(ptstr)
				if err != nil {
					return nil, err
				}
				ptmap[fmt.Sprintf("%d:%s", p, ttypes[idx])] = pts
			}
		}
		idx += 1
	}
	return ptmap, nil
}

func computeStandings(tourneys []Tournament, schemaFile string) ([]Standing, error) {

	return nil, nil
}
