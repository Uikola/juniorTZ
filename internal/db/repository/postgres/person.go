package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Uikola/juniorTZ/internal/entity"
	"github.com/Uikola/juniorTZ/internal/usecase/person_usecase"
	"github.com/Uikola/juniorTZ/pkg/filter"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r Repository) DeletePerson(ctx context.Context, personID int) error {
	const op = "PersonRepository.DeletePerson"
	query := `
	DELETE FROM people
	WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, personID)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}

func (r Repository) UpdatePerson(ctx context.Context, person entity.Person) (int, error) {
	const op = "PersonRepository.UpdatePerson"
	query := `
	UPDATE people
	SET name = $1,
	    surname = $2,
	    patronymic = $3,
	    age = $4,
	    gender = $5,
	    nation = $6
	WHERE id = $7
	RETURNING id`

	row := r.db.QueryRowContext(ctx, query, person.Name, person.Surname, person.Patronymic, person.Age, person.Gender, person.Nation, person.ID)

	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	return id, nil
}

func (r Repository) AddPerson(ctx context.Context, person entity.Person) (entity.Person, error) {
	const op = "PersonRepository.AddPerson"

	query := `
	INSERT INTO people(name, surname, patronymic, age, gender, nation)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, name, surname, patronymic, age, gender, nation`

	row := r.db.QueryRowContext(ctx, query, person.Name, person.Surname, person.Patronymic, person.Age, person.Gender, person.Nation)

	var name, surname, patronymic, gender, nation string
	var id, age int
	err := row.Scan(&id, &name, &surname, &patronymic, &age, &gender, &nation)
	if err != nil {
		return entity.Person{}, fmt.Errorf("%s:%w", op, err)
	}

	return entity.Person{
		ID:         id,
		Name:       name,
		Surname:    surname,
		Patronymic: patronymic,
		Age:        age,
		Gender:     gender,
		Nation:     nation,
	}, nil

}

func (r Repository) GetPeople(ctx context.Context, options person_usecase.GetPeopleOptions) ([]entity.Person, error) {
	const op = "PersonRepository.GetPeople"
	query := `
	SELECT *
	FROM people `
	var values []interface{}

	if options.FilterOptions.IsToApply() {
		query, values = r.addFilters(query, options.FilterOptions)
	}
	query = r.addPagination(query, len(options.FilterOptions.Fields()))
	values = append(values, options.Limit, options.Offset)

	rows, err := r.db.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	var name, surname, patronymic, gender, nation string
	var id, age int
	var persons []entity.Person
	for rows.Next() {
		err := rows.Scan(&id, &name, &surname, &patronymic, &age, &gender, &nation)
		if err != nil {
			return nil, fmt.Errorf("%s:%w", op, err)
		}
		person := entity.Person{
			ID:         id,
			Name:       name,
			Surname:    surname,
			Patronymic: patronymic,
			Age:        age,
			Gender:     gender,
			Nation:     nation,
		}
		persons = append(persons, person)
	}

	return persons, nil
}

func (r Repository) addFilters(query string, filterOptions person_usecase.Options) (string, []interface{}) {
	numCondMap := map[string]string{
		filter.OperatorEq:               "=",
		filter.OperatorNotEq:            "!=",
		filter.OperatorGreaterThan:      ">",
		filter.OperatorGreaterThanEqual: ">=",
		filter.OperatorLowerThan:        "<",
		filter.OperatorLowerThanEqual:   "<=",
	}
	var filterValues []interface{}

	filterFields := filterOptions.Fields()
	query += "WHERE "
	for i, filterField := range filterFields {
		if filterField.Type == filter.DataTypeStr {
			query += fmt.Sprintf("%s %s $%d AND ", filterField.Name, filterField.Operator, i+1)
		} else if filterField.Type == filter.DataTypeInt {
			query += fmt.Sprintf("%s %s $%d AND ", filterField.Name, numCondMap[filterField.Operator], i+1)
		}
		filterValues = append(filterValues, filterField.Value)
	}
	query = query[:len(query)-5]
	return query, filterValues
}

func (r Repository) addPagination(query string, ind int) string {
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", ind+1, ind+2)
	return query
}
