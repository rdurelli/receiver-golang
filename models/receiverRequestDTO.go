package models

type ReceiverRequestDTO struct {
	FileName       string `json:"fileName"`
	BucketName     string `json:"bucketName"`
	WhatZapNumber  string `json:"whatZapNumber"`
	DiscordWebHook string `json:"discordWebHook"`
}
