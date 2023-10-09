package dtos

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"server/application/use-cases/enums"
)

type SearchEmailsDto struct {
	Limit      int    `json:"limit"`
	Page       int    `json:"page"`
	Term       string `json:"term"`
	SearchType string `json:"searchType"`
}

var searchTypes = []string{
	enums.SearchTypeMatch,
	enums.SearchTypeMatchPhrase,
	enums.SearchTypeTerm,
	enums.SearchTypeQueryString,
	enums.SearchTypePrefix,
	enums.SearchTypeWildcard,
	enums.SearchTypeFuzzy,
}

const (
	SearchEmailsDtoPageLt   = 0
	SearchEmailsDtoLimitLt  = 0
	SearchEmailsDtoLimitGte = 50
)

const (
	SearchEmailsDtoPageLtError                = "limit must be greater than 0"
	SearchEmailsDtoLimitLtError               = "page must be greater than 0"
	SearchEmailsDtoLimitGteError              = "page must be less than or equal to 50"
	SearchEmailsDtoTermRequiredError          = "term is required"
	SearchEmailsDtoSearchTypeRequiredError    = "searchType is required"
	SearchEmailsDtoSearchTypeMustBeOneOfError = "searchType must be one of"
)

func (searchEmailsDto *SearchEmailsDto) Bind(r *http.Request) error {
	if searchEmailsDto.Page <= SearchEmailsDtoPageLt {
		return errors.New(SearchEmailsDtoPageLtError)
	}

	if searchEmailsDto.Limit <= SearchEmailsDtoLimitLt {
		return errors.New(SearchEmailsDtoLimitLtError)
	} else if searchEmailsDto.Limit > SearchEmailsDtoLimitGte {
		return errors.New(SearchEmailsDtoLimitGteError)
	}

	if searchEmailsDto.Term == "" {
		return errors.New(SearchEmailsDtoTermRequiredError)
	}

	if searchEmailsDto.SearchType == "" {
		return errors.New(SearchEmailsDtoSearchTypeRequiredError)
	}

	if !slices.Contains(searchTypes, searchEmailsDto.SearchType) {
		errorStr := fmt.Sprintf(SearchEmailsDtoSearchTypeMustBeOneOfError+" [%v]", strings.Join(searchTypes, ", "))

		return errors.New(errorStr)
	}

	return nil
}
