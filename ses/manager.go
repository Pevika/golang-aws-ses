//
// @Author: Geoffrey Bauduin <bauduin.geo@gmail.com>
//

package ses

import (
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"errors"
)

type profile struct {
	from			*string
	replyTo			[]*string
	returnPath		*string
	returnPathArn	*string
	sourceArn		*string
}

type Email struct {
	ses			*ses.SES
	profiles	map[string]*profile
}

// Creates a new manager
func NewEmail (awsAccessKey string, awsSecretKey string, region string) *Email {
	entity := new(Email)
	cred := credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, "")
	config := aws.NewConfig().WithRegion(region).WithCredentials(cred)
	sess := session.New(config)
	entity.ses = ses.New(sess)
	entity.profiles = map[string]*profile{}
	return entity
}

// Setup a profile to use with Send
func (this *Email) SetupProfile (name string, from string, replyTo []string, returnPath string, returnPathArn string, sourceArn string) bool {
	this.profiles[name] = &profile{
		from: aws.String(from),
		replyTo: []*string{},
		returnPath: aws.String(returnPath),
		returnPathArn: aws.String(returnPathArn),
		sourceArn: aws.String(sourceArn),
	}
	for _, d := range replyTo {
		this.profiles[name].replyTo = append(this.profiles[name].replyTo, aws.String(d))
	}
	return true
}

// Sends an email to the specified destination
func (this *Email) Send (p string, to []string, cc []string, bcc []string, subject string, htmlContent string, rawContent string, charset string) error {
	pr := this.profiles[p]
	if pr == nil {
		return errors.New("Cannot find profile " + p)
	}
	params := &ses.SendEmailInput{
		Destination: &ses.Destination{
			BccAddresses: []*string{},
			CcAddresses: []*string{},
			ToAddresses: []*string{},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Data: aws.String(htmlContent),
					Charset: aws.String(charset),
				},
				Text: &ses.Content{
					Data: aws.String(rawContent),
					Charset: aws.String(charset),
				},
			},
			Subject: &ses.Content{
				Data: aws.String(subject),
				Charset: aws.String(charset),
			},
		},
		Source: pr.from,
		ReplyToAddresses: pr.replyTo,
		ReturnPath: pr.returnPath,
		ReturnPathArn: pr.returnPathArn,
		SourceArn: pr.sourceArn,
	}
	for _, d := range to {
		params.Destination.ToAddresses = append(params.Destination.ToAddresses, aws.String(d))
	}
	for _, d := range cc {
		params.Destination.CcAddresses = append(params.Destination.CcAddresses, aws.String(d))
	}
	for _, d := range bcc {
		params.Destination.BccAddresses = append(params.Destination.BccAddresses, aws.String(d))
	}
	_, err :=  this.ses.SendEmail(params)
	return err
}
