package public

import (
	"github.com/donglei1234/platform/services/common/mq"
	"go.uber.org/zap"
	"time"
)

type Client struct {
	logger        *zap.Logger
	server        ChatServer
	profileID     string
	subscriptions map[string]mq.Subscription
	sendTime      map[string]time.Time
	chatInterval  time.Duration
}

func CreateClient(logger *zap.Logger, profileId string, server ChatServer, chatInterval time.Duration) *Client {
	return &Client{
		logger:        logger.With(zap.String("ProfileId", profileId)),
		profileID:     profileId,
		server:        server,
		subscriptions: make(map[string]mq.Subscription),
		sendTime:      make(map[string]time.Time),
		chatInterval:  chatInterval,
	}
}

func (c *Client) clear() {
	c.profileID = ""
}

func (c *Client) AddSubscription(topics string, sub mq.Subscription) {
	c.subscriptions[topics] = sub
}

func (c *Client) GetSubscription(topics string) (mq.Subscription, bool) {
	v, ok := c.subscriptions[topics]
	return v, ok
}

func (c *Client) Unsub(topics string) {
	if sub, ok := c.subscriptions[topics]; ok {
		err := sub.Unsubscribe()
		if err != nil {
			c.logger.Warn("Unsubscribe failed", zap.Error(err))
		}
		delete(c.subscriptions, topics)
	}
}

func (c *Client) UnsubAll() {
	for _, v := range c.subscriptions {
		err := v.Unsubscribe()
		if err != nil {
			c.logger.Warn("Unsubscribe failed", zap.Error(err))
		}
	}
	c.clear()
	c.logger.Info("unsub all")
}

func (c *Client) CheckInterval(key string) bool {
	if t, ok := c.sendTime[key]; ok {
		if time.Now().Sub(t) < c.chatInterval {
			return false
		}
	}
	c.sendTime[key] = time.Now()
	return true
}
