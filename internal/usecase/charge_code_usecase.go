// internal/usecase/charge_code_usecase
package usecase

type ChargeCode struct {
	ChargeCodeID int     `json:"charge_code_id"`
	Code         string  `json:"code" binding:"required"`
	MaxUses      int     `json:"max_uses" binding:"required"`
	CurrentUses  int     `json:"current_uses" binding:"required"`
	Amount       float64 `json:"amount" binding:"required"`
}

type ChargeCodeRepository interface {
	CreateChargeCode(chargeCode *ChargeCode) (*ChargeCode, error)
	GetChargeCodes() ([]*ChargeCode, error)
	GetChargeCodeByCode(code string) (*ChargeCode, error)
	GetChargeCodeByID(id int) (*ChargeCode, error)
	DeleteChargeCode(id int) error
	UpdateChargeCode(chargeCode *ChargeCode) (*ChargeCode, error)
	GetUserChargeCodes(userId int) ([]*ChargeCode, error)
}

type ChargeCodeUseCase struct {
	ChargeCodeRepository ChargeCodeRepository
}

func NewChargeCodeUseCase(chargeCodeRepo ChargeCodeRepository) *ChargeCodeUseCase {
	return &ChargeCodeUseCase{ChargeCodeRepository: chargeCodeRepo}
}

func (cu *ChargeCodeUseCase) GetChargeCodes() ([]*ChargeCode, error) {
	return cu.ChargeCodeRepository.GetChargeCodes()
}

func (cu *ChargeCodeUseCase) GetChargeCodeByCode(code string) (*ChargeCode, error) {
	return cu.ChargeCodeRepository.GetChargeCodeByCode(code)

}

func (cu *ChargeCodeUseCase) GetChargeCodeByID(id int) (*ChargeCode, error) {
	return cu.ChargeCodeRepository.GetChargeCodeByID(id)
}

func (uc *ChargeCodeUseCase) CreateChargeCode(chargeCode *ChargeCode) (*ChargeCode, error) {
	return uc.ChargeCodeRepository.CreateChargeCode(chargeCode)
}

func (cu *ChargeCodeUseCase) DeleteChargeCode(id int) error {
	return cu.ChargeCodeRepository.DeleteChargeCode(id)
}

func (cu *ChargeCodeUseCase) UpdateChargeCode(chargeCode *ChargeCode) (*ChargeCode, error) {
	return cu.ChargeCodeRepository.UpdateChargeCode(chargeCode)
}

func (cu *ChargeCodeUseCase) GetUserChargeCodes(userId int) ([]*ChargeCode, error) {
	return cu.ChargeCodeRepository.GetUserChargeCodes(userId)
}
