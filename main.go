package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gocolly/colly"
)

func ping(w http.ResponseWriter, r *http.Request) {
	log.Println("Ping")
	w.Write([]byte("ping"))
}

func main() {

	// create a service to listen on port 7272
	addr := ":7272"

	http.HandleFunc("/search", getIt)
	http.HandleFunc("/ping", ping)

	log.Println("listening on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func getIt(w http.ResponseWriter, r *http.Request) {
	//Verify URL
	URL := r.URL.Query().Get("url")
	if URL == "" {
		log.Println("missing URL argument")
		return
	}
	log.Println("visiting", URL)

	//Create a new colly collector
	c := colly.NewCollector()

	var httpResponse []string

	//check if the website html has any links if so get it
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		if link != "" {
			httpResponse = append(httpResponse, link)
		}
	})

	//run clever boy, run
	c.Visit(URL)

	//now parse our response to JSON

	b, err := json.Marshal(httpResponse)
	if err != nil {
		log.Println("failed to serialize response:", err)
		return
	}
	//Add some header and write the body
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}
