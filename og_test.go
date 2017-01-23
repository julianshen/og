package og

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTags(t *testing.T) {
	s := PageInfo{}
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
	urlStr := "https://techcrunch.com/2017/01/22/yahoo-hacking-sec/"

	pageInfo, e := GetPageInfoFromUrl(urlStr)

	assert.Nil(t, e)
	assert.NotEqual(t, "", pageInfo.Title)
	assert.NotEqual(t, 0, len(pageInfo.Images))

	assert.NotEmpty(t, pageInfo.Content)
}
