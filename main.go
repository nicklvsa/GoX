package main

import (
	"fmt"
	"gox/shared"
)

// TODO: Figure out reflect crashing issue when trying to marshal in utils.go

func main() {

	/*
		Examples of how to use the GoX cache library
	*/

	gox := &shared.GoxModule{} // Optionally configure the GoxCache object within the GoxModule (no need for .Init() then)
	// gox.Init("cache-id-000", "testing cache") // give the cache an id and name
	gox.InitWithSync("cache-id-000", "testing cache", "your-api-key")

	// gox.SetItem("hello", "pie", -1) // set an item within the cache

	val, err := gox.GetItem("hello", false) // get the item
	if err != nil {
		fmt.Println(fmt.Sprintf("Error while attempting to retrieve the item! Error: %s", err.Error()))
	}

	fmt.Println(fmt.Sprintf(`"hello" => "%s"`, val))

	items, count := gox.PurgeExpiredItems() // remove all expired items
	fmt.Println(fmt.Sprintf("Expired items: %+v | %d", items, count))

	select {} // testing - block forever
}
