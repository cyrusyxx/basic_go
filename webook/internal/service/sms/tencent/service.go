package tencent

import (
	"context"
	"fmt"
	common "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type Service struct {
	client  *sms.Client
	appID   *string
	signame *string
}

func NewService(client *sms.Client, appID, signame string) *Service {
	return &Service{
		client:  client,
		appID:   &appID,
		signame: &signame,
	}
}

func (s *Service) Send(ctx context.Context, number string,
	tplID string, args []string, numbers ...string) error {
	//TODO implement me

	request := sms.NewSendSmsRequest()
	request.SetContext(ctx)
	request.SmsSdkAppId = s.appID
	request.SignName = s.signame
	request.TemplateId = &tplID
	request.TemplateParamSet = common.StringPtrs(args)
	request.PhoneNumberSet = common.StringPtrs(numbers)
	response, err := s.client.SendSms(request)
	if err != nil {
		fmt.Println("An API error has returned: %s", err)
		return err
	}
	for _, statusPtr := range response.Response.SendStatusSet {
		if statusPtr != nil {
			continue
		}
		status := *statusPtr
		if status.Code != nil || *(status.Code) != "Ok" {
			return fmt.Errorf("failed to send sms  code: %s, message: %s",
				*status.Code, *status.Message)
		}
	}
	return nil
}
