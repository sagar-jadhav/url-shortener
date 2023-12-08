package utils

import (
	crypto_rand "crypto/rand"
	"encoding/base64"
	"encoding/binary"
	math_rand "math/rand"
	"regexp"
	"strings"
)

const (
	BUFFER_SIZE = 8
	LOCALHOST   = "localhost"
)

func init() {
	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("error in seeding math/rand package with cryptographically secure random number generator")
	}
	math_rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
}

// Generate the random string of the given length
func GenerateRandomString(length int) string {
	data := make([]byte, BUFFER_SIZE)
	math_rand.Read(data)
	return base64.URLEncoding.EncodeToString([]byte(data))[:length]
}

// ValidateURL validates the given URL
// It validate whether the given string is URL or not
// It validates whether the domain is not localhost
func ValidateURL(urlStr string) bool {
	r, _ := regexp.Compile("^(http:\\/\\/www\\.|https:\\/\\/www\\.|http:\\/\\/|https:\\/\\/|\\/|\\/\\/)?[A-z0-9_-]*?[:]?[A-z0-9_-]*?[@]?[A-z0-9]+([\\-\\.]{1}[a-z0-9]+)*\\.[a-z]{2,5}(:[0-9]{1,5})?(\\/.*)?$")
	return r.Match([]byte(urlStr))
}

// GetDomainName returns the domain name of the URL
func GetDomainName(url string) string {
	var domain string
	if strings.HasPrefix(url, "http://") {
		domain = strings.TrimPrefix(url, "http://")
	} else if strings.HasPrefix(url, "https://") {
		domain = strings.TrimPrefix(url, "https://")
	} else {
		domain = url
	}

	if strings.HasPrefix(domain, "www") {
		domain = strings.TrimPrefix(domain, "www.")
	}

	if strings.Contains(domain, "/") {
		return strings.Split(domain, "/")[0]
	}
	return domain
}
