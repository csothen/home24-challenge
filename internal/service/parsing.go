package service

import (
	"net/url"

	"github.com/csothen/htmlparser/pkg/models"
	"github.com/csothen/htmlparser/pkg/parsing"
	"github.com/csothen/htmlparser/pkg/request"
	"golang.org/x/net/html"
)

// ParsingService : Implements methods to handle the parsing of a website
type ParsingService struct{}

// NewParsingService : Creates a new instance of ParsingService that handles the business logic
func NewParsingService() *ParsingService {
	return &ParsingService{}
}

// Parse : Parses the website's contents and analyzes it, returning the result
func (service *ParsingService) Parse(url url.URL) (*models.Result, error) {

	response, err := request.Get(url.String())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	page := html.NewTokenizer(response.Body)

	result := parsing.Analyse(url, page)

	return result, nil
}
