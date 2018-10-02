package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/pkg/errors"
)

type TorchMetrics struct {
	client     *http.Client
	clientLock sync.Mutex
	host       string
}

func NewTrochMetrics(host string) (*TorchMetrics, error) {
	return &TorchMetrics{
		client: &http.Client{},
		host:   host,
	}, nil
}

type TorchMetricsSessionGrid struct {
	DisplayName string
	EntityId    int64
	GridSize    string
	BlocksCount int
	Mass        float64
	Position    struct {
		X float64
		Y float64
		Z float64
	}
	LinearSpeed      float64
	DistanceToPlayer float64
	OwnerSteamID     int64 `json:"OwnerSteamId"`
	OwnerDisplayName string
	OwnerFactionTag  string
	OwnerFactionName string
	IsPowered        bool
	PCU              int
	Concealed        bool
	DampenersEnabled bool
	IsStatic         bool
}

func (t *TorchMetrics) SessionGrids() ([]TorchMetricsSessionGrid, error) {
	res, err := t.client.Get(fmt.Sprintf("%s/metrics/v1/session/grids", t.host))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}
	var grids []TorchMetricsSessionGrid
	err = json.NewDecoder(res.Body).Decode(&grids)
	if err != nil {
		return nil, err
	}

	return grids, nil
}

type TorchMetricServer struct {
	Version            string
	ServerName         string
	WorldName          string
	IsReady            bool
	SimSpeed           float64
	SimulationCpuLoad  float64
	TotalTime          int64
	Players            int8
	UsedPCU            int
	MaxPlayers         int
	MaxFactionsCount   int
	MaxFloatingObjects int
	MaxGridSize        int
	MaxBlocksPerPlayer int
	BlockLimitEnabled  string
	TotalPCU           int
}

func (t *TorchMetrics) ServerInfo() (*TorchMetricServer, error) {
	res, err := t.client.Get(fmt.Sprintf("%s/metrics/v1/server", t.host))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}

	var server TorchMetricServer
	err = json.NewDecoder(res.Body).Decode(&server)
	if err != nil {
		return nil, err
	}

	return &server, nil
}
