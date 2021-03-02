package validate

import "net/url"

// IsValidURL : Checks if a string is a valid URL and if it is it returns
// a pointer to the URL and a bool with the result of the validation
func IsValidURL(input string) (*url.URL, bool) {
	u, err := url.Parse(input)
	return u, err == nil && u.Scheme != "" && u.Host != ""
}
