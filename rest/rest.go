package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

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
			URL: url("/blocks/{height}"),
			Method: "GET",
			Description: "searching for block's id",
		},
	}
	rw.Header().Add("Content-type", "application/json")
	json.NewEncoder(rw).Encode(data)
}

type addBlockBody struct{
	Data string `json:"data"`
}

func getBlocks(rw http.ResponseWriter, r *http.Request)  {
	rw.Header().Add("Content-type", "application/json")
	switch r.Method{
	case "GET": 
		json.NewEncoder(rw).Encode(blockchain.GetBlockChain().AllBlocks())
	case "POST": 
		var addblockbody addBlockBody
		utils.HandleError(json.NewDecoder(r.Body).Decode(&addblockbody))
		blockchain.GetBlockChain().AddBlock(addblockbody.Data)
		rw.WriteHeader(http.StatusCreated)
	}
}

func getBlock(rw http.ResponseWriter, r *http.Request)  {
	rw.Header().Add("Content-type", "application/json")
	
	id, err := strconv.Atoi(mux.Vars(r)["height"])
	utils.HandleError(err)

	block := blockchain.GetBlockChain().GetBlock(id)
	json.NewEncoder(rw).Encode(block)
}

func Start(aPort int)  {
	port = fmt.Sprintf(":%d",aPort)

	router := mux.NewRouter()
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/blocks", getBlocks).Methods("GET","POST")
	router.HandleFunc("/blocks/{height:[0-9]+}", getBlock).Methods("GET")

	fmt.Printf("Listening http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}