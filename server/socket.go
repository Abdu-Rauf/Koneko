package main

import(
	"net/http"
	"github.com/gorilla/websocket"
)

func SignalSocket(w http.ResponseWriter,r *http.Request)(conn *websocket.Conn,err error){
	conn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil,err
	}
	return conn,nil
}