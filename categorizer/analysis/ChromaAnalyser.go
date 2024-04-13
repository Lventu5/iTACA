package analysis

import (
	"categorizer/retrieve"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

/*	IN CASE CHROMA-GO GETS FIXED
type ChromaAnalyser struct {
	client     *chroma.Client
	collection *chroma.Collection
}

func NewChromaAnalyser(ctx context.Context, address string, port uint16, collection string, apiKey string) (*ChromaAnalyser, error) {

	finalAddress := fmt.Sprintf("http://%s:%d", address, port)

	cli := chroma.NewClient(finalAddress)

	coll, err := cli.GetCollection(ctx, collection, hf.NewHuggingFaceEmbeddingFunction(apiKey, "sentence-transformers/all-MiniLM-L6-v2"))
	if err != nil {
		return nil, err
	}

	return &ChromaAnalyser{client: cli, collection: coll}, nil
}
*/

type ChromaAnalyser struct {
	address    string
	port       uint16
	collection string
}

func NewChromaAnalyser(address string, port uint16, collection string) (*ChromaAnalyser, error) {
	r, err := regexp.Compile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`)
	if err != nil {
		return nil, err
	}
	if address != "localhost" && !r.MatchString(address) {
		return nil, fmt.Errorf("Invalid address")
	}
	return &ChromaAnalyser{address: address, port: port, collection: collection}, nil
}

func (a *ChromaAnalyser) Analyse(stream retrieve.Result, result chan<- StaticAnalysisResult) {
	/* IN CASE CHROMA-GO GETS FIXED
	qr, err := a.collection.Query(ctx, []string{stream.Stream}, 5, nil, nil, nil)
	if err != nil {
		fmt.Printf("Error querying: %v\n", err)
		cancel()
		return
	}

	var res StaticAnalysisResult
	for i, id := range qr.Ids[0] {
		if qr.Distances[0][i] > 1.0 {
			id = "SAFE"
		}
		res.MostLikelyCategories[i] = id
		res.SrcPort = stream.SrcPort
	}*/

	var res StaticAnalysisResult

	cmd := exec.Command("../analysis/QueryHandler.py", a.address, fmt.Sprintf("%d", a.port), a.collection, stream.Stream)
	qrRes, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error executing command: %v", err)
		return
	}
	strRes := strings.TrimSpace(string(qrRes))
	if strRes == "" {
		return
	}

	strRes = strings.ReplaceAll(strRes, "'", "")
	strRes = strings.ReplaceAll(strRes, "[", "")
	strRes = strings.ReplaceAll(strRes, "]", "")

	res.MostLikelyCategories = [5]string(strings.Split(strRes, ", "))
	res.SrcPort = stream.SrcPort

	result <- res
	return
}
