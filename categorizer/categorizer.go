package main

import (
	"categorizer/analysis"
	"categorizer/config"
	"categorizer/controllers"
	"categorizer/retrieve"
	"context"
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		args := os.Args[1:]

		if args[0] == "-h" {
			fmt.Println("Usage: ")
			fmt.Println("\nOptions:\nEdit the config.json file located in the config folder." +
				"\n- Retriever: \n\tallowed types: Caronte, Tulip (default Caronte)\n\taddress, port: address and port where the Caronte instance is running (usually localhost, 3333)" +
				"\n- Analyser: \n\tallowed types: ChromaDB\n\taddress, port: address and port where the ChromaDB server is running (usually localhost, 8000)\n\tapikey: insert your HF api key" +
				"\n- Log: true/false (default true)\n\tselect true if you want a persistent output over a log file")
			return
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	queue := make(chan retrieve.Result)
	results := make(chan analysis.StaticAnalysisResult, 10)
	exit := make(chan bool)
	var rtc *controllers.RetrieverController

	// Parse the config file
	cfg, err := config.ParseConfig("config/config.json")
	if err != nil {
		fmt.Printf("Error parsing config: %v", err)
		return
	}

	// Instantiate the RetrieverController
	switch cfg.Retriever.Type {
	case "Caronte":
		rtc = controllers.NewRetrieverController(ctx, queue, retrieve.NewCaronteRetriever(cfg.Retriever.Host, cfg.Retriever.Port))
	case "Tulip":
		rtc = controllers.NewRetrieverController(ctx, queue, retrieve.NewTulipRetriever(cfg.Retriever.Host, cfg.Retriever.Port))
	default:
		rtc = controllers.NewRetrieverController(ctx, queue, retrieve.NewCaronteRetriever(cfg.Retriever.Host, cfg.Retriever.Port))
	}

	otc := controllers.NewOutputController(ctx, results, cfg.Log)

	// Instantiate the ChromaAnalyser
	chr, err := analysis.NewChromaAnalyser(ctx, cfg.Analyser.Host, cfg.Analyser.Port, cfg.Analyser.Collection, cfg.Analyser.ApiKey)

	anc := controllers.NewAnalysisController(ctx, queue, chr)

	var stop string
	fmt.Println("Enter ^D to stop")

	go rtc.Start(exit, cancel)
	go anc.Start(exit, cancel)
	go otc.Start(exit, cancel)

	for {
		_, err := fmt.Scanln(&stop)
		if err == io.EOF {
			exit <- true
		}

		select {
		case <-ctx.Done():
			close(queue)
			close(results)
			close(exit)
			return
		}
	}
}
