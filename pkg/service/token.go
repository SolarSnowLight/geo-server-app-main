package service

import (
	"errors"
	userModel "main-server/pkg/model/user"
	"main-server/pkg/repository"

	"github.com/dgrijalva/jwt-go"
)

// Структура репозитория
type TokenService struct {
	role     repository.Role
	user     repository.User
	authType repository.AuthType
}

// Функция создания нового репозитория
func NewTokenService(role repository.Role,
	user repository.User,
	authType repository.AuthType,
) *TokenService {
	return &TokenService{
		role:     role,
		user:     user,
		authType: authType,
	}
}

/* Структура тела токена */
type tokenClaims struct {
	jwt.StandardClaims
	UsersId     string  `json:"users_id"`      // ID пользователя
	RolesId     string  `json:"roles_id"`      // Роль пользователя
	AuthTypesId string  `json:"auth_types_id"` // Тип аутентификации пользователя
	TokenApi    *string `json:"token_api"`     // Внешний токен доступа
}

func (s *TokenService) ParseToken(pToken, signingKey string) (userModel.TokenOutputParse, error) {
	token, err := jwt.ParseWithClaims(pToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})

	if !token.Valid {
		return userModel.TokenOutputParse{}, errors.New("token is not valid")
	}

	if err != nil {
		return userModel.TokenOutputParse{}, err
	}

	// Получение данных из токена (с преобразованием к указателю на tokenClaims)
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return userModel.TokenOutputParse{}, errors.New("token claims are not of type")
	}

	user, err := s.user.GetUser("uuid", claims.UsersId)

	if err != nil {
		return userModel.TokenOutputParse{}, err
	}

	role, err := s.role.GetRole("uuid", claims.RolesId)

	if err != nil {
		return userModel.TokenOutputParse{}, err
	}

	authType, err := s.authType.GetAuthType("uuid", claims.AuthTypesId)

	if err != nil {
		return userModel.TokenOutputParse{}, err
	}

	return userModel.TokenOutputParse{
		UsersId:  user.Id,
		RolesId:  role.Id,
		AuthType: authType,
		TokenApi: claims.TokenApi,
	}, nil
}

func (s *TokenService) ParseTokenWithoutValid(pToken, signingKey string) (userModel.TokenOutputParse, error) {
	token, err := jwt.ParseWithClaims(pToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})

	// Получение данных из токена (с преобразованием к указателю на tokenClaims)
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return userModel.TokenOutputParse{}, errors.New("token claims are not of type")
	}

	user, err := s.user.GetUser("uuid", claims.UsersId)

	if err != nil {
		return userModel.TokenOutputParse{}, err
	}

	role, err := s.role.GetRole("uuid", claims.RolesId)

	if err != nil {
		return userModel.TokenOutputParse{}, err
	}

	authType, err := s.authType.GetAuthType("uuid", claims.AuthTypesId)

	if err != nil {
		return userModel.TokenOutputParse{}, err
	}

	return userModel.TokenOutputParse{
		UsersId:  user.Id,
		RolesId:  role.Id,
		AuthType: authType,
		TokenApi: claims.TokenApi,
	}, nil
}
