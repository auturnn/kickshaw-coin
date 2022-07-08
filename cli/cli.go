package cli

import (
	"flag"

	"github.com/auturnn/kickshaw-coin/blockchain"
	"github.com/auturnn/kickshaw-coin/db"
	"github.com/auturnn/kickshaw-coin/rest"
)

func Start() {
	defer db.Close()
	// random port를 통해 네트워크 연결 충돌을 막고자 하였으나
	// 공유기 방화벽등의 문제로 인해 사용중지
	// randomPort, err := freeport.GetFreePort()
	// utils.HandleError(err)

	//하나의 컴퓨터에서 여러 테스트를 위해 적용됨
	var port int
	var status string
	flag.IntVar(&port, "port", 7120, "Set port of the someone server")
	flag.StringVar(&status, "network", "server", `Set application status in [local, server]`)

	//cmd shortcut
	flag.IntVar(&port, "p", 7120, "Set port of the someone server")
	flag.StringVar(&status, "n", "server", `Set application status in [local, server]`)
	flag.Parse()

	blockchain.Mempool()     //mempool Init
	db.InitDB()              //db init
	rest.Start(port, status) //rest start
}
