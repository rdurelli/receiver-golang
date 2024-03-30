package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"mime/multipart"
	"os"
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
	AwsService      services.AwsService
}

func New(receiverService services.ReceiverService, awsService services.AwsService) ReceiverController {
	return ReceiverController{ReceiverService: receiverService, AwsService: awsService}
}

func (receiver ReceiverController) ReceiverMp4Controller(c *fiber.Ctx) error {
	bucketName := os.Getenv("AWS_BUCKET")
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}
	zapNumber := c.FormValue("whatZapNumber")
	discordWebHook := c.FormValue("discordWebHook")

	err2, done := services.IsMp4(c, file)
	if done {
		return err2
	}

	key, err3, done2 := receiver.sendToBucket(c, err, file, bucketName)
	if done2 {
		return err3
	}

	err = receiver.sendToQueue(key+".mp4", zapNumber, discordWebHook, bucketName)
	if err != nil {
		return c.SendString("Error while sending to queue")
	}

	fileToConvert, id, error := receiver.saveToDataBase(key, zapNumber, discordWebHook)
	if error == true {
		return c.SendString("Error while inserting file")
	}

	fileToConvert.Id = id
	return c.JSON(fileToConvert)
}

func (receiver ReceiverController) sendToQueue(key string, zapNumber string, discordWebHook string, bucketName string) error {
	payload := models.ReceiverRequestDTO{
		FileName:       key,
		WhatZapNumber:  zapNumber,
		DiscordWebHook: discordWebHook,
		BucketName:     bucketName,
	}

	// Convert struct to JSON
	jsonData, err := json.Marshal(payload)
	err = receiver.AwsService.SendToQueue(string(jsonData))
	if err != nil {
		fmt.Println("Error while sending to queue")
		return err
	}
	return nil
}

func (receiver ReceiverController) sendToBucket(c *fiber.Ctx, err error, file *multipart.FileHeader, bucketName string) (string, error, bool) {
	// Get Buffer from file
	buffer, err := file.Open()

	if err != nil {
		return "", c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		}), true
	}
	defer buffer.Close()

	// Create a byte slice with the file size
	data := make([]byte, file.Size)

	// Read the file contents into the byte slice
	_, err = buffer.Read(data)
	if err != nil {
		fmt.Println("Error reading file:", err)
	}

	key, err := receiver.AwsService.UploadToS3(bucketName, data)
	if err != nil {
		return "", c.SendString("Error while uploading file to S3"), true
	}
	return key, nil, false
}

func (receiver ReceiverController) saveToDataBase(key string, zapNumber string, discordWebHook string) (models.FileToConvert, int, bool) {
	fileToConvert := models.FileToConvert{
		FileNameAsMp4:  key + ".mp4",
		Status:         SENT_TO_QUEUE,
		WhatZapNumber:  zapNumber,
		DiscordWebHook: discordWebHook,
		InsertedAt:     time.Now(),
		UpdatedAt:      time.Now(),
	}
	id, error := receiver.ReceiverService.Save(fileToConvert)
	return fileToConvert, id, error
}
