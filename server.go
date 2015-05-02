package main

import (
  "fmt"
  "flag"
  "time"
  "net/http"
  "os/exec"
  "github.com/gorilla/websocket"
)

func main(){
  port := flag.String("p", "3232", "specifies the port you want to be hosting on")

  flag.Parse()

  log( "Started server on port " + *port )

  upgrade := websocket.Upgrader{
    ReadBufferSize: 1024,
    WriteBufferSize: 1024,
  }

  http.HandleFunc("/websocket", func( res http.ResponseWriter, req *http.Request){
   handleConnection( res, req, upgrade )
  })

  http.Handle("/", http.FileServer(http.Dir("public/")))

  err := http.ListenAndServe(":" + *port, nil)
  if err != nil {
    panic( err )
  }
}

func handleConnection( res http.ResponseWriter, req *http.Request, upgrade websocket.Upgrader ){
  log("Started a stream!")

  conn, err := upgrade.Upgrade(res, req, nil)
  if err != nil {
    panic( err )
  }

  for{
    imageData, err := exec.Command("fswebcam", "-").Output()
    if err != nil{
      panic(err)
    }

    // Throttle this shit
    time.Sleep( time.Millisecond * 200 )

    conn.WriteMessage( websocket.BinaryMessage, imageData )
  }

}

func log( message string ){
  fmt.Printf( "[%2d:%2d] : %s \n", time.Now().Hour(), time.Now().Minute(), message )
}
