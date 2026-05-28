package handler

import (
	"encoding/json"
	"net/http"
	"payment_integration/internal/payment/usecases/create_invoice"
	"payment_integration/internal/transport/http_transport"
)

func (h *Handler) CreateInvoice() http.HandlerFunc {
	return http_transport.WrapHandler(h.logger, func(w http.ResponseWriter, r *http.Request) error {
		var dto create_invoice.CreateInvoiceRequest
		if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
			return err
		}
		if err := h.validate.Struct(dto); err != nil {
			return err
		}
		response, err := h.createInvoiceUseCase.Execute(r.Context(), dto)
		if err != nil {
			return err
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			h.logger.Error("Error encoding response", "error", err)
		}
		return nil
	})
}