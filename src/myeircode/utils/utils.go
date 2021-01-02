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

const Pretty = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Eircodes</title>
</head>
<body>
    <div id="myData"></div>
    <script>
        fetch('/api')
            .then(function (response) {
                return response.json();
            })
            .then(function (data) {
                appendData(data);
            })
            .catch(function (err) {
                console.log('error: ' + err);
            });

        function appendData(data) {
            var mainContainer = document.getElementById("myData");
            for (var i = 0; i < data.length; i++) {
                var div = document.createElement("div");
								div.innerHTML = '<b>Name:</b> ' + data[i].name + '<br><b>Code:</b> ' + data[i].code + '<br>';
                mainContainer.appendChild(div);
            }
        }
    </script>
</body>
</html>`

const PrettyForm = `<!DOCTYPE html>
<html>
<body>

<h3>New Code</h3>

<form method="POST" action="/new">
  <label for="fname">Who is it:</label><br>
  <input type="text" id="name" name="name" value=""><br>
  <label for="lname">Code:</label><br>
  <input type="text" id="code" name="code" value=""><br><br>
  <input type="submit" value="Submit">
</form>

</body>
</html>`

type Config struct {
	Bucket  string `yaml:"bucket"`
	MailKey string `yaml:"mailApiKey"`
	Admin   string `yaml:"adminMail"`
	Host    string `yaml:"host"`
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
		fmt.Printf(fmt.Sprint(response.StatusCode))
		fmt.Printf(response.Body)
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
		fmt.Errorf("storage.NewClient: %w", err)
	}

	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	wc := client.Bucket(bucket).Object(object).NewWriter(ctx)
	if err != nil {
		fmt.Errorf("Object(%q).NewReader: %w", object, err)
	}
	defer wc.Close()
	wc.ContentType = "text/plain"
	_, e := wc.Write(codes)
	if e != nil {
		fmt.Println(e)
	}

}
