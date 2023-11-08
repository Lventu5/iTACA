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

func NewRetrieverController(queue chan<- retrieve.Result, retriever retrieve.Retriever) *RetrieverController {
	return &RetrieverController{ctx: context.Background(), queue: queue, retriever: retriever}
}

func (c *RetrieverController) Start(exit <-chan bool) {
	go c.retriever.Retrieve(c.ctx, c.queue)
	select {
	case <-exit:
		fmt.Println("RetrieverController: task stopped")
		c.ctx.Done()
		return
	}
}
