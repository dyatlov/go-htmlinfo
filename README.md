Go HTML Info
===

Go HTML Info provides a simple interface to extract meaningful information from an html page.

source docs: http://godoc.org/github.com/dyatlov/go-htmlinfo/htmlinfo

Install: `go get github.com/dyatlov/go-htmlinfo/htmlinfo`

Use: `import "github.com/dyatlov/go-htmlinfo/htmlinfo"`

Example:

```go
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

	// if url can be nil too, just then we won't be able to fetch (and generate) oembed information
	err = info.Parse(resp.Body, &u)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Info:\n%s\n", info)

	fmt.Printf("Oembed information: %s\n", info.GenerateOembedFor(u))
}
```

Result would be:

_Info:_
```javascript
{"title":"iPad Pro Coming In November, Pricing Starts At $799  |  TechCrunch","description":"Apple unveiled its new iPad Pro today. If you're wondering when you can get your hands on it, and how much it will cost, here you go: Apple says the iPad Pro..","canonical_url":"http://techcrunch.com/2015/09/09/ipad-pro-coming-in-november-pricing-starts-at-799/","oembed_json_url":"https://public-api.wordpress.com/oembed/1.0/?format=json\u0026url=http%3A%2F%2Ftechcrunch.com%2F2015%2F09%2F09%2Fipad-pro-coming-in-november-pricing-starts-at-799%2F\u0026for=wpcom-auto-discovery","oembed_xml_url":"https://public-api.wordpress.com/oembed/1.0/?format=xml\u0026url=http%3A%2F%2Ftechcrunch.com%2F2015%2F09%2F09%2Fipad-pro-coming-in-november-pricing-starts-at-799%2F\u0026for=wpcom-auto-discovery","favicon_url":"https://s0.wp.com/wp-content/themes/vip/techcrunch-2013/assets/images/favicon.ico","image_src_url":"","main_content":"Apple unveiled its new iPad Pro today. If you’re wondering when you can get your hands on it, and how much it will cost, here you go: Apple says the iPad Pro and related accessories will be available in November.\nPricing will start at $799 with 32 gigabytes of memory and WiFi-only connectivity, with a $949 price tag for 128 GB, and $1,079 for 128 GB and a cellular connection. If you want the company’s new stylus, dubbed the Apple Pencil, that’ll cost you $99, and the Smart Keyboard will cost $169.\nThat may seem pretty pricey compared to other iPads — in fact, Apple said today that it’s dropping pricing on its iPad Mini 2, which its considers to be the entry-level iPad, to $269. What you’re paying for (among other things) is a 12.9-inch screen with resolution of 2,732 x 2,048 pixels, Apple’s A9X chip and four speakers.\nAnd, as the name and on-stage demos suggest, Apple seems to be pitching this for enterprise use and productivity, not for casual use.\n\t\t\t\n\t\t\t\t\n\t\t\t\tSay Hello To The Brand-New iPad Pro\n\t\t\t\n\t\t\t\t\t\t\t\t\n\t\t\t\t\t\t\t","opengraph":{"type":"article","url":"http://social.techcrunch.com/2015/09/09/ipad-pro-coming-in-november-pricing-starts-at-799/","title":"iPad Pro Coming In November, Pricing Starts At $799","description":"Apple unveiled its new iPad Pro today. If you're wondering when you can get your hands on it, and how much it will cost, here you go: Apple says the iPad Pro..","determiner":"","site_name":"TechCrunch","locale":"","locales_alternate":null,"images":[{"url":"https://tctechcrunch2011.files.wordpress.com/2015/09/screen-shot-2015-09-09-at-1-49-10-pm.png?w=560\u0026h=292\u0026crop=1","secure_url":"","type":"","width":0,"height":0}],"audios":null,"videos":null,"article":{"published_time":null,"modified_time":null,"expiration_time":null,"section":"","tags":null,"authors":null}},"oembed":{"type":"link","url":"http://techcrunch.com/2015/09/09/ipad-pro-coming-in-november-pricing-starts-at-799/","provider_url":"http://techcrunch.com","provider_name":"TechCrunch","title":"iPad Pro Coming In November, Pricing Starts At\u0026nbsp;$799","description":"","width":0,"height":0,"thumbnail_url":"https://i1.wp.com/tctechcrunch2011.files.wordpress.com/2015/09/screen-shot-2015-09-09-at-1-49-10-pm.png?fit=440%2C330","thumbnail_width":440,"thumbnail_height":218,"author_name":"\u003ca href=\"/author/anthony-ha/\" title=\"Posts by Anthony Ha\" onclick=\"s_objectID='river_author';\" rel=\"author\"\u003eAnthony Ha\u003c/a\u003e","author_url":"/author/anthony-ha/","html":"Apple \u003ca href=\"http://techcrunch.com/2015/09/09/apple-unveils-the-ipad-pro/\"\u003eunveiled its new iPad Pro today\u003c/a\u003e. If you're wondering when you can get your hands on it, and how much it will cost, here you go: Apple says the iPad Pro and related accessories will be available in November.\r\n\r\nPricing will start at $799 with 32 gigabytes of memory and WiFi-only connectivity, with a $949 price tag for 128 GB, and $1,079 for 128 GB and a cellular connection. If you want the company's new stylus, \u003ca href=\"http://techcrunch.com/2015/09/09/the-apple-pencil-is-the-ipad-pros-secret-weapon/#.91issd:LNXD\"\u003edubbed the Apple Pencil\u003c/a\u003e, that'll cost you $99, and the Smart Keyboard will cost $169.\r\n"}}
```

_Oembed information:_
```javascript
{"type":"link","url":"http://techcrunch.com/2015/09/09/ipad-pro-coming-in-november-pricing-starts-at-799/","provider_url":"http://techcrunch.com","provider_name":"TechCrunch","title":"iPad Pro Coming In November, Pricing Starts At\u0026nbsp;$799","description":"Apple unveiled its new iPad Pro today. If you're wondering when you can get your hands on it, and how much it will cost, here you go: Apple says the iPad Pro..","width":0,"height":0,"thumbnail_url":"https://i1.wp.com/tctechcrunch2011.files.wordpress.com/2015/09/screen-shot-2015-09-09-at-1-49-10-pm.png?fit=440%2C330","thumbnail_width":440,"thumbnail_height":218,"author_name":"\u003ca href=\"/author/anthony-ha/\" title=\"Posts by Anthony Ha\" onclick=\"s_objectID='river_author';\" rel=\"author\"\u003eAnthony Ha\u003c/a\u003e","author_url":"/author/anthony-ha/","html":"Apple \u003ca href=\"http://techcrunch.com/2015/09/09/apple-unveils-the-ipad-pro/\"\u003eunveiled its new iPad Pro today\u003c/a\u003e. If you're wondering when you can get your hands on it, and how much it will cost, here you go: Apple says the iPad Pro and related accessories will be available in November.\r\n\r\nPricing will start at $799 with 32 gigabytes of memory and WiFi-only connectivity, with a $949 price tag for 128 GB, and $1,079 for 128 GB and a cellular connection. If you want the company's new stylus, \u003ca href=\"http://techcrunch.com/2015/09/09/the-apple-pencil-is-the-ipad-pros-secret-weapon/#.91issd:LNXD\"\u003edubbed the Apple Pencil\u003c/a\u003e, that'll cost you $99, and the Smart Keyboard will cost $169.\r\n"}
```
