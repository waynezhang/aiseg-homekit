package aisegmanager

import (
	"os"

	"github.com/waynezhang/aiseg-hb/internal/ssdp"

	"github.com/PuerkitoBio/goquery"
	"github.com/waynezhang/aiseg-hb/internal/httpclient"
	"github.com/waynezhang/aiseg-hb/internal/log"
)

type AiSEGManager struct {
	Name    string
	Model   string
	client  *httpclient.HttpClient
	Devices []Device
}

type panel struct {
	title      string
	deviceType DeviceType
	link       string
}

// This API will block if no devices found
func DiscoverNewAiSEGManager() *AiSEGManager {
	device := ssdp.Discover()
	log.D("Found device %s (%s) at %s", device.Name, device.Model, device.Hostname)

	user := os.Getenv("AISEG_USER")
	if user == "" {
		log.F("AISEG_USER is not set")
	}
	password := os.Getenv("AISEG_PASSWORD")
	if password == "" {
		log.F("AISEG_PASSWORD is not set")
	}
	mgr := newManager(device.Hostname, user, password)
	mgr.Name = device.Name
	mgr.Model = device.Model
	mgr.Devices = mgr.findDevices(DeviceTypeFloorHeating | DeviceTypeLight)
	return mgr
}

func newManager(hostname string, username string, password string) *AiSEGManager {
	return &AiSEGManager{
		client:  httpclient.Client(hostname, username, password),
		Devices: []Device{},
	}
}

func (am *AiSEGManager) getDevicePageLink() *string {
	doc, err := am.client.Document("/")
	if err != nil {
		log.E("Failed to parse index page due to %s", err.Error())
		return nil
	}

	linkNode := doc.Find("#fmenu a[href^=\"/page/devices\"]").First()
	if linkNode == nil {
		return nil
	}

	href := linkNode.AttrOr("href", "")
	return &href
}

func (am *AiSEGManager) getPanelLinks(deviceLink string) []panel {
	log.D("Parsing %s", deviceLink)

	doc, err := am.client.Document(deviceLink)
	if err != nil {
		log.E("Failed to parse device page due to  %s", err)
		return []panel{}
	}

	panels := []panel{}
	doc.Find("#main .panel").Each(func(_ int, s *goquery.Selection) {
		title := s.Find(".kiki_title").Text()
		link := s.Find(".kiki_button a").FilterFunction(func(_ int, s *goquery.Selection) bool {
			// TODO figure a better way
			return s.Prev().Text() == "個別"
		}).First()
		if title == "照明" {
			panels = append(panels, panel{
				title,
				DeviceTypeLight,
				link.AttrOr("href", "#"),
			})
		} else if title == "床暖房" {
			panels = append(panels, panel{
				title,
				DeviceTypeFloorHeating,
				link.AttrOr("href", "#"),
			})
		} else {
			log.D("Found unsuppported device type %s", title)
		}
	})

	return panels
}

func (am *AiSEGManager) parseDevices(panel panel) []Device {
	if panel.deviceType == DeviceTypeLight {
		return am.parseLights(panel)
	}
	if panel.deviceType == DeviceTypeFloorHeating {
		return am.parseFloorHeating(panel)
	}
	return []Device{}
}
