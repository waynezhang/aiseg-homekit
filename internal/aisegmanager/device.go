package aisegmanager

import (
	"fmt"

	"github.com/waynezhang/aiseg-homekit/internal/log"
)

type Device struct {
	NodeId   string
	Name     string
	IsOn     bool
	Type     DeviceType
	eoj      string
	nodeType string
	link     string
	token    string
}

type DeviceType int32

const (
	DeviceTypeLight DeviceType = 1 << iota
	DeviceTypeFloorHeating
)

func (am *AiSEGManager) ToggleDevice(nodeId string) error {
	am.mutext.Lock()
	defer am.mutext.Unlock()

	ds := am.findDevice(nodeId)
	if len(ds) == 0 {
		return fmt.Errorf("Device %s not found", nodeId)
	}
	for _, d := range ds {
		log.D("Toggling device %s", nodeId)
		err := am.turnDevice(d, !d.IsOn)
		if err != nil {
			return err
		}
	}
	return nil
}

func (am *AiSEGManager) TurnDevice(nodeId string, on bool) error {
	am.mutext.Lock()
	defer am.mutext.Unlock()

	ds := am.findDevice(nodeId)
	if len(ds) == 0 {
		return fmt.Errorf("Device %s not found", nodeId)
	}

	for _, d := range ds {

		log.D("Turn device %s %t", nodeId, on)
		err := am.turnDevice(d, on)
		if err != nil {
			return err
		}
	}
	return nil
}

func (am *AiSEGManager) TurnAllDevices(deviceType DeviceType, on bool) error {
	am.mutext.Lock()
	defer am.mutext.Unlock()

	for _, d := range am.Devices {
		if d.Type == deviceType {
			if err := am.turnDevice(d, on); err != nil {
				return err
			}
		}
	}
	return nil
}

func (am *AiSEGManager) findDevice(nodeId string) []*Device {
	ds := []*Device{}
	for _, d := range am.Devices {
		if d.NodeId == nodeId {
			ds = append(ds, d)
		}
	}

	return ds
}

func (am *AiSEGManager) turnDevice(d *Device, on bool) error {
	if d.Type == DeviceTypeLight {
		return am.turnLight(d, on)
	} else if d.Type == DeviceTypeFloorHeating {
		return am.turnFloorHeating(d, on)
	}
	return fmt.Errorf("Unsupported device %d", d.Type)
}

func (am *AiSEGManager) findDevices(deviceType DeviceType) []*Device {
	deviceLink := am.getDevicePageLink()
	if deviceLink == nil {
		log.E("Failed to find device page link")
		return []*Device{}
	}

	devices := []*Device{}
	panels := am.getPanelLinks(*deviceLink)
	for _, panel := range panels {
		if panel.deviceType&deviceType != 0 {
			devices = append(devices, am.parseDevices(panel)...)
		}
	}

	return devices
}
