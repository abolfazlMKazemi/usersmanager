// internal/delivery/user_handler.go
package delivery

import (
	"chargeCode/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID          int     `json:"id"`
	PhoneNumber string  `json:"PhoneNumber" binding:"required"`
	Balance     float64 `json:"Balance"`
}

type UserHandler struct {
	UserUseCase *usecase.UserUseCase `json:"UserUseCase"`
}

func NewUserHandler(userUC *usecase.UserUseCase) *UserHandler {
	return &UserHandler{UserUseCase: userUC}

}

// GetUserByPhoneNumber godoc
// @Summary Get user by phoneNumber
// @Description Get a user by their unique phoneNumber.
// @Tags Users
// @ID get-user-by-phoneNumber
// @Produce json
// @Param phoneNumber path string true "User phoneNumber" Example: 09120000000
// @Success 200 {object} User
// @Router /api/v1/user/{phoneNumber} [get]
func (uh *UserHandler) GetUserByPhoneNumber(c *gin.Context) {
	userPhoneNumber := c.Param("phoneNumber")
	user, err := uh.UserUseCase.GetUserByPhoneNumber(userPhoneNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// Updateuser godoc
// @Summary Update a User
// @Description Update a User using the provided data.
// @Tags Users
// @Accept json
// @Produce json
// @Param User body User true "User object to update"
// @Success 200 {object} User
// @Router /api/v1/user [put]
func (uh *UserHandler) UpdateUser(c *gin.Context) {
	var user usecase.User

	// Parse the request body into a ChargeCode struct
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// At this point, user contains the data from the request body
	// You can use it as needed, such as passing it to your use case for update

	_, err := uh.UserUseCase.UpdateUser(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//c.JSON(http.StatusOK, updateduser)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// ListOfUserUsesChargeCode godoc
// @Summary Get List Of Users Use ChargeCode
// @Description Get a list of users who use a specific ChargeCode with pagination.
// @Tags Users
// @ID get-list-of-users-use-chargecode
// @Produce json
// @Param chargeCodeId path int true "ChargeCode ID" Example: 100
// @Param page query int false "Page number most start from 1"
// @Param pageSize query int false "Number of items per page"
// @Success 200 {object} ChargeCode
// @Router /api/v1/user/chargeCode/{chargeCodeId} [get]
func (uh *UserHandler) ListOfUsersUseChargeCode(c *gin.Context) {
	chargeCodeID, err := strconv.Atoi(c.Param("chargeCodeId"))
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "parsing error"})
		return
	}
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
	chargeCodes, err := uh.UserUseCase.ListOfUsersUseChargeCode(chargeCodeID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, chargeCodes)
}

// GetUserBalance godoc
// @Summary Get user balance by id
// @Description Get a user balance by their unique id.
// @Tags Users
// @ID get-user-balance-by-id
// @Produce json
// @Param userId path int true "User id" Example: 1
// @Success 200 {object} float64
// @Router /api/v1/user/balance/{userId} [get]
func (uh *UserHandler) GetUserBalance(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "parsing error"})
		return
	}
	user, err := uh.UserUseCase.GetUserBalance(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}
