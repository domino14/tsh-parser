module Main exposing (..)

import Browser exposing (Document, UrlRequest)
import Browser.Navigation as Nav
import Html exposing (..)
import ListTournaments
import Route exposing (Route)
import Url exposing (Url)


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
    }


type Page
    = NotFoundPage
    | TournamentListPage ListTournaments.Model
    | StandingsPage


type Msg
    = ListPageMsg ListTournaments.Msg
    | LinkClicked UrlRequest
    | UrlChanged Url


init : () -> Url -> Nav.Key -> ( Model, Cmd Msg )
init flags url navKey =
    let
        model =
            { route = Route.parseUrl url
            , page = NotFoundPage
            , navKey = navKey
            }
    in
    initCurrentPage ( model, Cmd.none )


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
                    ( NotFoundPage, Cmd.none )
    in
    ( { model | page = currentPage }
    , Cmd.batch [ existingCmds, mappedPageCmds ]
    )


view : Model -> Document Msg
view model =
    { title = "MGI Manager"
    , body = [ currentView model ]
    }


currentView : Model -> Html Msg
currentView model =
    case model.page of
        NotFoundPage ->
            notFoundView

        TournamentListPage pageModel ->
            ListTournaments.view pageModel
                |> Html.map ListPageMsg

        StandingsPage ->
            notFoundView


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

        ( _, _ ) ->
            ( model, Cmd.none )
