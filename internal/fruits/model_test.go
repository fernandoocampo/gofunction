package fruits_test

import (
	"context"
	"testing"

	"github.com/fernandoocampo/gofunction/internal/fruits"
	"github.com/stretchr/testify/assert"
)

func TestAudit(t *testing.T) {
	// Given
	newFruit := `{
	"source_id": "1d952b94-a5db-4d63-a500-b486dd96e8b2",
	"name": "lemon",
	"variety": "lima",
	"price": 2.50
}`

	expectedStoredFruits := 1
	expectedFruit := fruits.Fruit{
		SourceID: "1d952b94-a5db-4d63-a500-b486dd96e8b2",
		Name:     "lemon",
		Variety:  "lima",
		Price:    2.50,
	}

	ctx := context.TODO()

	store := &storerMock{}
	service := fruits.NewServiceWithStorer(store)

	// When
	service.Audit(ctx, newFruit)

	// Then
	assert.NotEmpty(t, store.data)
	assert.Equal(t, expectedStoredFruits, len(store.data))
	assert.NotEmpty(t, store.data[0].ID)
	store.data[0].ID = "" // guarantee idempotence
	assert.Equal(t, expectedFruit, store.data[0])
}

type storerMock struct {
	data []fruits.Fruit
}

func (s *storerMock) Save(ctx context.Context, fruit fruits.Fruit) error {
	s.data = append(s.data, fruit)

	return nil
}
