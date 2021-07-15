package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	rootPath              string
	kafkaBrokers          string
	amazonAccessKeyID     string
	amazonSecretAccessKey string
	amazonSesRegion       string
}

var config *Config

func LoadEnvVars() {
	serviceFile, err := filepath.Abs(os.Args[0])
	if err != nil {
		log.Fatalf("Não foi possivel resgatar o caminho do serviço %v", err)
	}
	rootPath := extractPath(serviceFile)
	if err != nil {
		log.Fatalf("Não foi possivel definir o caminho do executável: %v", err)
	}
	log.Print("RootPath: ", rootPath)

	if err = createTemp(rootPath); err != nil {
		log.Fatalf("Error when trying to create the temp folder [%v]", err)
	}

	config = &Config{
		rootPath:              rootPath,
		kafkaBrokers:          os.Getenv("KAFKA_BROKERS"),
		amazonAccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		amazonSecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		amazonSesRegion:       os.Getenv("AMAZON_SES_REGION"),
	}
}

func (c *Config) RootPath() string {
	return c.rootPath
}

func (c *Config) KafkaBrokers() string {
	return c.kafkaBrokers
}

func (c *Config) AmazonAccessKeyID() string {
	return c.amazonAccessKeyID
}

func (c *Config) AmazonSecretAccessKey() string {
	return c.amazonSecretAccessKey
}

func (c *Config) AmazonSESRegion() string {
	return c.amazonSesRegion
}

func GetConfig() *Config {
	return config
}

func extractPath(servicePath string) string {
	path := strings.Split(servicePath, `\`)
	return strings.Join(path[0:len(path)-1], `\`)
}

func createTemp(path string) error {
	_, err := os.Stat(fmt.Sprintf("%s/temp", path))
	if err != nil {
		err := os.Mkdir(fmt.Sprintf("%s/temp", path), 0666)
		if err != nil {
			return fmt.Errorf("cannot create temp folder: %s", err)
		}
	}

	return nil
}
