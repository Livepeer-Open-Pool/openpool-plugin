package models

type PoolEvent interface {
	GetTimestamp() int64
	GetData() string
	GetType() string
}

type DefaultPoolEvent struct {
	Timestamp int64
	Data      string
	Type      string
}

func (e DefaultPoolEvent) GetTimestamp() int64 {
	return e.Timestamp
}

func (e DefaultPoolEvent) GetData() string {
	return e.Data
}

func (e DefaultPoolEvent) GetType() string {
	return e.Type
}
