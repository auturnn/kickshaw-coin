package p2p

import (
	"fmt"
	"net/http"

	"github.com/auturnn/kickshaw-coin/utils"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func Upgrade(rw http.ResponseWriter, r *http.Request) {
	//port :3000 will upgrade the request from :4000
	//임시
	openPort := r.URL.Query().Get("openPort")
	ip := utils.Splitter(r.RemoteAddr, ":", 0)
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return openPort != "" && ip != ""
	}

	conn, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleError(err)
	initPeer(conn, ip, openPort)
}

func AddPeer(addr, port, openPort string) {
	// node:4000으로 들어오면 node:3000과 중계
	// peer -(request)> :4000 -(request(upgrade))> :3000
	// port :4000 is requesting an upgrade from the port :3000
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?openPort=%s", addr, port, openPort[1:]), nil)
	utils.HandleError(err)
	initPeer(conn, addr, port)

}
