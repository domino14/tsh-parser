module ListTournaments exposing (..)

import DateRange exposing (DateRange, dateRangeEncoder)
import Errors exposing (buildErrorMessage)
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onClick)
import Http
import Json.Encode as Encode
import RemoteData exposing (WebData)
import Tournament exposing (Tournament, tournamentsResponseDecoder)


type alias Model =
    { tournaments : WebData (List Tournament)
    , dateRange : DateRange
    , deleteError : Maybe String
    , jwt : String
    }


type Msg
    = FetchTournaments DateRange
    | TournamentsReceived (WebData (List Tournament))
    | DeleteTournament String
    | TournamentDeleted (Result Http.Error String)


init : String -> ( Model, Cmd Msg )
init jwt =
    let
        model =
            { tournaments = RemoteData.Loading
            , dateRange =
                { beginDate = "2022-01-01T00:00:00Z"
                , endDate = "2023-01-01T00:00:00Z"
                }
            , deleteError = Nothing
            , jwt = jwt
            }
    in
    ( model, fetchTournaments model.dateRange )


fetchTournaments : DateRange -> Cmd Msg
fetchTournaments dateRange =
    Http.post
        -- XXX: use some sort of env var on prod?
        { url = "http://localhost:8082/twirp/tshparser.TournamentRankerService/GetTournaments"
        , expect =
            tournamentsResponseDecoder
                |> Http.expectJson (RemoteData.fromResult >> TournamentsReceived)
        , body = Http.jsonBody (dateRangeEncoder dateRange)
        }


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        FetchTournaments dateRange ->
            ( { model | tournaments = RemoteData.Loading }, fetchTournaments dateRange )

        TournamentsReceived response ->
            ( { model | tournaments = response }, Cmd.none )

        DeleteTournament tid ->
            ( model, deleteTournament tid )

        TournamentDeleted (Ok _) ->
            ( model, fetchTournaments model.dateRange )

        TournamentDeleted (Err error) ->
            ( { model | deleteError = Just (buildErrorMessage error) }
            , Cmd.none
            )


deleteTournament : String -> Cmd Msg
deleteTournament tid =
    Http.post
        { url = "http://localhost:8082/twirp/tshparser.TournamentRankerService/RemoveTournament"
        , expect = Http.expectString TournamentDeleted -- "{}"
        , body = Http.jsonBody (tourneyIDEncoder tid)
        }


tourneyIDEncoder : String -> Encode.Value
tourneyIDEncoder tid =
    Encode.object
        [ ( "id", Encode.string tid )
        ]



-- VIEWS


view : Model -> Html Msg
view model =
    div []
        [ button [ class "button", onClick (FetchTournaments model.dateRange) ]
            [ text "Refresh tournaments" ]
        , br [] []
        , br [] []
        , a [ href "/tournaments/new" ]
            [ text "Add a new tournament" ]
        , br [] []
        , br [] []
        , a [ href "/standings" ]
            [ text "View standings to date" ]
        , br [] []
        , br [] []
        , viewTournaments model.tournaments
        , viewDeleteError model.deleteError
        ]


viewTournaments : WebData (List Tournament) -> Html Msg
viewTournaments tournaments =
    case tournaments of
        RemoteData.NotAsked ->
            text ""

        RemoteData.Loading ->
            h3 [ class "subtitle is-2" ] [ text "Loading..." ]

        RemoteData.Success actualTournaments ->
            div []
                [ h2 [ class "subtitle is-2" ] [ text "Tournaments" ]
                , table [ class "table" ]
                    (viewTableHeader :: List.map viewTournament actualTournaments)
                ]

        RemoteData.Failure httpError ->
            viewTournamentError (buildErrorMessage httpError)


viewTableHeader : Html Msg
viewTableHeader =
    tr []
        [ th [] [ text "Date" ]
        , th [] [ text "Category" ]
        , th [] [ text "Name" ]
        , th [] [ text "" ]
        ]


viewTournament : Tournament -> Html Msg
viewTournament tournament =
    tr []
        [ td [] [ text tournament.date ]
        , td [] [ text tournament.category ]
        , td [] [ text tournament.name ]
        , td []
            [ button [ class "button is-warning", type_ "button", onClick (DeleteTournament tournament.id) ]
                [ text "Delete" ]
            ]
        ]


viewTournamentError : String -> Html Msg
viewTournamentError errorMsg =
    let
        errorHeading =
            "Couldn't fetch tournaments at this time."
    in
    div []
        [ h3 [] [ text errorHeading ]
        , text ("Error: " ++ errorMsg)
        ]


viewDeleteError : Maybe String -> Html msg
viewDeleteError maybeError =
    case maybeError of
        Just error ->
            div []
                [ h3 [] [ text "Couldn't delete tournament at this time." ]
                , text ("Error: " ++ error)
                ]

        Nothing ->
            text ""
