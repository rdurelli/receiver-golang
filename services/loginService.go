package services

import (
	"context"
	"github.com/Nerzal/gocloak/v13"
	"github.com/gofiber/fiber/v2/log"
	"os"
)

var (
	keycloakRealm     = os.Getenv("KEYCLOAK_REALM")
	keycloakClientID  = os.Getenv("KEYCLOAK_CLIENT_ID")
	keycloakClientSec = os.Getenv("KEYCLOAK_CLIENT_SECRET")
)

type LoginService struct {
	keycloakClient *gocloak.GoCloak
	ctx            context.Context
}

func NewLoginService(keycloakClient *gocloak.GoCloak, ctx context.Context) LoginService {
	return LoginService{keycloakClient: keycloakClient, ctx: ctx}
}

func (login LoginService) Login(username string, password string) (*gocloak.JWT, error) {
	log.Info("LoginService.Login " + username + password)
	token, err := login.keycloakClient.Login(login.ctx, keycloakClientID, keycloakClientSec, keycloakRealm, username, password)
	if err != nil {
		log.Error("LoginService.Login " + err.Error())
		return nil, err
	}
	return token, nil
}
