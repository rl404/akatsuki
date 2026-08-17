package google

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Client) getPublisher(topic string) (*pubsub.Publisher, error) {
	if c.topicExist[topic] {
		return c.client.Publisher(topic), nil
	}

	// Check existing topic.
	if _, err := c.client.TopicAdminClient.GetTopic(context.Background(), &pubsubpb.GetTopicRequest{
		Topic: fmt.Sprintf("projects/%s/topics/%s", c.projectID, topic),
	}); err != nil {
		errStatus, ok := status.FromError(err)
		if !ok {
			return nil, err
		}

		if errStatus.Code() != codes.NotFound {
			return nil, err
		}

		// Create topic.
		if _, err := c.client.TopicAdminClient.CreateTopic(context.Background(), &pubsubpb.Topic{
			Name: fmt.Sprintf("projects/%s/topics/%s", c.projectID, topic),
		}); err != nil {
			return nil, err
		}

		// Also create subscription.
		if _, err := c.client.SubscriptionAdminClient.CreateSubscription(context.Background(), &pubsubpb.Subscription{
			Name:  fmt.Sprintf("projects/%s/subscriptions/%s-subscription", c.projectID, topic),
			Topic: fmt.Sprintf("projects/%s/topics/%s", c.projectID, topic),
		}); err != nil {
			return nil, err
		}
	}

	c.Lock()
	c.topicExist[topic] = true
	c.Unlock()

	return c.client.Publisher(topic), nil
}

func (c *Client) getSubscriber(topic string) (*pubsub.Subscriber, error) {
	if c.subscriptionExist[topic] != "" {
		return c.client.Subscriber(c.subscriptionExist[topic]), nil
	}

	// Check existing topic.
	if _, err := c.getPublisher(topic); err != nil {
		return nil, err
	}

	// Get existing subscriber.
	it := c.client.TopicAdminClient.ListTopicSubscriptions(context.Background(), &pubsubpb.ListTopicSubscriptionsRequest{
		Topic: fmt.Sprintf("projects/%s/topics/%s", c.projectID, topic),
	})
	for {
		subscriberName, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		c.Lock()
		c.subscriptionExist[topic] = subscriberName
		c.Unlock()

		return c.client.Subscriber(subscriberName), nil
	}

	// Still not exist.
	subscriberName := fmt.Sprintf("%s-subscription", topic)

	if _, err := c.client.SubscriptionAdminClient.CreateSubscription(context.Background(), &pubsubpb.Subscription{
		Name:  fmt.Sprintf("projects/%s/subscriptions/%s", c.projectID, subscriberName),
		Topic: fmt.Sprintf("projects/%s/topics/%s", c.projectID, topic),
	}); err != nil {
		return nil, err
	}

	c.Lock()
	c.subscriptionExist[topic] = subscriberName
	c.Unlock()

	return c.client.Subscriber(subscriberName), nil
}
