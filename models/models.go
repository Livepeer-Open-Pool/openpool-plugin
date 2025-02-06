package models

import (
	"encoding/json"
	"fmt"
	"time"
)

// EventLog represents the event_log table in the database.
type EventLog struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Type         string    `json:"type"`
	Data         string    `json:"data"`
	CreatedAt    time.Time `json:"created_at"`
	EndpointHash string    `json:"endpoint_hash"` // NEW: Identifies the endpoint that produced this event.
}

// PoolEvent represents each event in the pool.
type PoolEvent struct {
	ID            int            `json:"ID"`
	Payload       string         `json:"Payload"`
	Version       int            `json:"Version"`
	DT            time.Time      `json:"DT"`
	ParsedPayload PayloadWrapper `json:"-"`
}

// ParsePayload parses the Payload and populates the ParsedPayload field.
func (pe *PoolEvent) ParsePayload() error {
	return json.Unmarshal([]byte(pe.Payload), &pe.ParsedPayload)
}

// PayloadWrapper handles the dynamic payload based on event_type.
type PayloadWrapper struct {
	EventType string      `json:"event_type"`
	Payload   interface{} `json:"payload"`
}

// UnmarshalJSON implements custom unmarshaling for PayloadWrapper.
func (pw *PayloadWrapper) UnmarshalJSON(data []byte) error {
	var temp struct {
		EventType string          `json:"event_type"`
		Payload   json.RawMessage `json:"payload"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return fmt.Errorf("failed to unmarshal PayloadWrapper: %w", err)
	}

	pw.EventType = temp.EventType

	switch temp.EventType {
	case "orchestrator-reset":
		var pr OrchestratorResetPayload
		if err := json.Unmarshal(temp.Payload, &pr); err != nil {
			return fmt.Errorf("failed to unmarshal orchestrator-reset payload: %w", err)
		}
		pw.Payload = pr

	case "worker-connected":
		var wc RemoteWorker
		if err := json.Unmarshal(temp.Payload, &wc); err != nil {
			return fmt.Errorf("failed to unmarshal worker-connected payload: %w", err)
		}
		pw.Payload = wc

	case "worker-disconnected":
		var wc RemoteWorker
		if err := json.Unmarshal(temp.Payload, &wc); err != nil {
			return fmt.Errorf("failed to unmarshal worker-disconnected payload: %w", err)
		}
		pw.Payload = wc

	case "job-received":
		var jrp JobReceivedPayload
		if err := json.Unmarshal(temp.Payload, &jrp); err != nil {
			return fmt.Errorf("failed to unmarshal job-received payload: %w", err)
		}
		pw.Payload = jrp

	case "job-processed":
		var jpp JobProcessedPayload
		if err := json.Unmarshal(temp.Payload, &jpp); err != nil {
			return fmt.Errorf("failed to unmarshal job-processed payload: %w", err)
		}
		pw.Payload = jpp

	default:
		var unknown map[string]interface{}
		if err := json.Unmarshal(temp.Payload, &unknown); err != nil {
			return fmt.Errorf("failed to unmarshal unknown event_type payload: %w", err)
		}
		pw.Payload = unknown
	}

	return nil
}

// OrchestratorResetPayload represents the payload for the "orchestrator-reset" event.
type OrchestratorResetPayload struct{}

// PoolPayout represents the pool payout record.
type PoolPayout struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	EthAddress string    `json:"ethAddress"`
	TxHash     string    `json:"txHash"`
	Fees       int64     `json:"fees"`
	CreatedAt  time.Time `json:"createdAt" gorm:"autoUpdateTime"`
}

// RemoteWorker represents a worker. The unique composite key is built from EthAddress, NodeType, and Region.
type RemoteWorker struct {
	EthAddress   string    `json:"ethAddress" gorm:"primaryKey;not null"`
	NodeType     string    `json:"nodeType" gorm:"primaryKey;not null"`
	Region       string    `json:"region" gorm:"primaryKey;not null"`
	EndpointHash string    `json:"endpoint_hash" gorm:"primaryKey;not null;index"` // NEW: identifies the data source
	IsConnected  bool      `json:"is_connected"`
	PendingFees  int64     `json:"pending_fees"`
	PaidFees     int64     `json:"paid_fees"`
	LastUpdated  time.Time `json:"last_updated" gorm:"autoUpdateTime"`
	Connection   string    `json:"connection,omitempty"`
}

// JobReceivedPayload represents the payload for the "job-received" event.
type JobReceivedPayload struct {
	EthAddress string `json:"ethAddress"`
	ModelID    string `json:"modelID"`
	NodeType   string `json:"nodeType"`
	Pipeline   string `json:"pipeline"`
	RequestID  string `json:"requestID"`
	TaskID     int    `json:"taskID"`
}

// JobProcessedPayload represents the payload for the "job-processed" event.
type JobProcessedPayload struct {
	ComputeUnits        int64  `json:"computeUnits"`
	NodeType            string `json:"nodeType"`
	PricePerComputeUnit int64  `json:"pricePerComputeUnit"`
	Fees                int64  `json:"fees"`
	EthAddress          string `json:"ethAddress,omitempty"`
	ResponseTime        int64  `json:"responseTime,omitempty"`
	RequestID           string `json:"requestID,omitempty"`
	ModelID             string `json:"modelID,omitempty"`
	Pipeline            string `json:"pipeline,omitempty"`
}
