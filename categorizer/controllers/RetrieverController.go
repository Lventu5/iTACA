package controllers

import (
	"categorizer/retrieve"
	"context"
	"fmt"
)

type RetrieverController struct {
	ctx       context.Context
	queue     chan<- retrieve.Result
	retriever retrieve.Retriever
}

func NewRetrieverController(ctx context.Context, queue chan<- retrieve.Result, retriever retrieve.Retriever) *RetrieverController {
	return &RetrieverController{ctx: ctx, queue: queue, retriever: retriever}
}

func (c *RetrieverController) Start(exit <-chan bool, cancel context.CancelFunc) {
	go c.retriever.Retrieve(c.ctx, cancel, c.queue)
	select {
	case <-exit:
		fmt.Println("RetrieverController: task stopped")
		cancel()
		return
	case <-c.ctx.Done():
		fmt.Println("RetrieverController: task stopped")
		return
	}
}
