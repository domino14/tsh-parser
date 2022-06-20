module DateRange exposing (..)

import Json.Decode as Decode exposing (Decoder, string)
import Json.Decode.Pipeline exposing (required)
import Json.Encode as Encode


type alias DateRange =
    { beginDate : String
    , endDate : String
    }


dateRangeDecoder : Decoder DateRange
dateRangeDecoder =
    Decode.succeed DateRange
        |> required "begin_date" string
        |> required "end_date" string


dateRangeEncoder : DateRange -> Encode.Value
dateRangeEncoder req =
    Encode.object
        [ ( "begin_date", Encode.string req.beginDate )
        , ( "end_date", Encode.string req.endDate )
        ]
