module Login exposing (..)

import Browser.Navigation as Nav
import BulmaForm exposing (buttonInput, textInput)
import Errors exposing (buildErrorMessage)
import Html exposing (..)
import Html.Attributes exposing (..)
import Http
import Json.Decode as Decode exposing (Decoder, field, string)
import Json.Encode as Encode
import RemoteData
import Route exposing (Route(..))
import WebUtils exposing (buildExpect)


type alias LoginRequest =
    { email : String
    , password : String
    }


reqEncoder : LoginRequest -> Encode.Value
reqEncoder req =
    Encode.object
        [ ( "email", Encode.string req.email )
        , ( "password", Encode.string req.password )
        ]


type alias LoginResponse =
    { token : String }


loginResponseDecoder : Decoder LoginResponse
loginResponseDecoder =
    Decode.map LoginResponse
        (field "token" string)


type alias Model =
    { loginRequest : LoginRequest
    , loginError : Maybe String
    , navKey : Nav.Key
    }


type Msg
    = StoreEmail String
    | StorePassword String
    | Submit
    | LoggedIn (WebUtils.DetailedWebData LoginResponse)


init : Nav.Key -> ( Model, Cmd Msg )
init navKey =
    ( initialModel navKey, Cmd.none )


initialModel : Nav.Key -> Model
initialModel navKey =
    { navKey = navKey
    , loginRequest = LoginRequest "" ""
    , loginError = Nothing
    }


view : Model -> Html Msg
view model =
    div []
        [ h3 [] [ text "Log in" ]
        , loginForm
        , viewError model.loginError
        ]


loginForm : Html Msg
loginForm =
    Html.form []
        [ textInput
            { label = "Email"
            , placeholder = "my@email.com"
            , type_ = "text"
            , oninput = StoreEmail
            }
        , textInput
            { label = "Password"
            , placeholder = ""
            , type_ = "password"
            , oninput = StorePassword
            }
        , buttonInput
            { label = "Submit"
            , onclick = Submit
            }
        ]


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        StoreEmail email ->
            let
                req =
                    model.loginRequest

                updateEmail =
                    { req | email = email }
            in
            ( { model | loginRequest = updateEmail }, Cmd.none )

        StorePassword password ->
            let
                req =
                    model.loginRequest

                updatePassword =
                    { req | password = password }
            in
            ( { model | loginRequest = updatePassword }, Cmd.none )

        Submit ->
            ( model, requestJWT model.loginRequest )

        LoggedIn loginResponse ->
            case loginResponse of
                RemoteData.Success _ ->
                    ( { model | loginError = Nothing }
                    , Route.pushUrl Route.Tournaments model.navKey
                    )

                RemoteData.Failure detailedError ->
                    ( { model | loginError = Just (buildErrorMessage detailedError) }
                    , Cmd.none
                    )

                _ ->
                    ( model, Cmd.none )


requestJWT : LoginRequest -> Cmd Msg
requestJWT req =
    Http.post
        { url = "/twirp/tshparser.AuthenticationService/GetJWT"
        , body = Http.jsonBody (reqEncoder req)
        , expect = buildExpect loginResponseDecoder LoggedIn
        }


viewError : Maybe String -> Html msg
viewError maybeError =
    case maybeError of
        Just error ->
            div []
                [ h3 [] [ text "Couldn't log in at this time." ]
                , text ("Error: " ++ error)
                ]

        Nothing ->
            text ""
