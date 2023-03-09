package main

import (
    "bufio"
    "fmt"
    "os"
    "time"
    "github.com/gorilla/websocket"
)

// Save WebSocket data to file
func saveWebsocketData(url string, outputFile string) error {
    conn, _, err := websocket.DefaultDialer.Dial(url, nil)
    if err != nil {
        return err
    }
    defer conn.Close()

    f, err := os.Create(outputFile)
    if err != nil {
        return err
    }
    defer f.Close()

    for {
        _, message, err := conn.ReadMessage()
        if err != nil {
            return err
        }
        if _, err := f.WriteString(string(message) + "\n"); err != nil {
            return err
        }
    }

    return nil
}

// Replay WebSocket data from file
func replayWebsocketData(url string, inputFile string) error {
    conn, _, err := websocket.DefaultDialer.Dial(url, nil)
    if err != nil {
        return err
    }
    defer conn.Close()

    f, err := os.Open(inputFile)
    if err != nil {
        return err
    }
    defer f.Close()

    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        message := scanner.Text()
        err = conn.WriteMessage(websocket.TextMessage, []byte(message))
        if err != nil {
            return err
        }
        time.Sleep(100 * time.Millisecond) // delay between messages
    }

    return scanner.Err()
}


func main() {
     // Save WebSocket data to file
    saveWebsocketData("ws://echo.websocket.org/", "output.txt")

    // Replay WebSocket data from file
    err := replayWebsocketData("ws://echo.websocket.org/", "output.txt")
    if err != nil {
        fmt.Println("Error:", err)
    }
}
