package signature

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strings"
)

var (
	iniialSpace = regexp.MustCompile(`^\s+`)
	space       = regexp.MustCompile(`\s+`)
)

type SortSice []string

func (s SortSice) Len() int           { return len(s) }
func (s SortSice) Less(i, j int) bool { return s[i] < s[j] }
func (s SortSice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func sortMapKey(header map[string]string) []string {
	strSlice := make([]string, 0, 10)
	for key := range header {
		strSlice = append(strSlice, strings.ToLower(key))
	}
	sort.Sort(SortSice(strSlice))
	return strSlice
}

func normarizeHeader(header map[string]string) string {
	headerKeyMap := make(map[string]string)

	for key := range header {
		lowerKey := strings.ToLower(key)
		headerKeyMap[lowerKey] = key
	}
	sortedKey := sortMapKey(header)
	var buffer bytes.Buffer
	for _, str := range sortedKey {
		value := header[headerKeyMap[str]]
		f := iniialSpace.ReplaceAllString(value, "")
		sp := space.ReplaceAllString(f, ` `)
		l := strings.TrimSuffix(sp, " ")
		s := fmt.Sprintf("%s:%s\n", strings.ToLower(str), l)
		buffer.WriteString(s)
	}
	return buffer.String()
}

func linkSlice(strSlice []string) string {
	var buffer bytes.Buffer
	for n, str := range strSlice {
		buffer.WriteString(str)
		if n != len(strSlice)-1 {
			buffer.WriteString(";")
		}
	}
	return buffer.String()
}

func hashSHA256(payload string) string {
	hash := sha256.Sum256([]byte(payload))
	hexed := hex.EncodeToString(hash[:])
	return strings.ToLower(hexed)
}

// We should not encode URL for S3 request.
// I don't encode query for /ObjectName?uploads.
func canonicalRequest(method, URL, payload string, header map[string]string) string {
	HTTPRequestMethod := fmt.Sprintf("%s\n", method)
	u, _ := url.Parse(URL)
	canonicalURL := fmt.Sprintf("%s\n", u.EscapedPath())

	v := u.Query()
	canonicalQueryString := fmt.Sprintf("%s\n", v.Encode())

	nrm := normarizeHeader(header)
	canonicalHeaders := fmt.Sprintf("%s\n", nrm)

	sortKeySlice := sortMapKey(header)
	signedHeaders := fmt.Sprintf("%s\n", linkSlice(sortKeySlice))
	hash := hashSHA256(payload)
	return HTTPRequestMethod + canonicalURL + canonicalQueryString + canonicalHeaders + signedHeaders + hash
}

func stringToSign(ISODate, AWSRegion, hash string) string {
	algorithm := "AWS4-HMAC-SHA256\n"
	dateTime := fmt.Sprintf("%s\n", ISODate)
	date := ISODate[:8]
	credentialScope := fmt.Sprintf("%s/%s/s3/aws4_request\n", string(date), AWSRegion)
	return algorithm + dateTime + credentialScope + hash
}

func makeHMAC(key, msg []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(msg)
	return mac.Sum(nil)
}

func getSignatureKey(secret, date, region, service string) []byte {
	kSecret := secret
	kDate := makeHMAC([]byte("AWS4"+kSecret), []byte(date))
	kRegion := makeHMAC(kDate, []byte(region))
	kService := makeHMAC(kRegion, []byte(service))
	kSigning := makeHMAC(kService, []byte("aws4_request"))
	return kSigning
}

func getSignature(secret, date, region, service, stringToSign string) string {
	signatureKey := getSignatureKey(secret, date, region, service)
	signature := makeHMAC(signatureKey, []byte(stringToSign))
	return hex.EncodeToString(signature)
}
