package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/auturnn/kickshaw-coin/blockchain"
	"github.com/auturnn/kickshaw-coin/utils"
	"github.com/gorilla/mux"
)


type url string

func (u url) MashalText() ([]byte, error)  {
	return []byte(fmt.Sprintf("http://localhost%s%s", port, u)),nil
}

type urlDiscription struct{
	URL url `json:"url"`
	Method string `json:"method"`
	Description string `json:"description"`
	Payload string `json:"payload,omitempty"`
}

var port string

func documentation(rw http.ResponseWriter, r *http.Request)  {
	data := []urlDiscription{
		{
			URL: url("/"),
			Method: "GET",
			Description: "See Documentation",
		},
		{
			URL: url("/status"),
			Method: "GET",
			Description: "See the Status of the Blockchain",
		},
		{
			URL: url("/blocks"),
			Method: "GET",
			Description: "showing my blocks",
		},
		{
			URL: url("/blocks"),
			Method: "POST",
			Description: "Add blocks",
		},
		{
			URL: url("/blocks/{hash}"),
			Method: "GET",
			Description: "searching for block's id",
		},
	}
	json.NewEncoder(rw).Encode(data)
}

type addBlockBody struct{
	Data string `json:"data"`
}

func getBlocks(rw http.ResponseWriter, r *http.Request)  {
	switch r.Method{
	case "GET": 
		json.NewEncoder(rw).Encode(blockchain.BlockChain().Blocks())
	case "POST": 
		var addblockbody addBlockBody
		utils.HandleError(json.NewDecoder(r.Body).Decode(&addblockbody))
		blockchain.BlockChain().AddBlock(addblockbody.Data)
		rw.WriteHeader(http.StatusCreated)
	}
}
type errorResponse struct{
	ErrorMessage string `json:"errorMessage"`
}

func getBlock(rw http.ResponseWriter, r *http.Request)  {
	hash := mux.Vars(r)["hash"]
	block, err := blockchain.FindBlock(hash)
	encoder := json.NewEncoder(rw)
	if err == blockchain.ErrNotFound{
		encoder.Encode(errorResponse{fmt.Sprint(err)})
	} else{
		encoder.Encode(block)
	}
}

func getStatus(rw http.ResponseWriter, r *http.Request) {
	json.NewEncoder(rw).Encode(blockchain.BlockChain())
}

func Start(aPort int)  {
	port = fmt.Sprintf(":%d",aPort)

	router := mux.NewRouter()
	router.Use(jsonContentTypeMiddleware)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/status", getStatus).Methods("GET")
	router.HandleFunc("/blocks", getBlocks).Methods("GET","POST")
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", getBlock).Methods("GET")
	fmt.Printf("Listening http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler  {
	//adaptor pattern
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-type","application/json")
		next.ServeHTTP(rw, r)
	})
}