package logging

import (
	"categorizer/analysis"
	"context"
	"fmt"
	"time"
)

type ConsoleLogger struct {
	ctx context.Context
}

func NewConsoleLogger(ctx context.Context) *ConsoleLogger {
	return &ConsoleLogger{ctx: ctx}
}

func (l *ConsoleLogger) Log(result analysis.StaticAnalysisResult) {
	data := ""
	for i, category := range result.MostLikelyCategories {
		data += fmt.Sprintf("%d) %s ", i, category)
	}
	_, err := fmt.Printf("%v: port %d - %s\n", time.Now(), result.SrcPort, data)
	if err != nil {
		fmt.Printf("Error printing to console: %v", err)
	}
}
