package main

import (
  "fmt"
  "time"
  "net/http"
  "github.com/gorilla/websocket"
)

func main(){
  log( "Started server")

  upgrade := websocket.Upgrader{
    ReadBufferSize: 1024,
    WriteBufferSize: 1024,
  }

  http.HandleFunc("/websocket", func( res http.ResponseWriter, req *http.Request){
    handleConnection( res, req, upgrade )
  })

  http.Handle("/", http.FileServer(http.Dir("public/")))

  err := http.ListenAndServe(":8080", nil)
  if err != nil {
    panic( err )
  }
}

func handleConnection( res http.ResponseWriter, req *http.Request, upgrade websocket.Upgrader ){
  conn, err := upgrade.Upgrade(res, req, nil)
  if err != nil {
    panic( err )
  }

  for{
    messageType, message, err := conn.ReadMessage()
    if err != nil {
      panic( err )
    }

    messageStr := string(message[:])
    log( messageStr )

    if( messageStr == "ping" || messageStr == "PING" ){
      conn.WriteMessage( messageType, []byte("PONG") )
    }

    // Send the message to your business logic
  }
}

// Should probably shove this into a library
func log( message string ){
  fmt.Printf( "[%2d:%2d] : %s \n", time.Now().Hour(), time.Now().Minute(), message )
}
