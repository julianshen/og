// Package that provides functionality to parse open graph and twitter card meta data from a html page.
package og

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"strings"

	"strconv"

	"io/ioutil"

	"github.com/PuerkitoBio/goquery"
	readability "github.com/julianshen/go-readability"
)

var (
	ErrorType = errors.New("Should not be non-ptr or nil")
)

type OgImage struct {
	Url       string `meta:"og:image,og:image:url" json:"url,omitempty"`
	SecureUrl string `meta:"og:image:secure_url" json:"secureURL,omitempty"`
	Width     int    `meta:"og:image:width" json:"width,omitempty"`
	Height    int    `meta:"og:image:height" json:"height,omitempty"`
	Type      string `meta:"og:image:type" json:"type,omitempty"`
}

type OgVideo struct {
	Url       string `meta:"og:video,og:video:url" json:"url,omitempty"`
	SecureUrl string `meta:"og:video:secure_url" json:"secureURL,omitempty"`
	Width     int    `meta:"og:video:width" json:"width,omitempty"`
	Height    int    `meta:"og:video:height" json:"height,omitempty"`
	Type      string `meta:"og:video:type" json:"type,omitempty"`
}

type OgAudio struct {
	Url       string `meta:"og:audio,og:audio:url" json:"url,omitempty"`
	SecureUrl string `meta:"og:audio:secure_url" json:"secureURL,omitempty"`
	Type      string `meta:"og:audio:type" json:"type,omitempty"`
}

type TwitterCard struct {
	Card        string `meta:"twitter:card" json:"card,omitempty"`
	Site        string `meta:"twitter:site" json:"site,omitempty"`
	SiteId      string `meta:"twitter:site:id" json:"siteID,omitempty"`
	Creator     string `meta:"twitter:creator" json:"creator,omitempty"`
	CreatorId   string `meta:"twitter:creator:id" json:"creatorID,omitempty"`
	Description string `meta:"twitter:description" json:"description,omitempty"`
	Title       string `meta:"twitter:title" json:"title,omitempty"`
	Image       string `meta:"twitter:image,twitter:image:src" json:"image,omitempty"`
	ImageAlt    string `meta:"twitter:image:alt" json:"imageAlt,omitempty"`
	Url         string `meta:"twitter:url" json:"url,omitempty"`
	Player      struct {
		Url    string `meta:"twitter:player" json:"url,omitempty"`
		Width  int    `meta:"twitter:width" json:"width,omitempty"`
		Height int    `meta:"twitter:height" json:"height,omitempty"`
		Stream string `meta:"twitter:stream" json:"stream,omitempty"`
	}
	IPhone struct {
		Name string `meta:"twitter:app:name:iphone" json:"name,omitempty"`
		Id   string `meta:"twitter:app:id:iphone" json:"id,omitempty"`
		Url  string `meta:"twitter:app:url:iphone" json:"url,omitempty"`
	}
	IPad struct {
		Name string `meta:"twitter:app:name:ipad" json:"name,omitempty"`
		Id   string `meta:"twitter:app:id:ipad" json:"id,omitempty"`
		Url  string `meta:"twitter:app:url:ipad" json:"url,omitempty"`
	}
	Googleplay struct {
		Name string `meta:"twitter:app:name:googleplay" json:"name,omitempty"`
		Id   string `meta:"twitter:app:id:googleplay" json:"id,omitempty"`
		Url  string `meta:"twitter:app:url:googleplay" json:"url,omitempty"`
	}
}

type PageInfo struct {
	Title       string `meta:"og:title" json:"title,omitempty"`
	Type        string `meta:"og:type" json:"type,omitempty"`
	Url         string `meta:"og:url" json:"url,omitempty"`
	Site        string `meta:"og:site" json:"site,omitempty"`
	SiteName    string `meta:"og:site_name" json:"siteName,omitempty"`
	Description string `meta:"og:description" json:"description,omitempty"`
	Locale      string `meta:"og:locale" json:"locale,omitempty"`
	Images      []*OgImage
	Videos      []*OgVideo
	Audios      []*OgAudio
	Twitter     *TwitterCard
	Content     string
}

func GetPageDataFromUrl(urlStr string, data interface{}) error {
	doc, err := goquery.NewDocument(urlStr)

	if err != nil {
		return err
	}

	return GetPageData(doc, data)
}

func GetPageDataFromResponse(response *http.Response, data interface{}) error {
	doc, err := goquery.NewDocumentFromResponse(response)

	if err != nil {
		return err
	}

	return GetPageData(doc, data)
}

func GetPageDataFromHtml(html []byte, data interface{}) error {
	buf := bytes.NewBuffer(html)
	doc, err := goquery.NewDocumentFromReader(buf)

	if err != nil {
		return err
	}

	return GetPageData(doc, data)
}

func GetPageData(doc *goquery.Document, data interface{}) error {
	doc = goquery.CloneDocument(doc)
	return getPageData(doc, data)
}

func GetPageInfo(doc *goquery.Document) (*PageInfo, error) {
	info := PageInfo{}
	err := GetPageData(doc, &info)

	if err != nil {
		return nil, err
	}

	html, _ := doc.Html()
	r, err := readability.NewDocument(html)
	if err != nil {
		return nil, err
	}

	info.Content = r.Text()

	return &info, nil
}

func GetPageInfoFromResponse(response *http.Response) (*PageInfo, error) {
	info := PageInfo{}
	html, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	err = GetPageDataFromHtml(html, &info)

	if err != nil {
		return nil, err
	}

	r, err := readability.NewDocument(string(html))
	if err != nil {
		return nil, err
	}

	info.Content = r.Text()

	return &info, nil
}

func GetPageInfoFromHtml(html []byte) (*PageInfo, error) {
	info := PageInfo{}

	err := GetPageDataFromHtml(html, &info)

	if err != nil {
		return nil, err
	}

	r, err := readability.NewDocument(string(html))
	if err != nil {
		return nil, err
	}

	info.Content = r.Text()

	return &info, nil
}

func GetPageInfoFromUrl(urlStr string) (*PageInfo, error) {
	resp, err := http.Get(urlStr)

	if err != nil {
		return nil, err
	}
	return GetPageInfoFromResponse(resp)
}

func getPageData(doc *goquery.Document, data interface{}) error {
	var rv reflect.Value
	var ok bool
	if rv, ok = data.(reflect.Value); !ok {
		rv = reflect.ValueOf(data)
	}

	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return ErrorType
	}

	rt := rv.Type()

	for i := 0; i < rv.Elem().NumField(); i++ {
		fv := rv.Elem().Field(i)
		field := rt.Elem().Field(i)

		switch fv.Type().Kind() {
		case reflect.Ptr:
			if fv.IsNil() {
				fv.Set(reflect.New(fv.Type().Elem()))
			}
			e := getPageData(doc, fv)

			if e != nil {
				return e
			}
		case reflect.Struct:
			e := getPageData(doc, fv.Addr())

			if e != nil {
				return e
			}
		case reflect.Slice:
			if fv.IsNil() {
				fv.Set(reflect.MakeSlice(fv.Type(), 0, 0))
			}

			switch field.Type.Elem().Kind() {
			case reflect.Struct:
				last := reflect.New(field.Type.Elem())
				for {
					data := reflect.New(field.Type.Elem())
					e := getPageData(doc, data.Interface())

					if e != nil {
						return e
					}

					//Ugly solution (I can't remove nodes. Why?)
					if !reflect.DeepEqual(last.Elem().Interface(), data.Elem().Interface()) {
						fv.Set(reflect.Append(fv, data.Elem()))
						last.Elem().Set(data.Elem())

					} else {
						break
					}
				}
			case reflect.Ptr:
				last := reflect.New(field.Type.Elem().Elem())
				for {
					data := reflect.New(field.Type.Elem().Elem())
					e := getPageData(doc, data.Interface())

					if e != nil {
						return e
					}

					//Ugly solution (I can't remove nodes. Why?)
					if !reflect.DeepEqual(last.Elem().Interface(), data.Elem().Interface()) {
						fv.Set(reflect.Append(fv, data))
						last.Elem().Set(data.Elem())

					} else {
						break
					}
				}
			default:
				if tag, ok := field.Tag.Lookup("meta"); ok {
					tags := strings.Split(tag, ",")

					for _, t := range tags {
						contents := []reflect.Value{}

						processMeta := func(idx int, sel *goquery.Selection) {
							if c, existed := sel.Attr("content"); existed {
								if field.Type.Elem().Kind() == reflect.String {
									contents = append(contents, reflect.ValueOf(c))
								} else {
									i, e := strconv.Atoi(c)

									if e == nil {
										contents = append(contents, reflect.ValueOf(i))
									}
								}

								fv.Set(reflect.Append(fv, contents...))
							}
						}

						doc.Find(fmt.Sprintf("meta[property=\"%s\"]", t)).Each(processMeta)

						doc.Find(fmt.Sprintf("meta[name=\"%s\"]", t)).Each(processMeta)

						fv = reflect.Append(fv, contents...)
					}
				}
			}
		default:
			if tag, ok := field.Tag.Lookup("meta"); ok {

				tags := strings.Split(tag, ",")

				content := ""
				existed := false
				sel := (*goquery.Selection)(nil)
				for _, t := range tags {
					if sel = doc.Find(fmt.Sprintf("meta[property=\"%s\"]", t)).First(); sel.Size() > 0 {
						content, existed = sel.Attr("content")
					}

					if !existed {
						if sel = doc.Find(fmt.Sprintf("meta[name=\"%s\"]", t)).First(); sel.Size() > 0 {
							content, existed = sel.Attr("content")
						}
					}

					if existed {
						if fv.Type().Kind() == reflect.String {
							fv.Set(reflect.ValueOf(content))
						} else if fv.Type().Kind() == reflect.Int {
							if i, e := strconv.Atoi(content); e == nil {
								fv.Set(reflect.ValueOf(i))
							}
						}
						break
					}
				}
			}
		}
	}
	return nil
}
