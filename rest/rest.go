package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/auturnn/kickshaw-coin/blockchain"
	"github.com/auturnn/kickshaw-coin/p2p"
	"github.com/auturnn/kickshaw-coin/utils"
	"github.com/auturnn/kickshaw-coin/wallet"
	"github.com/gorilla/mux"
)

type url string

func (u url) MashalText() ([]byte, error) {
	return []byte(fmt.Sprintf("http://localhost%s%s", port, u)), nil
}

type urlDiscription struct {
	URL         url    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

var port string

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []urlDiscription{
		{
			URL:         url("/"),
			Method:      "GET",
			Description: "See Documentation",
		},
		{
			URL:         url("/status"),
			Method:      "GET",
			Description: "See the Status of the Blockchain",
		},
		{
			URL:         url("/blocks"),
			Method:      "GET",
			Description: "showing my blocks",
		},
		{
			URL:         url("/blocks"),
			Method:      "POST",
			Description: "Add blocks",
		},
		{
			URL:         url("/blocks/{hash}"),
			Method:      "GET",
			Description: "searching for block's id",
		},
		{
			URL:         url("/balance/{address}"),
			Method:      "GET",
			Description: "Get TxOuts for an Address",
		},
		{
			URL:         url("/mempool"),
			Method:      "GET",
			Description: "Get Mempool",
		},
		{
			URL:         url("/transactions"),
			Method:      "POST",
			Description: "Make a transaction",
		},
		{
			URL:         url("/ws"),
			Method:      "GET",
			Description: "Upgrade to websocket",
		},
	}
	json.NewEncoder(rw).Encode(data)
}

type addBlockBody struct {
	Data string `json:"data"`
}

func getBlocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		json.NewEncoder(rw).Encode(blockchain.Blocks(blockchain.BlockChain()))
	case "POST":
		newBlock := blockchain.BlockChain().AddBlock()
		p2p.BroadcastNewBlock(newBlock)
		rw.WriteHeader(http.StatusCreated)
	}
}

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

func getBlock(rw http.ResponseWriter, r *http.Request) {
	hash := mux.Vars(r)["hash"]
	block, err := blockchain.FindBlock(hash)
	encoder := json.NewEncoder(rw)
	if err == blockchain.ErrNotFound {
		utils.HandleError(encoder.Encode(errorResponse{fmt.Sprint(err)}))
	} else {
		utils.HandleError(encoder.Encode(block))
	}
}

func getStatus(rw http.ResponseWriter, r *http.Request) {
	blockchain.Status(blockchain.BlockChain(), rw)
}

type balanceResponse struct {
	Address string `json:"address"`
	Balance int    `json:"balance"`
}

func getBalance(rw http.ResponseWriter, r *http.Request) {
	address := mux.Vars(r)["address"]
	total := r.URL.Query().Get("total")
	switch total {
	case "true":
		amount := blockchain.BalanceByAddress(address, blockchain.BlockChain())
		utils.HandleError(json.NewEncoder(rw).Encode(balanceResponse{address, amount}))
	default:
		utils.HandleError(json.NewEncoder(rw).Encode(blockchain.UTxOutsByAddress(address, blockchain.BlockChain())))
	}
}

type addTxPayload struct {
	To     string `json:"to"`
	Amount int    `json:"amount"`
}

func getMempool(rw http.ResponseWriter, r *http.Request) {
	utils.HandleError(json.NewEncoder(rw).Encode(blockchain.Mempool().Txs))
}

func transactions(rw http.ResponseWriter, r *http.Request) {
	var payload addTxPayload
	utils.HandleError(json.NewDecoder(r.Body).Decode(&payload))
	tx, err := blockchain.Mempool().AddTx(payload.To, payload.Amount)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(errorResponse{err.Error()})
		return
	}
	go p2p.BroadcastNewTx(tx)
	rw.WriteHeader(http.StatusCreated)
}

type myWalletResponse struct {
	Address string `json:"address"`
}

func myWallet(rw http.ResponseWriter, r *http.Request) {
	address := wallet.Wallet().Address
	// json.NewEncoder(rw).Encode(struct{
	// 	Address string
	// }{Address: address})
	json.NewEncoder(rw).Encode(myWalletResponse{Address: address})
}

type addPeerPayload struct {
	Address, Port string
}

func peers(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var payload addPeerPayload
		json.NewDecoder(r.Body).Decode(&payload)
		p2p.AddPeer(payload.Address, payload.Port, port[1:], true)
		rw.WriteHeader(http.StatusOK)
	case "GET":
		json.NewEncoder(rw).Encode(p2p.AllPeers(&p2p.Peers))
	}
}

func Start(aPort int) {
	port = fmt.Sprintf(":%d", aPort)

	router := mux.NewRouter()
	router.Use(jsonContentTypeMiddleware)
	router.Use(loggerMiddleware)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/status", getStatus).Methods("GET")
	router.HandleFunc("/blocks", getBlocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", getBlock).Methods("GET")
	router.HandleFunc("/balance/{address}", getBalance).Methods("GET")
	router.HandleFunc("/mempool", getMempool).Methods("GET")
	router.HandleFunc("/wallet", myWallet).Methods("GET")
	router.HandleFunc("/transactions", transactions).Methods("POST")
	router.HandleFunc("/ws", p2p.Upgrade).Methods("GET")
	router.HandleFunc("/peers", peers).Methods("POST", "GET")
	fmt.Printf("Listening http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	//adaptor pattern
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL)
		next.ServeHTTP(rw, r)
	})
}
