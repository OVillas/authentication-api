package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var (
	ErrHashPassword             = errors.New("error trying hashed password")
	ErrUserAlreadyRegistered    = errors.New("there is already a registered user with this email")
	ErrCreateUser               = errors.New("error to create user")
	ErrGetUser                  = errors.New("error to get user")
	ErrConvertUserPayLoadToUser = errors.New("error to create id from new user")
	ErrInvalidId                = errors.New("the id passed is invalid")
	ErrUserNotFound             = errors.New("user not found")
	ErrDeleteUser               = errors.New("error to delete user")
	ErrSameEmail                = errors.New("the email cannot be the same as the previous one")
	ErrUserNotAuthorized        = errors.New("user not authorized to action")
	ErrPasswordNotMatch         = errors.New("invalid password")
	ErrGenToken                 = errors.New("error to generate new token jwt")
	ErrUnexpectedSigningMethod  = errors.New("unexpected signature method")
	ErrInvalidToken             = errors.New("token invalid")
	ErrIdNotFoundInPermissions  = errors.New("error to get id in token")
	ErrIdIsNotAString           = errors.New("'id' field value is not a string")
	ErrUpdatePassword           = errors.New("error to update password")
	ErrToSendConfirmationCode   = errors.New("error to send confirmation code")
	ErrInvalidOTP               = errors.New("wrong or expired OTP")
	ErrOTPNotFound              = errors.New("not found OTP from email")
	ErrUserIDMismatch           = errors.New("user ID mismatch")
)

type User struct {
	ID                  string    `gorm:"column:Id;type:char(36);primary_key"`
	Name                string    `gorm:"column:Name;type:varchar(75)"`
	Username            string    `gorm:"column:Username;type:varchar(255);unique_index"`
	Email               string    `gorm:"column:Email;type:varchar(255);unique_index"`
	Password            string    `gorm:"column:Password;type:varchar(255)"`
	EmailConfirmed      bool      `gorm:"column:EmailConfirmed;type:boolean"`
	TwoFactorAuthActive bool      `gorm:"column:TwoFactorAuthActive;type:boolean"`
	CreatedAt           time.Time `gorm:"column:CreatedAt"`
	UpdateAt            time.Time `gorm:"column:UpdateAt"`
}

func (User) TableName() string {
	return "user"
}

type UserPayLoad struct {
	Name     string `json:"name,omitempty" validate:"required,min=1,max=75"`
	Nick     string `json:"username,omitempty" validate:"required,min=1,max=75"`
	Email    string `json:"email,omitempty" validate:"required,email"`
	Password string `json:"password,omitempty" validate:"required,min=6,containsany=!@#&?"`
}

type UserUpdatePayLoad struct {
	Name     string `json:"name,omitempty" validate:"min=1,max=75"`
	Email    string `json:"email,omitempty" validate:"required,email"`
	Username string `json:"username,omitempty" validate:"required,min=6,max=75"`
}

type UserResponse struct {
	Id               string
	Name             string
	Email            string
	Username         string
	IsEmailConfirmed bool
	CreatedAt        string
	LastModified     string
}

type UserInfosResponse struct {
	Name string
	Nick string
}

type Login struct {
	Username string `json:"username,omitempty" validate:"required,min=6"`
	Password string `json:"password,omitempty" validate:"required"`
}

type UpdatePassword struct {
	Current string `json:"current,omitempty" validate:"required,min=6,containsany=!@#&?"`
	New     string `json:"new,omitempty" validate:"required,min=6,containsany=!@#&?"`
}

type ResetPassword struct {
	New     string `json:"new,omitempty" validate:"required,min=6,containsany=!@#&?"`
	Confirm string `json:"confirm,omitempty" validate:"required,min=6,containsany=!@#&?"`
}

type ConfirmationCode struct {
	Code       string
	ExpiryTime time.Time
}

type ConfirmCode struct {
	Email string `json:"email,omitempty" validate:"required,email"`
	Code  string `json:"code,omitempty" validate:"required"`
}

type RequestResetPassword struct {
	Email string `json:"email,omitempty" validate:"required,email"`
}

type UserHandler interface {
	Create(ctx echo.Context) error
	GetById(ctx echo.Context) error
	GetByNameOrNick(ctx echo.Context) error
	GetByEmail(ctx echo.Context) error
	GetAll(ctx echo.Context) error
	Update(ctx echo.Context) error
	Delete(ctx echo.Context) error
	Login(ctx echo.Context) error
	UpdatePassword(ctx echo.Context) error
	ConfirmEmail(ctx echo.Context) error
	ForgotPassword(ctx echo.Context) error
	ConfirmResetPasswordCode(ctx echo.Context) error
	ResetPassword(ctx echo.Context) error
}

type UserService interface {
	Create(userPayLoad UserPayLoad) error
	GetById(id string) (*UserResponse, error)
	GetByNameOrNick(nameOrNick string) ([]UserResponse, error)
	GetByEmail(email string) (*UserResponse, error)
	GetByUsername(username string) (*UserResponse, error)
	GetAll() ([]UserResponse, error)
	Update(id string, userUpdate UserUpdatePayLoad) error
	Delete(id string) error
	Login(login Login) (string, error)
	UpdatePassword(id string, updatePassword UpdatePassword) error
	SendConfirmationCode(email string) error
	ConfirmEmail(confirmCode ConfirmCode) error
	ConfirmResetPasswordCode(confirmCode ConfirmCode) (string, error)
	ResetPassword(userId string, resetPassword ResetPassword) error
	CheckUserIDMatch(idFromToken string) error
}

type UserRepository interface {
	Create(user User) error
	GetById(id string) (*User, error)
	GetByNameOrNick(nameOrNick string) ([]User, error)
	GetByEmail(email string) (*User, error)
	GetByUsername(username string) (*User, error)
	GetAll() ([]User, error)
	Update(id string, user User) error
	Delete(id string) error
	UpdatePassword(id string, password string) error
	ConfirmedEmail(id string) error
}

func (upl *UserPayLoad) Validate() error {
	validate := validator.New()
	return validate.Struct(upl)
}

func (uu *UserUpdatePayLoad) Validate() error {
	validate := validator.New()
	return validate.Struct(uu)
}

func (upl *UserPayLoad) ToUser(hashedPassword string) (*User, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &User{
		ID:       id.String(),
		Name:     strings.TrimSpace(upl.Name),
		Email:    strings.TrimSpace(upl.Email),
		Username: strings.TrimSpace(upl.Nick),
		Password: strings.TrimSpace(hashedPassword),
	}, nil
}

func (uu *UserUpdatePayLoad) ToUser() *User {
	return &User{
		Name:     uu.Name,
		Email:    uu.Email,
		Username: uu.Username,
	}
}

func (u *User) ToUserResponse() *UserResponse {
	return &UserResponse{
		Id:               u.ID,
		Name:             u.Name,
		Email:            u.Email,
		Username:         u.Username,
		IsEmailConfirmed: u.EmailConfirmed,
		CreatedAt:        u.CreatedAt.Format("2006-01-02 15:04:05"),
		LastModified:     u.UpdateAt.Format("2006-01-02 15:04:05"),
	}
}

func (u *User) ToUserInfosResponse() *UserInfosResponse {
	return &UserInfosResponse{
		Name: u.Name,
		Nick: u.Username,
	}
}

func (l *Login) Validate() error {
	validate := validator.New()
	return validate.Struct(l)
}

func (rrp *RequestResetPassword) Validate() error {
	validate := validator.New()
	return validate.Struct(rrp)
}

func (up *UpdatePassword) Validate() error {
	validate := validator.New()
	return validate.Struct(up)
}

func (ce *ConfirmCode) Validate() error {
	validate := validator.New()
	return validate.Struct(ce)
}

func (rp *ResetPassword) Validate() error {
	validate := validator.New()
	return validate.Struct(rp)
}
