package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Uikola/juniorTZ/internal/entity"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type UpdatePersonRequest struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
	Age        int    `json:"age"`
	Gender     string `json:"gender"`
	Nation     string `json:"nation"`
}

func (h Handler) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	h.log.Info().Msg("received an UpdatePerson PUT request")
	var req UpdatePersonRequest
	ctx := r.Context()
	id, err := strconv.Atoi(chi.URLParam(r, "person_id"))
	if err != nil {
		h.log.Error().Err(err).Msg("invalid person id")
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error().Err(err).Msg("error while decoding")
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	human := entity.Person{
		ID:         id,
		Name:       req.Name,
		Surname:    req.Surname,
		Patronymic: req.Patronymic,
		Age:        req.Age,
		Gender:     req.Gender,
		Nation:     req.Nation,
	}
	humanID, err := h.UseCase.UpdatePerson(ctx, human)
	if err != nil {
		h.log.Error().Err(err).Msg("error while updating a person")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	h.log.Info().Msg("the request was successfully completed")
	render.JSON(w, r, humanID)
}
