package notify

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"
	mediastore "video_processing_pipeline/internal/mediaStore"
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
  store mediastore.MediaStoreManager
}

func NewNotifyManager( cnfg ConfigSMTP, stre mediastore.MediaStore) *NotifyManager {
	return &NotifyManager{ cfg: cnfg,
   store: &stre,
  }
}

var  Notifier = (*NotifyManager)(nil)
//Go doc 
//func SendMail(addr string, a sasl.Client, from string, to []string, r io.Reader) error
func(N *NotifyManager)Notify(ctx context.Context, job *model.Jobs_Database) error {
  fmt.Println("Enter into notification service to notify users")
  authclient:=smtp.PlainAuth("",N.cfg.Username,N.cfg.Password,N.cfg.Host)
  fmt.Println("authn string is", authclient)
  urlstring,err:=N.store.PresignedGet(ctx,job.RawPath)
  if err!=nil{
    return err
  }
  msg := strings.Join([]string{
	"From: " + N.cfg.From,
	"To: " + strings.Join(N.cfg.To, ", "),
	"Subject: Video " + job.Status + ": " + job.FileName,
	"",
	"Job " + job.VideoId + " status=" + job.Status,
  "You can watch the video it will expire in 24 hour",
		urlstring,
}, "\r\n")
fmt.Println("msg is", msg)
  //func smtp.SendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error
  port:=N.cfg.Port
  if port==0{
    port=587
  }
  addr := fmt.Sprintf("%s:%d", N.cfg.Host, port)
  fmt.Println("addres along with post is",addr)
  err=smtp.SendMail(addr,authclient,N.cfg.From,N.cfg.To,[]byte(msg))
  if err!=nil{
    fmt.Printf("SMTP ERROR OCCURRED: %v\n", err)
    return err
  }
  fmt.Println("email sent successfully")
  return nil
}