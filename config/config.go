package config

import "os"

var Default = New()

func New() *Config {
	acceessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_ACCESS_KEY_SECRET")
	region := os.Getenv("AWS_DEFAULT_REGION")
	return &Config{AccessKeyID: acceessKeyID, SecretAccessKey: secretAccessKey, Region: region}
}

type Config struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
}

func (c *Config) AWSAccessKeyID() string {
	return c.AccessKeyID
}

func (c *Config) AWSSecretAccessKey() string {
	return c.SecretAccessKey
}

func (c *Config) AWSRegion() string {
	return c.Region
}
