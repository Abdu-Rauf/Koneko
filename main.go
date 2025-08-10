package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {

	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func loadpage(title string) (*Page, error) {

	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}
func create_loadpg() {
	p1 := &Page{Title: "pointers", Body: []byte("understanding referencing and pointers")}
	p1.save()
	p2, _ := loadpage(p1.Title)
	fmt.Println(string(p2.Body))
}

func myhandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "hello i know you are here niqqa %v", r.URL)
}

func main() {
	log.Println("Server starting on port 8080")
	http.HandleFunc("/ruffu", myhandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
