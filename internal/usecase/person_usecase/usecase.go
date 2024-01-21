package person_usecase

import (
	"context"
	"github.com/Uikola/juniorTZ/pkg/filter"

	"github.com/Uikola/juniorTZ/internal/entity"
)

type GetPeopleOptions struct {
	FilterOptions Options
	Limit         int
	Offset        int
}

type Options interface {
	IsToApply() bool
	AddField(name, operator, value, dType string) error
	Fields() []filter.Field
}

type Repository interface {
	DeletePerson(ctx context.Context, personID int) error
	UpdatePerson(ctx context.Context, person entity.Person) (int, error)
	AddPerson(ctx context.Context, person entity.Person) (entity.Person, error)
	GetPeople(ctx context.Context, options GetPeopleOptions) ([]entity.Person, error)
}

type UseCaseImp struct {
	repo Repository
}

func New(repo Repository) *UseCaseImp {
	return &UseCaseImp{repo: repo}
}
