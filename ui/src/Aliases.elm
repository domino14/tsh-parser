module Aliases exposing (..)

import Errors exposing (buildErrorMessage)
import Html exposing (..)
import Html.Attributes exposing (class, type_)
import Html.Events exposing (onClick)
import Http exposing (stringBody)
import Http.Detailed
import Json.Decode as Decode exposing (Decoder, field, list, string)
import Json.Encode as Encode
import RemoteData exposing (RemoteData)
import WebUtils exposing (DetailedWebData, buildExpect, twirpReq)


type alias Model =
    { playerAliases : DetailedWebData (List PlayerAlias)
    , deleteError : Maybe String
    }


type alias PlayerAlias =
    { correctName : String
    , alias : String
    }


playerAliasDecoder : Decoder PlayerAlias
playerAliasDecoder =
    Decode.map2 PlayerAlias
        (field "original_player" string)
        (field "alias" string)


playerAliasResponseDecoder : Decoder (List PlayerAlias)
playerAliasResponseDecoder =
    field "aliases" (list playerAliasDecoder)


type Msg
    = AliasesReceived (DetailedWebData (List PlayerAlias))
    | DeleteAlias String
    | AliasDeleted (Result (Http.Detailed.Error String) ( Http.Metadata, String ))


init : ( Model, Cmd Msg )
init =
    let
        model =
            { playerAliases = RemoteData.Loading
            , deleteError = Nothing
            }
    in
    ( model, fetchAliases )


fetchAliases : Cmd Msg
fetchAliases =
    twirpReq "TournamentRankerService"
        "ListPlayerAliases"
        (buildExpect playerAliasResponseDecoder AliasesReceived)
        (stringBody "application/json" "{}")


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        AliasesReceived response ->
            ( { model | playerAliases = response }, Cmd.none )

        DeleteAlias theAlias ->
            ( model, deleteAlias theAlias )

        -- fix
        AliasDeleted (Ok _) ->
            ( model, fetchAliases )

        AliasDeleted (Err error) ->
            ( model, Cmd.none )


deleteAlias : String -> Cmd Msg
deleteAlias alias_ =
    twirpReq
        "TournamentRankerService"
        "RemovePlayerAlias"
        (Http.Detailed.expectString AliasDeleted)
        (Http.jsonBody (aliasEncoder alias_))


aliasEncoder : String -> Encode.Value
aliasEncoder alias_ =
    Encode.object [ ( "alias", Encode.string alias_ ) ]



-- VIEWS


view : Model -> Html Msg
view model =
    case model.playerAliases of
        RemoteData.NotAsked ->
            text ""

        RemoteData.Loading ->
            h3 [ class "subtitle is-2" ] [ text "Loading..." ]

        RemoteData.Success aliases ->
            div []
                [ h2 [ class "subtitle is-2" ] [ text "Aliases" ]
                , table [ class "table" ]
                    (viewTableHeader :: List.map viewAlias aliases)
                ]

        RemoteData.Failure httpError ->
            div []
                [ h3 [] [ text "Couldn't fetch aliases at this time." ]
                , text ("Error: " ++ buildErrorMessage httpError)
                ]


viewTableHeader : Html Msg
viewTableHeader =
    tr []
        [ th [] [ text "Real name" ]
        , th [] [ text "Alias" ]
        , th [] [ text "" ]
        ]


viewAlias : PlayerAlias -> Html Msg
viewAlias playerAlias =
    tr []
        [ td [] [ text playerAlias.correctName ]
        , td [] [ text playerAlias.alias ]
        , td []
            [ button [ class "button is-warning", type_ "button", onClick (DeleteAlias playerAlias.alias) ]
                [ text "Delete" ]
            ]
        ]
