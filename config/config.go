package config

import "os"

// Default is package variable
var Default = New()

// New function create Config struct
func New() *Config {
	acceessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_ACCESS_KEY_SECRET")
	region := os.Getenv("AWS_DEFAULT_REGION")
	return &Config{AccessKeyID: acceessKeyID, SecretAccessKey: secretAccessKey, Region: region}
}

// Config represents AWS settings
type Config struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
}

// AWSAccessKeyID returns aws access key
func (c *Config) AWSAccessKeyID() string {
	return c.AccessKeyID
}

// AWSSecretAccessKey returns aws access secret
func (c *Config) AWSSecretAccessKey() string {
	return c.SecretAccessKey
}

// AWSRegion returns default region
func (c *Config) AWSRegion() string {
	return c.Region
}
