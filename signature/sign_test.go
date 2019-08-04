package signature

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortMapKey(t *testing.T) {
	cases := map[string]struct {
		testMap     map[string]string
		expectSlice []string
	}{
		"sort map key test": {
			testMap: map[string]string{
				"Host":         "iam.amazonaws.com\n",
				"content-type": "application/x-www-form-urlencoded; charset=utf-8\n",
				"My-header1":   `    a   b   c  ` + "\n",
				"X-Amz-Date":   "20150830T123600Z\n",
				"My-Header2":   `    "a   b   c"  ` + "\n",
			},
			expectSlice: []string{"content-type", "host", "my-header1", "my-header2", "x-amz-date"},
		},
	}

	for n, tc := range cases {
		tc := tc
		actualSlice := sortMapKey(tc.testMap)
		assert.Equal(t, tc.expectSlice, actualSlice, n)
	}
}
func TestNormraizeHeader(t *testing.T) {
	cases := map[string]struct {
		testHeader   map[string]string
		expectString string
	}{
		"normarize header test": {
			testHeader: map[string]string{
				"Host":         "iam.amazonaws.com\n",
				"content-type": "application/x-www-form-urlencoded; charset=utf-8\n",
				"My-header1":   `    a   b   c  ` + "\n",
				"X-Amz-Date":   "20150830T123600Z\n",
				"My-Header2":   `    "a   b   c"  ` + "\n",
			},
			expectString: `content-type:application/x-www-form-urlencoded; charset=utf-8
host:iam.amazonaws.com
my-header1:a b c
my-header2:"a b c"
x-amz-date:20150830T123600Z
`,
		},
		"example2 test": {
			testHeader: map[string]string{
				"Host":         "iam.amazonaws.com",
				"Content-Type": "application/x-www-form-urlencoded; charset=utf-8",
				"X-Amz-Date":   "20150830T123600Z",
			},
			expectString: `content-type:application/x-www-form-urlencoded; charset=utf-8
host:iam.amazonaws.com
x-amz-date:20150830T123600Z
`,
		},
	}

	for n, tc := range cases {
		tc := tc
		actualString := normarizeHeader(tc.testHeader)
		assert.Equal(t, tc.expectString, actualString, n)
	}
}

func TestLinkSlice(t *testing.T) {
	cases := map[string]struct {
		testSlice    []string
		expectString string
	}{
		"link slice test": {
			testSlice:    []string{"content-type", "host", "x-amz-date"},
			expectString: "content-type;host;x-amz-date",
		},
	}

	for n, tc := range cases {
		tc := tc
		actualString := linkSlice(tc.testSlice)
		assert.Equal(t, tc.expectString, actualString, n)
	}
}

func TestHashSHA256(t *testing.T) {
	cases := map[string]struct {
		testPayload string
		expectHash  string
	}{
		"hash test": {
			testPayload: "",
			expectHash:  "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		"normarize request hash test": {
			testPayload: `GET
/
Action=ListUsers&Version=2010-05-08
content-type:application/x-www-form-urlencoded; charset=utf-8
host:iam.amazonaws.com
x-amz-date:20150830T123600Z

content-type;host;x-amz-date
e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`,
			expectHash: "f536975d06c0309214f805bb90ccff089219ecd68b2577efef23edd43b7e1a59",
		},
	}

	for n, tc := range cases {
		tc := tc
		actualHash := hashSHA256(tc.testPayload)
		assert.Equal(t, tc.expectHash, actualHash, n)
	}
}

func TestCanonicalRequest(t *testing.T) {
	cases := map[string]struct {
		testMethod   string
		testURL      string
		testHeader   map[string]string
		testPayload  string
		expectString string
	}{
		"example canonical Request": {
			testMethod: "GET",
			testURL:    "https://iam.amazonaws.com/?Action=ListUsers&Version=2010-05-08",
			testHeader: map[string]string{
				"Host":         "iam.amazonaws.com",
				"Content-Type": "application/x-www-form-urlencoded; charset=utf-8",
				"X-Amz-Date":   "20150830T123600Z",
			},
			testPayload: "",
			expectString: `GET
/
Action=ListUsers&Version=2010-05-08
content-type:application/x-www-form-urlencoded; charset=utf-8
host:iam.amazonaws.com
x-amz-date:20150830T123600Z

content-type;host;x-amz-date
e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`,
		},
	}

	for n, tc := range cases {
		tc := tc
		actualString := canonicalRequest(tc.testMethod, tc.testURL, tc.testPayload, tc.testHeader)
		assert.Equal(t, tc.expectString, actualString, n)
	}
}
func TestStringToSign(t *testing.T) {
	cases := map[string]struct {
		testISODate string
		testRegion  string
		testHash    string
		expectSign  string
	}{
		"sign test": {
			testISODate: "20150830T123600Z",
			testRegion:  "us-east-1",
			testHash:    "f536975d06c0309214f805bb90ccff089219ecd68b2577efef23edd43b7e1a59",
			expectSign: `AWS4-HMAC-SHA256
20150830T123600Z
20150830/us-east-1/s3/aws4_request
f536975d06c0309214f805bb90ccff089219ecd68b2577efef23edd43b7e1a59`,
		},
	}

	for n, tc := range cases {
		tc := tc
		actualSign := stringToSign(tc.testISODate, tc.testRegion, tc.testHash)
		assert.Equal(t, actualSign, tc.expectSign, n)
	}
}
