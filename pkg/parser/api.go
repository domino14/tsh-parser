package parser

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	proto "github.com/domino14/tshparser/rpc"
	"github.com/rs/zerolog/log"
)

var URLRegex = regexp.MustCompile(`(https?://.+)html/([\w]+)-standings(?:-\d+)?.html`)

type Service struct {
	store      *SqliteStore
	schemaPath string
	secretKey  string
}

func NewService(store *SqliteStore, schemaPath, secretKey string) *Service {
	return &Service{store: store, schemaPath: schemaPath, secretKey: secretKey}
}

func adminRequired(ctx context.Context, secretKey string) error {
	u, err := userFromHTTPHeader(ctx, secretKey)
	if err != nil {
		return err
	}
	// Note: This doesn't check the database, but the IsAdmin is embedded in
	// the JWT. This is ok for this small app. Revisit this if this turns into
	// a problem.
	if !u.IsAdmin {
		return errors.New("user is not an admin")
	}
	return nil
}

func tshFileContents(turl string) ([]byte, error) {
	// turl looks like  https://tourneys.mindgamesincorporated.com/mspl122/html/A-standings-012.html
	// we need to get https://tourneys.mindgamesincorporated.com/mspl122/a.t
	matches := URLRegex.FindStringSubmatch(turl)
	if len(matches) != 3 {
		return nil, errors.New("turl not in right format")
	}
	resp, err := http.Get(matches[1] + "config.tsh")
	if err != nil {
		return nil, errors.New("config.tsh not found")
	}
	defer resp.Body.Close()
	cfgbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("could not read config.tsh: " + err.Error())
	}
	divFilename, err := GetDivisionFilename(cfgbody, matches[2])
	if err != nil {
		return nil, err
	}
	fullDivFilename := matches[1] + divFilename
	log.Info().Str("div-filename", fullDivFilename).Msg("fetching..")
	dresp, err := http.Get(fullDivFilename)
	if err != nil {
		return nil, errors.New(divFilename + " not found")
	}
	defer dresp.Body.Close()
	divcontent, err := ioutil.ReadAll(dresp.Body)
	if err != nil {
		return nil, errors.New("could not read division file: " + err.Error())
	}
	return divcontent, nil
}

func (s *Service) AddTournament(ctx context.Context, req *proto.AddTournamentRequest) (*proto.AddTournamentResponse, error) {
	err := adminRequired(ctx, s.secretKey)
	if err != nil {
		return nil, err
	}
	date, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		return nil, err
	}
	if req.Name == "" {
		return nil, errors.New("name required")
	}
	contents, err := tshFileContents(req.TshUrl)
	if err != nil {
		return nil, err
	}
	id, err := s.store.AddTournament(ctx, req.TournamentType, req.Name, date, contents)
	if err != nil {
		return nil, err
	}
	return &proto.AddTournamentResponse{Id: id}, nil
}

func (s *Service) RemoveTournament(ctx context.Context, req *proto.RemoveTournamentRequest) (*proto.RemoveTournamentResponse, error) {
	err := adminRequired(ctx, s.secretKey)
	if err != nil {
		return nil, err
	}
	err = s.store.RemoveTournament(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}
	return &proto.RemoveTournamentResponse{}, nil
}

func (s *Service) ComputeStandings(ctx context.Context, req *proto.ComputeStandingsRequest) (*proto.StandingsResponse, error) {
	dbegin, err := time.Parse(time.RFC3339, req.BeginDate)
	if err != nil {
		return nil, err
	}
	dend, err := time.Parse(time.RFC3339, req.EndDate)
	if err != nil {
		return nil, err
	}
	tourneys, err := s.store.GetTournaments(ctx, dbegin, dend)
	if err != nil {
		return nil, err
	}
	aliases, err := s.store.GetAllAliases(ctx)
	if err != nil {
		return nil, err
	}

	standings, err := computeStandings(tourneys, s.schemaPath, aliases)
	if err != nil {
		return nil, err
	}
	return &proto.StandingsResponse{Standings: standings}, nil
}

func (s *Service) GetTournaments(ctx context.Context, req *proto.GetTournamentsRequest) (*proto.TournamentsResponse, error) {
	dbegin, err := time.Parse(time.RFC3339, req.BeginDate)
	if err != nil {
		return nil, err
	}
	dend, err := time.Parse(time.RFC3339, req.EndDate)
	if err != nil {
		return nil, err
	}
	ts, err := s.store.GetTournaments(ctx, dbegin, dend)
	if err != nil {
		return nil, err
	}
	return &proto.TournamentsResponse{Tournaments: ts}, nil
}

func (s *Service) AddPlayerAlias(ctx context.Context, req *proto.PlayerAlias) (*proto.AddPlayerAliasResponse, error) {
	err := adminRequired(ctx, s.secretKey)
	if err != nil {
		return nil, err
	}
	origPlayer := strings.TrimSpace(req.OriginalPlayer)
	alias := strings.TrimSpace(req.Alias)
	if origPlayer == alias {
		return nil, errors.New("origPlayer must not be equal to alias")
	}
	if origPlayer == "" || alias == "" {
		return nil, errors.New("both origPlayer and alias must be specified")
	}
	err = s.store.AddPlayerAlias(ctx, origPlayer, alias)
	if err != nil {
		return nil, err
	}
	return &proto.AddPlayerAliasResponse{}, nil
}

func (s *Service) RemovePlayerAlias(ctx context.Context, req *proto.RemovePlayerAliasRequest) (*proto.RemovePlayerAliasResponse, error) {
	err := adminRequired(ctx, s.secretKey)
	if err != nil {
		return nil, err
	}
	err = s.store.RemovePlayerAlias(ctx, strings.TrimSpace(req.Alias))
	if err != nil {
		return nil, err
	}
	return &proto.RemovePlayerAliasResponse{}, nil
}

func (s *Service) ListPlayerAliases(ctx context.Context, req *proto.ListPlayerAliasesRequest) (*proto.PlayerAliasesResponse, error) {
	aliases, err := s.store.GetAllAliases(ctx)
	if err != nil {
		return nil, err
	}
	// aliases is a map; convert to a sorted list of aliases:
	aliasList := []*proto.PlayerAlias{}
	for alias, orig := range aliases {
		aliasList = append(aliasList, &proto.PlayerAlias{
			OriginalPlayer: orig,
			Alias:          alias,
		})
	}
	sort.Slice(aliasList, func(i, j int) bool {
		return aliasList[i].OriginalPlayer < aliasList[j].OriginalPlayer
	})

	return &proto.PlayerAliasesResponse{Aliases: aliasList}, nil
}
