package controllers

import (
	"categorizer/analysis"
	"categorizer/logging"
	"context"
	"fmt"
)

type OutputController struct {
	ctx           context.Context
	results       <-chan analysis.StaticAnalysisResult
	fileLogger    *logging.FileLogger
	consoleLogger *logging.ConsoleLogger
}

func NewOutputController(ctx context.Context, results <-chan analysis.StaticAnalysisResult, useFile bool) *OutputController {
	var fileLogger *logging.FileLogger
	var err error
	if useFile {
		fileLogger, err = logging.NewFileLogger(ctx)
		if err != nil {
			fmt.Println("Error creating file logger")
		}
	} else {
		fileLogger = nil
	}

	consoleLogger := logging.NewConsoleLogger(ctx)

	return &OutputController{ctx: ctx, results: results, fileLogger: fileLogger, consoleLogger: consoleLogger}
}

func (c *OutputController) Start(exit <-chan bool, cancel context.CancelFunc) {
	for {
		select {
		case <-exit:
			if c.fileLogger != nil {
				c.fileLogger.Close()
			}
			cancel()
			return
		case res := <-c.results:
			if c.fileLogger != nil {
				c.fileLogger.Log(res)
			}
			c.consoleLogger.Log(res)
		case <-c.ctx.Done():
			if c.fileLogger != nil {
				c.fileLogger.Close()
			}
			return
		}
	}
}
