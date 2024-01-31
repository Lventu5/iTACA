package controllers

import (
	"categorizer/analysis"
	"categorizer/retrieve"
	"context"
	"fmt"
)

type AnalysisController struct {
	ctx      context.Context
	queue    <-chan retrieve.Result
	results  chan<- analysis.StaticAnalysisResult
	analyser analysis.StaticAnalyser
}

func NewAnalysisController(ctx context.Context, queue <-chan retrieve.Result, analyser analysis.StaticAnalyser) *AnalysisController {
	return &AnalysisController{ctx: ctx, queue: queue, analyser: analyser}
}

func (a *AnalysisController) Start(exit <-chan bool) {
	for {
		select {
		case <-exit:
			fmt.Println("AnalysisController: task stopped")
			a.ctx.Done()
			return
		case <-a.ctx.Done():
			fmt.Println("RetrieverController: task stopped")
			return
		case stream := <-a.queue:
			go a.analyser.Analyse(a.ctx, stream, a.results)
		}
	}
}
