package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/jubobs/whocanibe/sites"
)

var client = sites.NewClient()

// var checkAll = sites.UniversalChecker(client, checkers)

type payload struct {
	Name string            `json:"username"`
	Res  map[string]string `json:"status"`
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	username := strings.TrimPrefix(r.URL.Path, "/")
	// res := checkAll(username)

	res2 := make(map[string]string)
	// for name, err := range res {
	// 	if err != nil {
	// 		switch {
	// 		// TODO: fix this
	// 		// case sites.IsInvalidUsernameError(err):
	// 		// 	res2[name] = "invalid"
	// 		// case sites.IsUnavailableUsernameError(err):
	// 		// 	res2[name] = "unavailable"
	// 		default:
	// 			res2[name] = "unknown"
	// 		}
	// 	} else {
	// 		res2[name] = "available"
	// 	}
	// }
	p := payload{
		Name: username,
		Res:  res2,
	}
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err := encoder.Encode(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
