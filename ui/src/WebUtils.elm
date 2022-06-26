module WebUtils exposing (..)

import Http.Detailed
import RemoteData exposing (RemoteData)


type alias DetailedWebData a =
    RemoteData (Http.Detailed.Error String) a
