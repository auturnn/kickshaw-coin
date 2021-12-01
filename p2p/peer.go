package p2p

import (
	"fmt"

	"github.com/gorilla/websocket"
)

var Peers map[string]*peer = make(map[string]*peer)

type peer struct {
	conn *websocket.Conn
}

func initPeer(conn *websocket.Conn, addr, port string) *peer {
	p := &peer{
		conn,
	}

	key := fmt.Sprintf("%s:%s", addr, port)
	Peers[key] = p

	return p
}
