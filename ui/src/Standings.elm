module Standings exposing (..)

import BulmaForm exposing (buttonInput, textInput)
import DateRange exposing (DateRange, dateRangeEncoder)
import Errors exposing (buildErrorMessage)
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onClick)
import Http
import Json.Encode as Encode
import RemoteData
import SingleStanding exposing (Standing, standingsResponseDecoder)
import WebUtils exposing (buildExpect, buildTextExpect, twirpReq)


type alias Model =
    { standings : WebUtils.DetailedWebData (List Standing)
    , dateRange : DateRange
    , modalVisible : Bool
    , potentialAlias : String -- the alias of a player we're potentially editing
    , potentialRealName : String
    , aliasCreationError : Maybe String
    }


type Msg
    = FetchStandings DateRange
    | StandingsReceived (WebUtils.DetailedWebData (List Standing))
    | OpenAliasModal String
    | CloseAliasModal
    | StoreRealName String
    | SubmitAlias
    | AliasCreated (WebUtils.DetailedWebData String)


init : ( Model, Cmd Msg )
init =
    let
        model =
            { standings = RemoteData.Loading
            , dateRange =
                { beginDate = "2022-01-01T00:00:00Z"
                , endDate = "2023-01-01T00:00:00Z"
                }
            , modalVisible = False
            , potentialAlias = ""
            , potentialRealName = ""
            , aliasCreationError = Nothing
            }
    in
    ( model, fetchStandings model.dateRange )


fetchStandings : DateRange -> Cmd Msg
fetchStandings dateRange =
    twirpReq
        "TournamentRankerService"
        "ComputeStandings"
        (buildExpect standingsResponseDecoder StandingsReceived)
        (Http.jsonBody (dateRangeEncoder dateRange))


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        FetchStandings dateRange ->
            ( { model | standings = RemoteData.Loading }, fetchStandings dateRange )

        StandingsReceived response ->
            ( { model | standings = response }, Cmd.none )

        OpenAliasModal playerName ->
            ( { model | modalVisible = True, potentialAlias = playerName }, Cmd.none )

        CloseAliasModal ->
            ( { model | modalVisible = False, potentialAlias = "" }, Cmd.none )

        StoreRealName realName ->
            ( { model | potentialRealName = realName }, Cmd.none )

        SubmitAlias ->
            ( model, createAlias model.potentialAlias model.potentialRealName )

        AliasCreated resp ->
            case resp of
                RemoteData.Success _ ->
                    ( { model | modalVisible = False }, fetchStandings model.dateRange )

                RemoteData.Failure detailedError ->
                    ( { model | aliasCreationError = Just (buildErrorMessage detailedError) }, Cmd.none )

                _ ->
                    ( model, Cmd.none )



-- VIEWS


view : Model -> Html Msg
view model =
    div []
        [ viewStandings model.standings
        , viewAliasModal model.modalVisible model.potentialAlias model.aliasCreationError
        ]


viewStandings : WebUtils.DetailedWebData (List Standing) -> Html Msg
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
        , td []
            [ span
                [ class "has-text-link is-clickable"
                , onClick (OpenAliasModal standing.playerName)
                ]
                [ text standing.playerName ]
            ]
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


viewAliasModal : Bool -> String -> Maybe String -> Html Msg
viewAliasModal visible potentialAlias aliasCreationError =
    let
        visibleModifier =
            if visible then
                "is-active"

            else
                ""

        modalClass =
            "modal " ++ visibleModifier
    in
    div [ class modalClass ]
        [ div [ class "modal-background", onClick CloseAliasModal ] []
        , div [ class "modal-card" ]
            [ header
                [ class "modal-card-head" ]
                [ p [ class "modal-card-title" ] [ text ("Editing alias: " ++ potentialAlias) ]
                , button [ class "delete", onClick CloseAliasModal ] []
                ]
            , section [ class "modal-card-body" ]
                [ p
                    []
                    [ text
                        ("Please type in a player name that you want "
                            ++ potentialAlias
                            ++ " to be an alias of. "
                        )
                    ]
                , p [] [ text "Please make sure to type it in **exactly**, with commas as needed." ]
                , p [] [ text "Note that if this player name does not exist in the standings, this alias won't work as expected. Both the alias and the player name must already be in the standings." ]
                , br [] []
                , br [] []
                , textInput
                    { label = "Real name"
                    , placeholder = "Richards, Nigel"
                    , type_ = "text"
                    , oninput = StoreRealName
                    }
                , buttonInput { label = "Submit", onclick = SubmitAlias }
                , br [] []
                , viewAliasError aliasCreationError
                ]
            , footer [ class "model-card-foot" ]
                []
            ]
        ]


viewAliasError : Maybe String -> Html Msg
viewAliasError maybeError =
    case maybeError of
        Just error ->
            div []
                [ h3 [] [ text "Couldn't create alias at this time." ]
                , text ("Error: " ++ error)
                ]

        Nothing ->
            text ""


type alias AliasRequest =
    { realname : String
    , alias : String
    }


aliasReqEncoder : AliasRequest -> Encode.Value
aliasReqEncoder req =
    Encode.object
        [ ( "original_player", Encode.string req.realname )
        , ( "alias", Encode.string req.alias )
        ]


createAlias : String -> String -> Cmd Msg
createAlias alias_ realname =
    twirpReq
        "TournamentRankerService"
        "AddPlayerAlias"
        (buildTextExpect AliasCreated)
        (Http.jsonBody (aliasReqEncoder (AliasRequest realname alias_)))



-- Http.post
--     { url = "http://localhost:8082/twirp/tshparser.TournamentRankerService/ComputeStandings"
--     , expect =
--         standingsResponseDecoder
--             |> Http.expectJson (RemoteData.fromResult >> StandingsReceived)
--     , body = Http.jsonBody (dateRangeEncoder dateRange)
--     }
