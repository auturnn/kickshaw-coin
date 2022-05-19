package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/auturnn/kickshaw-coin/rest"
)

func usage() {
	fmt.Printf("Welcome to kickshaw-coin\n\n")
	fmt.Printf("Please use the following flags:\n\n")
	fmt.Printf("-port     : Set port of the someone server\n")
	fmt.Printf("-id 	  : Input user name\n")
	fmt.Printf("-password : Input user password")
	os.Exit(0)
}

func Start() {
	if len(os.Args) == 1 {
		usage()
	}

	port := flag.Int("port", 8333, "Set port of the someone server")
	// name := flag.String("id", "anonymous", "Input user name")
	// password := flag.String("password", "0000", "Input user password")
	flag.Parse()
	rest.Start(*port)
}
