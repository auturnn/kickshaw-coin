package p2p

import (
	"fmt"

	"github.com/gorilla/websocket"
)

var Peers map[string]*peer = make(map[string]*peer)

type peer struct {
	key   string
	addr  string
	port  string
	conn  *websocket.Conn
	inbox chan []byte
}

func (p *peer) close() {
	p.conn.Close()
	delete(Peers, p.key)
}

func (p *peer) read() {
	defer p.close()
	//delete peer in case of error
	for {
		_, msg, err := p.conn.ReadMessage() //block
		if err != nil {
			break
		}
		fmt.Printf("%s", msg)
	}
}

func (p *peer) write() {
	defer p.close()
	for {
		msg, ok := <-p.inbox //block
		if !ok {
			break
		}
		p.conn.WriteMessage(websocket.TextMessage, msg)
	}
}

func initPeer(conn *websocket.Conn, addr, port string) *peer {
	key := fmt.Sprintf("%s:%s", addr, port)
	p := &peer{
		addr:  addr,
		port:  port,
		key:   key,
		conn:  conn,
		inbox: make(chan []byte),
	}

	go p.read()
	go p.write()
	Peers[key] = p
	return p
}
