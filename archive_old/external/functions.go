package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type CreateInvoiceRequest struct {
	ExternalID  string `json:"external_id"`
	Amount      int    `json:"amount"`
	PayerEmail  string `json:"payer_email"`
	Description string `json:"description"`
}

type CreateInvoiceResponse struct {
	ID                        string        `json:"id"`
	ExternalID                string        `json:"external_id"`
	UserID                    string        `json:"user_id"`
	Status                    string        `json:"status"`
	MerchantName              string        `json:"merchant_name"`
	MerchantProfilePictureUrl string        `json:"merchant_profile_picture_url"`
	Amount                    int           `json:"amount"`
	PayerEmail                string        `json:"payer_email"`
	Description               string        `json:"description"`
	ExpiryDate                string        `json:"expiry_date"`
	InvoiceURL                string        `json:"invoice_url"`
	AvailableBanks            []interface{} `json:"available_banks"`
}

type InvoiceResponse struct {
	ID                        string        `json:"id"`
	ExternalID                string        `json:"external_id"`
	UserID                    string        `json:"user_id"`
	PaymentMethod             string        `json:"payment_method"`
	Status                    string        `json:"status"`
	MerchantName              string        `json:"merchant_name"`
	MerchantProfilePictureUrl string        `json:"merchant_profile_picture_url"`
	Amount                    int           `json:"amount"`
	PaidAmout                 int           `json:"paid_amount"`
	PaidAt                    time.Time     `json:"paid_at"`
	PayerEmail                string        `json:"payer_email"`
	Description               string        `json:"description"`
	ExpiryDate                string        `json:"expiry_date"`
	InvoiceURL                string        `json:"invoice_url"`
	AvailableBanks            []interface{} `json:"available_banks"`
}

func CreateInvoice(externalId string, amount int, payerEmail string, description string) (InvoiceResponse, error) {
	request := CreateInvoiceRequest{
		ExternalID:  externalId,
		Amount:      amount,
		PayerEmail:  payerEmail,
		Description: description,
	}

	url := "https://api.xendit.co/v2/invoices"
	client := &http.Client{}

	reqBody, err := json.Marshal(request)
	if err != nil {
		return InvoiceResponse{}, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return InvoiceResponse{}, err
	}
	apiKey := os.Getenv("XENDIT_API_KEY")
	if apiKey == "" {
		log.Printf("XENDIT_API_KEY is not set in the environment")
		return InvoiceResponse{}, fmt.Errorf("XENDIT_API_KEY is not set in the environment")
	}
	req.SetBasicAuth(apiKey, "")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return InvoiceResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("failed to create invoice, status code: %d", resp.StatusCode)
		return InvoiceResponse{}, fmt.Errorf("failed to create invoice, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read response body: %v", err)
		return InvoiceResponse{}, err
	}

	var createInvoiceResponse InvoiceResponse
	if err := json.Unmarshal(body, &createInvoiceResponse); err != nil {
		log.Printf("failed to unmarshal response body: %v", err)
		return InvoiceResponse{}, err
	}

	return createInvoiceResponse, nil
}

func GetInvoice(invoiceID string) (InvoiceResponse, error) {
	url := "https://api.xendit.co/v2/invoices/" + invoiceID
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return InvoiceResponse{}, err
	}
	apiKey := os.Getenv("XENDIT_API_KEY")
	if apiKey == "" {
		log.Printf("XENDIT_API_KEY is not set in the environment")
		return InvoiceResponse{}, fmt.Errorf("XENDIT_API_KEY is not set in the environment")
	}
	req.SetBasicAuth(apiKey, "")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return InvoiceResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read response body: %v", err)
		return InvoiceResponse{}, err
	}

	var getInvoiceResponse InvoiceResponse
	if err := json.Unmarshal(body, &getInvoiceResponse); err != nil {
		log.Printf("failed to unmarshal response body: %v", err)
		return InvoiceResponse{}, err
	}

	return getInvoiceResponse, nil
}

func GetInvoiceStatus(invoiceID string) (string, error) {
	invoice, err := GetInvoice(invoiceID)
	if err != nil {
		return "", err
	}
	return invoice.Status, nil
}
