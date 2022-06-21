module SingleStanding exposing (..)

import Json.Decode exposing (Decoder, field, float, int, list, string)
import Json.Decode.Pipeline exposing (required)


type alias Standing =
    { playerName : String
    , points : Int
    , wins : Float
    , games : Int
    , spread : Int
    , tournamentsPlayed : Int
    , rank : Int
    }


standingsResponseDecoder : Decoder (List Standing)
standingsResponseDecoder =
    field "standings" (list standingDecoder)


standingDecoder : Decoder Standing
standingDecoder =
    Json.Decode.succeed Standing
        |> required "player_name" string
        |> required "points" int
        |> required "wins" float
        |> required "games" int
        |> required "spread" int
        |> required "tournaments_played" int
        |> required "rank" int
