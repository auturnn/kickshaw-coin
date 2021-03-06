package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/auturnn/kickshaw-coin/blockchain"
	"github.com/auturnn/kickshaw-coin/p2p"
	"github.com/auturnn/kickshaw-coin/utils"
	"github.com/auturnn/kickshaw-coin/wallet"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/kataras/golog"
	log "github.com/kataras/golog"
)

var port string

type url string

func (u url) MashalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

type urlDescription struct {
	URL         url    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

type addBlockBody struct {
	Data string `json:"data"`
}

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

type balanceResponse struct {
	Address string `json:"address"`
	Balance int    `json:"balance"`
}

type addTxPayload struct {
	To     string `json:"to"`
	Amount int    `json:"amount"`
}

type myWalletResponse struct {
	Address string `json:"address"`
}

type addPeerPayload struct {
	Address, Port string
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []urlDescription{
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
			Description: "See All Blocks",
		},
		{
			URL:         url("/blocks"),
			Method:      "POST",
			Description: "Add A Block",
			Payload:     "data:string",
		},
		{
			URL:         url("/blocks/{hash}"),
			Method:      "GET",
			Description: "See A Block",
		},
		{
			URL:         url("/balance/{address}"),
			Method:      "GET",
			Description: "Get TxOuts for an Address",
		},
		{
			URL:         url("/ws"),
			Method:      "GET",
			Description: "Upgrade to WebSockets",
		},
		{
			URL:         url("/peers"),
			Method:      "GET",
			Description: "Get all connecting Peer's address",
		},
	}
	json.NewEncoder(rw).Encode(data)
}

func getBlocks(rw http.ResponseWriter, r *http.Request) {
	json.NewEncoder(rw).Encode(blockchain.Blocks(blockchain.BlockChain()))
}

func postBlocks(rw http.ResponseWriter, r *http.Request) {
	newBlock := blockchain.BlockChain().AddBlock()
	p2p.BroadcastNewBlock(newBlock)
	rw.WriteHeader(http.StatusCreated)
}

func getBlock(rw http.ResponseWriter, r *http.Request) {
	hash := mux.Vars(r)["hash"]
	block, definedErr := blockchain.FindBlock(hash)
	encoder := json.NewEncoder(rw)
	if definedErr == blockchain.ErrNotFound {
		utils.HandleError(encoder.Encode(errorResponse{fmt.Sprint(definedErr)}), definedErr)
	} else {
		utils.HandleError(encoder.Encode(block), nil)
	}
}

func getStatus(rw http.ResponseWriter, r *http.Request) {
	blockchain.Status(blockchain.BlockChain(), rw)
}

func getBalance(rw http.ResponseWriter, r *http.Request) {
	address := mux.Vars(r)["address"]
	total := r.URL.Query().Get("total")
	switch total {
	case "true":
		amount := blockchain.BalanceByAddress(address, blockchain.BlockChain())
		utils.HandleError(json.NewEncoder(rw).Encode(balanceResponse{address, amount}), nil)
	default:
		utils.HandleError(json.NewEncoder(rw).Encode(blockchain.UTxOutsByAddress(address, blockchain.BlockChain())), nil)
	}
}

func getMempool(rw http.ResponseWriter, r *http.Request) {
	utils.HandleError(json.NewEncoder(rw).Encode(blockchain.Mempool().Txs), nil)
}

func transactions(rw http.ResponseWriter, r *http.Request) {
	var payload addTxPayload
	utils.HandleError(json.NewDecoder(r.Body).Decode(&payload), nil)

	tx, err := blockchain.Mempool().AddTx(payload.To, payload.Amount)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(errorResponse{err.Error()})
		return
	}

	p2p.BroadcastNewTx(tx)
	rw.WriteHeader(http.StatusCreated)
}

func myWallet(rw http.ResponseWriter, r *http.Request) {
	w := wallet.WalletLayer{}
	address := w.GetAddress()
	json.NewEncoder(rw).Encode(myWalletResponse{Address: address})
}

func getPeers(rw http.ResponseWriter, r *http.Request) {
	json.NewEncoder(rw).Encode(p2p.AllPeers(&p2p.Peers))
}

func P2PRouter(router *mux.Router) {
	router.HandleFunc("/peers", getPeers).Methods("GET")
	router.HandleFunc("/ws", p2p.Upgrade).Methods("GET")
}

func p2pServerConnect(ver, addr, sPort string) {
	fmt.Println(ver, addr, sPort)
	serverURL := fmt.Sprintf("%s://%s:%s", ver, addr, sPort)
	res, err := http.Get(serverURL + "/wallet")
	utils.HandleError(err, utils.ErrNetworkIsNotWork)
	var walletPayload struct {
		Address string `json:"address"`
	}
	json.NewDecoder(res.Body).Decode(&walletPayload)

	newPeer := []string{addr, sPort, walletPayload.Address[:5]}
	myInfo := []string{port[1:], wallet.WalletLayer{}.GetAddress()[:5]}
	p2p.AddPeer(newPeer, myInfo, true)
	log.Logf(log.InfoLevel, "p2p network Connecting...")
}

//wallet?????????????????? ????????? ?????? ????????? ????????? ?????????????????? ??????.
//?????????????????? ????????????????????? ??????????????? ?????????
func Start(p int, networkMode string) {
	port = fmt.Sprintf(":%s", strconv.Itoa(p))

	router := mux.NewRouter()
	router.Use(jsonContentTypeMiddleware, loggerMiddleware)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/status", getStatus).Methods("GET")
	router.HandleFunc("/blocks", getBlocks).Methods("GET")
	router.HandleFunc("/blocks", postBlocks).Methods("POST")
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", getBlock).Methods("GET")
	router.HandleFunc("/balance/{address}", getBalance).Methods("GET")
	router.HandleFunc("/mempool", getMempool).Methods("GET")
	router.HandleFunc("/wallet", myWallet).Methods("GET")
	router.HandleFunc("/transactions", transactions).Methods("POST")

	switch networkMode {
	case "server":
		P2PRouter(router)
		p2pServerConnect("http", "api.kickshaw-coin.com", "80")
		break
	case "local":
		P2PRouter(router)
		p2pServerConnect("http", "127.0.0.1", "8080")
		break
	case "alone":
		log.Info("alone mode start")
		break
	default:
		utils.HandleError(utils.ErrCMDNetwork, nil)
		break
	}

	cors := handlers.CORS()(router)
	log.Logf(golog.InfoLevel, "Listening http://localhost%s", port)
	recover := handlers.RecoveryHandler()(cors)
	log.Fatal(http.ListenAndServe(port, recover))
}
