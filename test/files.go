package test

import (
	"fmt"
	"io/fs"
	"os"
	"time"
)

// Practicing creating and reading files in Go

type Page struct {
	Title string
	Body  []byte
}

// specify the directory to use as the file system
var myFS = os.DirFS(".")

func (p *Page) save() error {

	filename := p.Title + ".txt"
	// 0600 means read and write permissions for the owner only (file permission in octal format)
	return os.WriteFile(filename, p.Body, 0600)
}

func Loadpage(title string) (*Page, error) {

	filename := title + ".txt"
	// read the file from the specified file system
	body, err := fs.ReadFile(myFS, filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}
func Create_loadpg() {
	p1 := &Page{Title: "pointers", Body: []byte("understanding referencing and pointers")}
	p1.save()
	p2, _ := Loadpage(p1.Title)
	fmt.Println(string(p2.Body))
}

// Practing data structures in Go

func Stringops() {
	var myst string = "hello world"
	fmt.Println(`string before`, myst)
	myst = `changed string`
	fmt.Println(`string after`, myst)

}

func ArrayandSlices() {
	//array
	var arr [5]int32
	arr = [5]int32{1, 2, 3, 4}
	fmt.Println("array", arr)
	arr[4] = 5
	fmt.Println("array after change", arr)

	// converting array to slice
	sliceconvert := arr[:]
	fmt.Println("length and capacity of slice before converting", len(sliceconvert), cap(sliceconvert))
	sliceconvert = append(sliceconvert, 6)
	fmt.Println("length and capacity of slice after converting", len(sliceconvert), cap(sliceconvert))
	fmt.Println("slice after", sliceconvert)

}

func SliceTime() {

	n := 1000000
	// defining slice without specifying si
	var t0 = time.Now()
	s1 := []int{}
	for len(s1) < n {
		s1 = append(s1, 1)
	}
	fmt.Println("Time taken without specifying size", time.Since(t0))

	// defining slice with specifying size
	var t1 = time.Now()
	s2 := make([]int, n)
	for len(s2) < n {
		s2 = append(s2, 1)
	}
	fmt.Println("Time taken wiht specifying time", time.Since(t1))

}

func Mapops() {

	// var mymap map[string]int
	// mymap[`Adam`]=25
	// mymap[`Eve`]=22
	// fmt.Println(mymap)

	mymap := make(map[string]int)
	mymap[`Adam`] = 25
	mymap[`Eve`] = 22
	fmt.Println(mymap)
	delete(mymap, `Eve`)
	fmt.Println(mymap)

}
