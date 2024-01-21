package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h Handler) DeletePerson(w http.ResponseWriter, r *http.Request) {
	h.log.Info().Msg("received an DeletePerson DELETE request")
	ctx := r.Context()

	humanID, err := strconv.Atoi(chi.URLParam(r, "person_id"))
	if err != nil {
		h.log.Error().Err(err).Msg("invalid person id")
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	err = h.UseCase.DeletePerson(ctx, humanID)
	if err != nil {
		h.log.Error().Err(err).Msg("error while deleting a person")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	h.log.Info().Msg("the request was successfully completed")
	w.WriteHeader(http.StatusNoContent)
}
