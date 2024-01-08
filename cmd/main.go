package main

import (
	"os"

	"github.com/taylormonacelli/hisfantasy"
)

func main() {
	code := hisfantasy.Execute()
	os.Exit(code)
}
