package tests

import (
	"categorizer/retrieve"
	"context"
	"testing"
)

func TestRetrieverController(t *testing.T) {
	ctx := context.Background()
	queue := make(chan retrieve.Result)
	retriever := retrieve.NewCaronteRetriever("0.0.0.0", 3333)
}
