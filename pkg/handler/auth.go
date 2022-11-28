package handler

import (
	userModel "main-server/pkg/model/user"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// @Summary SignUp
// @Tags auth
// @Description Регистрация пользователя
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body userModel.UserRegisterModel true "account info"
// @Success 200 {object} userModel.TokenAccessModel "data"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var input userModel.UserRegisterModel

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	data, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Добавление токена обновления в http only cookie
	c.SetCookie(viper.GetString("environment.refresh_token_key"), data.RefreshToken,
		30*24*60*60*1000, "/", viper.GetString("environment.domain"), false, true)

	c.JSON(http.StatusOK, userModel.TokenAccessModel{
		AccessToken: data.AccessToken,
	})
}

// @Summary SignIn
// @Tags auth
// @Description Авторизация пользователя
// @ID login
// @Accept  json
// @Produce  json
// @Param input body userModel.UserLoginModel true "credentials"
// @Success 200 {object} userModel.TokenAccessModel "data"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var input userModel.UserLoginModel

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	data, err := h.services.Authorization.LoginUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Добавление токена обновления в http only cookie
	c.SetCookie(viper.GetString("environment.refresh_token_key"), data.RefreshToken,
		30*24*60*60*1000, "/", viper.GetString("environment.domain"), false, true)

	c.JSON(http.StatusOK, userModel.TokenAccessModel{
		AccessToken: data.AccessToken,
	})
}

// @Summary SignInVK
// @Tags auth
// @Description Авторизация пользователя через VK
// @ID login_vk
// @Accept  json
// @Produce  json
// @Param input body userModel.UserLoginModel true "credentials"
// @Success 200 {object} userModel.TokenAccessModel "data"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-in/vk [post]
func (h *Handler) signInVK(c *gin.Context) {
	var input userModel.UserLoginModel

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	data, err := h.services.Authorization.LoginUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Добавление токена обновления в http only cookie
	c.SetCookie(viper.GetString("environment.refresh_token_key"), data.RefreshToken,
		30*24*60*60*1000, "/", viper.GetString("environment.domain"), false, true)

	c.JSON(http.StatusOK, userModel.TokenAccessModel{
		AccessToken: data.AccessToken,
	})
}

// @Summary SignInOAuth2
// @Tags auth
// @Description Авторизация пользователя через Google OAuth2
// @ID login_oauth2
// @Accept  json
// @Produce  json
// @Param input body userModel.GoogleOAuth2Code true "credentials"
// @Success 200 {object} userModel.TokenAccessModel "data"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-in/oauth2 [post]
func (h *Handler) signInOAuth2(c *gin.Context) {
	var input userModel.GoogleOAuth2Code

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	// For fast tests
	/*token, _ := configs.AppOAuth2Config.GoogleLogin.Exchange(c, input.Code)
	_, _ = google_oauth2.RevokeToken(token.AccessToken)
	return*/

	data, err := h.services.Authorization.LoginUserOAuth2(input.Code)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Добавление токена обновления в http only cookie
	c.SetCookie(viper.GetString("environment.refresh_token_key"), data.RefreshToken,
		30*24*60*60*1000, "/", viper.GetString("environment.domain"), false, true)

	c.JSON(http.StatusOK, userModel.TokenAccessModel{
		AccessToken: data.AccessToken,
	})
}

// @Summary Refresh
// @Tags auth
// @Description Обновление токена доступа и токена обновления
// @ID refresh
// @Accept  json
// @Produce  json
// @Param input body userModel.TokenRefreshModel true "credentials"
// @Success 200 {object} userModel.TokenAccessModel "data"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/refresh [post]
func (h *Handler) refresh(c *gin.Context) {
	refreshToken, err := c.Cookie(viper.GetString("environment.refresh_token_key"))

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	accessToken, _ := c.Get("access_token")
	authTypeValue, _ := c.Get("auth_type_value")
	tokenApi, _ := c.Get("token_api")

	data, err := h.services.Authorization.Refresh(userModel.TokenLogoutDataModel{
		AccessToken:   accessToken.(string),
		RefreshToken:  refreshToken,
		AuthTypeValue: authTypeValue.(string),
		TokenApi:      tokenApi.(*string),
	}, refreshToken)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.SetCookie(viper.GetString("environment.refresh_token_key"), data.RefreshToken,
		30*24*60*60*1000, "/", viper.GetString("environment.domain"), false, true)

	c.JSON(http.StatusOK, userModel.TokenAccessModel{
		AccessToken: data.AccessToken,
	})
}

type LogoutOutputModel struct {
	IsLogout bool `json:"is_logout"`
}

// @Summary Logout
// @Tags auth
// @Description Выход из аккаунта
// @ID logout
// @Accept  json
// @Produce  json
// @Success 200 {object} LogoutOutputModel "data"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/logout [post]
func (h *Handler) logout(c *gin.Context) {
	refreshToken, err := c.Cookie(viper.GetString("environment.refresh_token_key"))

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	accessToken, _ := c.Get("access_token")
	authTypeValue, _ := c.Get("auth_type_value")
	tokenApi, _ := c.Get("token_api")

	data, err := h.services.Authorization.Logout(userModel.TokenLogoutDataModel{
		AccessToken:   accessToken.(string),
		RefreshToken:  refreshToken,
		AuthTypeValue: authTypeValue.(string),
		TokenApi:      tokenApi.(*string),
	})

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if data {
		c.SetCookie(viper.GetString("environment.refresh_token_key"), "",
			30*24*60*60*1000, "/", viper.GetString("environment.domain"), false, true)
	}

	c.JSON(http.StatusOK, LogoutOutputModel{
		IsLogout: data,
	})
}

// @Summary Activate
// @Tags auth
// @Description Активация аккаунта по почте
// @ID activate
// @Accept  json
// @Produce  json
// @Success 200 {object} LogoutOutputModel "data"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/activate [get]
func (h *Handler) activate(c *gin.Context) {
	_, err := h.services.Activate(c.Params.ByName("link"))

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.HTML(http.StatusOK, "account_activate.html", gin.H{
		"title": "Подтверждение аккаунта",
	})
}
