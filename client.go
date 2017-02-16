package restuss

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type Client interface {
	GetScanTemplates() ([]*ScanTemplate, error)
	LaunchScan(scanId int64) error
	StopScan(scanId int64) error
	CreateScan(scan *Scan) (*PersistedScan, error)
	GetScans(lastModificationDate int64) ([]*PersistedScan, error)
	GetScanByID(id int64) (*ScanDetail, error)
}

type client struct {
	auth       AuthProvider
	url        string
	httpClient *http.Client
}

func (c *client) GetScanTemplates() ([]*ScanTemplate, error) {
	req, err := http.NewRequest("GET", c.url+"/editor/scan/templates", nil)

	if err != nil {
		return nil, errors.New("Unable to create request object: " + err.Error())
	}

	var data struct {
		Templates []*ScanTemplate `json:"templates"`
	}

	err = c.performCallAndReadResponse(req, &data)

	if err != nil {
		return nil, errors.New("Call failed: " + err.Error())
	}

	return data.Templates, nil
}

func (c *client) LaunchScan(scanId int64) error {
	url := c.url + "/scans/" + strconv.FormatInt(scanId, 10) + "/launch"
	req, err := http.NewRequest("POST", url, nil)

	if err != nil {
		return errors.New("Unable to create request object: " + err.Error())
	}

	err = c.performCallAndReadResponse(req, nil)

	if err != nil {
		return err
	}

	return nil
}

func (c *client) StopScan(scanId int64) error {
	url := c.url + "/scans/" + strconv.FormatInt(scanId, 10) + "/stop"
	req, err := http.NewRequest("POST", url, nil)

	if err != nil {
		return errors.New("Unable to create request object: " + err.Error())
	}

	err = c.performCallAndReadResponse(req, nil)

	if err != nil {
		return err
	}

	return nil
}

func (c *client) CreateScan(scan *Scan) (*PersistedScan, error) {
	jsonBody, err := json.Marshal(scan)

	if err != nil {
		return nil, errors.New("Unable to marshall request body" + err.Error())
	}

	req, err := http.NewRequest("POST", c.url+"/scans", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return nil, errors.New("Unable to create request object: " + err.Error())
	}

	var result struct {
		Scan PersistedScan `json:"scan"`
	}

	err = c.performCallAndReadResponse(req, &result)

	if err != nil {
		return nil, err
	}

	return &result.Scan, nil
}

func (c *client) GetScans(lastModificationDate int64) ([]*PersistedScan, error) {
	req, err := http.NewRequest("GET", c.url+"/scans", nil)

	if err != nil {
		return nil, errors.New("Unable to create request object: " + err.Error())
	}

	if lastModificationDate > 0 {
		q := req.URL.Query()
		q.Add("last_modification_date", strconv.FormatInt(lastModificationDate, 10))
		req.URL.RawQuery = q.Encode()
	}

	var data struct {
		Scans []*PersistedScan `json:"scans"`
	}

	err = c.performCallAndReadResponse(req, &data)

	if err != nil {
		return nil, err
	}

	return data.Scans, nil
}

func (c *client) GetScanByID(id int64) (*ScanDetail, error) {
	path := fmt.Sprintf("/scans/%d", id)

	req, err := http.NewRequest("GET", c.url+path, nil)

	if err != nil {
		return nil, errors.New("Unable to create request object: " + err.Error())
	}

	scanDetail := &ScanDetail{}

	err = c.performCallAndReadResponse(req, &scanDetail)

	if err != nil {
		return nil, err
	}

	scanDetail.ID = id

	return scanDetail, nil
}

func (c *client) performCallAndReadResponse(req *http.Request, data interface{}) error {
	c.auth.AddAuthHeaders(req)

	res, err := c.httpClient.Do(req)

	if err != nil {
		return errors.New("Failed call: " + err.Error())
	}

	defer res.Body.Close()

	if res.StatusCode > 299 {
		return errors.New("Call failed, status code: " + strconv.Itoa(res.StatusCode))
	}

	if data != nil {
		d := json.NewDecoder(res.Body)
		err = d.Decode(&data)

		if err != nil {
			return errors.New("Failed to read the response: " + err.Error())
		}
	}

	return nil
}

func NewClient(auth AuthProvider, url string, allowInsecureConnection bool) (Client, error) {
	var c *http.Client

	if allowInsecureConnection {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		c = &http.Client{Transport: tr}
	} else {
		c = &http.Client{}
	}

	err := auth.Prepare(url, c)

	if err != nil {
		return nil, errors.New("Failed to prepare auth provider: " + err.Error())
	}

	return &client{auth: auth, url: url, httpClient: c}, nil
}