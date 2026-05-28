package tiptop

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type TiptopService struct {
	Login    string
	Password string
	BaseURL  string
}

func NewTiptopService(login string, password string, baseURL string) *TiptopService {
	return &TiptopService{
		Login:    login,
		Password: password,
		BaseURL:  baseURL,
	}
}

type ChargeRequest struct {
	Amount int `json:"amount"`
	Currency string `json:"currency"`
	IpAddress string `json:"ip_address"`
	CardCryptogramPacket string `json:"card_cryptogram_packet"`
	Name string `json:"name"`
	InvoiceId string `json:"invoice_id"`
}

type ChargeResponse struct {
	Model *string `json:"model"`
	Success bool `json:"success"`
	Message *string `json:"message"`
}

func (s *TiptopService) Charge(ctx context.Context, request *ChargeRequest) (*ChargeResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", s.BaseURL+"/charge", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(s.Login, s.Password)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to charge: %s", resp.Status)
	}
	var response ChargeResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}