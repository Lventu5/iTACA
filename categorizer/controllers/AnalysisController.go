package controllers

import (
	"categorizer/analysis"
	"categorizer/retrieve"
	"context"
	"fmt"
	"strings"
)

type AnalysisController struct {
	ctx      context.Context
	queue    <-chan retrieve.Result
	results  chan<- analysis.StaticAnalysisResult
	analyser analysis.StaticAnalyser
}

func NewAnalysisController(ctx context.Context, queue <-chan retrieve.Result, results chan<- analysis.StaticAnalysisResult, analyser analysis.StaticAnalyser) *AnalysisController {
	return &AnalysisController{ctx: ctx, queue: queue, results: results, analyser: analyser}
}

func (a *AnalysisController) Start(exit <-chan bool, cancel context.CancelFunc) {
	ctrl := make(chan bool, 1)
	for {
		select {
		case <-exit:
			fmt.Println("AnalysisController: task stopped")
			cancel()
			return
		case <-a.ctx.Done():
			fmt.Println("RetrieverController: task stopped")
			return
		case stream := <-a.queue:
			streams := strings.Split(stream.Stream, "\n")
			for _, s := range streams {
				newStream := retrieve.Result{Stream: s, SrcPort: stream.SrcPort}
				go a.analyser.Analyse( /*a.ctx, cancel, */ newStream, a.results, ctrl)
			}
		}
	}
}
