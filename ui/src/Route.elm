module Route exposing (Route(..), parseUrl, pushUrl)

import Browser.Navigation as Nav
import Url exposing (Url)
import Url.Parser exposing (..)


type Route
    = NotFound
    | Tournaments
    | Standings
    | NewTournament
    | Login
    | Aliases


parseUrl : Url -> Route
parseUrl url =
    case parse matchRoute url of
        Just route ->
            route

        Nothing ->
            NotFound


matchRoute : Parser (Route -> a) a
matchRoute =
    oneOf
        [ map Tournaments top
        , map Tournaments (s "tournaments")
        , map Standings (s "standings")
        , map NewTournament (s "tournaments" </> s "new")
        , map Login (s "login")
        , map Aliases (s "aliases")
        ]


pushUrl : Route -> Nav.Key -> Cmd msg
pushUrl route navKey =
    routeToString route
        |> Nav.pushUrl navKey


routeToString : Route -> String
routeToString route =
    case route of
        NotFound ->
            "/not-found"

        Tournaments ->
            "/tournaments"

        Standings ->
            "/standings"

        NewTournament ->
            "/tournaments/new"

        Login ->
            "/login"

        Aliases ->
            "/aliases"
