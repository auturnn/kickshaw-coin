package main

import (
	"github.com/auturnn/kickshaw-coin/explorer"
	"github.com/auturnn/kickshaw-coin/rest"
)

func main()  {
	go explorer.Start(3000)
	rest.Start(8080)
}
