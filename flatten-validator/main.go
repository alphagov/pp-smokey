package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type DashboardConfigs struct {
	Items []DashboardConfigSummary
}

type DashboardConfigSummary struct {
	Slug  string
	Title string
}

type DashboardConfig struct {
	DashboardConfigSummary
	// TODO(mattrco): Handle module-type: tab and decode all tabs.
	Modules []Module
}

type Module struct {
	Slug       string
	Datasource DataSource `json:"data-source"`
}

type DataSource struct {
	DataGroup   string          `json:"data-group"`
	DataType    string          `json:"data-type"`
	QueryParams QueryParameters `json:"query-params"`
}

type QueryParameters struct {
	SortBy  string   `json:"sort_by"`
	Collect []string `json:"collect"`
	// TODO(mattrco): GroupBy may be an array of strings or a string.
	GroupBy  string   `json:"group_by"`
	FilterBy []string `json:"filter_by"`
	Limit    int      `json:"limit"`
}

func FetchDashboardConfigs(URL string) (DashboardConfigs, error) {

	var results DashboardConfigs
	resp, err := http.Get(URL)
	if err != nil {
		return results, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return results, err
	}

	err = json.Unmarshal(body, &results)
	return results, err
}

func FetchDashboardConfig(URL string) (DashboardConfig, error) {

	var dashboard DashboardConfig
	resp, err := http.Get(URL)
	if err != nil {
		return dashboard, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return dashboard, err
	}

	err = json.Unmarshal(body, &dashboard)
	return dashboard, err
}

func ConstructModuleURL(module Module) (string, error) {
	var err error

	baseURL := fmt.Sprintf(
		"https://www.performance.service.gov.uk/data/%s/%s?",
		module.Datasource.DataGroup,
		module.Datasource.DataType,
	)
	moduleParams := module.Datasource.QueryParams
	params := url.Values{}
	if moduleParams.SortBy != "" {
		params.Add("sort_by", module.Datasource.QueryParams.SortBy)
	}
	for _, collectParam := range moduleParams.Collect {
		params.Add("collect", collectParam)
	}
	if moduleParams.GroupBy != "" {
		params.Add("group_by", module.Datasource.QueryParams.GroupBy)
	}
	for _, filter := range moduleParams.FilterBy {
		params.Add("filter_by", filter)
	}
	if moduleParams.Limit != 0 {
		params.Add("limit", strconv.Itoa(moduleParams.Limit))
	}
	if len(params) == 0 {
		err = errors.New("Empty query string encoded for module")
	}

	moduleURL := baseURL + params.Encode()
	return moduleURL, err
}

func main() {
	baseURL := "https://stagecraft.production.performance.service.gov.uk/public/dashboards"
	dashConfs, err := FetchDashboardConfigs(baseURL)
	if err != nil {
		log.Fatal(err.Error())
	}

	modules := 0
	skipped := 0
	errored := 0
	flattenErrors := 0

	for _, dashConf := range dashConfs.Items {
		dash, err := FetchDashboardConfig(fmt.Sprintf("%s?slug=%s", baseURL, dashConf.Slug))
		if err != nil {
			log.Print(err.Error())
		}
		// For each module, print module URL
		for _, module := range dash.Modules {
			moduleURL, err := ConstructModuleURL(module)
			if err != nil {
				log.Printf("Skipping module because of: %s (%s)", err.Error(), dash.Slug)
				skipped += 1
			} else {
				resp, err := http.Get(moduleURL)
				if err != nil {
					log.Printf("Error fetching module: %s (%s)", err.Error(), moduleURL)
					errored += 1
				} else if resp.StatusCode != http.StatusOK {
					log.Printf("Got status %d for module %s", resp.StatusCode, moduleURL)
					errored += 1
				}
				resp.Body.Close()
				time.Sleep(time.Millisecond * 250)
				resp, err = http.Get(moduleURL + "&flatten=true")
				if err != nil {
					log.Printf("Error fetching module: %s (%s)", err.Error(), moduleURL)
					flattenErrors += 1
				} else if resp.StatusCode != http.StatusOK {
					log.Printf("Got status %d for module %s", resp.StatusCode, moduleURL)
					flattenErrors += 1
				}
				resp.Body.Close()
				time.Sleep(time.Millisecond * 250)
			}
			modules += 1
		}
	}
	log.Printf("Modules: %d, Skipped: %d, Errors: %d", modules, skipped, errored)
}
