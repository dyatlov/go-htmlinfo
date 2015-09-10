package main

import (
	"fmt"
	"net/http"

	"github.com/dyatlov/go-htmlinfo/htmlinfo"
)

func main() {
	u := "http://techcrunch.com/2015/09/09/ipad-pro-coming-in-november-pricing-starts-at-799/"

	resp, err := http.Get(u)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	info := htmlinfo.NewHTMLInfo()

	// if url is not provided it's fine too, just then we wont be able to fetch (and generate) oembed information
	err = info.Parse(resp.Body, &u)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Info:\n%s\n", info)

	fmt.Printf("Oembed information: %s\n", info.GenerateOembedFor(u))
}
