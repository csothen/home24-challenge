package parsing

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/csothen/htmlparser/pkg/models"
	"github.com/csothen/htmlparser/pkg/request"
	"golang.org/x/net/html"
)

var (
	tags map[string]string = map[string]string{
		"<!DOCTYPE HTML>":        "HTML 5",
		"HTML 4.01//EN":          "HTML 4.01 Strict",
		"HTML 4.01 TRANSITIONAL": "HTML 4.01 Transitional",
		"HTML 4.01 FRAMESET":     "HTML 4.01 Frameset",
		"XHTML 1.0 STRICT":       "XHTML 1.0 Strict",
		"XHTML 1.0 TRANSITIONAL": "XHTML 1.0 Transitional",
		"XHTML 1.0 FRAMESET":     "XHTML 1.0 Frameset",
		"XHTML 1.1":              "XHTML 1.1",
		"XHTML BASIC 1.1":        "XHTML Basic 1.1",
	}

	linkCache map[string]bool
)

// Analyse : Parses the website and returns the analysis result
func Analyse(url url.URL, page *html.Tokenizer) *models.Result {
	// Default values
	htmlVersion := "Failed to find Doctype declaration"
	pageTitle := "Failed to find the website's title"
	headingsCount := make(map[int]int)
	internalLinksCount := 0
	externalLinksCount := 0
	inaccessibleLinksCount := 0
	containsLoginForm := false

	depth := 0
	linkCache = make(map[string]bool)

	result := &models.Result{
		HTMLVersion:            htmlVersion,
		PageTitle:              pageTitle,
		HeadingsCount:          headingsCount,
		InternalLinksCount:     internalLinksCount,
		ExternalLinksCount:     externalLinksCount,
		InaccessibleLinksCount: inaccessibleLinksCount,
		ContainsLoginForm:      containsLoginForm,
	}

	for {
		tokenType := page.Next()
		if tokenType == html.ErrorToken {
			return result
		}

		token := page.Token()

		switch tokenType {
		case html.DoctypeToken:
			checkVersion(result, token)

		case html.StartTagToken:
			if token.DataAtom.String() != "meta" {
				depth++
			}

			checkTitle(result, token, page)

			checkHeadings(result, token, depth)

			checkLinks(result, token, url)

			checkLoginForm(result, token)

		case html.EndTagToken:
			depth--
		}
	}
}

func checkVersion(result *models.Result, token html.Token) {
	normalizedDoctype := strings.ToUpper(token.String())

	result.HTMLVersion = "Failed to match doctype declaration to a version"

	for key, value := range tags {
		if strings.Contains(normalizedDoctype, key) {
			result.HTMLVersion = value
			break
		}
	}
}

func checkTitle(result *models.Result, token html.Token, page *html.Tokenizer) {
	if token.DataAtom.String() == "title" {
		page.Next()
		titleToken := page.Token()
		result.PageTitle = titleToken.String()
	}
}

func checkHeadings(result *models.Result, token html.Token, depth int) {
	if len(token.DataAtom.String()) == 2 && token.DataAtom.String()[0] == 'h' {
		n := strings.Split(token.DataAtom.String(), "")[1]
		if _, err := strconv.Atoi(n); err == nil {
			count, ok := result.HeadingsCount[depth]
			if !ok {
				count = 0
			}
			count++
			result.HeadingsCount[depth] = count
		}
	}
}

func checkLinks(result *models.Result, token html.Token, inputURL url.URL) {
	if token.DataAtom.String() == "a" {
		for _, attr := range token.Attr {
			if attr.Key == "href" {
				url, err := url.Parse(attr.Val)

				if err != nil {
					result.InaccessibleLinksCount++
					continue
				}

				if url.Host != inputURL.Host && len(url.Host) != 0 {
					result.ExternalLinksCount++
				} else {
					if len(url.Host) == 0 {
						url.Scheme = inputURL.Scheme
						url.Host = inputURL.Host
					}
					result.InternalLinksCount++
				}

				urlString := url.String()
				accessible, ok := linkCache[urlString]

				if !ok {
					res, err := request.Get(urlString)
					accessible = err == nil && res.StatusCode < 400
					linkCache[urlString] = accessible
				}

				if !accessible {
					result.InaccessibleLinksCount++
				}

				break
			}
		}
	}
}

func checkLoginForm(result *models.Result, token html.Token) {
	if token.DataAtom.String() == "form" {
		for _, attr := range token.Attr {
			if attr.Key == "id" || attr.Key == "class" || attr.Key == "action" {
				normalizedID := strings.ToLower(attr.Val)
				if strings.Contains(normalizedID, "login") {
					result.ContainsLoginForm = true
				}
			}
		}
	}
}
