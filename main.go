package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	nexusURL, httpPort := parseArgs()

	http.HandleFunc("/badge", badgeHandler(nexusURL))

	port := strconv.Itoa(httpPort)
	fmt.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func parseArgs() (string, int) {
	nexusURL := flag.String("url", "", "Nexus url, including scheme and nexus base api path")
	httpPort := flag.Int("port", 8080, "Http port")

	flag.Parse()

	if *nexusURL == "" {
		flag.Usage()
		os.Exit(1)
	}
	return *nexusURL, *httpPort
}
