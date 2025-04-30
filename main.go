package main

import (
	"github.com/Ayobami6/schema_dump/cmd"
	_ "github.com/lib/pq"
)

func main() {

	if err := cmd.RootCmd.Execute(); err != nil {
		panic(err)
	}

}
