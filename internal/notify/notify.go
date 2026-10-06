package notify

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"
	"video_processing_pipeline/internal/model"
)
type ConfigSMTP struct{
  Host string
  Port int
  Username string
  Password string
  From string
  To []string
}
type NotifyManager struct {
	cfg ConfigSMTP
}

func NewNotifyManager( cnfg ConfigSMTP) *NotifyManager {
	return &NotifyManager{ cfg: cnfg}
}

var  Notifier = (*NotifyManager)(nil)
//Go doc 
//func SendMail(addr string, a sasl.Client, from string, to []string, r io.Reader) error
func(N *NotifyManager)Notify(ctx context.Context, job *model.Jobs_Database) error {
  fmt.Println("Enter into notification service to notify users")
  authclient:=smtp.PlainAuth("",N.cfg.Username,N.cfg.Password,N.cfg.Host)
  msg:= strings.Join([]string{
    "From: your-email@example.com",
    "To: recipient@example.com",
    "Subject: Hello there",
    "",
    "This is the email body.",
}, "\r\n")
  //func smtp.SendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error
  err:=smtp.SendMail(N.cfg.Host,authclient,N.cfg.From,N.cfg.To,[]byte(msg))
  if(err!=nil){
	return err
  }
  return nil
}