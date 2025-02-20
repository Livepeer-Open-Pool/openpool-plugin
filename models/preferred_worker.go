package models

type PreferredWorkerCriteria interface {
	GetNodeType() string
	GetRegion() string
	GetCriteria() any
}

// WorkerCriteria is a generic implementation of PreferredWorkerCriteria.
type WorkerCriteria struct {
	NodeType string
	Region   string
	// Additional filtering data can be stored in this map.
	Criteria map[string]any
}

func (c WorkerCriteria) GetNodeType() string {
	return c.NodeType
}

func (c WorkerCriteria) GetRegion() string {
	return c.Region
}

func (c WorkerCriteria) GetCriteria() any {
	return c.Criteria
}
