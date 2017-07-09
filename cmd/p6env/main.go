package main

import (
	"os"

	"github.com/Code-Hex/p6env"
)

func main() {
	os.Exit(p6env.New().Run())
}
