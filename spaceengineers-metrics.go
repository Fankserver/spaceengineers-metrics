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
	host = flag.String("host", "http://localhost:8080", "host url of the rcon server")
	//key        = flag.String("key", "", "rcon key")
	influxhost = flag.String("influxhost", "http://localhost:8086", "influxdb host")
	influxdb   = flag.String("influxdb", "spaceengineers", "influxdb database")
	influxuser = flag.String("influxuser", "", "influx username")
	influxpass = flag.String("influxpass", "", "influx password")
)

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	t, err := NewTrochMetrics(*host)
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
				info, err := t.ServerInfo()
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

				ready := 0
				if info.IsReady {
					ready++
				}
				pt, err := client.NewPoint(
					"server",
					map[string]string{
						"host":        *host,
						"server_name": info.ServerName,
						"version":     info.Version,
						"world_name":  info.WorldName,
						"block_limit": info.BlockLimitEnabled,
					},
					map[string]interface{}{
						"sim_speed":             info.SimSpeed,
						"players":               info.Players,
						"sim_cpu_load":          info.SimulationCpuLoad,
						"total_time":            info.TotalTime,
						"used_pcu":              info.UsedPCU,
						"ready":                 ready,
						"max_blocks_per_player": info.MaxBlocksPerPlayer,
						"max_factions_count":    info.MaxFactionsCount,
						"max_floating_objects":  info.MaxFloatingObjects,
						"max_grid_size":         info.MaxGridSize,
						"max_players":           info.MaxPlayers,
						"block_limit":           info.BlockLimitEnabled,
						"total_pcu":             info.TotalPCU,
					},
					time.Now(),
				)
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

				grids, err := t.SessionGrids()
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
					concealed := 0
					dampenersEnabled := 0
					isStatic := 0
					if grid.IsPowered {
						powered++
					}
					if grid.IsConcealed {
						concealed++
					}
					if grid.DampenersEnabled {
						dampenersEnabled++
					}
					if grid.IsStatic {
						isStatic++
					}
					pt, err := client.NewPoint(
						"grid",
						map[string]string{
							"host":                *host,
							"owner_steam_id":      fmt.Sprint(grid.OwnerSteamID),
							"owner_display_name":  grid.OwnerDisplayName,
							"owner_faction_tag":   grid.OwnerFactionTag,
							"owner_faction_name":  grid.OwnerFactionName,
							"display_name":        grid.DisplayName,
							"entity_id":           fmt.Sprint(grid.EntityId),
							"filter_is_powered":   toStringBool(grid.IsPowered),
							"grid_size":           grid.GridSize,
							"filter_is_concealed": toStringBool(grid.IsConcealed),
							"filter_is_static":    toStringBool(grid.IsStatic),
						},
						map[string]interface{}{
							"blocks_count":      grid.BlocksCount,
							"is_powered":        powered,
							"linear_speed":      grid.LinearSpeed,
							"mass":              grid.Mass,
							"pcu":               grid.PCU,
							"is_concealed":      concealed,
							"dampeners_enabled": dampenersEnabled,
							"is_static":         isStatic,
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

				asteroids, err := t.SessionAsteroids()
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

				for _, asteroid := range asteroids {
					pt, err := client.NewPoint(
						"asteroid",
						map[string]string{
							"host":         *host,
							"display_name": asteroid.DisplayName,
							"entity_id":    fmt.Sprint(asteroid.EntityId),
						},
						map[string]interface{}{
							"display_name": asteroid.DisplayName,
							"entity_id":    fmt.Sprint(asteroid.EntityId),
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

				planets, err := t.SessionPlanets()
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

				for _, planet := range planets {
					pt, err := client.NewPoint(
						"planet",
						map[string]string{
							"host":         *host,
							"display_name": planet.DisplayName,
							"entity_id":    fmt.Sprint(planet.EntityId),
						},
						map[string]interface{}{
							"display_name": planet.DisplayName,
							"entity_id":    fmt.Sprint(planet.EntityId),
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

				factions, err := t.SessionFactions()
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

				for _, faction := range factions {
					npconly := 0
					acceptHumans := 0
					autoAcceptMember := 0
					autoAcceptPeace := 0
					enableFriendlyFire := 0
					if faction.NPCOnly {
						npconly++
					}
					if faction.AcceptHumans {
						acceptHumans++
					}
					if faction.AutoAcceptMember {
						autoAcceptMember++
					}
					if faction.AutoAcceptPeace {
						autoAcceptPeace++
					}
					if faction.EnableFriendlyFire {
						enableFriendlyFire++
					}
					pt, err := client.NewPoint(
						"faction",
						map[string]string{
							"host":                        *host,
							"faction_id":                  fmt.Sprint(faction.FactionId),
							"founder_id":                  fmt.Sprint(faction.FounderId),
							"name":                        faction.Name,
							"tag":                         faction.Tag,
							"filter_accept_humans":        toStringBool(faction.AcceptHumans),
							"filter_auto_accept_member":   toStringBool(faction.AutoAcceptMember),
							"filter_auto_accept_peace":    toStringBool(faction.AutoAcceptPeace),
							"filter_enable_friendly_fire": toStringBool(faction.EnableFriendlyFire),
							"filter_npc_only":             toStringBool(faction.NPCOnly),
						},
						map[string]interface{}{
							"npc_only":             npconly,
							"auto_accept_humans":   acceptHumans,
							"auto_accept_member":   autoAcceptMember,
							"auto_accept_peace":    autoAcceptPeace,
							"enable_friendly_fire": enableFriendlyFire,
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

				floatingObjects, err := t.SessionFloatingObjects()
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

				for _, floatingObject := range floatingObjects {
					pt, err := client.NewPoint(
						"floating_object",
						map[string]string{
							"host":         *host,
							"entity_id":    fmt.Sprint(floatingObject.EntityId),
							"display_name": floatingObject.DisplayName,
							"kind":         floatingObject.Kind,
						},
						map[string]interface{}{
							"distance_to_player": floatingObject.DistanceToPlayer,
							"linear_speed":       floatingObject.LinearSpeed,
							"mass":               floatingObject.Mass,
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

func toStringBool(val bool) string {
	switch val {
	case true:
		return "yes"
	default:
		return "no"
	}
}
