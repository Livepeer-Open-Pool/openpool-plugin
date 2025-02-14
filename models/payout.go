package models

// Payout represents a recorded payout for a worker.
type Payout interface {
	GetID() string
	GetTimestamp() int64
	GetAmount() int64 //wei units
}

type DefaultPayout struct {
	ID        string
	Timestamp int64
	Amount    int64
}

func (e DefaultPayout) GetTimestamp() int64 {
	return e.Timestamp
}

func (e DefaultPayout) GetID() string {
	return e.ID
}

func (e DefaultPayout) GetAmount() int64 {
	return e.Amount
}
