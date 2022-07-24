module Errors exposing (..)

import Http
import Http.Detailed
import Json.Decode exposing (Decoder, field, string)


type alias TwirpError =
    { code : String
    , msg : String
    }


jsonErrorDecoder : Decoder TwirpError
jsonErrorDecoder =
    Json.Decode.map2 TwirpError
        (field "code" string)
        (field "msg" string)


buildErrorMessage : Http.Detailed.Error String -> String
buildErrorMessage httpError =
    case httpError of
        Http.Detailed.BadStatus _ body ->
            let
                res =
                    Json.Decode.decodeString jsonErrorDecoder body
            in
            case res of
                Ok twirpError ->
                    twirpError.msg

                Err _ ->
                    body

        _ ->
            "Other http error"



-- case httpError of
--     Http.BadUrl message ->
--         message
--     Http.Timeout ->
--         "Server is taking too long to respond. Please try again later."
--     Http.NetworkError ->
--         "Unable to reach server."
--     Http.BadStatus statusCode ->
--         "Request failed with status code: " ++ String.fromInt statusCode
--     Http.BadBody message ->
--         message
-- Debug.toString httpError
