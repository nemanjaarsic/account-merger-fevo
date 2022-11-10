package main

import (
	"os"
)

func main() {
	path := os.Args[1]
	acc := ImportAccounts(path)
	results := MergeAccounts(acc)
	PrintResults(results)
}
