package aisegmanager

import (
	"encoding/json"
	"fmt"

	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/waynezhang/aiseg-hb/internal/log"
)

func (am *AiSEGManager) parseLights(panel panel) []Device {
	return am.parseLightsPage([]Device{}, panel, 1, true)
}

func (am *AiSEGManager) parseLightsPage(last []Device, panel panel, page int, hasNext bool) []Device {
	if !hasNext {
		return last
	}

	_ = strings.Split(panel.link, "?")[0]
	path := fmt.Sprintf("/page/devices/device/%s&individual_page=%d", panel.link, page)
	log.D("Parsing lights at %s", path)
	doc, err := am.client.Document(path)
	if err != nil {
		log.E("Failed to parse light link %s due to %s", path, err.Error())
		return last
	}

	token := doc.Find(".setting_value").Text()
	link := strings.Split(panel.link, "?")[0]

	selection := doc.Find("#main .panel").Each(func(_ int, s *goquery.Selection) {
		nodeId := s.AttrOr("nodeid", "#")
		eoj := s.AttrOr("eoj", "#")
		nodeType := s.AttrOr("type", "#")
		title := s.Find(".lighting_title").Text()
		stateNode := s.Find(".lighting_state")
		isOn := stateNode.HasClass("on")
		last = append(last, Device{
			NodeId:   nodeId,
			Name:     title,
			IsOn:     isOn,
			Type:     DeviceTypeLight,
			eoj:      eoj,
			nodeType: nodeType,
			token:    token,
			link:     link,
		})
	})

	// TODO 8 per page
	hasNext = len(selection.Nodes) == 8
	return am.parseLightsPage(last, panel, page+1, hasNext)
}

func (am *AiSEGManager) turnLight(d *Device, on bool) error {
	// http://hostname/action/devices/device/xxx/change
	log.D("Toggling %s at %s", d.NodeId, d.link)

	path := fmt.Sprintf("/action/devices/device/%s/change", d.link)
	// From AiSEG
	// t.token = document.getElementsByClassName("setting_value")[0].innerHTML,
	// t.nodeId = q(e.parentElement).attribute("nodeId"),
	// t.eoj = q(e.parentElement).attribute("eoj"),
	// t.type = q(e.parentElement).attribute("type");
	// var a = "on" === q(e.parentElement).attribute("state") ? "0x31" : "0x30";
	// t.device = {
	//     onoff: a,
	//     modulate: "-"
	// },
	onoff := "0x30"
	if !on {
		onoff = "0x31"
	}
	values := map[string]interface{}{
		"token":  d.token,
		"nodeId": d.NodeId,
		"eoj":    d.eoj,
		"type":   d.nodeType,
		"device": map[string]interface{}{
			"onoff":    onoff,
			"modulate": "-",
		},
	}
	data, err := json.Marshal(values)
	if err != nil {
		return err
	}
	resp, err := am.client.PostForm(path, string(data))
	log.D("Response %s", resp)

	return err
}
