package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"go-postgres-template/internal/store"
)

func setupTestDB(t *testing.T) *store.Queries {
	ctx := context.Background()
	pgContainer, err := postgres.RunContainer(ctx,
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, pgContainer.Terminate(ctx))
	})

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	db, err := sql.Open("pgx", connStr)
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
		CREATE TABLE users (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`)
	require.NoError(t, err)

	return store.New(db)
}

func TestCreateUser(t *testing.T) {
	queries := setupTestDB(t)
	handler := NewUserHandler(queries)

	body := `{"name": "Alice", "email": "alice@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateUser(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var user store.User
	err := json.Unmarshal(w.Body.Bytes(), &user)
	require.NoError(t, err)
	require.Equal(t, "Alice", user.Name)
	require.Equal(t, "alice@example.com", user.Email)
	require.NotEqual(t, uuid.Nil, user.ID)
}
