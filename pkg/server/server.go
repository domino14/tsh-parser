package server

import (
	"fmt"
	"regexp"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/domino14/tshparser/pkg/parser"

	"encoding/json"
	"io/ioutil"
	"net/http"
)

var URLRegex = regexp.MustCompile(`(https?://.+)html/([\w]+)-standings-\d+.html`)

// an http server
type Server struct {
	service *parser.Service
}

func NewHTTPServer(service *parser.Service) *Server {
	return &Server{service: service}
}

type Request struct {
	Method string                 `json:"method"`
	Params map[string]interface{} `json:"params"`
}

func errResponse(w http.ResponseWriter, code int, err string) {
	w.WriteHeader(code)
	w.Write([]byte(`{"error": "` + err + `"}`))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	if r.Header.Get("Content-Type") != "application/json" {
		errResponse(w, 400, "bad content type")
		return
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errResponse(w, 500, err.Error())
		return
	}
	req := &Request{}
	err = json.Unmarshal(body, req)
	if err != nil {
		errResponse(w, 500, err.Error())
		return
	}
	log.Info().Interface("req", req).Msg("got request")
	if r.Method != "POST" {
		errResponse(w, 400, "must use POST method")
		return
	}
	// this is a terrible function and i should have used protobuf
	switch req.Method {
	case "add":
		// add tournament
		userdate, ok := req.Params["date"].(string)
		if !ok {
			errResponse(w, 400, "date not a string")
			return
		}
		ttype, ok := req.Params["ttype"].(string)
		if !ok {
			errResponse(w, 400, "ttype not a string")
			return
		}
		name, ok := req.Params["name"].(string)
		if !ok {
			errResponse(w, 400, "name not a string")
			return
		}

		date, err := time.Parse(time.RFC3339, userdate)
		if err != nil {
			errResponse(w, 400, err.Error())
			return
		}

		turl, ok := req.Params["tshurl"].(string)
		if !ok {
			errResponse(w, 400, "turl not a string")
			return
		}
		// turl looks like  https://tourneys.mindgamesincorporated.com/mspl122/html/A-standings-012.html
		// we need to get https://tourneys.mindgamesincorporated.com/mspl122/a.t
		matches := URLRegex.FindStringSubmatch(turl)
		if len(matches) != 3 {
			errResponse(w, 400, "turl not in right format")
			return
		}
		resp, err := http.Get(matches[1] + "config.tsh")
		if err != nil {
			errResponse(w, 400, "config.tsh not found")
			return
		}
		defer resp.Body.Close()
		cfgbody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			errResponse(w, 500, "could not read config.tsh: "+err.Error())
			return
		}
		divFilename, err := parser.GetDivisionFilename(cfgbody, matches[2])
		if err != nil {
			errResponse(w, 500, err.Error())
		}
		fullDivFilename := matches[1] + divFilename
		log.Info().Str("div-filename", fullDivFilename).Msg("fetching..")
		dresp, err := http.Get(fullDivFilename)
		if err != nil {
			errResponse(w, 400, divFilename+" not found")
			return
		}
		defer dresp.Body.Close()
		divcontent, err := ioutil.ReadAll(dresp.Body)
		if err != nil {
			errResponse(w, 500, "could not read division file: "+err.Error())
			return
		}
		id, err := s.service.AddTournament(r.Context(), parser.TournamentType(ttype), name, date, divcontent)
		if err != nil {
			errResponse(w, 500, "could not add tournament: "+err.Error())
			return
		}
		w.WriteHeader(201)
		w.Write([]byte(fmt.Sprintf(`{"id": %d}`, id)))
		return

	case "remove":
		// json parses numbers as float64s
		tid, ok := req.Params["tid"].(float64)
		if !ok {
			errResponse(w, 400, "tid not a number")
			return
		}
		if tid != float64(int(tid)) {
			errResponse(w, 400, "tid should be an integer")
			return
		}
		err := s.service.RemoveTournament(r.Context(), int(tid))
		if err != nil {
			errResponse(w, 500, "could not remove tournament: "+err.Error())
			return
		}

	case "standings":
		begin, ok := req.Params["begin"].(string)
		if !ok {
			errResponse(w, 400, "begin not a string")
			return
		}
		bdate, err := time.Parse(time.RFC3339, begin)
		if err != nil {
			errResponse(w, 400, err.Error())
			return
		}

		end, ok := req.Params["end"].(string)
		if !ok {
			errResponse(w, 400, "end not a string")
			return
		}
		edate, err := time.Parse(time.RFC3339, end)
		if err != nil {
			errResponse(w, 400, err.Error())
			return
		}

		st, err := s.service.ComputeStandings(r.Context(), bdate, edate)
		if err != nil {
			errResponse(w, 500, "error computing standings: "+err.Error())
			return
		}
		bts, err := json.Marshal(st)
		if err != nil {
			errResponse(w, 500, "error marshalling standings: "+err.Error())
			return
		}
		w.Write(bts)
		return

	case "tournaments":
		begin, ok := req.Params["begin"].(string)
		if !ok {
			errResponse(w, 400, "begin not a string")
			return
		}
		bdate, err := time.Parse(time.RFC3339, begin)
		if err != nil {
			errResponse(w, 400, err.Error())
			return
		}

		end, ok := req.Params["end"].(string)
		if !ok {
			errResponse(w, 400, "end not a string")
			return
		}
		edate, err := time.Parse(time.RFC3339, end)
		if err != nil {
			errResponse(w, 400, err.Error())
			return
		}
		ts, err := s.service.GetTournaments(r.Context(), bdate, edate)
		if err != nil {
			errResponse(w, 500, "error getting tournaments: "+err.Error())
			return
		}
		for ti := range ts {
			ts[ti].Standings, err = parser.SingleTourneyStandings(ts[ti].Contents)
			if err != nil {
				errResponse(w, 500, "error computing standings: "+err.Error())
				return
			}
			ts[ti].Contents = nil
		}
		bts, err := json.Marshal(ts)
		if err != nil {
			errResponse(w, 500, "error marshalling standings: "+err.Error())
			return
		}
		w.Write(bts)
		return

	case "addalias":
		origPlayer, ok := req.Params["origPlayer"].(string)
		if !ok {
			errResponse(w, 400, "origPlayer not a string")
			return
		}
		alias, ok := req.Params["alias"].(string)
		if !ok {
			errResponse(w, 400, "alias not a string")
			return
		}
		err := s.service.AddPlayerAlias(r.Context(), origPlayer, alias)
		if err != nil {
			errResponse(w, 500, "error adding alias: "+err.Error())
			return
		}

	case "removealias":
		origPlayer, ok := req.Params["origPlayer"].(string)
		if !ok {
			errResponse(w, 400, "origPlayer not a string")
			return
		}
		alias, ok := req.Params["alias"].(string)
		if !ok {
			errResponse(w, 400, "alias not a string")
			return
		}
		err := s.service.RemovePlayerAlias(r.Context(), origPlayer, alias)
		if err != nil {
			errResponse(w, 500, "error removing alias: "+err.Error())
			return
		}
	default:
		errResponse(w, 400, "method not handled")
		return
	}

	w.Write([]byte(`{"msg": "ok"}`))
}
