package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Uikola/juniorTZ/internal/usecase/person_usecase"
	"github.com/Uikola/juniorTZ/pkg/filter"
	"github.com/go-chi/render"
)

func (h Handler) GetPeople(w http.ResponseWriter, r *http.Request) {
	h.log.Info().Msg("received an GetPeople GET request")
	var limit, offset int
	ctx := r.Context()

	offset, _ = strconv.Atoi(r.URL.Query().Get("offset"))
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 1
	}

	filterOptions, err := FilterOptions(r)
	if err != nil {
		h.log.Error().Err(err).Msg("error while getting filter options")
		http.Error(w, "invalid options", http.StatusBadRequest)
		return
	}

	options := person_usecase.GetPeopleOptions{
		FilterOptions: filterOptions,
		Limit:         limit,
		Offset:        offset,
	}
	humans, err := h.UseCase.GetPeople(ctx, options)
	if err != nil {
		h.log.Error().Err(err).Msg("error while getting ")
		http.Error(w, "internal error", http.StatusBadRequest)
		return
	}

	h.log.Info().Msg("the request was successfully completed")
	render.JSON(w, r, humans)
}

func FilterOptions(r *http.Request) (Options, error) {
	filterOptions := filter.NewOptions()
	name := r.URL.Query().Get("name")
	if name != "" {
		err := filterOptions.AddField("name", filter.OperatorLike, name, filter.DataTypeStr)
		if err != nil {
			return nil, err
		}
	}

	surname := r.URL.Query().Get("surname")
	if surname != "" {
		err := filterOptions.AddField("surname", filter.OperatorLike, surname, filter.DataTypeStr)
		if err != nil {
			return nil, err
		}
	}

	patronymic := r.URL.Query().Get("patronymic")
	if patronymic != "" {
		err := filterOptions.AddField("patronymic", filter.OperatorLike, patronymic, filter.DataTypeStr)
		if err != nil {
			return nil, err
		}
	}

	age := r.URL.Query().Get("age")
	if age != "" {
		operator := filter.OperatorEq
		val := age
		if strings.Index(age, ":") != -1 {
			split := strings.Split(age, ":")
			operator = split[0]
			val = split[1]
		}
		err := filterOptions.AddField("age", operator, val, filter.DataTypeInt)
		if err != nil {
			return nil, err
		}
	}

	gender := r.URL.Query().Get("gender")
	if gender != "" {
		err := filterOptions.AddField("gender", filter.OperatorLike, gender, filter.DataTypeStr)
		if err != nil {
			return nil, err
		}
	}

	nation := r.URL.Query().Get("nation")
	if nation != "" {
		err := filterOptions.AddField("nation", filter.OperatorLike, nation, filter.DataTypeStr)
		if err != nil {
			return nil, err
		}
	}

	return filterOptions, nil
}
