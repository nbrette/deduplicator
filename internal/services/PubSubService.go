package services

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"cloud.google.com/go/pubsub"
)

const topicFullPathFormat = "projects/PROJECT/topics/TOPIC"

type PubSubService struct {
	topicsMap       map[string]*pubsub.Topic
	topicsMapMutex  sync.RWMutex
	ClientsMap      map[string]*pubsub.Client
	clientsMapMutex sync.RWMutex
}

func NewPubSubService() *PubSubService {
	topicsMap := make(map[string]*pubsub.Topic)
	clientsMap := make(map[string]*pubsub.Client)
	return &PubSubService{
		topicsMap:       topicsMap,
		topicsMapMutex:  sync.RWMutex{},
		ClientsMap:      clientsMap,
		clientsMapMutex: sync.RWMutex{},
	}
}

func (pubSubService *PubSubService) Publish(topicFullPath string, messagePubSub *pubsub.Message) error {
	splittedTopic := strings.Split(topicFullPath, "/")
	if len(splittedTopic) != len(strings.Split(topicFullPathFormat, "/")) {
		return fmt.Errorf("topic name should have this format: %s, got %s", topicFullPathFormat, topicFullPath)
	}
	project := splittedTopic[1]
	topicName := splittedTopic[3]

	pubSubService.topicsMapMutex.Lock()
	defer pubSubService.topicsMapMutex.Unlock()
	topic, ok := pubSubService.topicsMap[topicName]
	if ok {
		err := pubSubService.publish(topic, messagePubSub)
		if err != nil {
			return err
		}
		return nil
	}

	pubSubService.clientsMapMutex.Lock()
	defer pubSubService.clientsMapMutex.Unlock()
	client, ok := pubSubService.ClientsMap[project]
	if ok {
		topic = client.Topic(topicName)
		pubSubService.topicsMap[topicName] = topic
		err := pubSubService.publish(topic, messagePubSub)
		if err != nil {
			return err
		}
		return nil
	}
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, project)
	if err != nil {
		return err
	}
	pubSubService.ClientsMap[project] = client
	topic = client.Topic(topicName)
	pubSubService.topicsMap[topicName] = topic
	err = pubSubService.publish(topic, messagePubSub)
	if err != nil {
		return err
	}
	return nil
}

func (pubSubService *PubSubService) publish(topic *pubsub.Topic, messagePubSub *pubsub.Message) error {
	ctx := context.Background()
	result := topic.Publish(ctx, messagePubSub)
	_, err := result.Get(ctx)
	if err != nil {
		return err
	}
	return nil
}
