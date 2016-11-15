port module Main exposing (..)

import Html.App as App
import Html.Attributes as Attributes
import Html
import WebSocket
import Navigation
import Json.Decode exposing (..)


port coords : ( Float, Float ) -> Cmd msg


type alias Tweet =
    { text : String, location : ( Float, Float ) }


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
    Html.div
        [ Attributes.style
            [ ( "height", "400px" )
            , ( "overflow", "auto" )
            ]
        ]
        (List.concatMap viewTweet model.tweets)


viewTweet : Tweet -> List (Html.Html Msg)
viewTweet tweet =
    [ Html.text tweet.text, Html.br [] [] ]


subscriptions : Model -> Sub Msg
subscriptions model =
    let
        url =
            "ws://" ++ model.location.host ++ "/tweets"
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
