package service

import (
	userModel "main-server/pkg/model/user"
	"main-server/pkg/repository"

	"github.com/spf13/viper"
)

// Структура репозитория
type AuthService struct {
	repo         repository.Authorization
	tokenService TokenService
}

// Функция создания нового репозитория
func NewAuthService(repo repository.Authorization, tokenService TokenService) *AuthService {
	return &AuthService{
		repo:         repo,
		tokenService: tokenService,
	}
}

/*
*	Create user
 */
func (s *AuthService) CreateUser(user userModel.UserRegisterModel) (userModel.UserAuthDataModel, error) {
	return s.repo.CreateUser(user)
}

/*
*	Login user
 */
func (s *AuthService) LoginUser(user userModel.UserLoginModel) (userModel.UserAuthDataModel, error) {
	return s.repo.LoginUser(user)
}

/*
*	Login user with Google OAuth2
 */
func (s *AuthService) LoginUserOAuth2(code string) (userModel.UserAuthDataModel, error) {
	return s.repo.LoginUserOAuth2(code)
}

/*
*	Refresh user
 */
func (s *AuthService) Refresh(data userModel.TokenLogoutDataModel, refreshToken string) (userModel.UserAuthDataModel, error) {
	token, err := s.tokenService.ParseTokenWithoutValid(refreshToken, viper.GetString("token.signing_key_refresh"))

	if err != nil {
		return userModel.UserAuthDataModel{}, err
	}

	return s.repo.Refresh(data, refreshToken, token)
}

/*
*	Logout user
 */
func (s *AuthService) Logout(tokens userModel.TokenLogoutDataModel) (bool, error) {
	return s.repo.Logout(tokens)
}

/*
*	Активация аккаунта
 */
func (s *AuthService) Activate(link string) (bool, error) {
	return s.repo.Activate(link)
}

/*
*	Token parsing function
 */
func (s *AuthService) ParseToken(pToken, signingKey string) (userModel.TokenOutputParse, error) {
	return userModel.TokenOutputParse{}, nil
	/*token, err := jwt.ParseWithClaims(pToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
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

	user, err := s.repo.GetUser("uuid", claims.UsersId)

	if err != nil {
		return userModel.TokenOutputParse{}, err
	}

	role, err := s.repo.GetRole("uuid", claims.RolesId)

	if err != nil {
		return userModel.TokenOutputParse{}, err
	}

	return userModel.TokenOutputParse{
		UsersId: user.Id,
		RolesId: role.Id,
	}, nil*/
}
