package utils

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"cloud.google.com/go/storage"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Bucket    string `yaml:"bucket"`
	MailKey   string `yaml:"mailApiKey"`
	Admin     string `yaml:"adminMail"`
	Host      string `yaml:"host"`
	CookieKey string `yaml:"cookieKey"`
}

type Code struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

func SendMail(c Config, uuid string) {
	from := mail.NewEmail("My Eircode", c.Admin)
	subject := "Change waiting for approval"
	to := mail.NewEmail("Admin", c.Admin)
	plainTextContent := "This is a test message"
	htmlContent := "Click code to authorize https://" + c.Host + "/auth?id=" + uuid + " </strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(c.MailKey)

	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Print(fmt.Sprint(response.StatusCode))
		fmt.Print(response.Body)
	}
}

func (c *Config) LoadConfig() *Config {
	configFile, err := ioutil.ReadFile("config.yaml")

	if err != nil {
		log.Printf("Unable to open config file %v", err)
	}

	err = yaml.Unmarshal(configFile, c)

	if err != nil {
		log.Fatalf("Can't parse file %v", err)
	}

	return c
}

func DownloadFile(bucket string, object string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %w", err)
	}

	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	rc, err := client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("Object(%q).NewReader: %w", object, err)
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll: %w", err)
	}

	return data, nil
}

func SaveCodes(bucket string, object string, codes []byte) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	client, err := storage.NewClient(ctx)
	if err != nil {
		fmt.Printf("storage.NewClient: %+v", err)
	}

	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	wc := client.Bucket(bucket).Object(object).NewWriter(ctx)

	defer wc.Close()
	wc.ContentType = "text/plain"

	_, e := wc.Write(codes)
	if e != nil {
		fmt.Printf("Can't write to bucket %+v", e)
	}
}
