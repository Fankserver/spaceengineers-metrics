package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"flag"

	client "github.com/influxdata/influxdb/client/v2"
)

var (
	host       = flag.String("host", "http://localhost:8080", "host url of the rcon server")
	key        = flag.String("key", "", "rcon key")
	influxhost = flag.String("influxhost", "http://localhost:8086", "influxdb host")
	influxuser = flag.String("influxuser", "", "influx username")
	influxpass = flag.String("ifnluxpass", "", "influx password")
)

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	s, err := New(*host, *key)
	if err != nil {
		log.Fatal(err)
	}

	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     *influxhost,
		Username: *influxuser,
		Password: *influxpass,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		info, err := s.ServerInfo()
		if err != nil {
			log.Fatal(err)
		}

		grids, err := s.SessionGrids()
		if err != nil {
			log.Fatal(err)
		}

		// Create a new point batch
		bp, err := client.NewBatchPoints(client.BatchPointsConfig{
			Database:  "test",
			Precision: "s",
		})
		if err != nil {
			log.Fatal(err)
		}

		tags := map[string]string{
			"host":        *host,
			"server_name": info.ServerName,
			"server_id":   fmt.Sprint(info.ServerID),
			"version":     info.Version,
			"world_name":  info.WorldName,
		}

		ready := 0
		if info.IsReady {
			ready++
		}
		pt, err := client.NewPoint("spaceengineers", tags, map[string]interface{}{
			"sim_speed":    info.SimSpeed,
			"players":      info.Players,
			"sim_cpu_load": info.SimulationCpuLoad,
			"total_time":   info.TotalTime,
			"used_pcu":     info.UsedPCU,
			"ready":        ready,
		}, time.Now())
		if err != nil {
			log.Fatal(err)
		}
		bp.AddPoint(pt)

		for _, grid := range grids {
			powered := 0
			if grid.IsPowered {
				powered++
			}
			pt, err := client.NewPoint(
				"spaceengineers_grids",
				map[string]string{
					"host":               *host,
					"owner_steam_id":     fmt.Sprint(grid.OwnerSteamID),
					"owner_display_name": grid.OwnerDisplayName,
					"display_name":       grid.DisplayName,
					"entity_id":          fmt.Sprint(grid.EntityId),
					"is_powered":         fmt.Sprint(powered),
				},
				map[string]interface{}{
					"blocks_count": grid.BlocksCount,
					"grid_size":    grid.GridSize,
					"is_powered":   powered,
					"linear_speed": grid.LinearSpeed,
					"mass":         grid.Mass,
					"pcu":          grid.PCU,
				},
			)
			if err != nil {
				log.Fatal(err)
			}
			bp.AddPoint(pt)
		}

		// Write the batch
		if err := c.Write(bp); err != nil {
			log.Fatal(err)
		}
	}
}
