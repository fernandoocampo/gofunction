package fruits_test

import (
	"context"
	"flag"
	"testing"

	"github.com/fernandoocampo/gofunction/internal/fruits"
	"github.com/stretchr/testify/assert"
)

var integration = flag.Bool("integration", false, "run integration tests")

func TestAuditDynamodb(t *testing.T) {
	if !*integration {
		t.Skip("it's an integration test")
	}
	// Given
	newFruit := `{
	"source_id": "1d952b94-a5db-4d63-a500-b486dd96e8b2",
	"name": "lemon",
	"variety": "lima",
	"price": 2.50
}`
	region := "us-east-1"
	endpoint := "http://localhost:4566"

	ctx := context.TODO()

	service, err := fruits.NewService(region, endpoint)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	// When
	service.Audit(ctx, newFruit)

	// Then
	assert.True(t, false)
}
