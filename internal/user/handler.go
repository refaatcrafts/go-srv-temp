package user

import (
	"encoding/json"
	"net/http"

	"go-srv-temp/internal/httperr"
	"go-srv-temp/internal/middleware"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperr.RespondError(w, httperr.ErrBadRequest)
		return
	}

	res, err := h.svc.Signup(r.Context(), req)
	if err != nil {
		httperr.RespondError(w, err)
		return
	}

	httperr.RespondJSON(w, http.StatusCreated, res)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperr.RespondError(w, httperr.ErrBadRequest)
		return
	}

	res, err := h.svc.Login(r.Context(), req)
	if err != nil {
		httperr.RespondError(w, err)
		return
	}

	httperr.RespondJSON(w, http.StatusOK, res)
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		httperr.RespondError(w, httperr.ErrUnauthorized)
		return
	}

	u, err := h.svc.GetByID(r.Context(), userID)
	if err != nil {
		httperr.RespondError(w, err)
		return
	}

	httperr.RespondJSON(w, http.StatusOK, u)
}
