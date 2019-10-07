package data

import (
	"errors"
	"reflect"
	"testing"
)

func TestConsumer_HandleConsumeSuccess(t *testing.T) {
	tests := []struct {
		name          string
		existingAlert *Alert
		expectedAlert *Alert
	}{
		{"alert nil", nil, nil},
		{"alert present, not recoverable", &Alert{Recoverable: false}, &Alert{Recoverable: false}},
		{"alert present, recoverable", &Alert{Recoverable: true}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConsumer()
			c.Alert = tt.existingAlert
			c.HandleConsumeSuccess()
			if !reflect.DeepEqual(c.Alert, tt.expectedAlert) {
				t.Errorf("unexpected alert state after HandleConsumeSuccess(), want %v, got %v", tt.expectedAlert, c.Alert)
			}
		})
	}
}

func TestConsumer_HandleConsumeFailure(t *testing.T) {
	tests := []struct {
		name   string
		title  string
		err    error
		sample *Sample
		want   *Alert
	}{
		{
			"basic",
			"test",
			errors.New("test"),
			&Sample{
				Label: "test",
				Value: "test",
				Color: nil,
			},
			&Alert{
				Title:       "TEST",
				Text:        "test",
				Color:       nil,
				Recoverable: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConsumer()
			c.HandleConsumeFailure(tt.title, tt.err, tt.sample)
			got := <-c.AlertChannel
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("unexpected alert state after HandleConsumeFailure(), want %v, got %v", got, tt.want)
			}
		})
	}
}

func TestNewConsumer(t *testing.T) {
	tests := []struct {
		name      string
		checkFunc func(c *Consumer) bool
	}{
		{
			"initialized after creation",
			func(c *Consumer) bool {
				return c != nil
			},
		},
		{
			"alert is nil after creation",
			func(c *Consumer) bool {
				return c.Alert == nil
			},
		},
		{
			"command channel is initialized after creation",
			func(c *Consumer) bool {
				return c.CommandChannel != nil
			},
		},
		{
			"alert channel is initialized after creation",
			func(c *Consumer) bool {
				return c.AlertChannel != nil
			},
		},
		{
			"sample channel is initialized after creation",
			func(c *Consumer) bool {
				return c.SampleChannel != nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewConsumer(); !tt.checkFunc(got) {
				t.Errorf("unexpected consumer state after NewConsumer() = %v", got)
			}
		})
	}
}
