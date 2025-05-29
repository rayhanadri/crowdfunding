package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/rayhanadri/crowdfunding/api-gateway/model"
	"github.com/rayhanadri/crowdfunding/api-gateway/repository"
)

type DonationHandler interface {
	GetAllDonations(c echo.Context) error
	CreateDonation(c echo.Context) error
	GetDonationByID(c echo.Context) error
	UpdateDonation(c echo.Context) error
}

type donationHandler struct {
	donationRepo repository.DonationRepository
}

func NewDonationHandler(donationRepo repository.DonationRepository) DonationHandler {
	return &donationHandler{donationRepo: donationRepo}
}

// GetAllDonations godoc
// @Summary Get all donations for a user
// @Description Get all donations for an active user
// @Tags donations
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <access_token>"
// @Success 200 {object} model.Response
// @Router /donations/{id} [get]
func (h *donationHandler) GetAllDonations(c echo.Context) error {
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

	donations, err := h.donationRepo.GetAllDonations(userIdInt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{
			Status:  http.StatusInternalServerError,
			Message: "Internal Server Error",
		})
	}
	return c.JSON(200, model.Response{
		Status:  200,
		Message: "Success",
		Data:    donations,
	})
}

// CreateDonation godoc
// @Summary Create a new donation
// @Description Create a new donation for a user
// @Tags donations
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <access_token>"
// @Param model.Donation body model.Donation true "Donation object"
// @Success 201 {object} model.Response
// @Router /donations [post] // Updated the router path to use POST method
func (h *donationHandler) CreateDonation(c echo.Context) error {
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

	donation := new(model.Donation)
	if err := c.Bind(donation); err != nil {
		return c.JSON(400, model.Response{
			Status:  400,
			Message: "Bad Request, Invalid request body" + err.Error(),
		})
	}

	// Validate donation data
	if donation.Amount <= 0 {
		return c.JSON(400, model.Response{
			Status:  400,
			Message: "Invalid donation amount",
		})
	}

	// Validate campaign ID
	if donation.CampaignID <= 0 {
		return c.JSON(400, model.Response{
			Status:  400,
			Message: "Invalid campaign ID",
		})
	}

	donation, err := h.donationRepo.CreateDonation(userIdInt, donation)
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
		Data:    donation,
	})
}

// UpdateDonation godoc
// @Summary Update a donation based on the invoice status
// @Description Get details of a specific donation by its ID
// @Tags donations
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <access_token>"
// @Param id path int true "Donation ID"
// @Param model.Donation body model.Donation true "Donation object"
// @Success 200 {object} model.Response
// @Router /donations/{id} [put] // Updated the router path to use PUT method
func (h *donationHandler) UpdateDonation(c echo.Context) error {
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

	//get donation id from param
	donationID := c.Param("id")
	donationIdInt, err := strconv.Atoi(donationID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid donation ID",
		})
	}

	// get stored donation
	donation := new(model.Donation)
	if err := c.Bind(donation); err != nil {
		return c.JSON(400, model.Response{
			Status:  400,
			Message: "Bad Request, Invalid request body" + err.Error(),
		})
	}
	donation.ID = donationIdInt

	donation, err = h.donationRepo.UpdateDonation(userIdInt, donation)
	if err != nil {
		return c.JSON(http.StatusNotFound, model.Response{
			Status:  http.StatusNotFound,
			Message: "Donation not found",
		})
	}

	// return updated donation
	return c.JSON(200, model.Response{
		Status:  200,
		Message: "Success",
		Data:    donation,
	})
}

// GetDonationByID godoc
// @Summary Get Donation details by Donation ID
// @Description Get details of a specific donation by its ID
// @Tags donations
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <access_token>"
// @Param id path int true "Donation ID"
// @Success 200 {object} model.Response
// @Router /donations/{id} [get] // Updated the router path to include donation ID
func (h *donationHandler) GetDonationByID(c echo.Context) error {
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

	// get donation id from url param
	donationID := c.Param("id")
	// log.Println("donationID", donationID)
	if donationID == "" {
		return c.JSON(http.StatusBadRequest, model.Response{
			Status:  400,
			Message: "Donation ID is required",
		})
	}
	donationIDInt, err := strconv.Atoi(donationID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			Status:  400,
			Message: "Invalid donation ID",
		})
	}

	// get stored donation data
	donation := new(model.Donation)
	result, err := h.donationRepo.GetDonationByID(userIdInt, donationIDInt)
	if err != nil {
		return c.JSON(500, model.Response{
			Status:  500,
			Message: "Internal Server Error, Error when getting donation",
		})
	}

	donation = result

	return c.JSON(200, model.Response{
		Status:  200,
		Message: "Success",
		Data:    donation,
	})
}
