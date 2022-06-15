package cli

import (
	"flag"

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
	port := flag.Int("port", 8333, "Set port of the someone server")
	flag.Parse()

	db.InitDB(*port)
	rest.Start(*port)
}
