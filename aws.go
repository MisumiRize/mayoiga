package main

import (
	"bytes"
	"encoding/base64"

	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/s3"
)

var sess *session.Session

type s3Service interface {
	GetObject(*s3.GetObjectInput) (*s3.GetObjectOutput, error)
	PutObject(*s3.PutObjectInput) (*s3.PutObjectOutput, error)
}

var s3Svc s3Service

func getS3Svc() (s3Service, error) {
	if s3Svc == nil {
		config, err := readConfig()
		if err != nil {
			return nil, err
		}

		s3Svc = s3.New(sess, &aws.Config{Region: config.Region})
	}
	return s3Svc, nil
}

func s3GetObject(key *string) (*bytes.Buffer, error) {
	config, err := readConfig()
	if err != nil {
		return nil, err
	}

	svc, err := getS3Svc()
	if err != nil {
		return nil, err
	}

	params := &s3.GetObjectInput{
		Bucket: config.S3Bucket,
		Key:    key,
	}
	res, err := svc.GetObject(params)
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

func s3PutObject(key *string, body io.ReadSeeker) (err error) {
	config, err := readConfig()
	if err != nil {
		return
	}

	svc, err := getS3Svc()
	if err != nil {
		return
	}

	putParams := &s3.PutObjectInput{
		Bucket: config.S3Bucket,
		Key:    key,
		Body:   body,
	}
	_, err = svc.PutObject(putParams)
	return
}

type kmsService interface {
	Encrypt(*kms.EncryptInput) (*kms.EncryptOutput, error)
}

var kmsSvc kmsService

func getKmsSvc() (kmsService, error) {
	if kmsSvc == nil {
		config, err := readConfig()
		if err != nil {
			return nil, err
		}

		kmsSvc = kms.New(sess, &aws.Config{Region: config.Region})
	}
	return kmsSvc, nil
}

func kmsEncrypt(value string) (*string, error) {
	config, err := readConfig()
	if err != nil {
		return nil, err
	}

	svc, err := getKmsSvc()
	if err != nil {
		return nil, err
	}

	encParams := &kms.EncryptInput{
		KeyId:     config.KMSKeyID,
		Plaintext: []byte(value),
	}
	res, err := svc.Encrypt(encParams)
	if err != nil {
		return nil, err
	}

	v := base64.StdEncoding.EncodeToString(res.CiphertextBlob)
	return &v, nil
}
