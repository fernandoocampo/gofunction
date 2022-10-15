package fruits

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/google/uuid"
)

type NewFruit struct {
	SourceID string  `json:"source_id" dynamodbav:"source_id"`
	Name     string  `json:"name" dynamodbav:"name"`
	Variety  string  `json:"variety" dynamodbav:"variety"`
	Price    float32 `json:"price" dynamodbav:"price"`
}

type Fruit struct {
	ID       string  `json:"id" dynamodbav:"id"`
	SourceID string  `json:"source_id" dynamodbav:"source_id"`
	Name     string  `json:"name" dynamodbav:"name"`
	Variety  string  `json:"variety" dynamodbav:"variety"`
	Price    float32 `json:"price" dynamodbav:"price"`
}

type Storer interface {
	Save(ctx context.Context, fruit Fruit) error
}

type Service struct {
	repo Storer
}

var (
	errCreatingService    = errors.New("unable to create service")
	errUnmarshallingFruit = errors.New("unable to unmarsall fruit")
)

func NewService(region, endpoint string) (*Service, error) {
	ctx := context.Background()
	dbSetup := Setup{Region: region, Endpoint: endpoint}

	repo, err := newDynamoDBClient(ctx, dbSetup)
	if err != nil {
		return nil, errCreatingService
	}

	newService := Service{
		repo: repo,
	}

	return &newService, nil
}

func NewServiceWithStorer(storage Storer) *Service {
	newService := Service{
		repo: storage,
	}

	return &newService
}

func (s *Service) Audit(ctx context.Context, fruitstring string) {
	aNewfruit, err := newFruit(fruitstring)
	if err != nil {
		return
	}

	newid := uuid.New().String()

	fruit := aNewfruit.toFruit(newid)

	_ = s.repo.Save(ctx, fruit)
}

// toFruit transforms new fruit to a fruit.
func (n NewFruit) toFruit(fruitID string) Fruit {
	return Fruit{
		ID:       fruitID,
		SourceID: n.SourceID,
		Name:     n.Name,
		Variety:  n.Variety,
		Price:    n.Price,
	}
}

func newFruit(fruitstring string) (NewFruit, error) {
	var fruit NewFruit

	if fruitstring == "" {
		return fruit, nil
	}

	err := json.Unmarshal([]byte(fruitstring), &fruit)
	if err != nil {
		log.Println("unable to unmarshall fruit", "fruit", fruitstring, "error", err)

		return fruit, errUnmarshallingFruit
	}

	return fruit, nil
}
