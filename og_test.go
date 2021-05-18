package og_test

import (
	"log"
	"net"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/ringsaturn/og"
	"github.com/stretchr/testify/assert"
)

func TestTags(t *testing.T) {
	s := og.PageInfo{}
	st := reflect.TypeOf(s)
	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)

		tag, e := field.Tag.Lookup("og")
		if e {
			assert.True(t, true, strings.HasPrefix(tag, "og:"))
		}
	}
}

func TestGetPageInfo(t *testing.T) {
	url := "https://techcrunch.com/2017/01/22/yahoo-hacking-sec/"

	client := &http.Client{
		Timeout: time.Second * 5,
		Transport: &http.Transport{
			MaxIdleConns:        5,
			MaxIdleConnsPerHost: 5,
			IdleConnTimeout:     30 * time.Second,
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 5 * time.Second,
		},
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Panic(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Panic(err)
	}

	pageInfo, e := og.GetPageInfoFromResponse(resp)

	assert.Nil(t, e)
	assert.NotEqual(t, "", pageInfo.Title)
	assert.NotEqual(t, 0, len(pageInfo.Images))

	assert.NotEmpty(t, pageInfo.Content)
}
