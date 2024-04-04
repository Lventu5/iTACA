package logging

import (
	"categorizer/analysis"
	"context"
	"fmt"
	"os"
	"time"
)

type FileLogger struct {
	ctx  context.Context
	file *os.File
}

func NewFileLogger(ctx context.Context) (*FileLogger, error) {
	file, err := os.OpenFile("categorizer.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("Error opening file: %v", err)
		return nil, err
	}
	return &FileLogger{ctx: ctx, file: file}, nil
}

func (s *FileLogger) Log(result analysis.StaticAnalysisResult) {
	data := ""
	for i, category := range result.MostLikelyCategories {
		data += fmt.Sprintf("%d) %s ", i, category)
	}
	_, err := s.file.WriteString(fmt.Sprintf("%v: port %d - %s\n", time.Now(), result.SrcPort, data))
	if err != nil {
		fmt.Printf("Error writing to log file: %v", err)
	}
}

func (s *FileLogger) Close() {
	err := s.file.Close()
	if err != nil {
		fmt.Printf("Error closing file: %v", err)
	}
}
