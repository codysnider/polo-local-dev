package aws

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

const (
	ErrAuthExpiration = "aws authentication expired, please run aws-auth again"
)

func GetCallerIdentity() (string, error) {

	awsSession, awsSessionErr := session.NewSession(nil)
	if awsSessionErr != nil {
		return "", awsSessionErr
	}

	svc := sts.New(awsSession)
	input := &sts.GetCallerIdentityInput{}

	result, err := svc.GetCallerIdentity(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				return "", errors.New(ErrAuthExpiration)
			}
		} else {
			return "", err
		}
	}

	return fmt.Sprintf("account: %s", *result.Account), nil
}
