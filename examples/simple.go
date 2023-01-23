package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dyatlov/go-htmlinfo/htmlinfo"
)

func main() {
	ctx := context.Background()
	u := "http://techcrunch.com/2010/11/02/365-days-10-million-3-rounds-2-companies-all-with-5-magic-slides/"

	resp, err := http.Get(u)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	info := htmlinfo.NewHTMLInfo()
	info.AllowOembedFetching = true

	ct := resp.Header.Get("Content-Type")

	// if url and contentType are not provided it's fine too, just then we wont be able to fetch (and generate) oembed information
	err = info.ParseWithContext(ctx, resp.Body, &u, &ct)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Info:\n%s\n", info)

	fmt.Printf("Oembed information: %s\n", info.GenerateOembedFor(u))
}
