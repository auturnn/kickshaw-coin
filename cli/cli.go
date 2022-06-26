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

	//하나의 컴퓨터에서 여러 테스트를 위해 적용됨

	// randomPort, err := freeport.GetFreePort()
	// if err != nil{
	// 	utils.HandleError(errors.New("failed get port"))
	// }

	port := flag.Int("port", 8080, "Set port of the someone server")
	status := flag.Bool("network", true, `Set application status in [network, solo]`)
	flag.Parse()

	db.InitDB()
	rest.Start(*port, *status)
}
