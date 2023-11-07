package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func updateData(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	counter := 0
	for {
		currentTime := time.Now().Format("2006-01-02 15:04:05")
		counter++
		data := map[string]interface{}{"currentTime": currentTime, "counter": counter}

		if err := conn.WriteJSON(data); err != nil {
			fmt.Println(err)
			return
		}

		time.Sleep(5 * time.Second)
	}
}

func main() {
	http.HandleFunc("/update", updateData)
	http.Handle("/", http.FileServer(http.Dir("/app/static/tasks/")))

	fmt.Println("Server is running on :8080")
	http.ListenAndServe(":8080", nil)
}
