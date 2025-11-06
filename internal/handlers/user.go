package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"$PROJECT_NAME/internal/store"
)

type UserHandler struct {
	queries *store.Queries
	validate *validator.Validate
}

func NewUserHandler(queries *store.Queries) *UserHandler {
	return &UserHandler{
		queries: queries,
		validate: validator.New(),
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	user := struct {
		Name  string `validate:"required,min=2,max=100"`
		Email string `validate:"required,email"`
	}{input.Name, input.Email}

	if err := h.validate.Struct(user); err != nil {
		http.Error(w, "validation error: "+err.Error(), http.StatusBadRequest)
		return
	}

	dbUser, err := h.queries.CreateUser(r.Context(), store.CreateUserParams{
		Name:  input.Name,
		Email: input.Email,
	})
	if err != nil {
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"users_email_key\"" {
			http.Error(w, "email already exists", http.StatusConflict)
			return
		}
		log.Error().Err(err).Msg("failed to create user")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dbUser)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.queries.GetUserByID(r.Context(), userID)
	if err != nil {
		if err == store.ErrNoRows {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		log.Error().Err(err).Msg("failed to get user")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
