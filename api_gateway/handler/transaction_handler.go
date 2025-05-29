package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"crowdfund/external"
	"crowdfund/model"
	"crowdfund/repository"
)

type TransactionHandler interface {
	GetAllTransaction(c echo.Context) error
	CreateTransaction(c echo.Context) error
	GetTransactionByID(c echo.Context) error
	UpdateTransaction(c echo.Context) error
	CheckUpdateTransaction(c echo.Context) error
}

type transactionHandler struct {
	transactionRepo repository.TransactionRepository
}

func NewTransactionHandler(transactionRepo repository.TransactionRepository) TransactionHandler {
	return &transactionHandler{transactionRepo: transactionRepo}
}

// GetAllTransactions godoc
// @Summary Get all transactions for a user
// @Description Get all transactions for an active user
// @Tags transactions
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <access_token>"
// @Success 200 {object} model.Response
// @Router /transactions/{id} [get]
func (h *transactionHandler) GetAllTransaction(c echo.Context) error {
	//get user id from context
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, model.Response{
			Status:  http.StatusUnauthorized,
			Message: "User not authenticated",
		})
	}

	userIdFloat, ok := userID.(float64)
	if !ok {
		return c.JSON(http.StatusBadRequest, model.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid user ID",
		})
	}
	userIdInt := int(userIdFloat)

	transactions, err := h.transactionRepo.GetAllTransaction(userIdInt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{
			Status:  http.StatusInternalServerError,
			Message: "Internal Server Error",
		})
	}
	return c.JSON(200, model.Response{
		Status:  200,
		Message: "Success",
		Data:    transactions,
	})
}

// CreateTransaction godoc
// @Summary Create a new transaction
// @Description Create a new transaction for a user
// @Tags transactions
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <access_token>"
// @Param model.Transaction body model.Transaction true "Transaction object"
// @Success 201 {object} model.Response
// @Router /transactions [post] // Updated the router path to use POST method
func (h *transactionHandler) CreateTransaction(c echo.Context) error {
	//
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, model.Response{
			Status:  http.StatusUnauthorized,
			Message: "User not authenticated",
		})
	}

	userIdFloat, ok := userID.(float64)
	if !ok {
		return c.JSON(http.StatusBadRequest, model.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid user ID",
		})
	}

	userIdInt := int(userIdFloat)

	if userIdInt == 0 {
		return c.JSON(403, model.Response{
			Status:  403,
			Message: "Forbidden",
		})
	}

	transaction := new(model.Transaction)
	if err := c.Bind(transaction); err != nil {
		return c.JSON(400, model.Response{
			Status:  400,
			Message: "Bad Request, Invalid request body" + err.Error(),
		})
	}

	// Validate transaction data
	if transaction.Amount <= 0 {
		return c.JSON(400, model.Response{
			Status:  400,
			Message: "Invalid transaction amount",
		})
	}

	if transaction.DonationID <= 0 {
		return c.JSON(400, model.Response{
			Status:  400,
			Message: "Invalid donation ID",
		})
	}

	transaction, err := h.transactionRepo.CreateTransaction(userIdInt, transaction)
	if err != nil {
		return c.JSON(500, model.Response{
			Status:  500,
			Message: "Internal Server Error, " + err.Error(),
		})
	}
	//

	// Create Invoice
	// var invoice external.InvoiceResponse
	sendTransId := "transaction-" + strconv.Itoa(transaction.ID)

	// Create invoice using external API
	invoice, err := external.CreateInvoice(sendTransId, int(transaction.Amount), transaction.Donation.User.Email, transaction.InvoiceDescription)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to create invoice, " + err.Error(),
		})
	}

	transaction.InvoiceID = invoice.ID
	transaction.InvoiceURL = invoice.InvoiceURL
	transaction.Status = invoice.Status

	// Save updated transaction with invoice id, url, and status
	transaction, err = h.transactionRepo.UpdateTransaction(userIdInt, transaction)
	if err != nil {
		return c.JSON(500, model.Response{
			Status:  500,
			Message: "Internal Server Error, " + err.Error(),
		})
	}

	// return response
	return c.JSON(201, model.Response{
		Status:  201,
		Message: "Success",
		Data:    transaction,
	})
}

// UpdateTransaction godoc
// @Summary Update a transaction based on the invoice status
// @Description Get details of a specific transaction by its ID
// @Tags transactions
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <access_token>"
// @Param id path int true "Transaction ID"
// @Param model.Transaction body model.Transaction true "Transaction object"
// @Success 200 {object} model.Response
// @Router /transactions/{id} [put] // Updated the router path to use PUT method
func (h *transactionHandler) UpdateTransaction(c echo.Context) error {
	//get user id from context
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, model.Response{
			Status:  http.StatusUnauthorized,
			Message: "User not authenticated",
		})
	}

	userIdFloat, ok := userID.(float64)
	if !ok {
		return c.JSON(http.StatusBadRequest, model.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid user ID",
		})
	}
	userIdInt := int(userIdFloat)

	//get transaction id from param
	transactionID := c.Param("id")
	transactionIdInt, err := strconv.Atoi(transactionID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid transaction ID",
		})
	}

	// get stored transaction
	transaction := new(model.Transaction)
	if err := c.Bind(transaction); err != nil {
		return c.JSON(400, model.Response{
			Status:  400,
			Message: "Bad Request, Invalid request body" + err.Error(),
		})
	}
	transaction.ID = transactionIdInt

	transaction, err = h.transactionRepo.UpdateTransaction(userIdInt, transaction)
	if err != nil {
		return c.JSON(http.StatusNotFound, model.Response{
			Status:  http.StatusNotFound,
			Message: "Transaction not found",
		})
	}

	// return updated transaction
	return c.JSON(200, model.Response{
		Status:  200,
		Message: "Success",
		Data:    transaction,
	})
}

// CheckUpdateTransaction godoc
// @Summary Check and update a transaction based on the invoice status
// @Description Get details of a specific transaction by its ID and check its status
// @Tags transactions
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <access_token>"
// @Param id path int true "Transaction ID"
// @Param model.Transaction body model.Transaction true "Transaction object"
// @Success 200 {object} model.Response
// @Router /transactions/check-update-transaction/{id} [put] // Updated the router path to use PUT method
func (h *transactionHandler) CheckUpdateTransaction(c echo.Context) error {
	//get user id from context
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, model.Response{
			Status:  http.StatusUnauthorized,
			Message: "User not authenticated",
		})
	}

	userIdFloat, ok := userID.(float64)
	if !ok {
		return c.JSON(http.StatusBadRequest, model.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid user ID",
		})
	}
	userIdInt := int(userIdFloat)

	//get transaction id from param
	transactionID := c.Param("id")
	transactionIdInt, err := strconv.Atoi(transactionID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid transaction ID",
		})
	}

	// get stored transaction
	transaction := new(model.Transaction)
	transaction, err = h.transactionRepo.GetTransactionByID(userIdInt, transactionIdInt)
	if err != nil {
		return c.JSON(http.StatusNotFound, model.Response{
			Status:  http.StatusNotFound,
			Message: "Transaction not found",
		})
	}

	// check if stored transaction is paid
	if transaction.Status == "PENDING" {
		invoice, err := external.GetInvoice(transaction.InvoiceID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, model.Response{
				Status:  http.StatusInternalServerError,
				Message: "Error retrieving invoice status",
			})
		}
		// not paid
		if invoice.Status == "PENDING" {
			return c.JSON(200, model.Response{
				Status:  200,
				Message: "Transaction Status is Still PENDING",
				Data:    transaction,
			})
		}
		// paid
		if invoice.Status == "PAID" || invoice.Status == "SETTLED" {
			transaction.PaymentMethod = invoice.PaymentMethod
			transaction.Status = "PAID"
		}
		// expired
		if invoice.Status == "EXPIRED" {
			transaction.Status = "CANCELED"
		}

		// update after check
		if transaction.Status == "PAID" || transaction.Status == "CANCELED" {
			transaction, err = h.transactionRepo.CheckUpdateTransaction(userIdInt, transaction)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, model.Response{
					Status:  http.StatusInternalServerError,
					Message: "Error updating transaction status",
				})
			}
		}
		// return updated transaction
		return c.JSON(200, model.Response{
			Status:  200,
			Message: "Success Check and Update Transaction Status",
			Data:    transaction,
		})
	} else {
		// Already paid
		return c.JSON(200, model.Response{
			Status:  200,
			Message: "Transaction is Already Paid",
			Data:    transaction,
		})
	}

}

// GetTransactionByID godoc
// @Summary Get Transaction details by Transaction ID
// @Description Get details of a specific transaction by its ID
// @Tags transactions
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <access_token>"
// @Param id path int true "Transaction ID"
// @Success 200 {object} model.Response
// @Router /transactions/{id} [get] // Updated the router path to include transaction ID
func (h *transactionHandler) GetTransactionByID(c echo.Context) error {
	//get user id from context
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, model.Response{
			Status:  http.StatusUnauthorized,
			Message: "User not authenticated",
		})
	}

	userIdFloat, ok := userID.(float64)
	if !ok {
		return c.JSON(http.StatusBadRequest, model.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid user ID",
		})
	}
	userIdInt := int(userIdFloat)

	// get transaction id from url param
	transactionID := c.Param("id")
	// log.Println("transactionID", transactionID)
	if transactionID == "" {
		return c.JSON(http.StatusBadRequest, model.Response{
			Status:  400,
			Message: "Transaction ID is required",
		})
	}
	transactionIDInt, err := strconv.Atoi(transactionID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			Status:  400,
			Message: "Invalid transaction ID",
		})
	}

	// get stored transaction data
	transaction := new(model.Transaction)
	result, err := h.transactionRepo.GetTransactionByID(userIdInt, transactionIDInt)
	if err != nil {
		return c.JSON(500, model.Response{
			Status:  500,
			Message: "Internal Server Error, Error when getting transaction",
		})
	}

	transaction = result

	return c.JSON(200, model.Response{
		Status:  200,
		Message: "Success",
		Data:    transaction,
	})
}
