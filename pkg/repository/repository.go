package repository

import (
	rbacModel "main-server/pkg/model/rbac"
	userModel "main-server/pkg/model/user"

	"github.com/jmoiron/sqlx"
	"golang.org/x/oauth2"
)

type Authorization interface {
	CreateUser(user userModel.UserRegisterModel) (userModel.UserAuthDataModel, error)
	LoginUser(user userModel.UserLoginModel) (userModel.UserAuthDataModel, error)
	LoginUserOAuth2(code string) (userModel.UserAuthDataModel, error)
	CreateUserOAuth2(user userModel.UserRegisterOAuth2Model, token *oauth2.Token) (userModel.UserAuthDataModel, error)
	Refresh(data userModel.TokenLogoutDataModel, refreshToken string, token userModel.TokenOutputParse) (userModel.UserAuthDataModel, error)
	Logout(tokens userModel.TokenLogoutDataModel) (bool, error)
	Activate(link string) (bool, error)

	GetUser(column, value string) (userModel.UserModel, error)
	GetRole(column, value string) (rbacModel.RoleModel, error)
}

type Role interface {
	GetRole(column, value interface{}) (rbacModel.RoleModel, error)
}

type User interface {
	GetUser(column, value interface{}) (userModel.UserModel, error)
}

type AuthType interface {
	GetAuthType(column, value interface{}) (userModel.AuthTypeModel, error)
}

type Repository struct {
	Authorization
	Role
	User
	AuthType
}

func NewRepository(db *sqlx.DB) *Repository {
	user := NewUserPostgres(db)

	return &Repository{
		Authorization: NewAuthPostgres(db, *user),
		Role:          NewRolePostgres(db),
		User:          user,
		AuthType:      NewAuthTypePostgres(db),
	}
}
