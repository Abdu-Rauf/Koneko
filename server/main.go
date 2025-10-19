package main

import (
	"fmt"
	// "log"
	"net/http"

	"github.com/Abdu-Rauf/koneko/test"
)

func myhandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "hello world %v", r.URL)
}

func main() {
	// log.Println("Server starting on port 8080")
	// http.HandleFunc("/ruffu", myhandler)
	// log.Fatal(http.ListenAndServe(":8080", nil))

	test.Create_loadpg()
	test.SliceTime()
}
