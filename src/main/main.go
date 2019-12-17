package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"

	"../webhandler"
	"github.com/gorilla/mux"
)

var (
	address          = flag.String("addr", "0.0.0.0:8080", "Defines the server address")
	baseDirLogs      = flag.String("logs", "C:/Users/A64573642/DEV/workspace/irrp-testautomation/log/", "Directory which is queried for regression-logfiles")
	baseDirTestcases = flag.String("testcases", "C:/Users/A64573642/DEV/workspace/irrp-testautomation/testcases/", "Directoy which is queried for testcases")
)

func main() {
	flag.Parse()

	proctrl.RegressionLogDir = filepath.Clean(*baseDirLogs)
	proctrl.TestcaseDir = filepath.Clean(*baseDirTestcases)

	router := mux.NewRouter().StrictSlash(true)
	// lock handler
	router.HandleFunc("/api/lock", webhandler.LockHandler)

	// accepts post requests to trigger the tests
	router.HandleFunc("/api/post", webhandler.POSTHandler)
	// serves static logfiles in baseDirLogs
	router.PathPrefix("/logs").Handler(http.StripPrefix("/logs", http.FileServer(http.Dir(*baseDirLogs))))

	log.Printf("Listing on %s\n---", *address)
	if err := http.ListenAndServe(*address, router); err != nil {
		log.Fatal("Error starting server")
	}
}
