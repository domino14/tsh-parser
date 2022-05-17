package server

import (
	"regexp"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/domino14/tshparser/pkg/parser"

	"encoding/json"
	"io/ioutil"
	"net/http"
)

var URLRegex = regexp.MustCompile(`(https?://.+)html/(\w)+-standings-\d+.html`)

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
		err = s.service.AddTournament(r.Context(), parser.TournamentType(ttype), name, date, divcontent)
		if err != nil {
			errResponse(w, 500, "could not add tournament: "+err.Error())
		}
		return
	case "remove":

	case "standings":
	}

	w.WriteHeader(200)
	w.Write([]byte(`{"msg": "ok"}`))
}
