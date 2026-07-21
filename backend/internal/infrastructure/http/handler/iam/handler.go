package iam

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ruveydagundogan/llm-decision-score/backend/internal/application/iam/dto"
	iammodel "github.com/ruveydagundogan/llm-decision-score/backend/internal/domain/iam/model"
	"github.com/ruveydagundogan/llm-decision-score/backend/internal/shared/middleware"
	"github.com/ruveydagundogan/llm-decision-score/backend/internal/shared/response"
)

type RegisterUseCase interface {
	Execute(ctx context.Context, email, password, firstName, lastName string) (string, *iammodel.User, error)
}

type LoginUseCase interface {
	Execute(ctx context.Context, email, password string) (string, *iammodel.User, error)
}

type GetProfileUseCase interface {
	Execute(ctx context.Context, userID string) (*iammodel.User, error)
}

type UpdateProfileUseCase interface {
	Execute(ctx context.Context, userID, firstName, lastName string) error
}

type Handler struct {
	registerUC   RegisterUseCase
	loginUC      LoginUseCase
	getProfileUC GetProfileUseCase
	updateUC     UpdateProfileUseCase
}

func NewHandler(
	registerUC RegisterUseCase,
	loginUC LoginUseCase,
	getProfileUC GetProfileUseCase,
	updateUC UpdateProfileUseCase,
) *Handler {
	return &Handler{
		registerUC:   registerUC,
		loginUC:      loginUC,
		getProfileUC: getProfileUC,
		updateUC:     updateUC,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" {
		response.BadRequest(w, "email, password, first_name and last_name are required")
		return
	}

	if len(req.Password) < 6 {
		response.BadRequest(w, "password must be at least 6 characters")
		return
	}

	token, user, err := h.registerUC.Execute(r.Context(), req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Created(w, dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Status:    user.Status,
		},
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		response.BadRequest(w, "email and password are required")
		return
	}

	token, user, err := h.loginUC.Execute(r.Context(), req.Email, req.Password)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Success(w, dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Status:    user.Status,
		},
	})
}

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		response.Unauthorized(w, "authentication required")
		return
	}

	user, err := h.getProfileUC.Execute(r.Context(), userID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Success(w, dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Status:    user.Status,
	})
}

func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		response.Unauthorized(w, "authentication required")
		return
	}

	var req dto.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if err := h.updateUC.Execute(r.Context(), userID, req.FirstName, req.LastName); err != nil {
		response.Error(w, err)
		return
	}

	response.SuccessWithMessage(w, "profile updated")
}
