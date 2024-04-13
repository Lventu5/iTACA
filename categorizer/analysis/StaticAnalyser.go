package analysis

import (
	"categorizer/retrieve"
)

// StaticAnalyser : interface, defines the general method "analyse" which is used to analyse a tcp stream
type StaticAnalyser interface {
	Analyse(stream retrieve.Result, result chan<- StaticAnalysisResult)
}

type StaticAnalysisResult struct {
	MostLikelyCategories [5]string
	SrcPort              uint16
}
