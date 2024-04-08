package tests

import (
	"categorizer/retrieve"
	"context"
	"testing"
)

func TestTulipRetriever(t *testing.T) {
	address := "0.0.0.0"
	port := 3000

	ctx, cancel := context.WithCancel(context.Background())
	queue := make(chan retrieve.Result)

	retriever := retrieve.NewTulipRetriever(address, uint16(port))

	go retriever.Retrieve(ctx, cancel, queue)

	for i := 0; i < 10; i++ {
		select {
		case result := <-queue:
			t.Log(result)
		}
	}
	cancel()
	close(queue)
	return
}
