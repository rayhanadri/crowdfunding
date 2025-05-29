package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rayhanadri/crowdfunding/donation-service/model"

	"github.com/rayhanadri/crowdfunding/api-gateway/entity"
	"github.com/rayhanadri/crowdfunding/api-gateway/repository"
)

type TransactionHandler interface {
	GetAllTransaction(c echo.Context) error
	CreateTransaction(c echo.Context) error
	GetTransactionByID(c echo.Context) error
	UpdateTransaction(c echo.Context) error
	SyncTransaction(c echo.Context) error
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
// @Success 200 {object} entity.Response
// @Router /transactions [get]
func (h *transactionHandler) GetAllTransaction(c echo.Context) error {
	//get user id from context
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, entity.Response{
			Status:  http.StatusUnauthorized,
			Message: "User not authenticated",
		})
	}

	userIdFloat, ok := userID.(float64)
	if !ok {
		return c.JSON(http.StatusBadRequest, entity.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid user ID",
		})
	}
	userIdInt := int(userIdFloat)
	log.Println("userIdInt", userIdInt)

	transactions, err := h.transactionRepo.GetAllTransaction()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, entity.Response{
			Status:  http.StatusInternalServerError,
			Message: "Internal Server Error",
		})
	}
	return c.JSON(200, entity.Response{
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
// @Param entity.TransactionRequest body entity.TransactionRequest true "Transaction object"
// @Success 201 {object} entity.Response
// @Router /transactions [post] // Updated the router path to use POST method
func (h *transactionHandler) CreateTransaction(c echo.Context) error {
	//
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, entity.Response{
			Status:  http.StatusUnauthorized,
			Message: "User not authenticated",
		})
	}

	userIdFloat, ok := userID.(float64)
	if !ok {
		return c.JSON(http.StatusBadRequest, entity.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid user ID",
		})
	}

	userIdInt := int(userIdFloat)

	if userIdInt == 0 {
		return c.JSON(403, entity.Response{
			Status:  403,
			Message: "Forbidden",
		})
	}

	transaction := new(model.Transaction)

	// return response
	return c.JSON(201, entity.Response{
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
// @Param entity.TransactionRequest body entity.TransactionRequest true "Transaction object"
// @Success 200 {object} entity.Response
// @Router /transactions/{id} [put] // Updated the router path to use PUT method
func (h *transactionHandler) UpdateTransaction(c echo.Context) error {
	//get user id from context
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, entity.Response{
			Status:  http.StatusUnauthorized,
			Message: "User not authenticated",
		})
	}

	userIdFloat, ok := userID.(float64)
	if !ok {
		return c.JSON(http.StatusBadRequest, entity.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid user ID",
		})
	}
	userIdInt := int(userIdFloat)
	log.Println("userIdInt", userIdInt)

	//get transaction id from param
	transactionID := c.Param("id")
	transactionIdInt, err := strconv.Atoi(transactionID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, entity.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid transaction ID",
		})
	}

	// get stored transaction
	transaction := new(entity.Transaction)
	if err := c.Bind(transaction); err != nil {
		return c.JSON(400, entity.Response{
			Status:  400,
			Message: "Bad Request, Invalid request body" + err.Error(),
		})
	}
	transaction.ID = transactionIdInt

	// return response
	return c.JSON(201, entity.Response{
		Status:  201,
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
// @Success 200 {object} entity.Response
// @Router /transactions/sync-transaction/{id} [put] // Updated the router path to use PUT method
func (h *transactionHandler) SyncTransaction(c echo.Context) error {
	//get user id from context
	transaction := new(entity.Transaction)
	// return updated transaction
	return c.JSON(200, entity.Response{
		Status:  200,
		Message: "Success Check and Update Transaction Status",
		Data:    transaction,
	})
}

// GetTransactionByID godoc
// @Summary Get Transaction details by Transaction ID
// @Description Get details of a specific transaction by its ID
// @Tags transactions
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <access_token>"
// @Param id path int true "Transaction ID"
// @Success 200 {object} entity.Response
// @Router /transactions/{id} [get] // Updated the router path to include transaction ID
func (h *transactionHandler) GetTransactionByID(c echo.Context) error {
	//get user id from context
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, entity.Response{
			Status:  http.StatusUnauthorized,
			Message: "User not authenticated",
		})
	}

	userIdFloat, ok := userID.(float64)
	if !ok {
		return c.JSON(http.StatusBadRequest, entity.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid user ID",
		})
	}
	userIdInt := int(userIdFloat)
	log.Println("userIdInt", userIdInt)

	// get transaction id from url param
	transactionID := c.Param("id")
	// log.Println("transactionID", transactionID)
	if transactionID == "" {
		return c.JSON(http.StatusBadRequest, entity.Response{
			Status:  400,
			Message: "Transaction ID is required",
		})
	}
	transactionIDInt, err := strconv.Atoi(transactionID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, entity.Response{
			Status:  400,
			Message: "Invalid transaction ID",
		})
	}

	// get stored transaction data
	transaction := new(entity.Transaction)
	transaction.ID = transactionIDInt

	return c.JSON(200, entity.Response{
		Status:  200,
		Message: "Success",
		Data:    transaction,
	})
}
