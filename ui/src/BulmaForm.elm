module BulmaForm exposing (..)

import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onClick, onInput)


type alias TextInputParams msg =
    { label : String
    , placeholder : String
    , type_ : String
    , oninput : String -> msg
    }


textInput : TextInputParams msg -> Html msg
textInput params =
    div [ class "field" ]
        [ label [ class "label" ] [ text params.label ]
        , div [ class "control" ]
            [ input
                [ class "input"
                , type_ params.type_
                , placeholder params.placeholder
                , onInput params.oninput
                ]
                []
            ]
        ]


type alias Option =
    { value : String
    , display : String
    }


type alias SelectInputParams msg =
    { label : String
    , onchange : String -> msg
    , options : List Option
    }


optionfn : Option -> Html msg
optionfn opt =
    option [ value opt.value ] [ text opt.display ]


selectInput : SelectInputParams msg -> Html msg
selectInput params =
    div [ class "field" ]
        [ label [ class "label" ] [ text params.label ]
        , div [ class "control" ]
            [ div [ class "select " ]
                [ select [ onInput params.onchange ]
                    -- on "change" ?
                    (List.map optionfn params.options)

                -- [ option [] [ text "foo" ] ]
                ]
            ]
        ]


type alias ButtonParams msg =
    { label : String
    , onclick : msg
    }


buttonInput : ButtonParams msg -> Html msg
buttonInput params =
    div [ class "field" ]
        [ div [ class "control" ]
            [ button [ class "button is-link", onClick params.onclick, type_ "button" ]
                [ text params.label ]
            ]
        ]
