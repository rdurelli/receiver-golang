package services

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"os"
)

type AwsService struct {
	Session *session.Session
}

func NewAwsService() AwsService {
	s, _ := retrieveSession()
	return AwsService{Session: s}
}

func (awsService AwsService) UploadToS3(bucketName string, data []byte) (string, error) {
	log.Info("Uploading file to S3")
	log.Info("Bucket: " + bucketName)
	key := uuid.NewString()
	svc := s3.New(awsService.Session)
	putObject, err := svc.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(key + ".mp4"),
		Body:          bytes.NewReader(data),
		ContentLength: aws.Int64(int64(len(data))),
	})
	if err != nil {
		log.Error("Error uploading file to S3: " + err.Error())
		return "", err
	}
	log.Info("File uploaded to S3 " + putObject.String())
	return key, err
}

func (awsService AwsService) SendToQueue(payload string) error {
	sqsClient := sqs.New(awsService.Session)
	queueURL := os.Getenv("AWS_SQS_URL")
	_, err := sqsClient.SendMessage(&sqs.SendMessageInput{
		MessageBody:  aws.String(payload),
		QueueUrl:     &queueURL,
		DelaySeconds: aws.Int64(0),
	})
	return err
}

func retrieveSession() (*session.Session, error) {
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")
	endpoint := os.Getenv("AWS_ENDPOINT_URL")
	s3ForcePathStyle := os.Getenv("S3_FORCE_PATH_STYLE") == "true" // Expecting "true" or "false"

	// Initialize AWS session with credentials and region
	return session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
		Region:           aws.String(region),
		Endpoint:         aws.String(endpoint),
		S3ForcePathStyle: aws.Bool(s3ForcePathStyle),
	})
}
