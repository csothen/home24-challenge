package request

import "net/http"

// Get makes a GET request to a certain URL
func Get(url string) (*http.Response, error) {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", url, nil)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
