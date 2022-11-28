package handler

import (
	"errors"
	"main-server/pkg/service/google_oauth2"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

const (
	authorizationHeader = "Authorization"
	usersCtx            = "users_id"
	rolesCtx            = "roles_id"
	authTypeValueCtx    = "auth_type_value"
	accessTokenCtx      = "access_token"
	tokenApiCtx         = "token_api"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)

	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "Пустой заголовок авторизации!")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "Не корректный авторизационный заголовок!")
		return
	}

	data, err := h.services.Token.ParseToken(headerParts[1], viper.GetString("token.signing_key_access"))

	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	switch data.AuthType.Value {
	case "GOOGLE":
		if result, err := google_oauth2.VerifyAccessToken(*data.TokenApi); err != nil || result != true {
			newErrorResponse(c, http.StatusUnauthorized, "Не действительный токен доступа")
			return
		}
		break

	case "LOCAL":
		break
	}

	// Добавление к контексту дополнительных данных о пользователе
	c.Set(usersCtx, data.UsersId)
	c.Set(rolesCtx, data.RolesId)
	c.Set(authTypeValueCtx, data.AuthType.Value)
	c.Set(tokenApiCtx, data.TokenApi)
	c.Set(accessTokenCtx, headerParts[1])
}

func (h *Handler) userIdentityLogout(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)

	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "Пустой заголовок авторизации!")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "Не корректный авторизационный заголовок!")
		return
	}

	data, err := h.services.Token.ParseTokenWithoutValid(headerParts[1], viper.GetString("token.signing_key_access"))

	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	// Добавление к контексту дополнительных данных о пользователе
	c.Set(usersCtx, data.UsersId)
	c.Set(rolesCtx, data.RolesId)
	c.Set(authTypeValueCtx, data.AuthType.Value)
	c.Set(tokenApiCtx, data.TokenApi)
	c.Set(accessTokenCtx, headerParts[1])
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(usersCtx)
	if !ok {
		return 0, errors.New("user id not found")
	}

	idInt, ok := id.(int)
	if !ok {
		return 0, errors.New("user id is of invalid type")
	}

	return idInt, nil
}
