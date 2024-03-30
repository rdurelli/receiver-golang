package routes

import (
	"context"
	"github.com/Nerzal/gocloak/v13"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"net/http"
	"receiver/controller"
	"receiver/handlers"
	"receiver/repositories"
	"receiver/services"
)

var (
	keycloakURL = "http://auth:8080"
)

var (
	keycloakRealm     = "example-realm"
	keycloakClientID  = "test-client"
	keycloakClientSec = "AtU3FL4NDhfqiptgCZbwl9TkbK2s955S"
)

func SetupRoutes(app *fiber.App, handler handlers.Handler, ctx context.Context) {
	keycloakClient := gocloak.NewClient(keycloakURL)

	authMiddleware := authMiddleware(keycloakClient)

	receiverController := inject(handler)
	loginController := injectLogin(keycloakClient, ctx)
	app.Get("/api/v1/convert", authMiddleware, receiverController.ReceiverMp4Controller)
	app.Post("/api/v1/login", loginController.Login)
}

func authMiddleware(keycloakClient *gocloak.GoCloak) func(c *fiber.Ctx) error {
	// Middleware de autenticação
	authMiddleware := func(c *fiber.Ctx) error {
		// Verifica o token de acesso no cabeçalho da requisição
		token := c.Get("Authorization")
		if token == "" {
			log.Error("Token de acesso não fornecido")
			return c.Status(http.StatusUnauthorized).SendString("Token de acesso não fornecido")
		}

		// Valida o token com o Keycloak
		isValid, err := validateToken(token, keycloakClient)
		if err != nil {
			log.Error("Erro ao validar token: " + err.Error())
			return c.Status(http.StatusInternalServerError).SendString("Erro ao validar token")
		}
		if !isValid {
			log.Error("Token de acesso inválido")
			return c.Status(http.StatusUnauthorized).SendString("Token de acesso inválido")
		}

		log.Info("Token de acesso válido")
		// Se o token for válido, chama o próximo manipulador
		return c.Next()
	}
	return authMiddleware
}

// Função para validar o token de acesso com o Keycloak
func validateToken(token string, keycloakClient *gocloak.GoCloak) (bool, error) {
	ctx := context.Background()
	_, err := keycloakClient.RetrospectToken(ctx, token, keycloakClientID, keycloakClientSec, keycloakRealm)
	if err != nil {
		return false, err
	}
	return true, nil
}

func inject(handler handlers.Handler) controller.ReceiverController {
	receiverController := controller.New(services.New(repositories.New(handler)), services.AwsService{
		Session: services.NewAwsService().Session,
	})
	return receiverController
}

func injectLogin(keycloakClient *gocloak.GoCloak, ctx context.Context) controller.LoginController {
	loginController := controller.NewLoginController(services.NewLoginService(keycloakClient, ctx), ctx)
	return loginController
}
