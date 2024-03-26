package controller

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"receiver/services"
)

type LoginController struct {
	LoginService services.LoginService
	ctx          context.Context
}

func NewLoginController(loginService services.LoginService, ctx context.Context) LoginController {
	return LoginController{LoginService: loginService, ctx: ctx}
}

func (login LoginController) Login(c *fiber.Ctx) error {

	// Extrai as credenciais do corpo da requisição
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&credentials); err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	// Autentica as credenciais com o Keycloak
	token, err := login.LoginService.Login(credentials.Username, credentials.Password)
	if err != nil {
		return c.Status(http.StatusUnauthorized).SendString("Falha na autenticação")
	}

	return c.SendString("Login " + token.AccessToken)
}
