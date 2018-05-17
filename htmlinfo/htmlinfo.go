package htmlinfo

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"

	"github.com/dyatlov/go-oembed/oembed"
	"github.com/dyatlov/go-opengraph/opengraph"
	"golang.org/x/net/html"

	"github.com/dyatlov/go-readability"
)

// TouchIcon contains all icons parsed from page header, including Apple touch icons
type TouchIcon struct {
	URL        string `json:"url"`
	Type       string `json:"type"`
	Width      uint64 `json:"width"`
	Height     uint64 `json:"height"`
	IsScalable bool   `json:"is_scalable"`
}

// HTMLInfo contains information extracted from HTML page
type HTMLInfo struct {
	url *url.URL
	// http.Client instance to use, if nil then will be used default client
	Client *http.Client `json:"-"`
	// If it's true then parser will fetch oembed data from oembed url if possible
	AllowOembedFetching bool `json:"-"`
	// If it's true parser will extract main page content from html
	AllowMainContentExtraction bool `json:"-"`
	// We'll forward it to Oembed' fetchOembed method
	AcceptLanguage string `json:"-"`

	Title         string       `json:"title"`
	Description   string       `json:"description"`
	AuthorName    string       `json:"author_name"`
	CanonicalURL  string       `json:"canonical_url"`
	OembedJSONURL string       `json:"oembed_json_url"`
	OembedXMLURL  string       `json:"oembed_xml_url"`
	FaviconURL    string       `json:"favicon_url"`
	TouchIcons    []*TouchIcon `json:"touch_icons"`
	ImageSrcURL   string       `json:"image_src_url"`
	// Readability package is being used inside
	MainContent string               `json:"main_content"`
	OGInfo      *opengraph.OpenGraph `json:"opengraph"`
	OembedInfo  *oembed.Info         `json:"oembed"`
}

var (
	cleanHTMLTagsRegex    = regexp.MustCompile(`<.*?>`)
	replaceNewLinesRegex  = regexp.MustCompile(`[\r\n]+`)
	clearWhitespacesRegex = regexp.MustCompile(`\s+`)
	getImageRegex         = regexp.MustCompile(`(?i)<img[^>]+?src=("|')?(.*?)("|'|\s|>)`)
	linkWithIconsRegex    = regexp.MustCompile(`\b(icon|image_src)\b`)
	sizesRegex            = regexp.MustCompile(`(\d+)[^\d]+(\d+)`) // some websites use crazy unicode chars between height and width
)

// NewHTMLInfo return new instance of HTMLInfo
func NewHTMLInfo() *HTMLInfo {
	info := &HTMLInfo{AllowOembedFetching: true, AllowMainContentExtraction: true, OGInfo: opengraph.NewOpenGraph(), AcceptLanguage: "en-us"}
	return info
}

func (info *HTMLInfo) toAbsoluteURL(u string) string {
	if info.url == nil {
		return u
	}

	tu, _ := url.Parse(u)

	if tu != nil {
		if tu.Host == "" {
			tu.Scheme = info.url.Scheme
			tu.Host = info.url.Host
			tu.User = info.url.User
			tu.Opaque = info.url.Opaque
			if len(tu.Path) == 0 || tu.Path[0] != '/' {
				tu.Path = info.url.Path + tu.Path
			}
		} else if tu.Scheme == "" {
			tu.Scheme = info.url.Scheme
		}

		return tu.String()
	}

	return u
}

func (info *HTMLInfo) appendTouchIcons(url string, rel string, sizes []string) {
	for _, size := range sizes {
		icon := &TouchIcon{URL: url, Type: rel, IsScalable: (size == "any")}
		matches := sizesRegex.FindStringSubmatch(size)
		if len(matches) >= 3 {
			icon.Height, _ = strconv.ParseUint(matches[1], 10, 64)
			icon.Width, _ = strconv.ParseUint(matches[2], 10, 64)
		}
		info.TouchIcons = append(info.TouchIcons, icon)
	}
}

func (info *HTMLInfo) parseLinkIcon(attrs map[string]string) {
	rels := strings.Split(attrs["rel"], " ")
	url := info.toAbsoluteURL(attrs["href"])
	sizesString, present := attrs["sizes"]
	if !present {
		sizesString = "0x0"
	}
	sizes := strings.Split(sizesString, " ")

	for _, rel := range rels {
		if rel == "image_src" {
			info.ImageSrcURL = url
		} else if rel == "icon" {
			info.FaviconURL = url
			info.appendTouchIcons(url, rel, sizes)
		} else if rel == "apple-touch-icon" || rel == "apple-touch-icon-precomposed" {
			info.appendTouchIcons(url, rel, sizes)
		}
	}
}

func (info *HTMLInfo) parseHead(n *html.Node) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "title" {
			if c.FirstChild != nil {
				info.Title = c.FirstChild.Data
			}
		} else if c.Type == html.ElementNode && c.Data == "link" {
			m := make(map[string]string)
			for _, a := range c.Attr {
				m[a.Key] = a.Val
			}
			if m["rel"] == "canonical" {
				info.CanonicalURL = info.toAbsoluteURL(m["href"])
			} else if m["rel"] == "alternate" && m["type"] == "application/json+oembed" {
				info.OembedJSONURL = info.toAbsoluteURL(m["href"])
			} else if m["rel"] == "alternate" && m["type"] == "application/xml+oembed" {
				info.OembedXMLURL = info.toAbsoluteURL(m["href"])
			} else if linkWithIconsRegex.MatchString(m["rel"]) {
				info.parseLinkIcon(m)
			}
		} else if c.Type == html.ElementNode && c.Data == "meta" {
			m := make(map[string]string)
			for _, a := range c.Attr {
				m[a.Key] = a.Val
			}

			if m["name"] == "description" {
				info.Description = m["content"]
			} else if m["name"] == "author" {
				info.AuthorName = m["content"]
			}

			info.OGInfo.ProcessMeta(m)
		}
	}
}

func (info *HTMLInfo) parseBody(n *html.Node) {
	if !info.AllowMainContentExtraction {
		return
	}

	buf := new(bytes.Buffer)
	err := html.Render(buf, n)
	if err != nil {
		return
	}
	bufStr := buf.String()
	doc, err := readability.NewDocument(bufStr)
	if err != nil {
		return
	}

	doc.WhitelistTags = []string{"div", "p", "img"}
	doc.WhitelistAttrs["img"] = []string{"src", "title", "alt"}

	content := doc.Content()
	content = html.UnescapeString(content)

	info.MainContent = strings.Trim(content, "\r\n\t ")
}

// Parse return information about page
// @param s - contains page source
// @params pageURL - contains URL from where the data was taken [optional]
// @params contentType - contains Content-Type header value [optional]
// if no url is given then parser won't attempt to parse oembed info
func (info *HTMLInfo) Parse(s io.Reader, pageURL *string, contentType *string) error {
	contentTypeStr := "text/html"
	if contentType != nil && len(*contentType) > 0 {
		contentTypeStr = *contentType
	}
	utf8s, err := charset.NewReader(s, contentTypeStr)
	if err != nil {
		return err
	}

	if pageURL != nil {
		tu, _ := url.Parse(*pageURL)
		info.url = tu
	}

	doc, err := html.Parse(utf8s)
	if err != nil {
		return err
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode {
				if c.Data == "head" {
					info.parseHead(c)
					continue
				} else if c.Data == "body" {
					info.parseBody(c)
					continue
				}
			}
			f(c)
		}
	}
	f(doc)

	if info.AllowOembedFetching && pageURL != nil && len(info.OembedJSONURL) > 0 {
		pu, _ := url.Parse(info.OembedJSONURL)
		siteName := info.OGInfo.SiteName
		siteURL := strings.ToLower(pu.Scheme) + "://" + pu.Host

		if len(siteName) == 0 {
			siteName = pu.Host
		}

		oiItem := &oembed.Item{EndpointURL: info.OembedJSONURL, ProviderName: siteName, ProviderURL: siteURL, IsEndpointURLComplete: true}
		oi, _ := oiItem.FetchOembed(oembed.Options{URL: *pageURL, Client: info.Client, AcceptLanguage: info.AcceptLanguage})
		if oi != nil && oi.Status < 300 {
			info.OembedInfo = oi
		}
	}

	return nil
}

func (info *HTMLInfo) trimText(text string, maxLen int) string {
	var numRunes = 0
	runes := []rune(text)
	for index := range runes {
		numRunes++
		if numRunes > maxLen {
			return string(runes[:index-3]) + "..."
		}
	}
	return text
}

// GenerateOembedFor return Oembed Info for given url based on previously parsed data
// The returned oembed data is also updated in info.OembedInfo
// Example:
//
// info := NewHTMLInfo()
// info.Parse(dataReader, &sourceURL)
// oembed := info.GenerateOembedFor(sourceURL)
func (info *HTMLInfo) GenerateOembedFor(pageURL string) *oembed.Info {
	pu, _ := url.Parse(pageURL)

	if pu == nil {
		return nil
	}

	siteName := info.OGInfo.SiteName
	siteURL := strings.ToLower(pu.Scheme) + "://" + pu.Host

	if len(siteName) == 0 {
		siteName = pu.Host
	}

	title := info.OGInfo.Title
	if len(title) == 0 {
		title = info.Title
	}

	description := info.OGInfo.Description
	if len(description) == 0 {
		description = info.Description
		if len(description) == 0 {
			if len(info.MainContent) > 0 {
				description = cleanHTMLTagsRegex.ReplaceAllString(info.MainContent, " ")
				description = replaceNewLinesRegex.ReplaceAllString(description, " ")
				description = clearWhitespacesRegex.ReplaceAllString(description, " ")
				description = strings.Trim(description, " ")
				description = info.trimText(description, 200)
			}
		}
	}

	baseInfo := &oembed.Info{}

	baseInfo.Type = "link"
	baseInfo.URL = pageURL
	baseInfo.ProviderURL = siteURL
	baseInfo.ProviderName = siteName
	baseInfo.Title = title
	baseInfo.Description = description

	if len(info.ImageSrcURL) > 0 {
		baseInfo.ThumbnailURL = info.toAbsoluteURL(info.ImageSrcURL)
	}

	if len(info.OGInfo.Images) > 0 {
		baseInfo.ThumbnailURL = info.toAbsoluteURL(info.OGInfo.Images[0].URL)
		baseInfo.ThumbnailWidth = info.OGInfo.Images[0].Width
		baseInfo.ThumbnailHeight = info.OGInfo.Images[0].Height
	}

	if len(baseInfo.ThumbnailURL) == 0 && len(info.MainContent) > 0 {
		// get first image from body
		matches := getImageRegex.FindStringSubmatch(info.MainContent)
		if len(matches) > 0 {
			baseInfo.ThumbnailURL = info.toAbsoluteURL(matches[2])
		}
	}

	// first we check if there is link to oembed resource
	if info.OembedInfo != nil {
		info.OembedInfo.MergeWith(baseInfo)
		return info.OembedInfo
	}

	return baseInfo
}

// ToJSON return json represenation of structure, simple wrapper around json package
func (info *HTMLInfo) ToJSON() ([]byte, error) {
	return json.Marshal(info)
}

func (info *HTMLInfo) String() string {
	data, err := info.ToJSON()
	if err != nil {
		return err.Error()
	}
	return string(data[:])
}
