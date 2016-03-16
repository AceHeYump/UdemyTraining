package main

import (
	"bufio"
	"fmt"
	"github.com/pkg/browser"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	//"time"
)

//Set URL parts
var prefix = "http://stockcharts.com/def/servlet/SharpChartv05.ServletDriver?chart="
var suffix = ",pypawanrbo[pa][d][f1!3!2!!2!20]&pnf=y"
var filepath = "c:\\temp\\test\\"
var fileext = ".png"

var wg sync.WaitGroup

func worker(symbols <-chan string, output chan<- string) {
	for symbol := range symbols {
		//Get PnF image
		res, err := http.Get(prefix + symbol + suffix)
		if err != nil {
			log.Fatal(err)
		}

		//Read image data
		page, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		//Create image file
		ioutil.WriteFile(filepath+symbol+fileext, page, 644)
		res.Body.Close()

		output <- symbol + " file saved."
	}
	wg.Done()
}

func main() {
	//Create channels
	symbols := make(chan string)
	output := make(chan string)

	//Stage some workers
	for w := 1; w <= 10; w++ {
		wg.Add(1)
		go worker(symbols, output)
	}

	//Get symbol list from file
	f, _ := os.Open("C:\\temp\\test\\symbols.txt")
	defer f.Close()

	//Read file lines into slice
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	var symlist []string

	for scanner.Scan() {
		symlist = append(symlist, scanner.Text())
	}

	//Process each symbol
	go func() {
		for _, v := range symlist {
			symbols <- v
		}
		close(symbols)
	}()

	//Wait for workers to complete
	go func() {
		wg.Wait()
		close(output)
	}()

	//Spit out status as files complete
	for i := range output {
		fmt.Println(i)
	}

	fmt.Println("Done...")

	//...and now the for the web stuff...

	fns := template.FuncMap{
		"start": func(x int) bool {
			if (x+1)%3 == 1 {
				return true
			}
			return false
		},
		"finish": func(x, y int) bool {
			if (y+1) == x || (y+1)%3 == 0 {
				return true
			}
			return false
		},
	}

	//Parse template
	tmpl := template.Must(template.New("main").Funcs(fns).ParseGlob("*.gohtml"))

	//Function
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("c:\\temp\\test\\"))))
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		// execute template
		tmpl.ExecuteTemplate(res, "tmpl.gohtml", symlist)
	})

	//Call it...
	go browser.OpenURL("http://localhost:9002")

	// create server
	http.ListenAndServe(":9002", nil)
}
