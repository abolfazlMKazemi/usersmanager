package delivery

import (
	"chargeCode/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ChargeCode struct {
	ChargeCodeID int     `json:"charge_code_id"`
	Code         string  `json:"code" binding:"required"`
	MaxUses      int     `json:"max_uses" binding:"required"`
	CurrentUses  int     `json:"current_uses" binding:"required"`
	Amount       float64 `json:"amount" binding:"required"`
}

type CreateChargeCodeMode struct {
	Code        string  `json:"code" binding:"required"`
	MaxUses     int     `json:"max_uses" binding:"required"`
	CurrentUses int     `json:"current_uses" binding:"required"`
	Amount      float64 `json:"amount" binding:"required"`
}

type ChargeCodeHandler struct {
	ChargeCodeUseCase *usecase.ChargeCodeUseCase `json:"ChargeCodeUseCase"`
}

func NewChargeCodeHandler(chargeCodeUC *usecase.ChargeCodeUseCase) *ChargeCodeHandler {
	return &ChargeCodeHandler{ChargeCodeUseCase: chargeCodeUC}

}

// CreateChargeCode godoc
// @Summary Create a new chargeCode
// @Description Create a new chargeCode using the provided data.
// @Tags ChargeCode
// @Accept json
// @Produce json
// @Param CreateChargeCodeMode body CreateChargeCodeMode true "ChargeCode object to create"
// @Success 200 {object} ChargeCode
// @Router /api/v1/chargeCode [post]
func (cH *ChargeCodeHandler) CreateChargeCode(c *gin.Context) {
	var chargeCode usecase.ChargeCode

	// Parse the request body into a ChargeCode struct
	if err := c.ShouldBindJSON(&chargeCode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// At this point, chargeCode contains the data from the request body
	// You can use it as needed, such as passing it to your use case for creation

	_, err := cH.ChargeCodeUseCase.CreateChargeCode(&chargeCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//c.JSON(http.StatusOK, createdChargeCode)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// GetChargeCodes godoc
// @Summary Get chargeCodes
// @Description Get charge codes with pagination.
// @Tags ChargeCode
// @ID get-paginated-chargeCodes
// @Produce json
// @Param page query integer false "Page number most start from 1"
// @Param pageSize query integer false "Number of items per page"
// @Success 200 {object} ChargeCode
// @Router /api/v1/chargeCode [get]
func (cH *ChargeCodeHandler) GetChargeCodes(c *gin.Context) {
	// Parse the page and pageSize query parameters with default values
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "parsing error"})
		return
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "parsing error"})
		return
	}
	chargeCodes, err := cH.ChargeCodeUseCase.GetChargeCodes(page, pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, chargeCodes)
}

// GetChargeCodeByID godoc
// @Summary Get chargeCode by ID
// @Description Get a chargeCode by their unique ID.
// @Tags ChargeCode
// @ID get-chargeCode-by-id
// @Produce json
// @Param id path int true "chargeCode ID" Example: 123
// @Success 200 {object} ChargeCode
// @Router /api/v1/chargeCode/{id} [get]
func (cH *ChargeCodeHandler) GetChargeCodeByID(c *gin.Context) {
	chargeCodeID, _ := strconv.Atoi(c.Param("id"))
	chargeCode, err := cH.ChargeCodeUseCase.GetChargeCodeByID(chargeCodeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, chargeCode)
}

// GetChargeCodeByCode godoc
// @Summary Get chargeCode by Code
// @Description Get a chargeCode by their unique Code.
// @Tags ChargeCode
// @ID get-chargeCode-by-code
// @Produce json
// @Param code path string true "chargeCode Code" Example: c216
// @Success 200 {object} ChargeCode
// @Router /api/v1/chargeCode/code/{code} [get]
func (cH *ChargeCodeHandler) GetChargeCodeByCode(c *gin.Context) {
	chargeCodeParams := c.Param("code")
	chargeCode, err := cH.ChargeCodeUseCase.GetChargeCodeByCode(chargeCodeParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, chargeCode)
}

// DeleteChargeCode godoc
// @Summary Delete chargeCode by ID
// @Description Delete a chargeCode by their unique ID.
// @Tags ChargeCode
// @ID delete-chargeCode-by-id
// @Produce json
// @Param id path int true "chargeCode ID" Example: 123
// @Success 200 {object} string "OK"
// @Router /api/v1/chargeCode/{id} [delete]
func (cH *ChargeCodeHandler) DeleteChargeCodeByID(c *gin.Context) {
	chargeCodeID, _ := strconv.Atoi(c.Param("id"))
	err := cH.ChargeCodeUseCase.DeleteChargeCode(chargeCodeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "OK")
}

// UpdateChargeCode godoc
// @Summary Update a chargeCode
// @Description Update a chargeCode using the provided data.
// @Tags ChargeCode
// @Accept json
// @Produce json
// @Param chargeCode body ChargeCode true "ChargeCode object to update"
// @Success 200 {object} ChargeCode
// @Router /api/v1/chargeCode [put]
func (cH *ChargeCodeHandler) UpdateChargeCode(c *gin.Context) {
	var chargeCode usecase.ChargeCode

	// Parse the request body into a ChargeCode struct
	if err := c.ShouldBindJSON(&chargeCode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// At this point, chargeCode contains the data from the request body
	// You can use it as needed, such as passing it to your use case for update

	_, err := cH.ChargeCodeUseCase.UpdateChargeCode(&chargeCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//c.JSON(http.StatusOK, updatedChargeCode)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// GetUserChargeCodes godoc
// @Summary Get user chargeCodes with pagination
// @Description Get a user chargeCode by their unique userId with pagination support.
// @Tags ChargeCode
// @ID get-chargeCode-by-userId
// @Produce json
// @Param userId path int true "user id" Example: 123
// @Param page query int false "Page number most start from 1" Default: 1 Example: 1
// @Param pageSize query int false "page size" Default: 10 Example: 20
// @Success 200 {object} ChargeCode
// @Router /api/v1/chargeCode/user/{userId} [get]
func (cH *ChargeCodeHandler) GetUserChargeCodes(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}

	// Parse the page and pageSize query parameters with default values
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "parsing error"})
		return
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "parsing error"})
		return
	}

	chargeCodes, err := cH.ChargeCodeUseCase.GetUserChargeCodes(userId, page, pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, chargeCodes)
}
