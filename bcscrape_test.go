package bcscrape_test

import (
	"fmt"

	"github.com/thraxil/bcscrape"
)

func ExampleAlbumPage() {
	a := bcscrape.NewAlbumPage("https://solstafir.bandcamp.com/album/berdreyminn")
	a.Fetch()
	fmt.Println(a.Title)
	fmt.Println(a.Artist)
	// Output:
	// Berdreyminn
	// Sólstafir
}

func ExampleDeterminePage_album() {
	fmt.Println(bcscrape.DetermineType("https://solstafir.bandcamp.com/album/berdreyminn"))
	// Output: album https://solstafir.bandcamp.com/album/berdreyminn
}

func ExampleDeterminePage_releases() {
	fmt.Println(bcscrape.DetermineType("https://solstafir.bandcamp.com/releases"))
	// Output: album https://solstafir.bandcamp.com/album/berdreyminn
}

func ExampleDeterminePage_artist() {
	fmt.Println(bcscrape.DetermineType("https://solstafir.bandcamp.com/"))
	// Output: band https://solstafir.bandcamp.com
}
