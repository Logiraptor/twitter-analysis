module Main exposing (..)

import Html.App as App
import Html


type alias Model =
    {}


type Msg
    = NoOp


main : Program Never
main =
    App.program
        { init = ( model, Cmd.none )
        , view = view
        , update = update
        , subscriptions = subscriptions
        }


model : Model
model =
    {}


view : Model -> Html.Html Msg
view model =
    Html.text "Hello from elm!"


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        NoOp ->
            ( model, Cmd.none )
