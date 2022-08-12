package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"gocr/server"
	"gocr/utils"
)

func main() {
	var port string
	flag.StringVar(&port, "p", "2333", "server port")
	flag.StringVar(&port, "port", "2333", "server port")

	flag.Parse()

	host := utils.GetLocalIP()

	http.HandleFunc("/test", server.Test)
	http.HandleFunc("/test/", server.Test)

	http.HandleFunc("/file", server.File)
	http.HandleFunc("/file/", server.File)

	http.HandleFunc("/base64", server.Base64)
	http.HandleFunc("/base64/", server.Base64)

	log.Println(fmt.Sprintf("serving at: <0.0.0.0:%s>[%s]", port, host))

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

}
