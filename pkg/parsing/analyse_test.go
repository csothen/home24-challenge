package parsing

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/csothen/htmlparser/pkg/models"
	"golang.org/x/net/html"
)

type TestCase struct {
	input struct {
		url  url.URL
		page *html.Tokenizer
	}
	expected *models.Result
}

func TestAnalyse(t *testing.T) {
	websiteDoctype := `<!DOCTYPE html>`

	websiteHead := `
		<head>
    		<meta charset="utf-8" />
		    <meta name="referrer" content="origin-when-crossorigin" id="meta_referrer" />
			<title id="pageTitle">Facebook</title>
		</head>
	`

	websiteBody := `
		<body>
		    <div>
		        <div>
					<h1>Hello World!</h1>
					<a href="https://invalidlink_asd.com">Invalid external Link</a>
					<a href="https://google.com">Valid external link</a>
					<a href="/r.php">Valid internal link</a>
					<a href="/nasdkjfnakj">Invalid internal link</a>
					<div>
						<h2>Another header a level deeper</h2>
					</div>
				</div>
		    </div>
		</body>
		`

	websiteBodyWithForm := `
		<body>
		    <div>
		        <div>
					<h1>Hello World!</h1>
					<a href="https://invalidlink_asd.com">Invalid external Link</a>
					<a href="https://google.com">Valid external link</a>
					<a href="/r.php">Valid internal link</a>
					<a href="/nasdkjfnakj">Invalid internal link</a>
					<div>
						<h2>Another header a level deeper</h2>
					</div>
					<form id="login form" class="form" method="post">
					</form>
				</div>
		    </div>
		</body>
	`

	completeWebsite := fmt.Sprintf("%s\n<html>\n%s\n%s\n</html>", websiteDoctype, websiteHead, websiteBodyWithForm)
	noDoctypeWebsite := fmt.Sprintf("<html>\n%s\n%s\n</html>", websiteHead, websiteBodyWithForm)
	noTitleWebsite := fmt.Sprintf("%s\n<html>\n%s\n</html>", websiteDoctype, websiteBodyWithForm)
	noBodyWebsite := fmt.Sprintf("%s\n<html>\n%s\n</html>", websiteDoctype, websiteHead)
	noFormWebsite := fmt.Sprintf("%s\n<html>\n%s\n%s\n</html>", websiteDoctype, websiteHead, websiteBody)

	urlString := "https://facebook.com"
	urlObject, _ := url.ParseRequestURI(urlString)

	cases := []TestCase{
		{
			input: struct {
				url  url.URL
				page *html.Tokenizer
			}{
				url:  *urlObject,
				page: html.NewTokenizer(strings.NewReader(completeWebsite)),
			},
			expected: &models.Result{
				HTMLVersion:            "HTML 5",
				PageTitle:              "Facebook",
				HeadingsCount:          map[int]int{5: 1, 6: 1},
				InternalLinksCount:     2,
				ExternalLinksCount:     2,
				InaccessibleLinksCount: 2,
				ContainsLoginForm:      true,
			},
		},
		{
			input: struct {
				url  url.URL
				page *html.Tokenizer
			}{
				url:  *urlObject,
				page: html.NewTokenizer(strings.NewReader(noDoctypeWebsite)),
			},
			expected: &models.Result{
				HTMLVersion:            "Failed to find Doctype declaration",
				PageTitle:              "Facebook",
				HeadingsCount:          map[int]int{5: 1, 6: 1},
				InternalLinksCount:     2,
				ExternalLinksCount:     2,
				InaccessibleLinksCount: 2,
				ContainsLoginForm:      true,
			},
		},
		{
			input: struct {
				url  url.URL
				page *html.Tokenizer
			}{
				url:  *urlObject,
				page: html.NewTokenizer(strings.NewReader(noTitleWebsite)),
			},
			expected: &models.Result{
				HTMLVersion:            "HTML 5",
				PageTitle:              "Failed to find the website's title",
				HeadingsCount:          map[int]int{5: 1, 6: 1},
				InternalLinksCount:     2,
				ExternalLinksCount:     2,
				InaccessibleLinksCount: 2,
				ContainsLoginForm:      true,
			},
		},
		{
			input: struct {
				url  url.URL
				page *html.Tokenizer
			}{
				url:  *urlObject,
				page: html.NewTokenizer(strings.NewReader(noFormWebsite)),
			},
			expected: &models.Result{
				HTMLVersion:            "HTML 5",
				PageTitle:              "Facebook",
				HeadingsCount:          map[int]int{5: 1, 6: 1},
				InternalLinksCount:     2,
				ExternalLinksCount:     2,
				InaccessibleLinksCount: 2,
				ContainsLoginForm:      false,
			},
		},
		{
			input: struct {
				url  url.URL
				page *html.Tokenizer
			}{
				url:  *urlObject,
				page: html.NewTokenizer(strings.NewReader(noBodyWebsite)),
			},
			expected: &models.Result{
				HTMLVersion:            "HTML 5",
				PageTitle:              "Facebook",
				HeadingsCount:          map[int]int{},
				InternalLinksCount:     0,
				ExternalLinksCount:     0,
				InaccessibleLinksCount: 0,
				ContainsLoginForm:      false,
			},
		},
	}

	for _, c := range cases {
		got := Analyse(c.input.url, c.input.page)
		if !reflect.DeepEqual(got, c.expected) {
			prettyGot, _ := json.MarshalIndent(got, "", "    ")
			prettyExpected, _ := json.MarshalIndent(c.expected, "", "    ")
			t.Log(fmt.Sprintf("should be\n%s but got\n%s", prettyExpected, prettyGot))
			t.Fail()
		}
	}
}
