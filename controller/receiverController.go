package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"receiver/models"
	"receiver/services"
	"time"
)

const (
	SAVED_TO_DATABASE             = iota
	SENT_TO_QUEUE                 = iota
	CONVERTED_TO_MP3_SUCCESSFULLY = iota
	CONVERTED_TO_MP3_WITH_FAILED  = iota
)

type ReceiverController struct {
	ReceiverService services.ReceiverService
}

func New(receiverService services.ReceiverService) ReceiverController {
	return ReceiverController{ReceiverService: receiverService}
}

func (receiver ReceiverController) ReceiverMp4Controller(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}
	zapNumber := c.FormValue("whatZapNumber")
	discordWebHook := c.FormValue("discordWebHook")
	fmt.Println(zapNumber, discordWebHook)

	err2, done := services.IsMp4(c, file)
	if done {
		return err2
	}

	fileToConvert := models.FileToConvert{
		FileNameAsMp4: file.Filename,
		Status:        SAVED_TO_DATABASE,
		InsertedAt:    time.Now(),
		UpdatedAt:     time.Now(),
	}
	id, error := receiver.ReceiverService.Save(fileToConvert)
	if error == true {
		return c.SendString("Error while inserting file")
	}

	fileToConvert.Id = id
	return c.JSON(fileToConvert)
}
