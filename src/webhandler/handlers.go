package webhandler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"../ctrl"
)

type event struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	Task  string `json:"task"`
}

type ResponseWrapper struct {
	ResponseList []Response
}

type Response struct {
	Type    string          `json:"type"`
	Title   string          `json:"title"`
	Task    string          `json:"task"`
	Results []ctrl.Feedback `json:"results"`
}

var (
	authMap = map[string]string{
		"tester1": "7abcddbb2c74e4c0789c2c0aa6abcf5172e82e9f4916bc6409fc3989ed673e08",
		"tester2": "7cd477192d54ceb8673be093f443b8622c612896880f6879c7f8ec16fa7ba114",
	}
	lockOwner string
)

func setLoggerPrefix(r *http.Request) {
	if src := r.Header.Get("source"); src != "" {
		log.SetPrefix(fmt.Sprintf("%s ", src))
	}

}

func LockHandler(w http.ResponseWriter, r *http.Request) {

	if err := checkRequestMethod(r, http.MethodGet); err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusUnauthorized)
		log.Print(err)
		return
	}
	err := checkAuthLock(r)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusUnauthorized)
		log.Print(err)
		return
	}
	w.WriteHeader(200)
}

//POSTHandler is the superduper function
func POSTHandler(w http.ResponseWriter, r *http.Request) {

	defer log.Print("---")
	setLoggerPrefix(r)

	if err := checkRequestMethod(r, http.MethodPost); err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusUnauthorized)
		log.Print(err)
		return
	}

	if err := checkAuth(r); err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusUnauthorized)
		log.Print(err)
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request-body: JSON-object expected", http.StatusBadRequest)
		return
	}
	// check if the sent JSON object has the right fields and content
	newEvent, err := validateJSON(reqBody)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		log.Print(err)
		return
	}
	log.Printf("Got event: %v", newEvent)
	jsonObj, err := newEvent.start()

	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonObj)
}
