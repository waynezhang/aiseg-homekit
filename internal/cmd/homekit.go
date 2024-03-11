package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/brutella/hap"
	"github.com/brutella/hap/accessory"
	"github.com/spf13/cobra"
	"github.com/waynezhang/aiseg-hb/internal/aisegmanager"
	"github.com/waynezhang/aiseg-hb/internal/log"
)

const (
	refreshInterval = 15 * time.Minute // 15 mins
)

var HKServeCmd = func() *cobra.Command {
	fn := func(cmd *cobra.Command, args []string) {
		serve()
	}

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start HomeKit server. AISEG_USER and AIUSEG_PASSWORD are required as environment variables. PIN code (00102003 by default) can be configured by AISEG_PIN.",
		Run:   fn,
	}

	return cmd
}()

func serve() {
	mgr, bridge, accessories := discoverAccessories()
	if len(accessories) == 0 {
		log.E("No accessories found")
		os.Exit(1)
	}
	log.D("Found %d accessories", len(accessories))

	store := hap.NewFsStore("./db")

	server, err := hap.NewServer(store, bridge.A, accessories...)
	if err != nil {
		log.F("Failed to create server due to %s", err.Error())
	}

	pin := os.Getenv("AISEG_PIN")
	if pin == "" {
		pin = "00102003"
	}
	server.Pin = pin

	startRefresh(mgr)

	fmt.Printf("Starting server with PIN code %s...\n", server.Pin)
	if err = server.ListenAndServe(context.Background()); err != nil {
		log.F("Failed to start server due to %s", err.Error())
	}
}

func discoverAccessories() (*aisegmanager.AiSEGManager, *accessory.Bridge, []*accessory.A) {
	log.D("Discovering devices")

	accessories := []*accessory.A{}
	mgr := aisegmanager.DiscoverNewAiSEGManager()
	for idx, d := range mgr.Devices {
		switch d.Type {
		case aisegmanager.DeviceTypeLight:
			a := accessory.NewLightbulb(accessory.Info{
				Name: d.Name,
			})
			a.Lightbulb.On.OnValueRemoteUpdate(func(on bool) {
				_ = mgr.TurnDevice(d.NodeId, on)
			})
			a.Id = uint64(10000 + idx)
			a.Lightbulb.On.SetValue(d.IsOn)
			accessories = append(accessories, a.A)
		case aisegmanager.DeviceTypeFloorHeating:
			a := accessory.NewFan(accessory.Info{
				Name: d.Name,
			})
			a.Fan.On.OnSetRemoteValue(func(on bool) error {
				return mgr.TurnDevice(d.NodeId, on)
			})
			a.Id = uint64(10000 + idx)
			a.Fan.On.SetValue(d.IsOn)
			accessories = append(accessories, a.A)
		}
	}

	for _, a := range accessories {
		log.D("Created accessory %12d %s", a.Id, a.Name())
	}

	bridge := accessory.NewBridge(accessory.Info{
		Name:         mgr.Name,
		Model:        mgr.Model,
		Manufacturer: "Panasonic",
	})
	bridge.Id = 1

	return mgr, bridge, accessories
}

func startRefresh(mgr *aisegmanager.AiSEGManager) {
	go func() {
		time.Sleep(refreshInterval)
		log.D("Refreshing tokens")
		mgr.Refresh()
	}()
}
