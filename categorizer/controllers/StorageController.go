package controllers

import (
	"categorizer/analysis"
	"context"
	"fmt"
	"os"
	"time"
)

type StorageController struct {
	ctx     context.Context
	file    string
	results <-chan analysis.StaticAnalysisResult
}

func NewStorageController(ctx context.Context, results <-chan analysis.StaticAnalysisResult) *StorageController {
	return &StorageController{ctx: ctx, file: "categorizer.log", results: results}
}

func (s *StorageController) Start(exit <-chan bool) {
	for {
		file, err := os.OpenFile(s.file, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			fmt.Printf("Error opening file: %v", err)
		}

		select {
		case <-exit:
			err := file.Close()
			if err != nil {
				fmt.Printf("Error closing log file: %v", err)
				return
			}
			s.ctx.Done()
			return
		case <-s.ctx.Done():
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
