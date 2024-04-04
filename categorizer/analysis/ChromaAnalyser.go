package analysis

import (
	"categorizer/retrieve"
	"context"
	"fmt"
	"github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/hf"
)

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

func (a *ChromaAnalyser) Analyse(ctx context.Context, cancel context.CancelFunc, stream retrieve.Result, result chan<- StaticAnalysisResult) {
	qr, err := a.collection.Query(ctx, []string{stream.Stream}, 5, nil, nil, nil)
	if err != nil {
		fmt.Printf("Error querying: %v\n", err)
		cancel()
		return
	}

	var res StaticAnalysisResult
	for i, id := range qr.Ids[0] {
		res.MostLikelyCategories[i] = id
		res.SrcPort = stream.SrcPort
	}

	result <- res
	return
}
