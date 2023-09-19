// internal/usecase/user_usecase.go
package usecase

type User struct {
	ID          int     `json:"id" binding:"required"`
	PhoneNumber string  `json:"PhoneNumber" binding:"required"`
	Balance     float64 `json:"Balance" binding:"required"`
}

type UserRepository interface {
	GetUserByPhoneNumber(phoneNumber string) (*User, error)
	UpdateUser(user *User) (*User, error)
	ListOfUsersUseChargeCode(chargeCodeId int, page int, pageSize int) ([]*User, error)
	GetUserBalance(userId int) (float64, error)
}

type UserUseCase struct {
	UserRepository UserRepository
}

func NewUserUseCase(userRepo UserRepository) *UserUseCase {
	return &UserUseCase{UserRepository: userRepo}
}

func (uc *UserUseCase) GetUserByPhoneNumber(phoneNumber string) (*User, error) {
	return uc.UserRepository.GetUserByPhoneNumber(phoneNumber)
}

func (uc *UserUseCase) UpdateUser(user *User) (*User, error) {
	return uc.UserRepository.UpdateUser(user)
}

func (uc *UserUseCase) ListOfUsersUseChargeCode(chargeCodeId int, page int, pageSize int) ([]*User, error) {
	return uc.UserRepository.ListOfUsersUseChargeCode(chargeCodeId, page, pageSize)
}

func (uc *UserUseCase) GetUserBalance(userId int) (float64, error) {
	return uc.UserRepository.GetUserBalance(userId)
}
