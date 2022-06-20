module ListTournaments exposing (..)

import DateRange exposing (DateRange, dateRangeEncoder)
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onClick)
import Http
import RemoteData exposing (WebData)
import Tournament exposing (Tournament, tournamentsDecoder)


type alias Model =
    { tournaments : WebData (List Tournament)
    , dateRange : DateRange
    }


type Msg
    = FetchTournaments DateRange
    | TournamentsReceived (WebData (List Tournament))


init : () -> ( Model, Cmd Msg )
init _ =
    ( { tournaments = RemoteData.NotAsked
      , dateRange = { beginDate = "2022-01-01", endDate = "2023-01-01" }
      }
    , Cmd.none
    )


fetchTournaments : DateRange -> Cmd Msg
fetchTournaments dateRange =
    Http.post
        -- XXX: use some sort of env var on prod?
        { url = "http://localhost:8082/twirp/tshparser.TournamentRankerService/GetTournaments"
        , expect =
            tournamentsDecoder
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



-- VIEWS


view : Model -> Html Msg
view model =
    div []
        [ button [ onClick (FetchTournaments model.dateRange) ]
            [ text "Refresh tournaments" ]
        , viewTournaments model.tournaments
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
        ]


viewTournament : Tournament -> Html Msg
viewTournament tournament =
    tr []
        [ td [] [ text tournament.date ]
        , td [] [ text tournament.category ]
        , td [] [ text tournament.name ]
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
