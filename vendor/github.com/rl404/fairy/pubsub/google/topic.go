package google

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/iterator"
)

func (c *Client) getTopic(topic string) (*pubsub.Topic, error) {
	if c.topicExist[topic] {
		return c.client.Topic(topic), nil
	}

	// Check for first time get topic.
	t := c.client.Topic(topic)

	isExist, err := t.Exists(context.Background())
	if err != nil {
		return nil, err
	}

	if !isExist {
		// Create new topic.
		t, err = c.client.CreateTopic(context.Background(), topic)
		if err != nil {
			return nil, err
		}

		// Also create new subscription.
		subID := fmt.Sprintf("%s-subscription", topic)
		if _, err = c.client.CreateSubscription(context.Background(), subID, pubsub.SubscriptionConfig{
			Topic: t,
		}); err != nil {
			return nil, err
		}
	}

	c.Lock()
	c.topicExist[topic] = true
	c.Unlock()

	return t, nil
}

func (c *Client) getSubscription(topic string) (*pubsub.Subscription, error) {
	if c.subscriptionExist[topic] != "" {
		return c.client.Subscription(c.subscriptionExist[topic]), nil
	}

	// Check for first time.
	subID := fmt.Sprintf("%s-subscription", topic)
	s := c.client.Subscription(subID)

	isExist, err := s.Exists(context.Background())
	if err != nil {
		return nil, err
	}

	if isExist {
		c.Lock()
		c.subscriptionExist[topic] = subID
		c.Unlock()

		return s, nil
	}

	// Get existing subscription from topic.
	t, err := c.getTopic(topic)
	if err != nil {
		return nil, err
	}

	subs := t.Subscriptions(context.Background())
	for {
		sub, err := subs.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		subID = sub.ID()

		c.Lock()
		c.subscriptionExist[topic] = subID
		c.Unlock()

		return sub, nil
	}

	// Still not exist.
	if s, err = c.client.CreateSubscription(context.Background(), subID, pubsub.SubscriptionConfig{
		Topic: t,
	}); err != nil {
		return nil, err
	}

	c.Lock()
	c.subscriptionExist[topic] = subID
	c.Unlock()

	return s, nil
}
