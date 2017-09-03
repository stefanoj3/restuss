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

// Client expose the methods callable on Nessus Api
type Client interface {
	GetScanTemplates() ([]*ScanTemplate, error)
	LaunchScan(scanId int64) error
	StopScan(scanId int64) error
	CreateScan(scan *Scan) (*PersistedScan, error)
	GetScans(lastModificationDate int64) ([]*PersistedScan, error)
	GetScanByID(id int64) (*ScanDetail, error)
	GetPluginByID(id int64) (*Plugin, error)
}

type NessusClient struct {
	auth       AuthProvider
	url        string
	httpClient *http.Client
}

// NewClient returns a new NessusClient
func NewClient(auth AuthProvider, url string, allowInsecureConnection bool) (*NessusClient, error) {
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

	return &NessusClient{auth: auth, url: url, httpClient: c}, nil
}

// GetScanTemplates retrieves Scan Templates
func (c *NessusClient) GetScanTemplates() ([]*ScanTemplate, error) {
	req, err := http.NewRequest(http.MethodGet, c.url+"/editor/scan/templates", nil)

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

// LaunchScan launch spe scan with the specified scanID
func (c *NessusClient) LaunchScan(scanID int64) error {
	path := "/scans/" + strconv.FormatInt(scanID, 10) + "/launch"
	req, err := http.NewRequest(http.MethodPost, c.url+path, nil)

	if err != nil {
		return errors.New("Unable to create request object: " + err.Error())
	}

	err = c.performCallAndReadResponse(req, nil)

	if err != nil {
		return err
	}

	return nil
}

// StopScan stops the scan with the given scanID
func (c *NessusClient) StopScan(scanID int64) error {
	path := "/scans/" + strconv.FormatInt(scanID, 10) + "/stop"
	req, err := http.NewRequest(http.MethodPost, c.url+path, nil)

	if err != nil {
		return errors.New("Unable to create request object: " + err.Error())
	}

	err = c.performCallAndReadResponse(req, nil)
	if err != nil {
		return err
	}

	return nil
}

// CreateScan creates a scan
func (c *NessusClient) CreateScan(scan *Scan) (*PersistedScan, error) {
	jsonBody, err := json.Marshal(scan)
	if err != nil {
		return nil, errors.New("Unable to marshall request body" + err.Error())
	}

	req, err := http.NewRequest(http.MethodPost, c.url+"/scans", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, errors.New("Unable to create request object: " + err.Error())
	}

	req.Header.Set("Content-Type", "application/json")

	var result struct {
		Scan PersistedScan `json:"scan"`
	}

	err = c.performCallAndReadResponse(req, &result)
	if err != nil {
		return nil, err
	}

	return &result.Scan, nil
}

// GetScans get a list of scan matching the provided lastModificationDate (check Nessus documentation)
func (c *NessusClient) GetScans(lastModificationDate int64) ([]*PersistedScan, error) {
	req, err := http.NewRequest(http.MethodGet, c.url+"/scans", nil)
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

// GetScanByID retrieve a scan by ID
func (c *NessusClient) GetScanByID(ID int64) (*ScanDetail, error) {
	path := fmt.Sprintf("/scans/%d", ID)

	req, err := http.NewRequest(http.MethodGet, c.url+path, nil)
	if err != nil {
		return nil, errors.New("Unable to create request object: " + err.Error())
	}

	scanDetail := &ScanDetail{}

	err = c.performCallAndReadResponse(req, &scanDetail)
	if err != nil {
		return nil, err
	}

	scanDetail.ID = ID

	return scanDetail, nil
}

// GetPluginByID retrieves a plugin by ID
func (c *NessusClient) GetPluginByID(ID int64) (*Plugin, error) {
	path := fmt.Sprintf("/plugins/plugin/%d", ID)

	req, err := http.NewRequest(http.MethodGet, c.url+path, nil)
	if err != nil {
		return nil, errors.New("Unable to create request object: " + err.Error())
	}

	p := &Plugin{}

	err = c.performCallAndReadResponse(req, p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (c *NessusClient) performCallAndReadResponse(req *http.Request, data interface{}) error {
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
