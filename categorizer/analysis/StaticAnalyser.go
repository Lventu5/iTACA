package analysis

import (
	"errors"
	chroma "github.com/amikos-tech/chroma-go"
	hf "github.com/amikos-tech/chroma-go/hf"
	"os"
)

// StaticAnalyser : interface, defines the general method "analyse" which is used to analyse a tcp stream
type StaticAnalyser interface {
	Analyse(stream []byte, result chan<- StaticAnalysisResult)
}

type StaticAnalysisResult struct {
	mostLikelyCategory string
	otherCategories    [4]string
}

type ChromaAnalyser struct {
	client     *chroma.Client
	collection *chroma.Collection
}

func NewChromaAnalyser(params ...string) (*ChromaAnalyser, error) {
	if len(params) < 1 || len(params) > 2 {
		err := errors.New("invalid number of parameters")
		return nil, err
	}

	cli := chroma.NewClient(params[0])
	var apiKey string = ""

	if len(params) == 1 {
		apiKey = os.Getenv("HF_API_KEY")
		if apiKey == "" {
			err := errors.New("no api key found in environment variables")
			return nil, err
		}
	} else {
		apiKey = params[1]
	}

	coll, err := cli.GetCollection("payloads", hf.NewHuggingFaceEmbeddingFunction(apiKey, "sentence-transformers/all-MiniLM-L6-v2"))
	if err != nil {
		return nil, err
	}

	return &ChromaAnalyser{client: cli, collection: coll}, nil
}

func (a *ChromaAnalyser) Analyse(stream []byte, result chan<- StaticAnalysisResult) {
	
}
