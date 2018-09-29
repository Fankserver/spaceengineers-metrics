package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"flag"

	"context"

	"github.com/influxdata/influxdb/client/v2"
	"golang.org/x/sync/errgroup"
)

var (
	host       = flag.String("host", "http://localhost:8080", "host url of the rcon server")
	key        = flag.String("key", "", "rcon key")
	influxhost = flag.String("influxhost", "http://localhost:8086", "influxdb host")
	influxdb   = flag.String("influxdb", "spaceengineers", "influxdb database")
	influxuser = flag.String("influxuser", "", "influx username")
	influxpass = flag.String("influxpass", "", "influx password")
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

	errWg, gCtx := errgroup.WithContext(context.Background())
	errWg.Go(func() error {
		running := false
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-gCtx.Done():
				return gCtx.Err()
			case <-ticker.C:
				if running {
					continue
				}
				running = true
				info, err := s.ServerInfo()
				if err != nil {
					return err
				}

				// Create a new point batch
				bp, err := client.NewBatchPoints(client.BatchPointsConfig{
					Database:  *influxdb,
					Precision: "s",
				})
				if err != nil {
					return err
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
					return err
				}
				bp.AddPoint(pt)

				// Write the batch
				if err := c.Write(bp); err != nil {
					return err
				}
				running = false
			}
		}

		return nil
	})
	errWg.Go(func() error {
		running := false
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-gCtx.Done():
				return gCtx.Err()
			case <-ticker.C:
				if running {
					continue
				}
				running = true

				grids, err := s.SessionGrids()
				if err != nil {
					return err
				}

				// Create a new point batch
				bp, err := client.NewBatchPoints(client.BatchPointsConfig{
					Database:  *influxdb,
					Precision: "s",
				})
				if err != nil {
					return err
				}

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
							"grid_size":          grid.GridSize,
						},
						map[string]interface{}{
							"blocks_count": grid.BlocksCount,
							"is_powered":   powered,
							"linear_speed": grid.LinearSpeed,
							"mass":         grid.Mass,
							"pcu":          grid.PCU,
						},
					)
					if err != nil {
						return err
					}
					bp.AddPoint(pt)
				}

				// Write the batch
				if err := c.Write(bp); err != nil {
					return err
				}
				running = false
			}
		}

		return nil
	})

	err = errWg.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
