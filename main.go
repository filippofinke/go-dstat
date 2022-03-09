package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
var clients []*websocket.Conn = []*websocket.Conn{}

func ws(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	clients = append(clients, c)
}

func main() {

	ticker := time.NewTicker(time.Second)
	var count int64 = 0

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	http.HandleFunc("/ws", ws)
	http.HandleFunc("/dstat", func(rw http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&count, 1)
		rw.WriteHeader(http.StatusOK)
	})

	go func() {
		for {
			<-ticker.C
			for i, c := range clients {
				err := c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprint(count)))
				if err != nil {
					c.Close()
					clients = append(clients[:i], clients[i+1:]...)
				}
			}
			count = 0
		}
	}()

	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}
