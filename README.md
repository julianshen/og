# og - An open graph parser for Golang

This is for parsing open graph (and Twitter card) meta tags in a html page

## Data structures

### PageInfo

```go
type PageInfo struct {
	Title    string `meta:"og:title"`
	Type     string `meta:"og:type"`
	Url      string `meta:"og:url"`
	Site     string `meta:"og:site"`
	SiteName string `meta:"og:site_name"`
	Images   []*OgImage
	Videos   []*OgVideo
	Audios   []*OgAudio
	Twitter  *TwitterCard
	Content  string
}
```

Except "Content", all fields come from og and twitter meta tags. And "Content" is readable text in this page.

### OgImage/OgVideo/OgAudio

```go
type OgImage struct {
	Url       string `meta:"og:image,og:image:url"`
	SecureUrl string `meta:"og:image:secure_url"`
	Width     int    `meta:"og:image:width"`
	Height    int    `meta:"og:image:height"`
	Type      string `meta:"og:image:type"`
}

type OgVideo struct {
	Url       string `meta:"og:video,og:video:url"`
	SecureUrl string `meta:"og:video:secure_url"`
	Width     int    `meta:"og:video:width"`
	Height    int    `meta:"og:video:height"`
	Type      string `meta:"og:video:type"`
}

type OgAudio struct {
	Url       string `meta:"og:audio,og:audio:url"`
	SecureUrl string `meta:"og:audio:secure_url"`
	Type      string `meta:"og:audio:type"`
}
```

### TwitterCard

```go
type TwitterCard struct {
	Card        string `meta:"twitter:card"`
	Site        string `meta:"twitter:site"`
	SiteId      string `meta:"twitter:site:id"`
	Creator     string `meta:"twitter:creator"`
	CreatorId   string `meta:"twitter:creator:id"`
	Description string `meta:"twitter:description"`
	Title       string `meta:"twitter:title"`
	Image       string `meta:"twitter:image,twitter:image:src"`
	ImageAlt    string `meta:"twitter:image:alt"`
	Url         string `meta:"twitter:url"`
	Player      struct {
		Url    string `meta:"twitter:player"`
		Width  int    `meta:"twitter:width"`
		Height int    `meta:"twitter:height"`
		Stream string `meta:"twitter:stream"`
	}
	IPhone struct {
		Name string `meta:"twitter:app:name:iphone"`
		Id   string `meta:"twitter:app:id:iphone"`
		Url  string `meta:"twitter:app:url:iphone"`
	}
	IPad struct {
		Name string `meta:"twitter:app:name:ipad"`
		Id   string `meta:"twitter:app:id:ipad"`
		Url  string `meta:"twitter:app:url:ipad"`
	}
	Googleplay struct {
		Name string `meta:"twitter:app:name:googleplay"`
		Id   string `meta:"twitter:app:id:googleplay"`
		Url  string `meta:"twitter:app:url:googleplay"`
	}
}
```

## GetPageInfo

You can get page information direct from the url like this:

```go
urlStr := "https://techcrunch.com/2017/01/22/yahoo-hacking-sec/"
pageInfo, e := GetPageInfoFromUrl(urlStr)
```

## GetPageData

Another way to retrieve page related data is `GetPageData`. For example,
if you want to get `og:image` information just do this:

```go
onImage := OgImage{}
urlStr := "https://techcrunch.com/2017/01/22/yahoo-hacking-sec/"
GetPageDataFromUrl(urlStr, &onImage)
```

If you have custom attributes that is not supported, you could also define your own data structure with struct field tag `meta:`. And pass your data structure to `GetPageData`.
