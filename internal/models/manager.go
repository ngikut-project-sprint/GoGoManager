package models

import (
	"time"

	"github.com/ngikut-project-sprint/GoGoManager/internal/utils"
	"github.com/ngikut-project-sprint/GoGoManager/internal/validators"
)

type Manager struct {
	ID              int
	Email           *string
	Password        *string
	Name            *string
	UserImageUri    *string
	CompanyName     *string
	CompanyImageUri *string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}

func (m Manager) ValidEmail() bool {
	err := validators.ValidateEmail(*m.Email)
	return err == nil
}

func (m Manager) ValidPassword() bool {
	lenght := len(*m.Password)
	return lenght >= 8 && lenght <= 32
}

func (m Manager) ValidName() bool {
	lenght := len(*m.Name)
	return lenght >= 4 && lenght <= 52
}

func (m Manager) ValidImageURI() bool {
	err := validators.ValidateURI(*m.UserImageUri)
	return err == nil
}

func (m Manager) ValidCompanyName() bool {
	lenght := len(*m.CompanyName)
	return lenght >= 4 && lenght <= 52
}

func (m Manager) ValidCompanyImageURI() bool {
	err := validators.ValidateURI(*m.CompanyImageUri)
	return err == nil
}

func (m Manager) ToManagerResponse() utils.ManagerResponse {
	var (
		email           string
		name            string
		userImageUri    string
		companyName     string
		companyImageUri string
	)

	if m.Email != nil {
		email = *m.Email
	} else {
		email = ""
	}

	if m.Name != nil {
		name = *m.Name
	} else {
		name = ""
	}
	if m.UserImageUri != nil {
		userImageUri = *m.UserImageUri
	} else {
		userImageUri = ""
	}
	if m.CompanyName != nil {
		companyName = *m.CompanyName
	} else {
		companyName = ""
	}

	if m.CompanyImageUri != nil {
		companyImageUri = *m.CompanyImageUri
	} else {
		companyImageUri = ""
	}

	return utils.ManagerResponse{
		Email:           email,
		Name:            name,
		UserImageUri:    userImageUri,
		CompanyName:     companyName,
		CompanyImageUri: companyImageUri,
	}
}
