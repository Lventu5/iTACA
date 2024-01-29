package controllers

import (
	"categorizer/analysis"
	"categorizer/retrieve"
	"context"
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
	go a.analyser.Analyse(a.ctx, a.queue, a.results)
	select {
	case <-exit:
		a.ctx.Done()
		return
	case <-a.ctx.Done():
		return
	}
}
