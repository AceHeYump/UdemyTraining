package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	res, err := http.Get("http://stockcharts.com/def/servlet/SharpChartv05.ServletDriver?c=%24SILVER,PLTCDANRBO[PA][D][F1!3!!!4!20]&r=9926&pnf=y")
	if err != nil {
		log.Fatal(err)
	}

	page, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	ioutil.WriteFile("c:\\temp\\silver.gif", page, 644)
	res.Body.Close()
}
