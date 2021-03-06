package p2p

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/auturnn/kickshaw-coin/blockchain"
	"github.com/auturnn/kickshaw-coin/utils"
	log "github.com/kataras/golog"
)

type MessageKind int

const (
	MessageNewestBlock MessageKind = iota
	MessageAllBlocksrequest
	MessageAllBlocksResponse
	MessageNewBlockNotify
	MessageNewTxNotify
	MessageNewPeerNotify
)

type Message struct {
	Kind    MessageKind
	Payload []byte
}

func makeMessage(kind MessageKind, payload interface{}) []byte {
	m := Message{
		Kind:    kind,
		Payload: utils.ToJSON(payload),
	}
	return utils.ToJSON(m)
}

func sendNewestBlock(p *peer) {
	logf(log.InfoLevel, "Peer %s - Sending newest block", p.key)
	b, definedErr := blockchain.FindBlock(blockchain.BlockChain().NewestHash)
	utils.HandleError(definedErr, nil)
	m := makeMessage(MessageNewestBlock, b)
	p.inbox <- m
}

func requestAllBlocks(p *peer) {
	m := makeMessage(MessageAllBlocksrequest, nil)
	p.inbox <- m
}

func sendAllBlocks(p *peer) {
	m := makeMessage(MessageAllBlocksResponse, blockchain.Blocks(blockchain.BlockChain()))
	p.inbox <- m
}

func notifyNewBlock(b *blockchain.Block, p *peer) {
	m := makeMessage(MessageNewBlockNotify, b)
	p.inbox <- m
}

func notifyNewTx(tx *blockchain.Tx, p *peer) {
	m := makeMessage(MessageNewTxNotify, tx)
	p.inbox <- m
}

func notifyNewPeer(addr string, p *peer) {
	m := makeMessage(MessageNewPeerNotify, addr)
	p.inbox <- m
}

func handlerMsg(m *Message, p *peer) {
	switch m.Kind {
	case MessageNewestBlock:
		logf(log.InfoLevel, "Peer %s - Received the newest block", p.key)
		var payload blockchain.Block
		utils.HandleError(json.Unmarshal(m.Payload, &payload), nil)
		block, definedErr := blockchain.FindBlock(blockchain.BlockChain().NewestHash)
		utils.HandleError(definedErr, nil)

		if payload.Hash == block.Hash && payload.Height == block.Height {
			return
		}

		if payload.Height >= block.Height {
			logf(log.InfoLevel, "Peer %s - Requesting all blocks", p.key)
			requestAllBlocks(p)
		} else {
			sendNewestBlock(p)
		}

	case MessageAllBlocksrequest:
		logf(log.InfoLevel, "Peer %s - wants all the blocks", p.key)
		sendAllBlocks(p)

	case MessageAllBlocksResponse:
		logf(log.InfoLevel, "Peer %s - Received all the blocks", p.key)
		var payload []*blockchain.Block
		utils.HandleError(json.Unmarshal(m.Payload, &payload), nil)
		blockchain.BlockChain().Replace(payload)

	case MessageNewBlockNotify:
		logf(log.InfoLevel, "Peer %s - NewBlockNotify!", p.key)
		var payload *blockchain.Block
		utils.HandleError(json.Unmarshal(m.Payload, &payload), nil)
		blockchain.BlockChain().AddPeerBlock(payload)

	case MessageNewTxNotify:
		logf(log.InfoLevel, "Peer %s - NewTxNotify!", p.key)
		var payload *blockchain.Tx
		utils.HandleError(json.Unmarshal(m.Payload, &payload), nil)
		blockchain.Mempool().AddPeerTx(payload)

	case MessageNewPeerNotify:
		var payload string
		// {???????????????peerAddr : ???????????????peerPort : ???????????????peerWallet}
		// :{???????????????peerAddr: ???????????????peerPort : ???????????????peerWallet}
		utils.HandleError(json.Unmarshal(m.Payload, &payload), nil)
		parts := strings.Split(payload, ":")
		logf(log.InfoLevel, "Peer %s - NewPeerNotify!", parts[:3])
		server, _ := strconv.ParseBool(parts[5])
		AddPeer(parts[0:3], parts[3:5], server)
	}
}
