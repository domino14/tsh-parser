module Main exposing (..)

import Aliases
import Browser exposing (Document, UrlRequest)
import Browser.Navigation as Nav
import Html exposing (..)
import Html.Attributes exposing (class, href)
import Http exposing (stringBody)
import Json.Decode as Decode exposing (Decoder, field, string)
import Json.Decode.Pipeline exposing (required)
import Jwt
import ListTournaments
import Login
import NewTournament
import RemoteData exposing (RemoteData)
import Route exposing (Route(..))
import SingleStanding exposing (Standing)
import Standings
import Url exposing (Url)
import WebUtils exposing (DetailedWebData, buildExpect, twirpReq)


main : Program () Model Msg
main =
    Browser.application
        { init = init
        , view = view
        , update = update
        , subscriptions = \_ -> Sub.none
        , onUrlRequest = LinkClicked
        , onUrlChange = UrlChanged
        }


type alias Model =
    { route : Route
    , page : Page
    , navKey : Nav.Key
    , myuser : DetailedWebData User
    }


type alias User =
    { email : String }


type Page
    = NotFoundPage
    | TournamentListPage ListTournaments.Model
    | StandingsPage Standings.Model
    | NewTournamentPage NewTournament.Model
    | LoginPage Login.Model
    | AliasesPage Aliases.Model


type Msg
    = ListPageMsg ListTournaments.Msg
    | LinkClicked UrlRequest
    | UrlChanged Url
    | NewTournamentPageMsg NewTournament.Msg
    | StandingsPageMsg Standings.Msg
    | LoginPageMsg Login.Msg
    | AliasesPageMsg Aliases.Msg
    | WhoAmIReceived (DetailedWebData User)


init : () -> Url -> Nav.Key -> ( Model, Cmd Msg )
init flags url navKey =
    let
        model =
            { route = Route.parseUrl url
            , page = NotFoundPage
            , navKey = navKey
            , myuser = RemoteData.Loading
            }
    in
    initCurrentPage ( model, fetchWhoAmI )


whoAmIResponseDecoder : Decoder User
whoAmIResponseDecoder =
    Decode.succeed User
        |> required "email" string


fetchWhoAmI : Cmd Msg
fetchWhoAmI =
    twirpReq "AuthenticationService"
        "WhoAmI"
        (buildExpect whoAmIResponseDecoder WhoAmIReceived)
        (stringBody "application/json" "{}")


initCurrentPage : ( Model, Cmd Msg ) -> ( Model, Cmd Msg )
initCurrentPage ( model, existingCmds ) =
    let
        ( currentPage, mappedPageCmds ) =
            case model.route of
                Route.NotFound ->
                    ( NotFoundPage, Cmd.none )

                Route.Tournaments ->
                    let
                        ( pageModel, pageCmds ) =
                            ListTournaments.init
                    in
                    ( TournamentListPage pageModel, Cmd.map ListPageMsg pageCmds )

                Route.Standings ->
                    let
                        ( pageModel, pageCmds ) =
                            Standings.init
                    in
                    ( StandingsPage pageModel, Cmd.map StandingsPageMsg pageCmds )

                Route.NewTournament ->
                    let
                        ( pageModel, pageCmd ) =
                            NewTournament.init model.navKey
                    in
                    ( NewTournamentPage pageModel, Cmd.map NewTournamentPageMsg pageCmd )

                Route.Login ->
                    let
                        ( pageModel, pageCmd ) =
                            Login.init model.navKey
                    in
                    ( LoginPage pageModel, Cmd.map LoginPageMsg pageCmd )

                Route.Aliases ->
                    let
                        ( pageModel, pageCmd ) =
                            Aliases.init
                    in
                    ( AliasesPage pageModel, Cmd.map AliasesPageMsg pageCmd )
    in
    ( { model | page = currentPage }
    , Cmd.batch [ existingCmds, mappedPageCmds ]
    )


view : Model -> Document Msg
view model =
    { title = "MGI Manager"
    , body =
        [ div [ class "container" ]
            [ loggedInBar model
            , currentView model
            ]
        ]
    }


jwtDecoder : Decode.Decoder String
jwtDecoder =
    Decode.field "sub" Decode.string


emailFromJwt : String -> Result Jwt.JwtError String
emailFromJwt jwt =
    Jwt.decodeToken jwtDecoder jwt


userOrLogin : DetailedWebData User -> Html Msg
userOrLogin user =
    case user of
        RemoteData.NotAsked ->
            text ""

        RemoteData.Loading ->
            text "Loading..."

        RemoteData.Success actualuser ->
            span [] [ text ("Logged in as " ++ actualuser.email) ]

        RemoteData.Failure _ ->
            a [ href "/login" ] [ text "Log in" ]


loggedInBar : Model -> Html Msg
loggedInBar model =
    nav [ class "navbar" ]
        [ div [ class "navbar-brand" ]
            [ text "MGI Management Portal" ]
        , div [ class "navbar-end" ]
            [ div [ class "navbar-item" ]
                [ userOrLogin model.myuser ]
            ]
        ]


currentView : Model -> Html Msg
currentView model =
    case model.page of
        NotFoundPage ->
            notFoundView

        TournamentListPage pageModel ->
            ListTournaments.view pageModel
                |> Html.map ListPageMsg

        StandingsPage pageModel ->
            Standings.view pageModel
                |> Html.map StandingsPageMsg

        NewTournamentPage pageModel ->
            NewTournament.view pageModel
                |> Html.map NewTournamentPageMsg

        LoginPage pageModel ->
            Login.view pageModel
                |> Html.map LoginPageMsg

        AliasesPage pageModel ->
            Aliases.view pageModel
                |> Html.map AliasesPageMsg


notFoundView : Html msg
notFoundView =
    h3 [] [ text "Oops! The page you requested was not found!" ]


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case ( msg, model.page ) of
        ( ListPageMsg subMsg, TournamentListPage pageModel ) ->
            let
                ( updatedPageModel, updatedCmd ) =
                    ListTournaments.update subMsg pageModel
            in
            ( { model | page = TournamentListPage updatedPageModel }
            , Cmd.map ListPageMsg updatedCmd
            )

        ( StandingsPageMsg subMsg, StandingsPage pageModel ) ->
            let
                ( updatedPageModel, updatedCmd ) =
                    Standings.update subMsg pageModel
            in
            ( { model | page = StandingsPage updatedPageModel }
            , Cmd.map StandingsPageMsg updatedCmd
            )

        ( NewTournamentPageMsg subMsg, NewTournamentPage pageModel ) ->
            let
                ( updatedPageModel, updatedCmd ) =
                    NewTournament.update subMsg pageModel
            in
            ( { model | page = NewTournamentPage updatedPageModel }
            , Cmd.map NewTournamentPageMsg updatedCmd
            )

        ( LoginPageMsg subMsg, LoginPage pageModel ) ->
            let
                ( updatedPageModel, updatedCmd ) =
                    Login.update subMsg pageModel
            in
            ( { model
                | page = LoginPage updatedPageModel
              }
            , Cmd.map LoginPageMsg updatedCmd
            )

        ( AliasesPageMsg subMsg, AliasesPage pageModel ) ->
            let
                ( updatedPageModel, updatedCmd ) =
                    Aliases.update subMsg pageModel
            in
            ( { model
                | page = AliasesPage updatedPageModel
              }
            , Cmd.map AliasesPageMsg updatedCmd
            )

        ( LinkClicked urlRequest, _ ) ->
            case urlRequest of
                Browser.Internal url ->
                    ( model
                    , Nav.pushUrl model.navKey (Url.toString url)
                    )

                Browser.External url ->
                    ( model
                    , Nav.load url
                    )

        ( UrlChanged url, _ ) ->
            let
                newRoute =
                    Route.parseUrl url
            in
            ( { model | route = newRoute }, Cmd.none )
                |> initCurrentPage

        ( WhoAmIReceived response, _ ) ->
            ( { model | myuser = response }, Cmd.none )

        ( _, _ ) ->
            ( model, Cmd.none )
