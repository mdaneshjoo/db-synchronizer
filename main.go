package main

import (
	"fmt"

	"github.com/mdaneshjoo/db-synchronizer/cdc"
)


func main() {
	cdcCh := make(chan cdc.DebeziumPayload)

	go func() {
		defer close(cdcCh)
		cdc.Capture(100, cdcCh)
	}()
	for payload := range cdcCh {
		fmt.Println("message recived",payload)
	}
}
