module Tournament exposing (..)

import Json.Decode as Decode exposing (Decoder, int, list, string)
import Json.Decode.Pipeline exposing (required)


type alias Tournament =
    { id : Int
    , name : String
    , category : String -- or "type"
    , date : String
    }


tournamentsDecoder : Decoder (List Tournament)
tournamentsDecoder =
    list tournamentDecoder


tournamentDecoder : Decoder Tournament
tournamentDecoder =
    Decode.succeed Tournament
        |> required "id" int
        |> required "name" string
        |> required "tournament_type" string
        |> required "date" string
