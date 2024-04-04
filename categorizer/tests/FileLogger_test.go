package tests

import (
	"categorizer/analysis"
	"categorizer/logging"
	"context"
	"testing"
)

func TestFileController(t *testing.T) {
	ctx := context.Background()

	flc, err := logging.NewFileLogger(ctx)
	if err != nil {
		t.Errorf("Error creating FileLogger")
	}

	var test [5]analysis.StaticAnalysisResult
	test[0] = analysis.StaticAnalysisResult{MostLikelyCategories: [5]string{"a", "b", "c", "d", "e"}, SrcPort: 9999}
	test[1] = analysis.StaticAnalysisResult{MostLikelyCategories: [5]string{"XSS", "SQL", "CINJ", "LFI", "RFI"}, SrcPort: 1234}
	test[2] = analysis.StaticAnalysisResult{MostLikelyCategories: [5]string{"PWN", "PWN", "PWN", "SQL", "PWN"}, SrcPort: 2345}
	test[3] = analysis.StaticAnalysisResult{MostLikelyCategories: [5]string{"PRV", "PRV", "PRV", "PRV", "PRV"}, SrcPort: 2222}
	test[4] = analysis.StaticAnalysisResult{MostLikelyCategories: [5]string{"CRYPTO", "CRYPTO", "CRYPTO", "CRYPTO", "CRYPTO"}, SrcPort: 10456}

	for _, res := range test {
		flc.Log(res)
	}

	flc.Close()
}
