package responses

import (
	"server/domain/models"
)

type SearchEmailsResponse struct {
	Hits   int                 `json:"hits"`
	Emails []models.EmailModel `json:"emails"`
}
