package main

import (
	"bytes"

	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/s3"
)

type s3Service interface {
	GetObject(*s3.GetObjectInput) (*s3.GetObjectOutput, error)
	PutObject(*s3.PutObjectInput) (*s3.PutObjectOutput, error)
}

var s3Svc s3Service

func s3GetObject() (*bytes.Buffer, error) {
	s3Bucket, err := getStringConfig("S3Bucket")
	if err != nil {
		return nil, err
	}

	s3Key, err := getStringConfig("S3Key")
	if err != nil {
		return nil, err
	}

	params := &s3.GetObjectInput{
		Bucket: aws.String(*s3Bucket),
		Key:    aws.String(*s3Key),
	}
	res, err := s3Svc.GetObject(params)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	buf := new(bytes.Buffer)
	if _, err = buf.ReadFrom(res.Body); err != nil {
		return nil, err
	}

	return buf, nil
}

func s3PutObject(body io.ReadSeeker) (err error) {
	s3Bucket, err := getStringConfig("S3Bucket")
	if err != nil {
		return
	}

	s3Key, err := getStringConfig("S3Key")
	if err != nil {
		return
	}

	putParams := &s3.PutObjectInput{
		Bucket: aws.String(*s3Bucket),
		Key:    aws.String(*s3Key),
		Body:   body,
	}
	_, err = s3Svc.PutObject(putParams)
	return
}

type kmsService interface {
	Encrypt(*kms.EncryptInput) (*kms.EncryptOutput, error)
}

var kmsSvc kmsService

func kmsEncrypt(value string) (*string, error) {
	keyID, err := getStringConfig("KMSKeyId")
	if err != nil {
		return nil, err
	}

	encParams := &kms.EncryptInput{
		KeyId:     aws.String(*keyID),
		Plaintext: []byte(value),
	}
	res, err := kmsSvc.Encrypt(encParams)
	if err != nil {
		return nil, err
	}

	v := string(res.CiphertextBlob)
	return &v, nil
}
