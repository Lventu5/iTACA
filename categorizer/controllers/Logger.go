package controllers

import (
	"categorizer/analysis"
	"context"
	"fmt"
	"time"
)

type Logger struct {
	ctx     context.Context
	results <-chan analysis.StaticAnalysisResult
}

func NewLogger(ctx context.Context, results <-chan analysis.StaticAnalysisResult) *Logger {
	return &Logger{ctx: ctx, results: results}
}

func (l *Logger) Start(exit <-chan bool) {
	for {
		select {
		case <-exit:
			l.ctx.Done()
			return
		case <-l.ctx.Done():
			return
		case result := <-l.results:
			data := ""
			for i, category := range result.MostLikelyCategories {
				data += fmt.Sprintf("%d) %s ", i, category)
			}
			_, err := fmt.Printf("%v: port %d - %s\n", time.Now(), result.SrcPort, data)
			if err != nil {
				fmt.Printf("Error writing to log file: %v", err)
			}
		}
	}
}
