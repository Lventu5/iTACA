package tests

import (
	"categorizer/retrieve"
	"context"
	"testing"
)

func TestCaronteRetriever(t *testing.T) {
	address := "0.0.0.0"
	port := 3333

	ctx, cancel := context.WithCancel(context.Background())
	queue := make(chan retrieve.Result)

	retriever := retrieve.NewCaronteRetriever(address, uint16(port))

	go retriever.Retrieve(ctx, cancel, queue)

	for i := 0; i < 1000; i++ {
		select {
		case result := <-queue:
			t.Log(result)
		}
	}
	cancel()
	close(queue)
	return
}
