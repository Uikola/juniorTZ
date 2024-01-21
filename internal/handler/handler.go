package handler

import (
	"context"
	"database/sql"
	"github.com/Uikola/juniorTZ/pkg/filter"

	"github.com/Uikola/juniorTZ/internal/db/repository/postgres"
	"github.com/Uikola/juniorTZ/internal/entity"
	"github.com/Uikola/juniorTZ/internal/usecase/person_usecase"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

type UseCase interface {
	AddPerson(ctx context.Context, human entity.Person) (entity.Person, error)
	DeletePerson(ctx context.Context, personID int) error
	GetPeople(ctx context.Context, options person_usecase.GetPeopleOptions) ([]entity.Person, error)
	UpdatePerson(ctx context.Context, human entity.Person) (int, error)
}

type Options interface {
	IsToApply() bool
	AddField(name, operator, value, dType string) error
	Fields() []filter.Field
}

type Handler struct {
	UseCase UseCase
	log     zerolog.Logger
}

func New(useCase UseCase, log zerolog.Logger) *Handler {
	return &Handler{
		UseCase: useCase,
		log:     log,
	}
}

func Router(db *sql.DB, router chi.Router, log zerolog.Logger) {
	repo := postgres.New(db)
	useCase := person_usecase.New(repo)
	handler := New(useCase, log)

	router.Get("/people", handler.GetPeople)
	router.Post("/person", handler.AddPerson)
	router.Put("/person/{person_id}", handler.UpdatePerson)
	router.Delete("/person/{person_id}", handler.DeletePerson)
}
