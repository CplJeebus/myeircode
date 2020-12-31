package utils

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"cloud.google.com/go/storage"
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

type Config struct {
	Bucket string `yaml:"bucket"`
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
