module NewTournament exposing (..)

import Browser.Navigation as Nav
import BulmaForm exposing (Option, TextInputParams, buttonInput, selectInput, textInput)
import Errors exposing (buildErrorMessage)
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onClick, onInput)
import Http
import Json.Decode as Decode exposing (Decoder, field, map, string)
import Json.Encode as Encode
import ListTournaments exposing (Msg)
import RemoteData exposing (RemoteData)
import Route exposing (Route)
import WebUtils exposing (buildExpect, twirpReq)


type alias TournamentRequest =
    { date : String
    , name : String
    , category : String
    , tshURL : String
    }


reqEncoder : TournamentRequest -> Encode.Value
reqEncoder req =
    Encode.object
        [ ( "date", Encode.string req.date )
        , ( "tournament_type", Encode.string req.category )
        , ( "name", Encode.string req.name )
        , ( "tsh_url", Encode.string req.tshURL )
        ]


type alias Model =
    { navKey : Nav.Key
    , tournamentRequest : TournamentRequest
    , createError : Maybe String
    }


type Msg
    = StoreName String
    | StoreDate String
    | StoreTshURL String
    | StoreCategory String
    | Submit
    | TournamentCreated (WebUtils.DetailedWebData TournamentCreationResponse)


init : Nav.Key -> ( Model, Cmd Msg )
init navKey =
    ( initialModel navKey, Cmd.none )


initialModel : Nav.Key -> Model
initialModel navKey =
    { navKey = navKey
    , tournamentRequest = TournamentRequest "" "" "MGILeague-Div1" ""
    , createError = Nothing
    }


view : Model -> Html Msg
view model =
    div []
        [ h3 [ class "subtitle is-2" ] [ text "Add new tournament" ]
        , newTournamentForm
        , viewError model.createError
        ]


newTournamentForm : Html Msg
newTournamentForm =
    Html.form []
        [ textInput
            { label = "Name"
            , placeholder = "My tournament"
            , type_ = "text"
            , oninput = StoreName
            }
        , selectInput
            { label = "Category"
            , onchange = StoreCategory
            , options =
                [ Option "MGILeague-Div1" "MGI League Div 1"
                , Option "MGILeague-Div2" "MGI League Div 2"
                , Option "Masters" "Masters"
                , Option "Open" "Open"
                , Option "Intermediate" "Intermediate"
                ]
            }
        , textInput
            { label = "Date"
            , placeholder = ""
            , type_ = "date"
            , oninput = StoreDate
            }
        , textInput
            { label = "TSH Standings URL"
            , placeholder = "https://tsh.com/14-a.html"
            , type_ = "text"
            , oninput = StoreTshURL
            }
        , buttonInput
            { label = "Submit"
            , onclick = Submit
            }
        ]


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        StoreCategory category ->
            let
                req =
                    model.tournamentRequest

                updateCategory =
                    { req | category = category }
            in
            ( { model | tournamentRequest = updateCategory }, Cmd.none )

        StoreDate date ->
            let
                req =
                    model.tournamentRequest

                updateDate =
                    { req | date = date ++ "T00:00:00Z" }
            in
            ( { model | tournamentRequest = updateDate }, Cmd.none )

        StoreName name ->
            let
                oldTourney =
                    model.tournamentRequest

                updateName =
                    { oldTourney | name = name }
            in
            ( { model | tournamentRequest = updateName }, Cmd.none )

        StoreTshURL turl ->
            let
                oldTourney =
                    model.tournamentRequest

                updateUrl =
                    { oldTourney | tshURL = turl }
            in
            ( { model | tournamentRequest = updateUrl }, Cmd.none )

        Submit ->
            ( model, createTournament model.tournamentRequest )

        TournamentCreated resp ->
            case resp of
                RemoteData.Success _ ->
                    ( { model | createError = Nothing }
                    , Route.pushUrl Route.Tournaments model.navKey
                    )

                RemoteData.Failure detailedError ->
                    ( { model | createError = Just (buildErrorMessage detailedError) }
                    , Cmd.none
                    )

                _ ->
                    ( model, Cmd.none )


createTournament : TournamentRequest -> Cmd Msg
createTournament req =
    twirpReq
        "TournamentRankerService"
        "AddTournament"
        (buildExpect creationDecoder TournamentCreated)
        (Http.jsonBody (reqEncoder req))


type alias TournamentCreationResponse =
    { id : String }


creationDecoder : Decoder TournamentCreationResponse
creationDecoder =
    Decode.map TournamentCreationResponse
        (field "id" string)


viewError : Maybe String -> Html msg
viewError maybeError =
    case maybeError of
        Just error ->
            div []
                [ h3 [] [ text "Couldn't create a tournament at this time." ]
                , text ("Error: " ++ error)
                ]

        Nothing ->
            text ""
