package retrieve

import (
	"bytes"
	"categorizer/storage"
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"slices"
)

// TulipRetriever : implements Retriever for Tulip, fetching reconstructed TCP streams from Tulip exposed API
// address and port must be the address of the machine hosting the service and the port exposed by the service for API interactions
type TulipRetriever struct {
	address string
	port    uint16
}

// NewTulipRetriever : creates a new instance of TulipRetriever
func NewTulipRetriever(address string, port uint16) *TulipRetriever {
	return &TulipRetriever{address: address, port: port}
}

// Retrieve : retrieves tcp streams from a TulipDB database
func (r *TulipRetriever) Retrieve(ctx context.Context, cancel context.CancelFunc, results chan<- Result) {
	var visited []primitive.ObjectID

	for {
		select {
		case <-ctx.Done():
			return
		default:
			addr := fmt.Sprintf("http://%s:%d/api/query", r.address, r.port)
			var data = []byte(`{}`)

			req, err := http.NewRequest(http.MethodPost, addr, bytes.NewBuffer(data))
			if err != nil {
				fmt.Printf("client: could not create request: %s\n", err)
				cancel()
				return
			}

			req.Header.Add("Accept", "*/*")
			req.Header.Add("Accept-Language", "en-GB,en;q=0.9,it-IT;q=0.8,it;q=0.7,en-US;q=0.6")
			req.Header.Add("Cache-Control", "no-cache")
			req.Header.Add("Connection", "keep-alive")
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("DNT", "1")
			req.Header.Add("Pragma", "no-cache")
			req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Printf("client: error making http request: %s\n", err)
				cancel()
				return
			}

			var flows []storage.FlowEntry
			err = json.NewDecoder(res.Body).Decode(&flows)
			if err != nil {
				fmt.Printf("client: error decoding response body: %s\n", err)
				continue
			}

			var IDs = map[primitive.ObjectID]uint16{}
			for _, flow := range flows {
				IDs[flow.Id] = uint16(flow.Dst_port)
			}

			for id, port := range IDs {
				if slices.Contains(visited, id) {
					continue
				}

				visited = append(visited, id)

				req.URL.Path = fmt.Sprintf("/api/flow/%v", id.Hex())
				req.Method = http.MethodGet
				req.ContentLength = 0

				res, err = http.DefaultClient.Do(req)
				if err != nil {
					fmt.Printf("client: error making http request: %s\n", err)
					cancel()
					return
				}

				var singleFlow storage.FlowEntry
				err = json.NewDecoder(res.Body).Decode(&singleFlow)
				if err != nil {
					fmt.Printf("client: error decoding response body: %s\n", err)
					continue
				}

				reconstructedStream := ""

				for _, flow := range singleFlow.Flow {
					if flow.From == "s" {
						continue
					}
					// fmt.Printf(" Retrieved stream from port %d\n%s\n", port, v.Content)
					reconstructedStream += flow.Data
				}

				results <- Result{Stream: reconstructedStream, SrcPort: uint16(port)}
			}
		}
	}
}
