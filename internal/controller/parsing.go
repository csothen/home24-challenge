package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/csothen/htmlparser/internal/service"
	"github.com/csothen/htmlparser/pkg/models"
	"github.com/csothen/htmlparser/pkg/validate"
)

// ParsingController : Handles parsing related HTTP requests
type ParsingController struct {
	l       *log.Logger
	service *service.ParsingService
}

// NewParsingController : Creates a new instance of ParsingController
func NewParsingController(l *log.Logger) *ParsingController {
	service := service.NewParsingService()
	return &ParsingController{l, service}
}

// ParseWebsite : Handles the Parsing of a website based on a certain URL
func (pc *ParsingController) ParseWebsite(rw http.ResponseWriter, r *http.Request) {

	// Decode request body
	input := new(models.ParsingRequest)

	d := json.NewDecoder(r.Body)
	if err := d.Decode(input); err != nil {
		http.Error(rw, "Unable to parse request body", 400)
		return
	}

	// Validate URL
	valid := validate.IsValidURL(input.URL)

	if !valid {
		http.Error(rw, "The URL is invalid", 400)
		return
	}

	// Process request
	pc.l.Println("Parsing URL", input.URL)

	result, err := pc.service.Parse(input.URL)
	if err != nil {
		http.Error(rw, "Failed to parse website", 500)
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	e := json.NewEncoder(rw)
	err = e.Encode(result)
	if err != nil {
		http.Error(rw, "Failed to encode result", 500)
		return
	}
}
