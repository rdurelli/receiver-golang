package models

import "github.com/gofiber/fiber/v2"

type ReceiverRequestDTO struct {
	File           fiber.FormFile `json:"file"`
	WhatZapNumber  string         `json:"whatZapNumber"`
	DiscordWebHook string         `json:"discordWebHook"`
}
