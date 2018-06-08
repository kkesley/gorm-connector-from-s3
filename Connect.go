package gconn

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // mysql
)

//MysqlConnect to mysql db based on DB_BUCKET
func MysqlConnect() (*gorm.DB, error) {
	config := aws.Config{
		Region: aws.String("ap-southeast-2"),
	}
	sess := session.Must(session.NewSession(&config))

	svc := s3.New(sess)
	fmt.Println("accessing bucket: " + os.Getenv("DB_BUCKET") + "/" + os.Getenv("DB_BUCKET_KEY_FINAL"))
	s3Output, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("DB_BUCKET")),
		Key:    aws.String(os.Getenv("DB_BUCKET_KEY_FINAL")),
	})

	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(s3Output.Body)
	conBytes := buf.Bytes()

	var connection DB
	json.Unmarshal(conBytes, &connection)

	fmt.Println("accessing database")

	db, err := gorm.Open("mysql", connection.User+":"+connection.Password+"@("+connection.Host+":"+connection.Port+")"+"/"+connection.Db+"?charset=utf8&parseTime=True")
	if err != nil {
		return nil, err
	}
	return db, nil
}
