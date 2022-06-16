package cli

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/auturnn/kickshaw-coin/db"
	"github.com/auturnn/kickshaw-coin/rest"
)

// func usage() {
// 	fmt.Printf("Welcome to kickshaw-coin\n\n")
// 	fmt.Printf("Please use the following flags:\n\n")
// 	fmt.Printf("-port     : Set port of the someone server\n")
// 	os.Exit(0)
// }

func Start() {
	defer db.Close()

	envPort, err := strconv.Atoi(os.Getenv("KICKSHAW_PORT"))
	if err != nil{
		log.Fatal("port environment has not setting")
		return
	}
	
	//하나의 컴퓨터에서 여러 테스트를 위해 적용됨
	port := flag.Int("port", envPort, "Set port of the someone server")

	flag.Parse()

	db.InitDB(*port)
	rest.Start(*port)
}
