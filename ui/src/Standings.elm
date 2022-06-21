module Standings exposing (..)

import DateRange exposing (DateRange, dateRangeEncoder)
import Errors exposing (buildErrorMessage)
import Html exposing (..)
import Html.Attributes exposing (..)
import Http
import RemoteData exposing (WebData)
import SingleStanding exposing (Standing, standingsResponseDecoder)


type alias Model =
    { standings : WebData (List Standing)
    , dateRange : DateRange
    }


type Msg
    = FetchStandings DateRange
    | StandingsReceived (WebData (List Standing))


init : ( Model, Cmd Msg )
init =
    let
        model =
            { standings = RemoteData.Loading
            , dateRange =
                { beginDate = "2022-01-01T00:00:00Z"
                , endDate = "2023-01-01T00:00:00Z"
                }
            }
    in
    ( model, fetchStandings model.dateRange )


fetchStandings : DateRange -> Cmd Msg
fetchStandings dateRange =
    Http.post
        { url = "http://localhost:8082/twirp/tshparser.TournamentRankerService/ComputeStandings"
        , expect =
            standingsResponseDecoder
                |> Http.expectJson (RemoteData.fromResult >> StandingsReceived)
        , body = Http.jsonBody (dateRangeEncoder dateRange)
        }


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        FetchStandings dateRange ->
            ( { model | standings = RemoteData.Loading }, fetchStandings dateRange )

        StandingsReceived response ->
            ( { model | standings = response }, Cmd.none )



-- VIEWS


view : Model -> Html Msg
view model =
    div []
        [ viewStandings model.standings ]


viewStandings : WebData (List Standing) -> Html Msg
viewStandings standings =
    case standings of
        RemoteData.NotAsked ->
            text ""

        RemoteData.Loading ->
            h3 [ class "subtitle is-2" ] [ text "Loading..." ]

        RemoteData.Success actualStandings ->
            div []
                [ h2 [ class "subtitle is-2" ] [ text "Standings" ]
                , table [ class "table" ]
                    (viewTableHeader :: List.map viewStanding actualStandings)
                ]

        RemoteData.Failure httpError ->
            viewStandingsError (buildErrorMessage httpError)


viewTableHeader : Html Msg
viewTableHeader =
    tr []
        [ th []
            [ abbr [ title "Position" ]
                [ text "Pos" ]
            ]
        , th []
            [ text "Name" ]
        , th []
            [ text "Points" ]
        , th []
            [ text "Wins" ]
        , th []
            [ text "Spread"
            ]
        , th []
            [ text "Games played"
            ]
        , th []
            [ text "Tournaments"
            ]
        ]


viewStanding : Standing -> Html Msg
viewStanding standing =
    tr []
        [ td [] [ text (String.fromInt standing.rank) ]
        , td [] [ text standing.playerName ]
        , td [] [ text (String.fromInt standing.points) ]
        , td [] [ text (String.fromFloat standing.wins) ]
        , td [] [ text (String.fromInt standing.spread) ]
        , td [] [ text (String.fromInt standing.games) ]
        , td [] [ text (String.fromInt standing.tournamentsPlayed) ]
        ]


viewStandingsError : String -> Html Msg
viewStandingsError errorMsg =
    let
        errorHeading =
            "Couldn't fetch standings at this time."
    in
    div []
        [ h3 [] [ text errorHeading ]
        , text ("Error: " ++ errorMsg)
        ]
