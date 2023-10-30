package retrieve

import (
	c "categorizer"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"slices"
	"sync"
)

// Retriever : interface which defines the general method "retrieve", which is used to retrieve tcp streams from various sources
type Retriever interface {
	Retrieve(ctx context.Context, wg sync.WaitGroup, results chan<- string)
}

// CaronteRetriever : implements Retriever for Caronte service, fetching reconstructed TCP streams from a service exposing a REST API
// address must be the address of the machine which is hosting the service, port must be the port which exposes the service
type CaronteRetriever struct {
	address string
	port    uint16
}

func (r *CaronteRetriever) Retrieve(ctx context.Context, wg sync.WaitGroup, results chan<- string) {
	var visited []c.RowID

	for {
		case<- ctx.Done():
			fmt.Println("Retriever: task stopped")
			exit(0)
		default:
			addr := fmt.Sprintf("%s:%d/api/connections?limit=50", r.address, r.port)
			req, err := http.NewRequest(http.MethodGet, addr, nil)
			if err != nil {
				fmt.Printf("client: could not create request: %s\n", err)
				os.Exit(1)
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
				os.Exit(1)
			}

			var connection []c.Connection
			err = json.NewDecoder(res.Body).Decode(&connection)
			if err != nil {
				fmt.Printf("client: error decoding response body: %s\n", err)
				os.Exit(1)
			}

			var Ids []c.RowID
			for _, conn := range connection {
				Ids = append(Ids, conn.ID)
			}

			for _, id := range Ids {
				if slices.Contains(visited, id) {
					continue
				}

				visited = append(visited, id)
				req.URL.Path = fmt.Sprintf("/api/streams/%s/download?format=default", id)

				res, err = http.DefaultClient.Do(req)
				if err != nil {
					fmt.Printf("client: error making http request: %s\n", err)
					os.Exit(1)
				}

				resBody, err := ioutil.ReadAll(res.Body)
				if err != nil {
					fmt.Printf("client: error reading response body: %s\n", err)
					os.Exit(1)
				}
				results <- string(resBody)
		}
	}

}
