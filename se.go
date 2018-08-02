package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"strings"

	"github.com/pkg/errors"
)

type SpaceEngineers struct {
	client *http.Client
	host   string
	key    []byte
	random *rand.Rand
}

func New(host string, key string) (*SpaceEngineers, error) {
	decodedKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}

	return &SpaceEngineers{
		client: &http.Client{},
		host:   host,
		key:    decodedKey,
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
	}, nil
}

type SessionAsteroid struct {
	DisplayName *string
	EntityId    int64
	Position    struct {
		X float64
		Y float64
		Z float64
	}
}

func (s *SpaceEngineers) SessionAsteroids() ([]SessionAsteroid, error) {
	req, err := s.createRequest("/v1/session/asteroids", http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	res, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}
	var responseData struct {
		Data struct {
			Asteroids []SessionAsteroid
		} `json:"data"`
	}
	err = json.NewDecoder(res.Body).Decode(&responseData)
	if err != nil {
		return nil, err
	}

	return responseData.Data.Asteroids, nil
}

type SessionGrid struct {
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
	IsPowered        bool
	PCU              int
}

func (s *SpaceEngineers) SessionGrids() ([]SessionGrid, error) {
	req, err := s.createRequest("/v1/session/grids", http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	res, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}
	var responseData struct {
		Data struct {
			Grids []SessionGrid
		} `json:"data"`
	}
	err = json.NewDecoder(res.Body).Decode(&responseData)
	if err != nil {
		return nil, err
	}

	return responseData.Data.Grids, nil
}

//func (s *SpaceEngineers) SessionFloatingObjects() (*string, error) {
//	req, err := s.createRequest("/v1/session/floatingObjects", http.MethodGet, nil)
//	if err != nil {
//		return nil, err
//	}
//
//	res, err := s.client.Do(req)
//	if err != nil {
//		return nil, err
//	}
//
//	defer res.Body.Close()
//	if res.StatusCode != http.StatusOK {
//		return nil, errors.New(res.Status)
//	}
//	//var responseData struct {
//	//	Data SessionAsteroids `json:"data"`
//	//}
//	b, _ := ioutil.ReadAll(res.Body)
//	log.Fatal(string(b))
//	//err = json.NewDecoder(res.Body).Decode(&responseData)
//	//if err != nil {
//	//	return nil, err
//	//}
//
//	return nil, nil
//}

type ServerInfo struct {
	Game              string
	IsReady           bool
	Players           int8
	ServerID          int64 `json:"ServerId"`
	ServerName        string
	SimSpeed          float64
	SimulationCpuLoad float64
	TotalTime         int64
	UsedPCU           int64
	Version           string
	WorldName         string
}

func (s *SpaceEngineers) ServerInfo() (*ServerInfo, error) {
	req, err := s.createRequest("/v1/server", http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	res, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}
	var responseData struct {
		Data ServerInfo `json:"data"`
	}
	err = json.NewDecoder(res.Body).Decode(&responseData)
	if err != nil {
		return nil, err
	}

	return &responseData.Data, nil
}

func (s *SpaceEngineers) createRequest(resourceLink string, method string, queryParams map[string]string) (*http.Request, error) {
	url := "/vrageremote" + resourceLink
	req, err := http.NewRequest(method, s.host+url, nil)
	if err != nil {
		return nil, err
	}

	// Space engineers requires a GMT date, UTC won't work.
	t := strings.Replace(time.Now().UTC().Format(time.RFC1123), "UTC", "GMT", 1)

	req.Header.Add("Date", t)

	randomNumber := s.random.Uint32()
	message := req.URL.Path + "\r\n"
	message += fmt.Sprint(randomNumber) + "\r\n"
	message += t + "\r\n"

	mac := hmac.New(sha1.New, s.key)
	mac.Write([]byte(message))

	req.Header.Add("Authorization", fmt.Sprint(randomNumber)+":"+base64.StdEncoding.EncodeToString(mac.Sum(nil)))

	return req, nil
}
