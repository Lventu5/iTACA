package analysis

import (
	"categorizer/retrieve"
	"context"
	"errors"
	"fmt"
	chroma "github.com/amikos-tech/chroma-go"
	hf "github.com/amikos-tech/chroma-go/hf"
	"os"
)

// StaticAnalyser : interface, defines the general method "analyse" which is used to analyse a tcp stream
type StaticAnalyser interface {
	Analyse(ctx context.Context, stream retrieve.Result, result chan<- StaticAnalysisResult)
}

type StaticAnalysisResult struct {
	MostLikelyCategories [5]string
	SrcPort              uint16
}
type ChromaAnalyser struct {
	client     *chroma.Client
	collection *chroma.Collection
}

func NewChromaAnalyser(ctx context.Context, params ...string) (*ChromaAnalyser, error) {
	if len(params) < 1 || len(params) > 2 {
		err := errors.New("invalid number of parameters")
		return nil, err
	}

	cli := chroma.NewClient(params[0])
	apiKey := ""

	if len(params) == 1 {
		apiKey = os.Getenv("HF_API_KEY")
		if apiKey == "" {
			err := errors.New("no api key found in environment variables")
			return nil, err
		}
	} else {
		apiKey = params[1]
	}

	coll, err := cli.GetCollection(ctx, "payloads", hf.NewHuggingFaceEmbeddingFunction(apiKey, "sentence-transformers/all-MiniLM-L6-v2"))
	if err != nil {
		return nil, err
	}

	return &ChromaAnalyser{client: cli, collection: coll}, nil
}

func (a *ChromaAnalyser) Analyse(ctx context.Context, stream retrieve.Result, result chan<- StaticAnalysisResult) {
	qr, err := a.collection.Query(ctx, []string{stream.Stream}, 5, nil, nil, nil)
	if err != nil {
		fmt.Printf("Error querying: %v\n", err)
		ctx.Done()
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
