module Aliases exposing (..)

import Html exposing (..)
import Http exposing (stringBody)
import Http.Detailed
import Json.Decode as Decode exposing (Decoder, field, list, string)
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
            ( model, Cmd.none )

        -- fix
        AliasDeleted (Ok _) ->
            ( model, Cmd.none )

        AliasDeleted (Err error) ->
            ( model, Cmd.none )



-- VIEWS


view : Model -> Html Msg
view model =
    Html.text "foo"
