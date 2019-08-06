package signature

import (
	"encoding/hex"
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
		t.Run(n, func(t *testing.T) {
			actualSlice := sortMapKey(tc.testMap)
			assert.Equal(t, tc.expectSlice, actualSlice, n)
		})
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
		t.Run(n, func(t *testing.T) {
			actualString := normarizeHeader(tc.testHeader)
			assert.Equal(t, tc.expectString, actualString, n)
		})
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
		t.Run(n, func(t *testing.T) {
			actualString := linkSlice(tc.testSlice)
			assert.Equal(t, tc.expectString, actualString, n)
		})
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
		t.Run(n, func(t *testing.T) {
			actualHash := hashSHA256(tc.testPayload)
			assert.Equal(t, tc.expectHash, actualHash, n)
		})
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
		t.Run(n, func(t *testing.T) {
			actualString := canonicalRequest(tc.testMethod, tc.testURL, tc.testPayload, tc.testHeader)
			assert.Equal(t, tc.expectString, actualString)
		})

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
		t.Run(n, func(t *testing.T) {
			actualSign := stringToSign(tc.testISODate, tc.testRegion, tc.testHash)
			assert.Equal(t, actualSign, tc.expectSign, n)
		})
	}
}

func TestGetSignatureKey(t *testing.T) {
	cases := map[string]struct {
		testSecret         string
		testDate           string
		testRegion         string
		testService        string
		expectSignatureKey string
	}{
		"test get signature key": {
			testSecret:         "wJalrXUtnFEMI/K7MDENG+bPxRfiCYEXAMPLEKEY",
			testDate:           "20150830",
			testRegion:         "us-east-1",
			testService:        "iam",
			expectSignatureKey: "c4afb1cc5771d871763a393e44b703571b55cc28424d1a5e86da6ed3c154a4b9",
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			actualSignatureKey := getSignatureKey(tc.testSecret, tc.testDate, tc.testRegion, tc.testService)
			assert.Equal(t, tc.expectSignatureKey, hex.EncodeToString(actualSignatureKey))
		})
	}
}

func TestGetSignature(t *testing.T) {
	cases := map[string]struct {
		testSecret       string
		testDate         string
		testRegion       string
		testService      string
		testStringToSign string
		expectSignature  string
	}{
		"test get signature key": {
			testSecret:  "wJalrXUtnFEMI/K7MDENG+bPxRfiCYEXAMPLEKEY",
			testDate:    "20150830",
			testRegion:  "us-east-1",
			testService: "iam",
			testStringToSign: `AWS4-HMAC-SHA256
20150830T123600Z
20150830/us-east-1/iam/aws4_request
f536975d06c0309214f805bb90ccff089219ecd68b2577efef23edd43b7e1a59`,
			expectSignature: "5d672d79c15b13162d9279b0855cfba6789a8edb4c82c400e06b5924a6f2b5d7",
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			actualSignature := getSignature(tc.testSecret, tc.testDate, tc.testRegion, tc.testService, tc.testStringToSign)
			assert.Equal(t, tc.expectSignature, actualSignature)
		})
	}
}

func TestGetAuthrization(t *testing.T) {
	cases := map[string]struct {
		testSecretAccessKey string
		testCredentialScope string
		testSignedHeaders   string
		testSignature       string
		expectAuthorization string
	}{
		"get authorization test": {
			testSecretAccessKey: "AKIDEXAMPLE",
			testCredentialScope: "20150830/us-east-1/iam/aws4_request",
			testSignedHeaders:   "content-type;host;x-amz-date",
			testSignature:       "5d672d79c15b13162d9279b0855cfba6789a8edb4c82c400e06b5924a6f2b5d7",
			expectAuthorization: "AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20150830/us-east-1/iam/aws4_request, SignedHeaders=content-type;host;x-amz-date, Signature=5d672d79c15b13162d9279b0855cfba6789a8edb4c82c400e06b5924a6f2b5d7",
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			actualAuthorization := getAuthorization(tc.testSecretAccessKey, tc.testCredentialScope, tc.testSignedHeaders, tc.testSignature)
			assert.Equal(t, tc.expectAuthorization, actualAuthorization)
		})
	}
}
