module Session exposing (..)

import Http


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
