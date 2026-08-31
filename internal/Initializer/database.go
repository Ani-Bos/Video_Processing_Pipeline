package initializer

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DbInitializer struct {
	DB *gorm.DB
}

func ConnectDB(d *DbInitializer) {
  fmt.Println("Entering into connecting to db")
  db_host:=os.Getenv("DATABASE_URL")
  fmt.Println("dbstring is",db_host)
  if db_host==""{
   fmt.Println("Failed to get db host from the environment")
   return
  }
  var err error
  d.DB,err=gorm.Open(postgres.Open(db_host), &gorm.Config{})
  if err!=nil{
     fmt.Println("Got error while loadign gorm postgres",err)
	 return
  }
  fmt.Println("connected to postgres successfully")
}