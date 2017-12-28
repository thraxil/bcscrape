// Package bcscrape implements a simple scraper for bandcamp.com pages
package bcscrape

import (
	"net/url"
	"strings"

	"github.com/gocolly/colly"
)

// Page represents an unknown type of bandcamp page
type Page struct {
	URL          string
	Type         string
	NormativeURL string
}

// NewPage is a Page constructor
func NewPage(url string) *Page {
	return &Page{URL: url}
}

func trim(s string) string {
	return strings.Trim(s, " \t\n\r")
}

// Fetch makes the HTTP request and populates the Type and NormativeURL fields
func (p *Page) Fetch() {
	c := colly.NewCollector()

	c.OnHTML("meta[property=\"og:type\"]", func(e *colly.HTMLElement) {
		p.Type = e.Attr("content")
	})

	c.OnHTML("meta[property=\"og:url\"]", func(e *colly.HTMLElement) {
		p.NormativeURL = e.Attr("content")
	})

	c.Visit(p.URL)
}

// TrackPage represents a track page
type TrackPage struct {
	URL       string
	Title     string
	Artist    string
	ArtistURL string
	Album     string
	AlbumURL  string
	CoverSRC  string
	Published string
}

// NewTrackPage is the constructor for TrackPage
func NewTrackPage(url string) *TrackPage {
	return &TrackPage{URL: url}
}

// Fetch makes the HTTP request and populates the data fields
func (t *TrackPage) Fetch() {
	c := colly.NewCollector()

	c.OnHTML("span[itemprop=\"byArtist\"] a", func(e *colly.HTMLElement) {
		t.Artist = trim(e.Text)
		t.ArtistURL = trim(e.Attr("href"))
	})

	c.OnHTML("span[itemprop=\"inAlbum\"] a", func(e *colly.HTMLElement) {
		t.Album = trim(e.Text)
		t.AlbumURL = trim(e.Attr("href"))
		if strings.HasPrefix(t.AlbumURL, "/") {
			// deal with relative URL
			u, _ := url.Parse(t.URL)
			u.Path = t.AlbumURL
			t.AlbumURL = u.String()
		}
	})

	c.OnHTML("h2[itemprop=\"name\"]", func(e *colly.HTMLElement) {
		t.Title = trim(e.Text)
	})

	c.OnHTML("#tralbumArt img[itemprop=\"image\"]", func(e *colly.HTMLElement) {
		t.CoverSRC = e.Attr("src")
	})

	c.OnHTML("meta[itemprop=\"datePublished\"]", func(e *colly.HTMLElement) {
		// TODO: parse into actual date struct
		t.Published = e.Attr("content")
	})

	// Tags

	c.Visit(t.URL)
}

// AlbumPage represents an album page
type AlbumPage struct {
	URL         string
	Artist      string
	ArtistURL   string
	Title       string
	Description string
	CoverSRC    string
	Published   string
}

// NewAlbumPage is the constructor for AlbumPage
func NewAlbumPage(url string) *AlbumPage {
	return &AlbumPage{URL: url}
}

// Fetch makes the HTTP request and populates the data fields
func (a *AlbumPage) Fetch() {
	c := colly.NewCollector()
	c.OnHTML("span[itemprop=\"byArtist\"] a", func(e *colly.HTMLElement) {
		a.Artist = trim(e.Text)
		a.ArtistURL = e.Attr("href")
	})

	c.OnHTML("h2[itemprop=\"name\"]", func(e *colly.HTMLElement) {
		a.Title = trim(e.Text)
	})

	c.OnHTML("#tralbumArt img[itemprop=\"image\"]", func(e *colly.HTMLElement) {
		a.CoverSRC = e.Attr("src")
	})

	c.OnHTML("div.tralbumData[itemprop=\"description\"]", func(e *colly.HTMLElement) {
		a.Description = trim(e.Text)
	})

	c.OnHTML("meta[itemprop=\"datePublished\"]", func(e *colly.HTMLElement) {
		// TODO: parse into actual date struct
		a.Published = e.Attr("content")
	})

	// Tracks
	// Tags

	c.Visit(a.URL)
}

// DetermineType takes a url and determines what type of page it refers to
// and also returns a normalized URL (eg, an artist '/releases' page often
// contains the data of the most recently released album page, so it normalizes
// to that album page's URL)
func DetermineType(url string) (string, string) {
	if strings.Contains(url, "bandcamp.com/album/") {
		return "album", url
	}
	if strings.Contains(url, "bandcamp.com/track/") {
		return "track", url
	}
	p := NewPage(url)
	p.Fetch()
	return p.Type, p.NormativeURL
}
