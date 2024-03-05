package analysis

import (
	"categorizer/retrieve"
	"context"
)

// StaticAnalyser : interface, defines the general method "analyse" which is used to analyse a tcp stream
type StaticAnalyser interface {
	Analyse(ctx context.Context, stream retrieve.Result, result chan<- StaticAnalysisResult)
}

type StaticAnalysisResult struct {
	MostLikelyCategories [5]string
	SrcPort              uint16
}
