package main

import (
	"MrScraper/scrapers"
	"log"
)

func main() {
	scr := scrapers.NewHabrScraper()
	_, err := scr.Scrap("microservices")
	if err != nil {
		log.Printf("ОШИБКА: ", err)
		return
	}

}
