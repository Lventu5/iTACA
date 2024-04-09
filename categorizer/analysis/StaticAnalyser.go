package analysis

import (
	"categorizer/retrieve"
)

// StaticAnalyser : interface, defines the general method "analyse" which is used to analyse a tcp stream
type StaticAnalyser interface {
	Analyse( /*ctx context.Context, cancel context.CancelFunc, */ stream retrieve.Result, result chan<- StaticAnalysisResult, ctrl chan bool)
}

type StaticAnalysisResult struct {
	MostLikelyCategories [5]string
	SrcPort              uint16
}
