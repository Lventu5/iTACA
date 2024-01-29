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

func (c *RetrieverController) Start(exit <-chan bool) {
	go c.retriever.Retrieve(c.ctx, c.queue)
	select {
	case <-exit:
		fmt.Println("RetrieverController: task stopped")
		c.ctx.Done()
		return
	case <-c.ctx.Done():
		fmt.Println("RetrieverController: task stopped")
		return
	}
}
