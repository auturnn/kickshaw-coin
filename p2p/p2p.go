package p2p

import (
	"fmt"
	"log"
	"net/http"

	"github.com/auturnn/kickshaw-coin/blockchain"
	"github.com/auturnn/kickshaw-coin/utils"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func Upgrade(rw http.ResponseWriter, r *http.Request) {
	//port :3000 will upgrade the request from :4000
	openPort := r.URL.Query().Get("openPort")
	wAddr := r.URL.Query().Get("wAddr")
	ip := utils.Splitter(r.RemoteAddr, ":", 0)
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return openPort != "" && ip != ""
	}
	fmt.Printf("%s wants an upgrade\n", openPort)

	conn, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleError(err)
	initPeer(conn, ip, openPort, wAddr)
}

func AddPeer(addr, port, openPort, wAddr string, broadcast bool) {
	log.Printf("%s:%s:%s wants to connect to port %s\n", addr, openPort, wAddr, port)
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?openPort=%s&wAddr=%s", addr, port, openPort, wAddr), nil)
	utils.HandleError(err)
	p := initPeer(conn, addr, port, wAddr)
	sendNewestBlock(p)
}

func BroadcastNewBlock(b *blockchain.Block) {
	for _, p := range Peers.v {
		notifyNewBlock(b, p)
	}
}

func BroadcastNewTx(tx *blockchain.Tx) {
	for _, p := range Peers.v {
		notifyNewTx(tx, p)
	}
}
