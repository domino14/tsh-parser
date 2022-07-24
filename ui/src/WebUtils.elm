module WebUtils exposing (..)

import Http
import Http.Detailed
import Json.Decode exposing (Decoder)
import RemoteData exposing (RemoteData)


type alias DetailedWebData a =
    RemoteData (Http.Detailed.Error String) a


twirpReq : String -> String -> Http.Expect msg -> Http.Body -> Cmd msg
twirpReq service method expect body =
    Http.post
        { url = "/twirp/tshparser." ++ service ++ "/" ++ method
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


buildTextExpect : (DetailedWebData String -> msg) -> Http.Expect msg
buildTextExpect msg =
    Http.Detailed.expectString
        (\result ->
            (case result of
                Err err ->
                    Err err

                Ok ( _, str ) ->
                    Ok str
            )
                |> RemoteData.fromResult
                |> msg
        )
