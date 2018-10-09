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

func NewTorchMetrics(host string) (*TorchMetrics, error) {
	return &TorchMetrics{
		client: &http.Client{},
		host:   host,
	}, nil
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
	ModCount           int
	SaveDuration       int64
}

func (t *TorchMetrics) Server() (*TorchMetricServer, error) {
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

type TorchMetricsProcess struct {
	PrivateMemorySize64        int64
	VirtualMemorySize64        int64
	WorkingSet64               int64
	NonpagedSystemMemorySize64 int64
	PagedMemorySize64          int64
	PagedSystemMemorySize64    int64
	PeakPagedMemorySize64      int64
	PeakVirtualMemorySize64    int64
	PeakWorkingSet64           int64
	GCLatencyMode              int8
	GCTotalMemory              int64
	GCMaxGeneration            int
	GCCollectionCount0         int
	GCCollectionCount1         int
	GCCollectionCount2         int
}

type TorchMetricsLoad struct {
	ServerCPULoad          float64
	ServerCPULoadSmooth    float64
	ServerSimulationRatio  float64
	ServerThreadLoad       float64
	ServerThreadLoadSmooth float64
	MillisecondsInThePast  int64
}

func (t *TorchMetrics) Load() ([]TorchMetricsLoad, error) {
	res, err := t.client.Get(fmt.Sprintf("%s/metrics/v1/load", t.host))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}
	var loads []TorchMetricsLoad
	err = json.NewDecoder(res.Body).Decode(&loads)
	if err != nil {
		return nil, err
	}

	return loads, nil
}

func (t *TorchMetrics) Process() (*TorchMetricsProcess, error) {
	res, err := t.client.Get(fmt.Sprintf("%s/metrics/v1/process", t.host))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}

	var server TorchMetricsProcess
	err = json.NewDecoder(res.Body).Decode(&server)
	if err != nil {
		return nil, err
	}

	return &server, nil
}

type TorchMetricsEvent struct {
	Type             string
	Text             string
	Tags             []string
	SecondsInThePast int
}

func (t *TorchMetrics) Events() ([]TorchMetricsEvent, error) {
	res, err := t.client.Get(fmt.Sprintf("%s/metrics/v1/events", t.host))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}
	var events []TorchMetricsEvent
	err = json.NewDecoder(res.Body).Decode(&events)
	if err != nil {
		return nil, err
	}

	return events, nil
}

type TorchMetricsSessionGrid struct {
	DisplayName      string
	EntityId         int64
	GridSize         string
	BlocksCount      int
	Mass             float64
	LinearSpeed      float64
	DistanceToPlayer float64
	OwnerSteamID     int64 `json:"OwnerSteamId"`
	OwnerDisplayName string
	OwnerFactionTag  string
	OwnerFactionName string
	IsPowered        bool
	PCU              int
	IsConcealed      bool
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

type TorchMetricsSessionAsteroidOrPlanet struct {
	DisplayName string
	EntityId    int64
}

func (t *TorchMetrics) SessionAsteroids() ([]TorchMetricsSessionAsteroidOrPlanet, error) {
	res, err := t.client.Get(fmt.Sprintf("%s/metrics/v1/session/asteroids", t.host))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}
	var grids []TorchMetricsSessionAsteroidOrPlanet
	err = json.NewDecoder(res.Body).Decode(&grids)
	if err != nil {
		return nil, err
	}

	return grids, nil
}

func (t *TorchMetrics) SessionPlanets() ([]TorchMetricsSessionAsteroidOrPlanet, error) {
	res, err := t.client.Get(fmt.Sprintf("%s/metrics/v1/session/planets", t.host))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}
	var grids []TorchMetricsSessionAsteroidOrPlanet
	err = json.NewDecoder(res.Body).Decode(&grids)
	if err != nil {
		return nil, err
	}

	return grids, nil
}

type TorchMetricsSessionFloatingObject struct {
	DisplayName      string
	EntityId         int64
	Kind             string
	Mass             float64
	LinearSpeed      float64
	DistanceToPlayer float64
	TypeDisplayName  string
}

func (t *TorchMetrics) SessionFloatingObjects() ([]TorchMetricsSessionFloatingObject, error) {
	res, err := t.client.Get(fmt.Sprintf("%s/metrics/v1/session/floatingObjects", t.host))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}
	var grids []TorchMetricsSessionFloatingObject
	err = json.NewDecoder(res.Body).Decode(&grids)
	if err != nil {
		return nil, err
	}

	return grids, nil
}

type TorchMetricsSessionFaction struct {
	AcceptHumans       bool
	AutoAcceptMember   bool
	AutoAcceptPeace    bool
	EnableFriendlyFire bool
	FactionId          int64
	FounderId          int64
	MemberCount        int
	Name               string
	Tag                string
	NPCOnly            bool
}

func (t *TorchMetrics) SessionFactions() ([]TorchMetricsSessionFaction, error) {
	res, err := t.client.Get(fmt.Sprintf("%s/metrics/v1/session/factions", t.host))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}
	var grids []TorchMetricsSessionFaction
	err = json.NewDecoder(res.Body).Decode(&grids)
	if err != nil {
		return nil, err
	}

	return grids, nil
}
