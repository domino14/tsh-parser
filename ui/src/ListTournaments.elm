module ListTournaments exposing (..)

import DateRange exposing (DateRange, dateRangeEncoder)
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
    }


type Msg
    = FetchTournaments DateRange
    | TournamentsReceived (WebData (List Tournament))
    | DeleteTournament String
    | TournamentDeleted (Result Http.Error String)


init : ( Model, Cmd Msg )
init =
    let
        model =
            { tournaments = RemoteData.Loading
            , dateRange =
                { beginDate = "2022-01-01T00:00:00Z"
                , endDate = "2023-01-01T00:00:00Z"
                }
            , deleteError = Nothing
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
        [ button [ onClick (FetchTournaments model.dateRange) ]
            [ text "Refresh tournaments" ]
        , viewTournaments model.tournaments
        , viewDeleteError model.deleteError
        ]


viewTournaments : WebData (List Tournament) -> Html Msg
viewTournaments tournaments =
    case tournaments of
        RemoteData.NotAsked ->
            text ""

        RemoteData.Loading ->
            h3 [] [ text "Loading..." ]

        RemoteData.Success actualTournaments ->
            div []
                [ h3 [] [ text "Tournaments" ]
                , table []
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
            [ button [ type_ "button", onClick (DeleteTournament tournament.id) ]
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


buildErrorMessage : Http.Error -> String
buildErrorMessage httpError =
    case httpError of
        Http.BadUrl message ->
            message

        Http.Timeout ->
            "Server is taking too long to respond. Please try again later."

        Http.NetworkError ->
            "Unable to reach server."

        Http.BadStatus statusCode ->
            "Request failed with status code: " ++ String.fromInt statusCode

        Http.BadBody message ->
            message
