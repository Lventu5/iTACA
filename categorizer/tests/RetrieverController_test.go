package tests

import (
	"categorizer/controllers"
	"categorizer/retrieve"
	"context"
	"testing"
)

func TestRetrieverController(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	queue := make(chan retrieve.Result)
	retriever := retrieve.NewCaronteRetriever("0.0.0.0", 3333)
	exit := make(chan bool)

	controller := controllers.NewRetrieverController(ctx, queue, retriever)
	go controller.Start(exit, cancel)

	for i := 0; i < 100; i++ {
		select {
		case result := <-queue:
			t.Log(result)
		}
	}
	exit <- true
	close(queue)
	close(exit)
	return
}
