package main

import (
    "fmt"
    twitch "github.com/Onestay/go-new-twitch"
)

func main() {
    client := twitch.NewClient("jp4kqyf9c6l5ll0dgghbbgqrf6daj8")
    input := twitch.GetStreamsInput{
        GameID: []string{"33214"},
        Language: []string{"en"},
    }

    resp, err := client.GetStreams(input)
    if err != nil {
        fmt.Println("An error occured while getting the streams: %v", err)
    }

    for k := range resp {
        user_resp, err := client.GetUsersByID(resp[k].UserID)
        if err != nil {
            fmt.Println("An error occurred while getting the user: %v", err)
        }
        fmt.Println(user_resp[0].DisplayName)
    }
}
