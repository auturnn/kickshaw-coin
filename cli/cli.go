package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/auturnn/kickshaw-coin/explorer"
	"github.com/auturnn/kickshaw-coin/rest"
)

func usage() {
	fmt.Printf("Welcome to kickshaw-coin\n\n")
	fmt.Printf("Please use the following flags:\n\n")
	fmt.Printf("-mode  :   Choose between 'html' or 'rest' or 'all'\n")
	fmt.Printf("-port  :   Set port of the someone server\n")
	os.Exit(0)
}

func Start() {
	if len(os.Args) == 1 {
		usage()
	}

	mode := flag.String("mode", "rest", "Choose between 'html' or 'rest' or 'all'")
	port := flag.Int("port", 8333, "Set port of the someone server")
	flag.Parse()

	switch *mode {
	case "rest":
		rest.Start(*port)
	case "html":
		explorer.Start(*port)
	case "all":
		go explorer.Start(4000)
		rest.Start(8333)
	default:
		usage()
	}

}
