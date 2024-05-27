package main

import (
	dbase "MrScraper/db"
)

func main() {

	db, err := dbase.NewDatabase()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//scr := scrapers.NewSpringerOpenScraper()
	//_, err := scr.Scrap("microservices")
	//if err != nil {
	//	log.Printf("ОШИБКА: ", err)
	//	return
	//}

}
