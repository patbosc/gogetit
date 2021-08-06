package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/gocolly/colly"
)

//use: curl http://localhost:7272/product?url=https://www.amazon.com/s?k=playstation+5&ref=nb_sb_noss_2

func main() {

	// create a service to listen on port 7272
	addr := ":7272"

	//registerd functions
	http.HandleFunc("/search", getIt)
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/datetime", datetime)
	http.HandleFunc("/product", product)

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

func ping(w http.ResponseWriter, r *http.Request) {
	log.Println("Ping")
	w.Write([]byte("ping"))
}

func datetime(w http.ResponseWriter, r *http.Request) {
	log.Println("datetime")
	w.Write([]byte(time.Now().String()))
}

func product(w http.ResponseWriter, r *http.Request) {
	//Verify URL
	URL := r.URL.Query().Get("url")
	if URL == "" {
		log.Println("missing URL argument")
		return
	}
	log.Println("visiting", URL)

	//Create a new colly collector
	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	var httpResponse []string

	c.OnHTML("div.s-result-list.s-search-results.sg-row", func(e *colly.HTMLElement) {
		e.ForEach("div.a-section.a-spacing-medium", func(_ int, e *colly.HTMLElement) {
			var productName, stars, price string

			productName = e.ChildText("span.a-size-medium.a-color-base.a-text-normal")

			if productName == "" {
				// If we can't get any name, we return and go directly to the next element
				return
			}

			stars = e.ChildText("span.a-icon-alt")
			FormatStars(&stars)

			price = e.ChildText("span.a-price > span.a-offscreen")
			FormatPrice(&price)

			prods := fmt.Sprintf("Product Name: %s \nStars: %s \nPrice: %s \n", productName, stars, price)

			httpResponse = append(httpResponse, prods)

			//now parse our response to JSON

			b, err := json.Marshal(httpResponse)
			if err != nil {
				log.Println("failed to serialize response:", err)
				return
			}
			//Add some header and write the body
			w.Header().Add("Content-Type", "application/json")
			w.Write(b)
		})
	})

	c.Visit(URL)

}

func FormatPrice(price *string) {
	r := regexp.MustCompile(`\$(\d+(\.\d+)?).*$`)

	newPrices := r.FindStringSubmatch(*price)

	if len(newPrices) > 1 {
		*price = newPrices[1]
	} else {
		*price = "Unknown"
	}

}

func FormatStars(stars *string) {
	if len(*stars) >= 3 {
		*stars = (*stars)[0:3]
	} else {
		*stars = "Unknown"
	}
}
