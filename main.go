package main

import (
	"MrScraper/scrapers"
	"log"
)

func main() {
	scr := scrapers.NewSpringerOpenScraper()
	_, err := scr.Scrap("microservices")
	if err != nil {
		log.Printf("ОШИБКА: ", err)
		return
	}

}
