package normal

import (
	"github.com/apache/rocketmq-clients/golang/v5"
	"github.com/apache/rocketmq-clients/golang/v5/credentials"
)

type ProducerConfig struct {
	Endpoint string
	//GroupName string
	AccessKey   string
	SecretKey   string
	Region      string
	Topic       string
	MaxAttempts int32
}

func NewNormalProducer(c *ProducerConfig) (golang.Producer, error) {
	if c.MaxAttempts == 0 {
		c.MaxAttempts = 3
	}
	//os.Setenv("mq.consoleAppender.enabled", "true")
	golang.ResetLogger()
	producer, err := golang.NewProducer(&golang.Config{
		Endpoint: c.Endpoint,
		//Region:   Region,
		Credentials: &credentials.SessionCredentials{
			AccessKey:    c.AccessKey,
			AccessSecret: c.SecretKey,
		},
	},
		golang.WithTopics(c.Topic),
		golang.WithMaxAttempts(c.MaxAttempts),
	)
	if err != nil {
		return nil, err
	}
	return producer, nil
}
