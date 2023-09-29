package dtos

import (
	"errors"
	"net/http"
)

const (
	SearchEmailsDtoLimitLt  = 0
	SearchEmailsDtoLimitGte = 50
	SearchEmailsDtoPageLt   = 0
)

const (
	SearchEmailsDtoLimitLtError      = "page must be greater than 0"
	SearchEmailsDtoLimitGteError     = "page must be less than or equal to 50"
	SearchEmailsDtoPageLtError       = "limit must be greater than 0"
	SearchEmailsDtoTermRequiredError = "term is required"
)

type SearchEmailsDto struct {
	Limit int    `json:"limit"`
	Page  int    `json:"page"`
	Term  string `json:"term"`
}

func (searchEmailsDto *SearchEmailsDto) Bind(r *http.Request) error {
	if searchEmailsDto.Limit <= SearchEmailsDtoLimitLt {
		return errors.New(SearchEmailsDtoLimitLtError)
	} else if searchEmailsDto.Limit > SearchEmailsDtoLimitGte {
		return errors.New(SearchEmailsDtoLimitGteError)
	}

	if searchEmailsDto.Page <= SearchEmailsDtoPageLt {
		return errors.New(SearchEmailsDtoPageLtError)
	}

	if searchEmailsDto.Term == "" {
		return errors.New(SearchEmailsDtoTermRequiredError)
	}

	return nil
}
