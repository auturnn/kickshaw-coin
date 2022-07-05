package p2p

import (
	"fmt"
	"net/http"

	"github.com/auturnn/kickshaw-coin/blockchain"
	"github.com/auturnn/kickshaw-coin/utils"
	"github.com/auturnn/kickshaw-coin/wallet"
	"github.com/gorilla/websocket"
	log "github.com/kataras/golog"
)

var upgrader = websocket.Upgrader{}

var logf = log.Logf

func Upgrade(rw http.ResponseWriter, r *http.Request) {
	//port :3000 will upgrade the request from :4000
	//AddPeer에서는 기존의 peer에 저장되있던 node들의 정보를 port와 waddr 쿼리로 보내지만
	//Upgrade를 받는 쪽에서는 해당 노드들이 새로 연결을 요청하는 쪽이기 때문에 newPeer가 된다.
	myWalletAddr := wallet.WalletLayer{}.GetAddress()[:5]
	if r.URL.Query().Get("nwddr") != myWalletAddr {
		//정말 드물게 자신의 IP, Port가 같을 수도 있다고 생각하여 구분사항을 추가.
		return
	}

	newPeerPort := r.URL.Query().Get("port")
	newPeerWddr := r.URL.Query().Get("wddr")

	ip := utils.Splitter(r.RemoteAddr, ":", 0)
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return newPeerPort != "" && ip != ""
	}
	logf(log.InfoLevel, "%s wants an upgrade", newPeerPort)

	conn, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleError(err)
	initPeer(conn, ip, newPeerPort, newPeerWddr, false)
}

// AddPeer 매개변수 => 연결해오는peerAddr / 연결해오는peerPort / 연결해오는peerWallet / 연결되있는peerPort / 연결되있는peerWallet / server유무
func AddPeer(newPeerAddr, newPeerPort, newPeerWAddr, existPeerPort, existPeerWAddr string, server bool) {
	logf(log.InfoLevel, "%s:%s:%s - wants to connect to port - %s:%s", newPeerAddr, newPeerPort, newPeerWAddr, existPeerPort, existPeerWAddr)
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?nwddr=%s&port=%s&wddr=%s", newPeerAddr, newPeerPort, newPeerWAddr, existPeerPort, existPeerWAddr), nil)
	utils.HandleError(err)
	p := initPeer(conn, newPeerAddr, newPeerPort, newPeerWAddr, server)
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
