package main

import (
    "bufio"
    "fmt"
    "os"
    "time"
    "net/http"
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
func replayWebsocketData(url string, inputFile string, port int) error {
    upgrader := websocket.Upgrader{}
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        conn, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
            fmt.Println("Error upgrading to WebSocket:", err)
            return
        }
        defer conn.Close()

        f, err := os.Open(inputFile)
        if err != nil {
            fmt.Println("Error opening input file:", err)
            return
        }
        defer f.Close()

        scanner := bufio.NewScanner(f)
        for scanner.Scan() {
            message := scanner.Text()
            err = conn.WriteMessage(websocket.TextMessage, []byte(message))
            if err != nil {
                fmt.Println("Error sending WebSocket message:", err)
                return
            }
            time.Sleep(100 * time.Millisecond) // delay between messages
        }

        if err := scanner.Err(); err != nil {
            fmt.Println("Error reading input file:", err)
            return
        }
    })

    err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
    if err != nil {
        return err
    }

    return nil
}

func main() {
    if len(os.Args) < 4 {
        fmt.Println("Usage: go run main.go [save|replay] [url] [filename] [port]")
        os.Exit(1)
    }

    command := os.Args[1]
    url := os.Args[2]
    filename := os.Args[3]

    if command == "save" {
        err := saveWebsocketData(url, filename)
        if err != nil {
            fmt.Println("Error:", err)
            os.Exit(1)
        }
    } else if command == "replay" {
        if len(os.Args) < 5 {
            fmt.Println("Usage: go run main.go replay [url] [filename] [port]")
            os.Exit(1)
        }
        port := os.Args[4]
        err := replayWebsocketData(url, filename, port)
        if err != nil {
            fmt.Println("Error:", err)
            os.Exit(1)
        }
    } else {
        fmt.Println("Invalid command:", command)
        os.Exit(1)
    }
}
