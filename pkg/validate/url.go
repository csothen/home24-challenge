package validate

import "net/url"

func IsValidURL(input string) (*url.URL, bool) {
	u, err := url.Parse(input)
	return u, err == nil && u.Scheme != "" && u.Host != ""
}
