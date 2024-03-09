package aisegmanager

import (
	"encoding/json"
	"fmt"

	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/waynezhang/aiseg-hb/internal/log"
)

func (am *AiSEGManager) parseFloorHeating(panel panel) []Device {
	path := fmt.Sprintf("/page/devices/device/%s", panel.link)
	doc, err := am.client.Document(path)
	if err != nil {
		log.E("Failed to parse floor heating link %s due to %s", path, err.Error())
		return []Device{}
	}

	link := strings.Split(panel.link, "?")[0]

	devices := []Device{}
	doc.Find(".main .panel").Each(func(_ int, s *goquery.Selection) {
		nodeId := s.AttrOr("nodeid", "#")
		eoj := s.AttrOr("eoj", "#")
		nodeType := s.AttrOr("type", "#")
		title := s.Find(".kiki_title").Text()
		stateNode := s.Find(".kiki_state").Children().First()
		token := s.Find(".kiki_button .control").AttrOr("token", "")
		isOn := stateNode.HasClass("on")
		devices = append(devices, Device{
			NodeId:   nodeId,
			Name:     title,
			IsOn:     isOn,
			Type:     DeviceTypeFloorHeating,
			eoj:      eoj,
			nodeType: nodeType,
			token:    token,
			link:     link,
		})
	})
	return devices
}

func (am *AiSEGManager) turnFloorHeating(d *Device, on bool) error {
	// http://hostname/action/devices/device/xxx/change
	log.D("Toggling %s at %s", d.NodeId, d.link)

	path := fmt.Sprintf("/action/devices/device/%s/change", d.link)
	// From AiSEG
	// if (a.hasClass("panel")) {
	//      var s = a.attribute("nodeId"),
	//          n = a.attribute("eoj"),
	//          i = a.attribute("type");
	//      return {
	//          nodeId: s,
	//          eoj: n,
	//          type: i,
	//          state: a.attribute("state")
	//      }

	// On/Off is the oppsite of Light device
	onoff := "0x31"
	if !on {
		onoff = "0x30"
	}
	values := map[string]interface{}{
		"token":  d.token,
		"nodeId": d.NodeId,
		"eoj":    d.eoj,
		"type":   d.nodeType,
		"state":  onoff,
	}
	data, err := json.Marshal(values)
	if err != nil {
		return err
	}
	resp, err := am.client.PostForm(path, string(data))
	log.D("Response %s", resp)

	return err
}
