module Tournament exposing (..)

import Json.Decode as Decode exposing (Decoder, field, list, string)
import Json.Decode.Pipeline exposing (required)


type alias Tournament =
    { id : String -- it's actually an int64
    , name : String
    , category : String -- or "type"
    , date : String
    }


tournamentsResponseDecoder : Decoder (List Tournament)
tournamentsResponseDecoder =
    field "tournaments" (list tournamentDecoder)


tournamentDecoder : Decoder Tournament
tournamentDecoder =
    Decode.succeed Tournament
        |> required "id" string
        |> required "name" string
        |> required "tournament_type" string
        |> required "date" string
