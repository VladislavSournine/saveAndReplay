# saveAndReplay
Save websocket data to file, and replay from file to websocket

This code first calls the saveWebsocketData function to save WebSocket data to a file. Then it reads the data from the file using a scanner, sends each message to the WebSocket server using conn.WriteMessage, and adds a delay between messages using time.Sleep. Note that this code assumes that the messages in the input file are in the same format as the messages that the server expects. If the message format needs to be modified, the code should be adjusted accordingly.
