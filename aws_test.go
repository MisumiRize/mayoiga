package main

import (
	"testing"

	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/golang/mock/gomock"
)

func TestS3GetObject(t *testing.T) {
	resetConfig(t)

	createJSON(t, map[string]string{
		"S3Bucket": "bucket",
		"S3Key":    "key",
	})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockS3 := NewMockS3API(ctrl)
	mockS3.EXPECT().GetObject(&s3.GetObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("key"),
	}).Return(&s3.GetObjectOutput{
		Body: ioutil.NopCloser(strings.NewReader("string")),
	}, nil)
	s3Svc = mockS3

	buf, err := s3GetObject()
	if err != nil {
		t.Fatalf(err.Error())
	}

	if buf.String() != "string" {
		t.Fatalf("assertion failed. actual: %s", buf.String())
	}
}
