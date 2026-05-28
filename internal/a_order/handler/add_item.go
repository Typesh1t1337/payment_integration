package handler

import (
	"encoding/json"
	"net/http"
	"payment_integration/internal/a_order"
	"payment_integration/internal/transport/http_transport"

	"github.com/google/uuid"
)

func (h *Handler) AddItemHandler() http.HandlerFunc {
	return http_transport.WrapHandler(h.logger, func(w http.ResponseWriter, r *http.Request) error {
		var dto a_order.AddItemRequest
		if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
			return http_transport.NewErrorResponse(
				http.StatusUnprocessableEntity, "Request body is invalid",
			)
		}
		ctx := r.Context()
		userID := ctx.Value(http_transport.UserIDContextKey).(uuid.UUID)

		err := h.AddItemUseCase.Execute(ctx, userID, dto)
		if err != nil {
			switch err {
			case a_order.ErrOrderBeingProcessed:
				return http_transport.NewErrorResponse(http.StatusBadRequest, err.Error())
			case a_order.ErrProductNotFound:
				return http_transport.NewErrorResponse(http.StatusNotFound, err.Error())
			default:
				return &http_transport.ErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Unexpected error",
				}
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(map[string]string{
			"message": "Ok",
		}); err != nil {
			h.logger.Error("Error encoding response", "error", err)
		}
		return nil
	})
}
