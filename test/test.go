package test

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
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

func ArrayandSlices() {
	//array
	arr := [5]int32{1, 2, 3, 4}
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

	mymap := make(map[string]int, 10)
	mymap[`Adam`] = 25
	mymap[`Eve`] = 22
	mymap[`Charlie`] = 30
	// fmt.Println(mymap)
	// delete(mymap, `Eve`)
	fmt.Println(mymap)
	for name, age := range mymap {
		fmt.Printf("Name:%v Age:%v\n", name, age)
	}
}

func Stringops() {
	st := "résumé"
	// st = [114,233,115,117,109,233] // variable length encoding through utf-8
	fmt.Printf("%v\n", st[0]) // utf encoding of r is printed
	for i, v := range st {
		fmt.Printf("%v,%v\n", i, v) // index is skipped for multibyte characters
	}
	fmt.Println("bytes in string:", len(st))

	st1 := []rune(st) // converting string to rune slice to get actual characters
	for i, v := range st1 {
		fmt.Printf("%v,%v\n", i, v) // index is not skipped for multibyte characters
	}
	fmt.Println("runes in string:", len(st1))

	var strBuilder strings.Builder
	strslice := []string{"r", "e", "s", "u", "m", "e"}
	// str := ""
	for i := range strslice {
		strBuilder.WriteString(strslice[i])
	}
	str := strBuilder.String()
	fmt.Println("constructed string:", str)
}

// method cant be defined inside a function, defining outside

type gasengine struct {
	mpg     uint16
	gallons uint16
}

type electricengine struct {
	kwatt           uint16
	batterycapacity uint16
}

func (g gasengine) milesleft() uint16 {
	return g.mpg * g.gallons
}

func (e electricengine) milesleft() uint16 {
	return e.kwatt * e.batterycapacity
}

type engine interface {
	milesleft() uint16
}

func distancePossble(e engine, miles uint16) {
	if miles <= e.milesleft() {
		fmt.Println("can reach destination")
	} else {
		fmt.Println("cannot reach destination")
	}
}

func Interfaceops() {
	var gasCar gasengine = gasengine{mpg: 30, gallons: 10}
	var electricCar electricengine = electricengine{kwatt: 5, batterycapacity: 50}
	fmt.Printf("Gas car can go %v miles\n", gasCar.milesleft())
	fmt.Printf("Electric car can go %v miles\n", electricCar.milesleft())
	distancePossble(gasCar, 300)
	distancePossble(electricCar, 300)

}

func Sqaurevals(vals *[5]int64) [5]int64 {
	fmt.Printf("The memory address of slice in function is %p\n", vals)
	for i := range vals {
		vals[i] = vals[i] * vals[i]
	}
	return *vals
}
func Malmain() {

	var slice = [5]int64{1, 2, 3, 4, 5}
	fmt.Printf("The memory address of original slice is %p\n", &slice)
	var result = Sqaurevals(&slice)
	fmt.Println(result)
}

func Task1() {
    time.Sleep(1 * time.Second)
    fmt.Println("Task 1")
}

func Task2() {
    time.Sleep(1 * time.Second)
    fmt.Println("Task 2")
}

