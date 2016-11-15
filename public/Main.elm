port module Main exposing (..)

import Html.App as App
import Html.Attributes as Attributes
import Html.Events as Events
import Html
import WebSocket
import Navigation
import Json.Decode exposing (..)
import String


port coords : ( Float, Float ) -> Cmd msg


port resetView : () -> Cmd msg


port mapView : (( Float, Float, Float, Float ) -> msg) -> Sub msg


type alias Tweet =
    { text : String, location : ( Float, Float ) }


type alias Model =
    { location : Navigation.Location
    , tweets : List Tweet
    , paused : Bool
    , boundingBox : ( Float, Float, Float, Float )
    }


type Msg
    = NewLocation Navigation.Location
    | NewTweet Tweet
    | TogglePause
    | FilterToView
    | NewView ( Float, Float, Float, Float )


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
    ( { location = loc, tweets = [], paused = True, boundingBox = ( -142.29531250000002, 21.07961382717576, -54.40468750000002, 54.082101510457534 ) }, Cmd.none )


view : Model -> Html.Html Msg
view model =
    Html.div []
        [ Html.button [ Events.onClick TogglePause ]
            [ Html.text
                (if model.paused then
                    "Start"
                 else
                    "Stop"
                )
            ]
        , Html.button [ Events.onClick FilterToView ]
            [ Html.text "Filter tweets not in current view"
            ]
        , Html.div
            [ Attributes.style
                [ ( "height", "400px" )
                , ( "overflow", "auto" )
                ]
            ]
            (List.concatMap viewTweet model.tweets)
        ]


viewTweet : Tweet -> List (Html.Html Msg)
viewTweet tweet =
    [ Html.text tweet.text, Html.br [] [] ]


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.batch [ tweetStream model, mapView NewView ]


tweetStream : Model -> Sub Msg
tweetStream model =
    case model.paused of
        True ->
            Sub.none

        False ->
            let
                ( a, b, c, d ) =
                    model.boundingBox

                location =
                    [ a, b, c, d ]
                        |> List.map toString
                        |> String.join ","

                url =
                    "ws://" ++ model.location.host ++ "/tweets?locations=" ++ location
            in
                WebSocket.listen url
                    (decodeString decodeTweet
                        >> Result.withDefault { text = "Error", location = ( 0, 0 ) }
                        >> NewTweet
                    )


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        NewLocation location ->
            ( model, Cmd.none )

        NewTweet tweet ->
            let
                newTweets =
                    (tweet :: model.tweets)
                        |> List.take 100
            in
                ( { model | tweets = newTweets }, coords tweet.location )

        TogglePause ->
            ( { model | paused = not model.paused }, Cmd.none )

        FilterToView ->
            ( model, resetView () )

        NewView v ->
            ( { model | boundingBox = v, tweets = [] }, Cmd.none )


urlUpdate : Navigation.Location -> Model -> ( Model, Cmd Msg )
urlUpdate loc model =
    ( { model | location = loc }, Cmd.none )


decodeTweet : Decoder Tweet
decodeTweet =
    object2 Tweet
        ("text" := string)
        ("computed_coords" := decodeCoords)


decodeCoords : Decoder ( Float, Float )
decodeCoords =
    at [ "coordinates" ] (tuple2 (,) float float)
