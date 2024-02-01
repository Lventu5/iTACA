package categorizer

import (
	"categorizer/analysis"
	"categorizer/controllers"
	"categorizer/retrieve"
	"context"
	"fmt"
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
		fmt.Println("Usage: categorizer [-h] [-r] [-f] [args]")
		fmt.Println("Options:\n-r: Select which source streams must be retrieved from: [Caronte, Tulip, pcap, file] (default Caronte)" +
			"\n\t- Caronte: categorizer -r Caronte address port" +
			"\n\t- Tulip: categorizer -r Tulip dbAddress dbPort" +
			"\n\t- pcap: categorizer -r pcap /path/to/dir/containing/pcaps" +
			"\n\t- file: categorizer -r file /path/to/dir/containing/files" +
			"\n\t- default: categorizer [options] address port" +
			"\n-f: Save results to a log file" +
			"\n-h: Help")
		return
	}

	ctx := context.Background()
	queue := make(chan retrieve.Result)
	results := make(chan analysis.StaticAnalysisResult)
	exit := make(chan bool)
	var stc *controllers.StorageController
	var rtc *controllers.RetrieverController

	if IndexOf(args, "-f") != -1 {
		stc = controllers.NewStorageController(ctx, results)
	}

	i := IndexOf(args, "-r")
	if i == -1 {
		if len(args) < 2 {
			fmt.Println("Invalid arguments. Usage: categorizer [-h] [-f] [-r Caronte|Tulip|pcap|file args]")
			return
		}

		match, err := regexp.MatchString("^(?:[0-9]{1,3}\\.){3}[0-9]{1,3}$", args[len(args)-2])
		if err != nil {
			fmt.Printf("Error checking address validity: %v", err)
			return
		}

		if !match {
			fmt.Println("Invalid address")
			return
		}

		port, err := strconv.ParseUint(args[len(args)-1], 10, 32)
		if err != nil {
			fmt.Printf("Error parsing port: %v", err)
			return
		}

		rtc = controllers.NewRetrieverController(ctx, queue, retrieve.NewCaronteRetriever(args[len(args)-2], uint16(port)))
	}
	else {
		if len(args) < i+3 {
			fmt.Println("Invalid arguments. Usage: categorizer [-h] [-f] [-r Caronte|Tulip|pcap|file args]")
			return
		}
	}
}
