module Main exposing (..)

import Html.App as App
import Html
import WebSocket
import Navigation


type alias Tweet =
    String


type alias Model =
    { location : Navigation.Location
    , tweets : List Tweet
    }


type Msg
    = NewLocation Navigation.Location
    | NewTweet Tweet


main : Program Never
main =
    Navigation.program
        (Navigation.makeParser identity)
        { init = model
        , view = view
        , update = update
        , urlUpdate = urlUpdate
        , subscriptions = subscriptions
        }


model : Navigation.Location -> ( Model, Cmd Msg )
model loc =
    ( { location = loc, tweets = [] }, Cmd.none )


view : Model -> Html.Html Msg
view model =
    Html.div [] (List.map viewTweet model.tweets)


viewTweet : Tweet -> Html.Html Msg
viewTweet tweet =
    Html.text tweet


subscriptions : Model -> Sub Msg
subscriptions model =
    let
        url =
            "ws://" ++ model.location.host ++ "/tweets"
    in
        WebSocket.listen url NewTweet


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        NewLocation location ->
            ( model, Cmd.none )

        NewTweet tweet ->
            let
                newTweets =
                    (tweet :: model.tweets)
                        |> List.take 10
            in
                ( { model | tweets = newTweets }, Cmd.none )


urlUpdate : Navigation.Location -> Model -> ( Model, Cmd Msg )
urlUpdate loc model =
    ( { model | location = loc }, Cmd.none )
