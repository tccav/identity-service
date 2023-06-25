package kafka

import (
	"context"
	"encoding/json"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Producer struct {
	client *kgo.Client
}

func NewProducer(client *kgo.Client) Producer {
	return Producer{
		client: client,
	}
}

type produceInput struct {
	topic   string
	headers map[string]string
	event   event
}

func (p Producer) produce(ctx context.Context, input produceInput) error {
	eventValue, err := json.Marshal(input.event)
	if err != nil {
		return err
	}

	record := kgo.Record{
		Headers: recordHeadersFromHeaders(input.headers),
		Value:   eventValue,
		Topic:   input.topic,
	}
	result := p.client.ProduceSync(ctx, &record)
	if err = result.FirstErr(); err != nil {
		return err
	}
	return nil
}

func recordHeadersFromHeaders(headers map[string]string) []kgo.RecordHeader {
	var recordHeaders []kgo.RecordHeader
	for k, v := range headers {
		recordHeaders = append(recordHeaders, kgo.RecordHeader{
			Key:   k,
			Value: []byte(v),
		})
	}
	return recordHeaders
}
