package p2p

import (
	"encoding/json"
	"fmt"

	"github.com/auturnn/kickshaw-coin/blockchain"
	"github.com/auturnn/kickshaw-coin/utils"
)

type MessageKind int

const (
	MessageNewestBlock MessageKind = iota
	MessageAllBlocksrequest
	MessageAllBlocksResponse
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
	fmt.Printf("Sending newest block to %s\n", p.key)
	b, err := blockchain.FindBlock(blockchain.BlockChain().NewestHash)
	utils.HandleError(err)
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

func handlerMsg(m *Message, p *peer) {
	switch m.Kind {
	case MessageNewestBlock:
		{
			fmt.Printf("Received the newst block from %s\n", p.key)
			var payload blockchain.Block
			utils.HandleError(json.Unmarshal(m.Payload, &payload))

			b, err := blockchain.FindBlock(blockchain.BlockChain().NewestHash)
			utils.HandleError(err)

			if payload.Height >= b.Height {
				fmt.Printf("Requesting all blocks from %s\n", p.key)
				requestAllBlocks(p)
			} else {
				sendNewestBlock(p)
			}
		}

	case MessageAllBlocksrequest:
		{
			fmt.Printf("%s wants all the blocks\n", p.key)
			sendAllBlocks(p)
		}

	case MessageAllBlocksResponse:
		{
			fmt.Printf("Received all the blocks from %s\n", p.key)
			var payload []*blockchain.Block
			utils.HandleError(json.Unmarshal(m.Payload, &payload))
			blockchain.BlockChain().Replace(payload)
		}

	}
}
