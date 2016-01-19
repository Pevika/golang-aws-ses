//
// @Author: Geoffrey Bauduin <bauduin.geo@gmail.com>
//

package main

import (
    "github.com/pevika/golang-aws-ses/ses"
    "log"
)

func main () {
    sender := ses.NewEmail("access key", "secret key", "region")
    sender.SetupProfile("default", "Geoffrey Bauduin <bauduin.geo@gmail.com>", []string{"bauduin.geo@gmail.com"}, "return path", "return path arn", "source arn")
    err := sender.Send("default", []string{"recipient@gmail.com"}, []string{"cc@gmail.com"}, []string{"bcc@gmail.com"}, "Subject", "Html content", "Raw content", "encoding")
    if err != nil {
        log.Println(err)
    }
}