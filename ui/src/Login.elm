module Login exposing (..)

import Browser.Navigation as Nav
import BulmaForm exposing (buttonInput, textInput)
import Errors exposing (buildErrorMessage)
import Html exposing (..)
import Html.Attributes exposing (..)
import Http
import Http.Detailed
import Json.Decode as Decode exposing (Decoder, field, string)
import Json.Encode as Encode
import Route exposing (Route(..))
import Session exposing (buildExpect)


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
    , token : String
    }


type Msg
    = StoreEmail String
    | StorePassword String
    | Submit
    | LoggedIn (Result (Http.Detailed.Error String) ( Http.Metadata, String ))


init : Nav.Key -> ( Model, Cmd Msg )
init navKey =
    ( initialModel navKey, Cmd.none )


initialModel : Nav.Key -> Model
initialModel navKey =
    { navKey = navKey
    , loginRequest = LoginRequest "" ""
    , loginError = Nothing
    , token = ""
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

        LoggedIn (Ok loginResponse) ->
            ( { model | loginError = Nothing, token = loginResponse.token }
            , Route.pushUrl Route.Tournaments model.navKey
            )

        LoggedIn (Err error) ->
            ( { model | loginError = Just (buildErrorMessage error) }
            , Cmd.none
            )


requestJWT : LoginRequest -> Cmd Msg
requestJWT req =
    Http.post
        { url = "http://localhost:8082/twirp/tshparser.AuthenticationService/GetJWT"
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
