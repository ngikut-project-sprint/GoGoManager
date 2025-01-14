package utils

import (
	"github.com/ngikut-project-sprint/GoGoManager/internal/validators"
)

type ManagerResponse struct {
	Email           string `json:"email"`
	Name            string `json:"name"`
	UserImageUri    string `json:"userImageUri"`
	CompanyName     string `json:"companyName"`
	CompanyImageUri string `json:"companyImageUri"`
}
type ManagerRequest struct {
	ID              int
	Email           *string `json:"email"`
	Password        *string `json:"password"`
	Name            *string `json:"name"`
	UserImageUri    *string `json:"userImageUri"`
	CompanyName     *string `json:"companyName"`
	CompanyImageUri *string `json:"companyImageUri"`
}

func (m ManagerRequest) ValidEmail() bool {
	if m.Email == nil {
		return false
	}
	err := validators.ValidateEmail(*m.Email)
	return err == nil
}

func (m ManagerRequest) ValidPassword() bool {
	if m.Password == nil {
		return false
	}
	lenght := len(*m.Password)
	return lenght >= 8 && lenght <= 32
}

func (m ManagerRequest) ValidName() bool {
	if m.Name == nil {
		return false
	}

	lenght := len(*m.Name)
	return m.Name != nil && lenght >= 4 && lenght <= 52
}

func (m ManagerRequest) ValidImageURI() bool {
	if m.UserImageUri == nil {
		return false
	}

	err := validators.ValidateURI(*m.UserImageUri)
	return m.UserImageUri != nil && err == nil
}

func (m ManagerRequest) ValidCompanyName() bool {
	if m.CompanyName == nil {
		return false
	}

	lenght := len(*m.CompanyName)
	return m.CompanyName != nil && lenght >= 4 && lenght <= 52
}

func (m ManagerRequest) ValidCompanyImageURI() bool {
	if m.CompanyImageUri == nil {
		return false
	}

	err := validators.ValidateURI(*m.CompanyImageUri)
	return m.CompanyImageUri != nil && err == nil
}
