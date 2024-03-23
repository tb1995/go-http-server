package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const webPage = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <title>Simple Web Page</title>
  </head>
  <body>
    <h1 style="color: red;">Test Web Page</h1>
    <p>My web server served this page!</p>
  </body>
</html>
`

func getRoot(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    fmt.Fprint(w, webPage)
    fmt.Printf("got / Request\n")
}

func main() {
    // load .env
    err:= godotenv.Load(".env")
    if err!=nil {
        log.Fatal("Error loading .env file: %s", err)
    }
    
    publicDirectory := os.Getenv("PUBLIC_DIRECTORY_PATH")
    fmt.Printf(publicDirectory)

    http.HandleFunc("/", getRoot)

    err = http.ListenAndServe(":3333", nil)

    if errors.Is(err, http.ErrServerClosed) {
        fmt.Printf("Server closed \n")
    } else if err!= nil{
        fmt.Printf("Error starting server: %s\n", err)
    }

}