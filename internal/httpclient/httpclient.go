package httpclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/icholy/digest"
	"github.com/waynezhang/aiseg-hb/internal/log"
)

type HttpClient struct {
	hostname string
	client   *http.Client
}

func Client(hostname string, username string, password string) *HttpClient {
	return &HttpClient{
		hostname,
		&http.Client{
			Transport: &digest.Transport{
				Username: username,
				Password: password,
			},
		},
	}
}

func (hc *HttpClient) Get(path string) (string, error) {
	url := hc.url(path)

	log.D("Request GET: %s", url)

	resp, err := hc.client.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// {"result":"0","acceptId":"108946","errorInfo":"-"}
type deviceChangeResponse struct {
	Result    string `json:"result"`
	ErrorInfo string `json:"errorInfo"`
}

func (hc *HttpClient) RequestChange(path string, data string) error {
	body, err := hc.postForm(path, data)
	if err != nil {
		return err
	}
	log.D("Response %s", body)

	var resp deviceChangeResponse
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		return err
	}
	if resp.Result != "0" {
		return fmt.Errorf("Request failure %s", resp.ErrorInfo)
	}

	return nil
}

func (hc *HttpClient) postForm(path string, data string) (string, error) {
	url := hc.url(path)

	reqBody := fmt.Sprintf("data=%s", data)

	log.D("Request Post: %s Body: %s", url, data)

	resp, err := hc.client.Post(url, "application/x-www-form-urlencoded", strings.NewReader(reqBody))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (hc *HttpClient) Document(path string) (*goquery.Document, error) {
	html, err := hc.Get(path)
	if err != nil {
		return nil, err
	}

	return goquery.NewDocumentFromReader(strings.NewReader(html))
}

func (hc *HttpClient) url(path string) string {
	return fmt.Sprintf("http://%s%s", hc.hostname, path)
}
