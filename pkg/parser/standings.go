package parser

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"

	proto "github.com/domino14/tshparser/rpc"
)

var (
	singlePosRegex = regexp.MustCompile(`(\d+)`)
	rangePosRegex  = regexp.MustCompile(`(\d+)-(\d+)`)

	tFileFirstFieldRegex = regexp.MustCompile(`([\pL,\.\s\d]+)[\s]{2,}[\d]+(.+)`)
)

// t-files are iso-8859-1 by default. let's ignore this for now.
// config.tsh and archive.tsh are UTF8

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

func deleteByes(standings map[string]*proto.Standing) {
	toDelete := []string{}
	for s := range standings {
		if strings.Contains(s, "bye") ||
			strings.Contains(s, "Bye") ||
			strings.Contains(s, "Gracious Picker") ||
			strings.Contains(s, "GraciousPicker") {
			toDelete = append(toDelete, s)
		}
	}
	for _, t := range toDelete {
		delete(standings, t)
	}
	log.Info().Interface("toDelete", toDelete).Msg("deleting from standings")
}

func computeStandings(tourneys []*proto.Tournament, schemaFile string, aliases map[string]string) ([]*proto.Standing, error) {
	ptmap, err := createPtMap(schemaFile)
	if err != nil {
		return nil, err
	}

	playerStandings := make(map[string]*proto.Standing)

	// Now, iterate through all the tournaments and assign scores accordingly.
	for _, t := range tourneys {
		sts, err := SingleTourneyStandings(t.TfileContents)
		if err != nil {
			return nil, err
		}
		for si, s := range sts {
			// assign pts for this tournament:
			pts, ok := ptmap[fmt.Sprintf("%d:%s", si+1, t.TournamentType)]
			if !ok {
				pts = 100
				log.Warn().Msgf("place %d had no entry for tournament type %s ... defaulting to 100 pts", si+1, t.TournamentType)
			}
			s.Points = int32(pts)
			playerStandings[s.PlayerName] = aggregate(playerStandings[s.PlayerName], s)
		}
	}

	deleteByes(playerStandings)

	// now aggregate and remove aliases
	for alias, realPlayer := range aliases {
		if _, ok := playerStandings[alias]; ok {
			playerStandings[realPlayer] = aggregate(playerStandings[realPlayer], playerStandings[alias])
			delete(playerStandings, alias)
			// make sure the real player name gets passed through correctly,
			// as the aggregate function doesn't do this.
			playerStandings[realPlayer].PlayerName = realPlayer
		}
	}
	vals := []*proto.Standing{}
	for _, v := range playerStandings {
		vals = append(vals, v)
	}
	sort.Slice(vals, func(i, j int) bool {
		if vals[i].Points == vals[j].Points {
			if vals[i].Wins == vals[j].Wins {
				if vals[i].Spread == vals[j].Spread {
					return vals[i].PlayerName < vals[j].PlayerName
				}
				return vals[i].Spread > vals[j].Spread
			}
			return vals[i].Wins > vals[j].Wins
		}
		return vals[i].Points > vals[j].Points
	})

	// re-assign ranks
	for idx := range vals {
		vals[idx].Rank = int32(idx + 1)
	}
	return vals, nil
}

func SingleTourneyStandings(tfileContents []byte) ([]*proto.Standing, error) {
	reader := bytes.NewReader(tfileContents)
	bufReader := bufio.NewReader(reader)

	standings := []*proto.Standing{}
	matchups := [][]int{}
	allScores := [][]int{}
	for {
		line, err := bufReader.ReadString('\n')
		if err != nil {
			if line == "" {
				break
			}
		}
		// ignore blank lines
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}
		fields := strings.Split(line, ";")

		if len(fields) < 2 {
			return nil, errors.New("bad line")
		}
		firstField := tFileFirstFieldRegex.FindStringSubmatch(fields[0])
		if len(firstField) != 3 {
			return nil, fmt.Errorf("badly formatted first field: %v", fields[0])
		}
		s := strings.Fields(fields[1])
		if len(s) == 0 {
			return nil, errors.New("badly formatted scores")
		}
		pname := strings.TrimSpace(firstField[1])
		//  I don't know just some characters here that might be erroneously at the beginning
		// or end of a name:
		pname = strings.Trim(pname, `.,;:<>"-=+`)
		m := strings.Fields(strings.TrimSpace(firstField[2]))
		if len(m) != len(s) {
			return nil, fmt.Errorf("matchups don't match scores (is the tournament still going?) %v %v", m, s)
		}
		numericMatchups := make([]int, len(m))
		scores := make([]int, len(s))

		for i := 0; i < len(m); i++ {
			numericMatchups[i], err = strconv.Atoi(m[i])
			if err != nil {
				return nil, err
			}
			scores[i], err = strconv.Atoi(s[i])
			if err != nil {
				return nil, err
			}
		}
		matchups = append(matchups, numericMatchups)
		allScores = append(allScores, scores)
		standings = append(standings, &proto.Standing{PlayerName: pname})
	}
	// Now compute wins and spread
	for i := range standings {
		for mi, m := range matchups[i] {
			var theirScore int
			ourScore := allScores[i][mi]
			if m == 0 {
				// it's a bye or forfeit
				theirScore = 0
			} else {
				theirScore = allScores[m-1][mi] // player numbers are 1-indexed
			}
			if ourScore > theirScore {
				standings[i].Wins += 1
			} else if ourScore == theirScore {
				standings[i].Wins += 0.5
			} // otherwise, we lost, but the other player's standing will take care of this.
			standings[i].Spread += int32(ourScore - theirScore)
			standings[i].Games += 1 // XXX: what if it's a bye?
		}
		standings[i].TournamentsPlayed = 1
	}
	// Sort standings by wins then spread.
	sort.Slice(standings, func(i, j int) bool {
		if standings[i].Wins == standings[j].Wins {
			if standings[i].Spread == standings[j].Spread {
				return standings[i].PlayerName < standings[j].PlayerName
			}
			return standings[i].Spread > standings[j].Spread
		}
		return standings[i].Wins > standings[j].Wins
	})

	for idx := range standings {
		standings[idx].Rank = int32(idx + 1)
	}

	return standings, nil
}

func aggregate(origStanding, toAdd *proto.Standing) *proto.Standing {
	if origStanding == nil {
		// Use the zero-values if origStanding is nil.
		origStanding = &proto.Standing{}
	}
	// "add" the standing to origStanding
	st := &proto.Standing{
		PlayerName:        toAdd.PlayerName,
		Points:            origStanding.Points + toAdd.Points,
		Wins:              origStanding.Wins + toAdd.Wins,
		Spread:            origStanding.Spread + toAdd.Spread,
		TournamentsPlayed: origStanding.TournamentsPlayed + toAdd.TournamentsPlayed,
		Games:             origStanding.Games + toAdd.Games,
	}
	return st
}
