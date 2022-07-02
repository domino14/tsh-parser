module Session exposing (..)

import Http
import Http.Detailed
import Json.Decode exposing (Decoder)
import RemoteData exposing (RemoteData)
import WebUtils exposing (DetailedWebData)


type alias Session =
    { jwt : String
    }


twirpReq : String -> String -> Http.Expect msg -> Http.Body -> Cmd msg
twirpReq service method expect body =
    Http.post
        { url = "http://localhost:8082/twirp/tshparser." ++ service ++ "/" ++ method
        , body = body
        , expect = expect
        }


buildExpect : Decoder a -> (DetailedWebData a -> msg) -> Http.Expect msg
buildExpect decoder msg =
    decoder
        |> Http.Detailed.expectJson
            (\result ->
                (case result of
                    Err err ->
                        Err err

                    Ok ( _, a ) ->
                        Ok a
                )
                    |> RemoteData.fromResult
                    |> msg
            )
