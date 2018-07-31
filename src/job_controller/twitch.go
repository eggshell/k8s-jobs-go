package job_controller

import (
    "fmt"
    "os"
    twitch "github.com/Onestay/go-new-twitch"
)

func UpdateStreamsList() []string {
    client := twitch.NewClient(os.Getenv("client_id"))
    input := twitch.GetStreamsInput{
        GameID: []string{"33214"},
        Language: []string{"en"},
    }

    resp, err := client.GetStreams(input)
    if err != nil {
        fmt.Println("An error occured while getting the streams: %v", err)
    }

    size := len(resp)
    user_ids := make([]string, size)

    for k := range resp {
        user_resp, err := client.GetUsersByID(resp[k].UserID)
        if err != nil {
            fmt.Println("An error occurred while getting the user: %v", err)
        } else {
            user_ids[k] = user_resp[0].DisplayName
        }
    }

    return user_ids
}
