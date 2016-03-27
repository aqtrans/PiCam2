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

  imageData := make(chan []byte)
  go camThread(imageData)

  log( "Started server on port " + *port )

  upgrade := websocket.Upgrader{
    ReadBufferSize: 1024,
    WriteBufferSize: 1024,
  }

  http.HandleFunc("/websocket", func( res http.ResponseWriter, req *http.Request){
   handleConnection( res, req, upgrade, imageData )
  })

  http.Handle("/", http.FileServer(http.Dir("public/")))

  err := http.ListenAndServe(":" + *port, nil)
  if err != nil {
    panic( err )
  }
}

func handleConnection( res http.ResponseWriter, req *http.Request, upgrade websocket.Upgrader, imageData chan []byte ){
  log("Started a stream!")

  conn, err := upgrade.Upgrade(res, req, nil)
  if err != nil {
    panic( err )
  }


  for{
    conn.WriteMessage( websocket.BinaryMessage, <-imageData )
  }
}

func camThread(imageData chan []byte ){
  for{
    currentFrame, err := exec.Command("/opt/vc/bin/raspistill", "-n", "-vf", "-hf", "-ex", "night", "-w", "1024", "-h", "768", "-q", "50", "-t", "1000", "-o", "-").Output()
    if err != nil{
      panic( err )
    }

    imageData<-currentFrame
    time.Sleep( time.Millisecond * 1000 )
  }
}

func log( message string ){
  fmt.Printf( "[%2d:%2d] : %s \n", time.Now().Hour(), time.Now().Minute(), message )
}
