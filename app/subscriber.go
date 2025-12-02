package main

import "github.com/codecrafters-io/redis-starter-go/app/protocol"

type ChannelDetails struct {
	protocol *protocol.Protocol
}

type Subscriber struct {
	channels map[string][]ChannelDetails
}

var SubscriberInstance = NewSubsciber()

func NewSubsciber() *Subscriber {
	return &Subscriber{
		channels: make(map[string][]ChannelDetails),
	}
}

type SubscribeCount struct {
	channelName string
	count       int
}

func (subscriber *Subscriber) Subscribe(channelNames []string, protocol *protocol.Protocol) []SubscribeCount {
	counts := []SubscribeCount{}
	for _, channelName := range channelNames {
		subscriber.subscribe(channelName, protocol)
		counts = append(counts, SubscribeCount{
			channelName: channelName,
			count:       protocol.SubscribeCount,
		})
	}

	return counts
}

func (subscriber *Subscriber) subscribe(channelName string, protocol *protocol.Protocol) {
	channel := subscriber.channels[channelName]

	// Check for existing subscription
	for _, details := range channel {
		if details.protocol.GetClientId() == protocol.GetClientId() {
			return
		}
	}

	// Add new subscription
	subscriber.channels[channelName] = append(channel, ChannelDetails{protocol: protocol})
	protocol.SubscribeCount++
}

func (subscriber *Subscriber) unsubscribe(channelName string, clientId string) {
	channel := subscriber.channels[channelName]

	for i, details := range channel {
		if details.protocol.GetClientId() == clientId {
			details.protocol.SubscribeCount--
			subscriber.channels[channelName] = append(channel[:i], channel[i+1:]...)
			return
		}
	}
}

func (subscriber Subscriber) publish(channelName string, message string) {
	channels := subscriber.channels[channelName]

	for _, channel := range channels {
		channel.protocol.Write([]protocol.RESPValue{protocol.BulkString{Value: message}})
	}
}
