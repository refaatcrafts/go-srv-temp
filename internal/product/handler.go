package product

import (
	"encoding/json"
	"net/http"

	"go-srv-temp/internal/httperr"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperr.RespondError(w, httperr.ErrBadRequest)
		return
	}

	p, err := h.svc.Create(r.Context(), req)
	if err != nil {
		httperr.RespondError(w, err)
		return
	}

	httperr.RespondJSON(w, http.StatusCreated, p)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httperr.RespondError(w, httperr.New(400, "invalid id"))
		return
	}

	p, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		httperr.RespondError(w, err)
		return
	}

	httperr.RespondJSON(w, http.StatusOK, p)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	var categoryID *uuid.UUID
	if cid := r.URL.Query().Get("category_id"); cid != "" {
		id, err := uuid.Parse(cid)
		if err != nil {
			httperr.RespondError(w, httperr.New(400, "invalid category_id"))
			return
		}
		categoryID = &id
	}

	pp, err := h.svc.List(r.Context(), categoryID)
	if err != nil {
		httperr.RespondError(w, err)
		return
	}

	httperr.RespondJSON(w, http.StatusOK, pp)
}
