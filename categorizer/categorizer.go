package main

import (
	"categorizer/analysis"
	"categorizer/controllers"
	"categorizer/retrieve"
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
)

func IndexOf[T comparable](collection []T, el T) int {
	for i, x := range collection {
		if x == el {
			return i
		}
	}
	return -1
}

func main() {
	var args []string
	if len(os.Args) > 1 {
		args = os.Args[1:]
	}

	if args[0] == "-h" {
		fmt.Println("Usage: categorizer [-h] [-f] [-r Caronte|Tulip address port] [-a apiKey]")
		fmt.Println("Options:\n-r: Select which source streams must be retrieved from: [Caronte, Tulip] (default Caronte)" +
			"\n\t- Caronte: categorizer -r Caronte address port [-a apiKey]" +
			"\n\t- Tulip: categorizer -r Tulip dbAddress dbPort [-a apiKey]" +
			"\n\t- default: categorizer [-f] [-a apiKey] address port" +
			"\n-f: Save results to a log file" +
			"\n-a: Specify the API key as an argument rather than using an environment variable" +
			"\n-h: Help")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	queue := make(chan retrieve.Result)
	results := make(chan analysis.StaticAnalysisResult)
	exit := make(chan bool)
	var ftc *controllers.FileController
	var rtc *controllers.RetrieverController
	const CHROMASERVER = "http://localhost:8000"

	if IndexOf(args, "-f") != -1 {
		ftc = controllers.NewFileController(ctx, results)
	}

	i := IndexOf(args, "-r")
	if i == -1 {
		if len(args) < 2 {
			fmt.Println("Invalid arguments. Usage: categorizer [-h] [-f] [-r Caronte|Tulip args]")
			return
		}

		match, err := regexp.MatchString("^(?:[0-9]{1,3}\\.){3}[0-9]{1,3}$", args[len(args)-2])
		if err != nil {
			fmt.Printf("Error checking address validity: %v", err)
			return
		}

		if !match && args[len(args)-2] != "localhost" {
			fmt.Println("Invalid address")
			return
		}

		port, err := strconv.ParseUint(args[len(args)-1], 10, 32)
		if err != nil {
			fmt.Printf("Error parsing port: %v", err)
			return
		}

		rtc = controllers.NewRetrieverController(ctx, queue, retrieve.NewCaronteRetriever(args[len(args)-2], uint16(port)))
	} else {
		if len(args) != i+4 {
			fmt.Println("Invalid arguments. Usage: categorizer [-h] [-f] [-r Caronte|Tulip args]")
			return
		}
		if args[i+1] == "Caronte" {
			match, err := regexp.MatchString("^(?:[0-9]{1,3}\\.){3}[0-9]{1,3}$", args[i+2])
			if err != nil {
				fmt.Printf("Error checking address validity: %v", err)
				return
			}

			if !match && args[i+2] != "localhost" {
				fmt.Println("Invalid address")
				return
			}

			port, err := strconv.ParseUint(args[i+3], 10, 32)
			if err != nil {
				fmt.Printf("Error parsing port: %v", err)
				return
			}

			rtc = controllers.NewRetrieverController(ctx, queue, retrieve.NewCaronteRetriever(args[i+2], uint16(port)))
		} else if args[i+1] == "Tulip" { // TODO: Implement TulipRetriever
			/*match, err := regexp.MatchString("^(?:[0-9]{1,3}\\.){3}[0-9]{1,3}$", args[i+2])
			if err != nil {
				fmt.Printf("Error checking address validity: %v", err)
				return
			}

			if !match && args[i+2] != "localhost" {
				fmt.Println("Invalid address")
				return
			}

			port, err := strconv.ParseUint(args[i+3], 10, 32)
			if err != nil {
				fmt.Printf("Error parsing port: %v", err)
				return
			}

			rtc = controllers.NewRetrieverController(ctx, queue, retrieve.NewTulipRetriever(args[i+2], uint16(port)))*/
			return
		} else {
			fmt.Println("Invalid retriever")
			cancel()
			return
		}
	}

	i = IndexOf(args, "-a")
	if i != -1 {
		if len(args) >= i+2 {
			err := os.Setenv("HF_API_KEY", args[i+1])
			if err != nil {
				fmt.Printf("Error setting API key: %v", err)
				return
			}
		} else {
			fmt.Println("Invalid API key")
			return
		}
	}

	// Instantiate the ChromaAnalyser. The second parameter is the address of the Chroma server which must be running
	chr, err := analysis.NewChromaAnalyser(ctx, CHROMASERVER)
	if err != nil {
		fmt.Printf("Error creating ChromaAnalyser: %v", err)
		return
	}

	anc := controllers.NewAnalysisController(ctx, queue, chr)
	log := controllers.NewLogger(ctx, results)

	var stop string
	fmt.Println("Enter ^D to stop")

	go rtc.Start(exit, cancel)
	go anc.Start(exit, cancel)
	if ftc != nil {
		go ftc.Start(exit, cancel)
	}
	go log.Start(exit)

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
