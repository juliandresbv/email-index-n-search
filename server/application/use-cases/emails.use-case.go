package usecases

import (
	"encoding/json"

	customlogger "pkg/custom-logger"
	searchdbclient "pkg/zinc-search-db/db-client/search-db-client"
	"pkg/zinc-search-db/db-client/search-db-client/dtos"

	controllersemailsdtos "server/adapter/controllers/emails/dtos"
	"server/adapter/controllers/emails/responses"
	"server/domain/models"
)

var logger = customlogger.NewLogger()

type EmailsUseCase struct {
	searchDbClient searchdbclient.ISearchDbClient
}

type IEmailsUseCase interface {
	SearchEmails(searchEmailsDto *controllersemailsdtos.SearchEmailsDto) (responses.SearchEmailsResponse, error)
}

func NewEmailsUseCase(searchDbClient searchdbclient.ISearchDbClient) IEmailsUseCase {
	return &EmailsUseCase{
		searchDbClient,
	}
}

func (emailsUseCase *EmailsUseCase) SearchEmails(
	searchEmailsDto *controllersemailsdtos.SearchEmailsDto,
) (responses.SearchEmailsResponse, error) {
	indexName := "emails"

	term := searchEmailsDto.Term
	from := (searchEmailsDto.Page - 1) * searchEmailsDto.Limit
	maxResults := searchEmailsDto.Limit
	searchType := searchEmailsDto.SearchType

	searchSearchV1Dto := dtos.SearchSearchV1Dto{
		SearchType: searchType,
		Query: dtos.QuerySearchSearchV1Dto{
			Term: term,
		},
		From:       from,
		MaxResults: maxResults,
		SortFields: []string{
			"-Date",
		},
	}

	response, err := emailsUseCase.searchDbClient.SearchV1(indexName, searchSearchV1Dto)

	if err != nil {
		logger.Println("Error on searchV1 ZincSearch API request: ", err)

		return responses.SearchEmailsResponse{
			Hits:   0,
			Emails: []models.EmailModel{},
		}, err
	}

	emails := []models.EmailModel{}
	hits := response.Data.Hits.Total.Value

	for _, hit := range response.Data.Hits.Hits {
		marshaledSource, _ := json.Marshal(hit.Source)

		var email models.EmailModel
		json.Unmarshal(marshaledSource, &email)
		email.Id = hit.Id

		emails = append(emails, email)
	}

	return responses.SearchEmailsResponse{
		Hits:   hits,
		Emails: emails,
	}, nil
}
