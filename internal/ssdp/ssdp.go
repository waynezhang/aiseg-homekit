package ssdp

import (
	"time"

	"github.com/waynezhang/aiseg-homekit/internal/log"

	"github.com/huin/goupnp"
)

type SSDPDevice struct {
	Name     string
	Model    string
	Hostname string
}

func Discover() *SSDPDevice {
	serviceId := "urn:panasonic-com:service:p60AiSeg2DataService:1"

	for {
		devices, err := goupnp.DiscoverDevices(serviceId)
		if err != nil {
			log.E("Failed to discover device due to %s", err.Error())
			time.Sleep(1 * time.Second)
			continue
		}

		if len(devices) > 0 {
			root := devices[0].Root
			return &SSDPDevice{
				root.Device.ModelName,
				root.Device.ModelNumber,
				root.URLBase.Hostname(),
			}
		}
	}
}
