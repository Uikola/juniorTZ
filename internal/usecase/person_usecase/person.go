package person_usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Uikola/juniorTZ/internal/entity"
)

func (uc UseCaseImp) AddPerson(ctx context.Context, person entity.Person) (entity.Person, error) {
	age, err := EnrichAge(person.Name, "https://api.agify.io")
	if err != nil {
		return entity.Person{}, err
	}

	person.Age = age
	gender, err := EnrichGender(person.Name, "https://api.genderize.io")
	if err != nil {
		return entity.Person{}, err
	}

	person.Gender = gender
	nation, err := EnrichNation(person.Name, "https://api.nationalize.io")
	if err != nil {
		return entity.Person{}, err
	}

	person.Nation = nation
	return uc.repo.AddPerson(ctx, person)
}

func (uc UseCaseImp) DeletePerson(ctx context.Context, personID int) error {
	return uc.repo.DeletePerson(ctx, personID)
}

func (uc UseCaseImp) GetPeople(ctx context.Context, options GetPeopleOptions) ([]entity.Person, error) {
	return uc.repo.GetPeople(ctx, options)
}

func (uc UseCaseImp) UpdatePerson(ctx context.Context, human entity.Person) (int, error) {
	return uc.repo.UpdatePerson(ctx, human)
}

type EnrichmentResult struct {
	Country []Country `json:"country"`
	Gender  string    `json:"gender"`
	Age     int       `json:"age"`
}
type Country struct {
	CountryID string `json:"country_id"`
}

func EnrichNation(name string, url string) (string, error) {
	var result EnrichmentResult
	err := EnrichData(name, url, &result)
	if err != nil {
		return "", err
	}
	return result.Country[0].CountryID, nil
}

func EnrichGender(name string, url string) (string, error) {
	var result EnrichmentResult
	err := EnrichData(name, url, &result)
	if err != nil {
		return "", err
	}
	return result.Gender, nil
}

func EnrichAge(name string, url string) (int, error) {
	var result EnrichmentResult
	err := EnrichData(name, url, &result)
	if err != nil {
		return 0, err
	}
	return result.Age, nil
}

func EnrichData(name string, url string, result *EnrichmentResult) error {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?name=%s", url, name), nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return err
	}

	return nil
}
