module Session exposing (..)

import Http
import Http.Detailed
import Json.Decode exposing (Decoder)
import RemoteData exposing (RemoteData)
import WebUtils exposing (DetailedWebData)


type alias Session =
    { jwt : String
    }


twirpReq : Session -> String -> String -> Http.Expect msg -> Http.Body -> Cmd msg
twirpReq sess service method expect body =
    Http.request
        { method = "POST"
        , headers = [ Http.header "Authorization" ("Bearer " ++ sess.jwt) ]
        , url = "http://localhost:8082/twirp/tshparser." ++ service ++ "/" ++ method
        , body = body
        , expect = expect
        , timeout = Nothing
        , tracker = Nothing
        }


buildExpect : Decoder a -> (DetailedWebData a -> msg) -> Http.Expect msg
buildExpect decoder msg =
    decoder
        |> Http.Detailed.expectJson
            (\result ->
                (case result of
                    Err err ->
                        Err err

                    Ok ( metadata, a ) ->
                        Ok a
                )
                    |> RemoteData.fromResult
                    |> msg
            )
