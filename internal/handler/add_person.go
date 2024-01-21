package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Uikola/juniorTZ/internal/entity"
	"github.com/go-chi/render"
)

type AddPersonRequest struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
}

func (h Handler) AddPerson(w http.ResponseWriter, r *http.Request) {
	h.log.Info().Msg("received an AddPerson POST request")

	var req AddPersonRequest
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error().Err(err).Msg("error while decoding")
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	human := entity.Person{
		Name:       req.Name,
		Surname:    req.Surname,
		Patronymic: req.Patronymic,
	}

	person, err := h.UseCase.AddPerson(ctx, human)
	if err != nil {
		h.log.Error().Err(err).Msg("error while adding a person")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	h.log.Info().Msg("the request was successfully completed")
	render.JSON(w, r, person)
}
