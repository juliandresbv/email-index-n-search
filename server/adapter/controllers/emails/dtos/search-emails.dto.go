package dtos

import (
	"errors"
	"net/http"
)

type SearchEmailsDto struct {
	Limit int    `json:"limit"`
	Page  int    `json:"page"`
	Term  string `json:"term"`
}

const (
	SearchEmailsDtoPageLt   = 0
	SearchEmailsDtoLimitLt  = 0
	SearchEmailsDtoLimitGte = 50
)

const (
	SearchEmailsDtoPageLtError       = "limit must be greater than 0"
	SearchEmailsDtoLimitLtError      = "page must be greater than 0"
	SearchEmailsDtoLimitGteError     = "page must be less than or equal to 50"
	SearchEmailsDtoTermRequiredError = "term is required"
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

	return nil
}
