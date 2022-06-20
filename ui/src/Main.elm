module Main exposing (..)

import Browser
import Dict exposing (Dict)
import Html exposing (..)
import Html.Attributes exposing (style)
import Html.Events exposing (..)
import Http
import Json.Decode exposing (Decoder, field, float, int, list, map5, string)
import Json.Encode as Encode


main : Program () Model Msg
main =
    Browser.element
        { init = init
        , update = update
        , subscriptions = subscriptions
        , view = view
        }


type alias Model =
    { standings : List Standing
    , playerAliases : List PlayerAlias

    -- dates for filtering standings
    , beginDate : String
    , endDate : String
    }


type alias PlayerAlias =
    { origName : String
    , alias : String
    }


type alias StandingsRequest =
    { beginDate : String, endDate : String }


type alias Standing =
    { player_name : String
    , points : Int
    , games : Int
    , wins : Float
    , spread : Int
    , tournaments_played : Int
    , rank : Int
    }


type alias Standings =
    List Standing


init : () -> ( Model, Cmd Msg )
init _ =
    ( { standings = []
      }
    , getStandings { beginDate = "2022-01-01", endDate = "2022-04-01" }
    )



-- UPDATE


type Msg
    = GotStandings (Result Http.Error Standings)


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        GotStandings result ->
            case result of
                Ok standings ->
                    ( Success standings, Cmd.none )

                Err _ ->
                    ( Failure, Cmd.none )



-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none



-- VIEW


view : Model -> Html Msg
view model =
    div []
        [ h2 [] [ text "Random Quotes" ]
        , viewQuote model
        ]


viewQuote : Model -> Html Msg
viewQuote model =
    case model of
        Failure ->
            div []
                [ text "I could not load a random quote for some reason. "
                , button [ onClick MorePlease ] [ text "Try Again!" ]
                ]

        Loading ->
            text "Loading..."

        Success quote ->
            div []
                [ button [ onClick MorePlease, style "display" "block" ] [ text "More Please!" ]
                , blockquote [] [ text quote.quote ]
                , p [ style "text-align" "right" ]
                    [ text "— "
                    , cite [] [ text quote.source ]
                    , text (" by " ++ quote.author ++ " (" ++ String.fromInt quote.year ++ ")")
                    ]
                ]



-- HTTP
-- getRandomQuote : Cmd Msg
-- getRandomQuote =
--     Http.get
--         { url = "https://elm-lang.org/api/random-quotes"
--         , expect = Http.expectJson GotQuote quoteDecoder
--         }
-- quoteDecoder : Decoder Quote
-- quoteDecoder =
--     map4 Quote
--         (field "quote" string)
--         (field "source" string)
--         (field "author" string)
--         (field "year" int)


standingsDecoder : Decoder Standings
standingsDecoder =
    list
        (map5
            Standing
            (field "player_name" string)
            (field "points" int)
            (field "wins" float)
            (field "spread" int)
            (field "tournaments_played" int)
        )


getStandings : StandingsRequest -> Cmd Msg
getStandings req =
    Http.post
        { url = "localhost:8082/twirp/tshparser.TournamentRankerService/ComputeStandings"
        , body = Http.jsonBody (standingsReqEncoder req)
        , expect = Http.expectJson GotStandings standingsDecoder
        }


standingsReqEncoder : StandingsRequest -> Encode.Value
standingsReqEncoder req =
    Encode.object
        [ ( "method", Encode.string "standings" )
        , ( "params"
          , Encode.dict identity
                Encode.string
                -- Is this the best way to convert the req into a Dict?
                (Dict.fromList
                    [ ( "begin", req.beginDate )
                    , ( "end", req.endDate )
                    ]
                )
          )
        ]
