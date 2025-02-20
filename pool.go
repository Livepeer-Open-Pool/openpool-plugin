package pool

import (
	"github.com/Livepeer-Open-Pool/openpool-plugin/config"
	"github.com/Livepeer-Open-Pool/openpool-plugin/models"
	"time"
)

// PluginInterface defines the methods all plugins must implement.
type PluginInterface interface {
	Init(cfg config.Config, store StorageInterface)
	Start()
}

// StorageInterface defines the methods a storage plugin must implement.
type StorageInterface interface {
	Init(config *config.Config)
	GetLastEventTimestamp() (time.Time, error)
	AddEvent(event models.PoolEvent) error
	GetWorkers() ([]models.Worker, error)
	GetPreferredWorkers(criteria models.PreferredWorkerCriteria) ([]models.Worker, error)
	AddPaidFees(ethAddress string, amount int64, txhash string, region string, nodeType string) error //wei values
	UpdateWorkerStatus(ethAddress string, online bool, region string, nodeType string) error
	AddPendingFees(ethAddress string, amount int64, region string, nodeType string) error //wei values
	ResetWorkersOnlineStatus(region string, nodeType string) error
	GetPendingFees() (float64, error) //eth value
	GetPaidFees() (float64, error)    //eth value
}
