package routes

import (
	"context"
	"github.com/Nerzal/gocloak/v13"
	"github.com/gofiber/fiber/v2"
	"log/slog"
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

func SetupRoutes(app *fiber.App, handler handlers.Handler, ctx context.Context, logger *slog.Logger) {
	keycloakClient := gocloak.NewClient(keycloakURL)

	authMiddleware := authMiddleware(keycloakClient, logger)

	receiverController := inject(handler, logger)
	loginController := injectLogin(keycloakClient, ctx, logger)
	app.Get("/api/v1/convert", authMiddleware, receiverController.ReceiverMp4Controller)
	app.Post("/api/v1/login", loginController.Login)
}

func authMiddleware(keycloakClient *gocloak.GoCloak, logger *slog.Logger) func(c *fiber.Ctx) error {
	// Middleware de autenticação
	authMiddleware := func(c *fiber.Ctx) error {
		// Verifica o token de acesso no cabeçalho da requisição
		token := c.Get("Authorization")
		if token == "" {
			logger.Error("Token de acesso não fornecido")
			return c.Status(http.StatusUnauthorized).SendString("Token de acesso não fornecido")
		}

		// Valida o token com o Keycloak
		isValid, err := validateToken(token, keycloakClient, logger)
		if err != nil {
			logger.Error("Erro ao validar token: " + err.Error())
			return c.Status(http.StatusInternalServerError).SendString("Erro ao validar token")
		}
		if !isValid {
			logger.Error("Token de acesso inválido")
			return c.Status(http.StatusUnauthorized).SendString("Token de acesso inválido")
		}

		logger.Info("Token de acesso válido")
		// Se o token for válido, chama o próximo manipulador
		return c.Next()
	}
	return authMiddleware
}

// Função para validar o token de acesso com o Keycloak
func validateToken(token string, keycloakClient *gocloak.GoCloak, logger *slog.Logger) (bool, error) {
	ctx := context.Background()
	logger.Info("Validando token de acesso")
	_, err := keycloakClient.RetrospectToken(ctx, token, keycloakClientID, keycloakClientSec, keycloakRealm)
	if err != nil {
		logger.Error("Erro ao validar token: " + err.Error())
		return false, err
	}
	return true, nil
}

func inject(handler handlers.Handler, logger *slog.Logger) controller.ReceiverController {
	logger.Info("Injecting receiver controller")
	receiverController := controller.New(services.New(repositories.New(handler)), services.AwsService{
		Session: services.NewAwsService().Session,
	}, logger)
	logger.Info("Receiver controller injected")
	return receiverController
}

func injectLogin(keycloakClient *gocloak.GoCloak, ctx context.Context, logger *slog.Logger) controller.LoginController {
	logger.Info("Injecting login controller")
	loginController := controller.NewLoginController(services.NewLoginService(keycloakClient, ctx), ctx)
	logger.Info("Login controller injected")
	return loginController
}
