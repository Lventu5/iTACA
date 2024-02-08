package controllers

import (
	"categorizer/analysis"
	"context"
	"fmt"
	"os"
	"time"
)

type FileController struct {
	ctx     context.Context
	file    string
	results <-chan analysis.StaticAnalysisResult
}

func NewFileController(ctx context.Context, results <-chan analysis.StaticAnalysisResult) *FileController {
	return &FileController{ctx: ctx, file: "categorizer.log", results: results}
}

func (s *FileController) Start(exit <-chan bool, cancel context.CancelFunc) {
	for {
		file, err := os.OpenFile(s.file, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			fmt.Printf("Error opening file: %v", err)
		}

		select {
		case <-exit:
			fmt.Println("FileController: task stopped")
			err := file.Close()
			if err != nil {
				fmt.Printf("Error closing log file: %v", err)
				return
			}
			cancel()
			return
		case <-s.ctx.Done():
			fmt.Println("FileController: task stopped")
			err := file.Close()
			if err != nil {
				fmt.Printf("Error closing log file: %v", err)
				return
			}
			return
		case result := <-s.results:
			data := ""
			for i, category := range result.MostLikelyCategories {
				data += fmt.Sprintf("%d) %s ", i, category)
			}
			_, err := file.WriteString(fmt.Sprintf("%v: port %d - %s\n", time.Now(), result.SrcPort, data))
			if err != nil {
				fmt.Printf("Error writing to log file: %v", err)
			}
		}
	}
}
