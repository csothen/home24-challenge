package models

/**
HTMLVersion : Version being used in the website
PageTitle : Title of the website
HeadingsCount : Count of headings used by level of Html
InternalLinksCount : Count of links that lead to the same domain
ExternalLinksCount : Count of links that lead to a different domain
InaccessibleLinksCount : Count of links that are not accessible to an unauthorized user
ContainsLoginForm : Check if the website contains a login form
*/

// Result : The result of the parsing and analyzation of the data in the website
type Result struct {
	HTMLVersion            string      `json:"htmlVersion"`
	PageTitle              string      `json:"pageTitle"`
	HeadingsCount          map[int]int `json:"headingsCount"`
	InternalLinksCount     int         `json:"internalLinksCount"`
	ExternalLinksCount     int         `json:"externalLinksCount"`
	InaccessibleLinksCount int         `json:"inaccessibleLinksCount"`
	ContainsLoginForm      bool        `json:"containsLoginForm"`
}

// ParsingRequest : Holds the input URL that will be used for the request
type ParsingRequest struct {
	URL *string `json:"url"`
}
