package cli

import (
	"flag"
	"fmt"
	"runtime"

	"github.com/auturnn/kickshaw-coin/explorer"
	"github.com/auturnn/kickshaw-coin/rest"
)

func usage()  {
	fmt.Printf("Welcome to kickshaw-coin\n\n")
	fmt.Printf("Please use the following flags:\n\n")
	runtime.Goexit()
}

func Start()  {
	mode := flag.String("mode", "all", "Choose between 'html' or 'rest' or 'all'")
	port := flag.Int("port", 8080, "Set port of the someone server")
	hport := flag.Int("hport", 3000, "Set port of the html server for mode 'all' only")
	rport := flag.Int("rport", 8080, "Set port of the REST API server for mode 'all' only")
	flag.Parse()

	switch *mode {
	case "rest" :
		rest.Start(*port)
	case "html" :
		explorer.Start(*port)
	case "all" :
		go explorer.Start(*hport)
		rest.Start(*rport)
	default: 
		usage()
	}
}
