package emails

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"server/adapter/controllers/emails/dtos"
	"server/adapter/controllers/emails/responses"
	"server/adapter/controllers/types"
	usecases "server/application/use-cases"
)

type EmailsController struct {
	emailsUseCase usecases.IEmailsUseCase
}

type IEmailsController interface {
	SearchEmails(http.ResponseWriter, *http.Request)
}

func NewEmailsController(server *chi.Mux, emailsUseCase usecases.IEmailsUseCase) {
	controller := &EmailsController{
		emailsUseCase,
	}

	server.Route("/emails", func(r chi.Router) {
		r.Post("/search", controller.SearchEmails)
	})
}

// SearchEmails godoc
//
//	@Summary			Search emails
//	@Description	Endpoint to search emails by term, limit, and page
//	@Tags 				search emails
//	@Accept				json
//	@Produce			json
//	@Param				searchEmailsDto body dtos.SearchEmailsDto true "Search emails DTO"
//	@Success			200	{object}	types.SuccessResponse[responses.SearchEmailsResponse]{}
//	@Failure			400	{object}	types.ErrorResponse{}
//	@Failure			500	{object}	types.ErrorResponse{}
//	@Router				/emails/search	[post]

func (controller *EmailsController) SearchEmails(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength <= 0 {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, types.ErrorResponse{
			Error: "body request is missing",
		})

		return
	}

	searchEmailsDto := &dtos.SearchEmailsDto{}

	err := render.Bind(r, searchEmailsDto)

	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, types.ErrorResponse{
			Error: err.Error(),
		})

		return
	}

	searchEmails, err := controller.emailsUseCase.SearchEmails(searchEmailsDto)

	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, types.ErrorResponse{
			Error: err.Error(),
		})

		return
	}

	render.JSON(w, r, types.SuccessResponse[responses.SearchEmailsResponse]{
		Data: responses.SearchEmailsResponse{
			Hits:   searchEmails.Hits,
			Emails: searchEmails.Emails,
		},
	})
}
