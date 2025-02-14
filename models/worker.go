package models

// Worker represents a remote worker with payout tracking.
type Worker interface {
	GetID() string
	GetPendingFees() int64 //wei values
	GetPaidFees() int64    //wei values
	GetNodeType() string
	GetRegion() string
}

type DefaultWorker struct {
	ID          string
	PendingFees int64
	PaidFees    int64
	NodeType    string
	Online      bool
	Region      string
}

func (e DefaultWorker) GetPendingFees() int64 {
	return e.PendingFees
}
func (e DefaultWorker) GetPaidFees() int64 {
	return e.PaidFees
}

func (e DefaultWorker) GetNodeType() string {
	return e.NodeType
}

func (e DefaultWorker) GetID() string {
	return e.ID
}

func (e DefaultWorker) GetRegion() string {
	return e.Region
}
