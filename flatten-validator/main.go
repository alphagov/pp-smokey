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
	Modules []Module
}

type Module struct {
	Slug       string
	Datasource DataSource `json:"data-source"`
	ModuleType string     `json:"module-type"`
	Tabs       []Module   `json:"tabs"`
}

type DataSource struct {
	DataGroup   string          `json:"data-group"`
	DataType    string          `json:"data-type"`
	QueryParams QueryParameters `json:"query-params"`
}

type QueryParameters struct {
	SortBy  string   `json:"sort_by"`
	Collect []string `json:"collect"`
	// GroupBy may be an array of strings or a string.
	GroupBy  interface{} `json:"group_by"`
	FilterBy []string    `json:"filter_by"`
	Limit    int         `json:"limit"`
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

	// TODO(mrc): use net/url to create the URL.
	baseURL := fmt.Sprintf(
		"https://www.staging.performance.service.gov.uk/data/%s/%s?",
		module.Datasource.DataGroup,
		module.Datasource.DataType,
	)
	moduleParams := module.Datasource.QueryParams
	params := url.Values{}
	if moduleParams.SortBy != "" {
		params.Add("sort_by", moduleParams.SortBy)
	}
	for _, collectParam := range moduleParams.Collect {
		params.Add("collect", collectParam)
	}

	// Since GroupBy may be a string or a string array, switch on the type.
	switch moduleParams.GroupBy.(type) {
	case string:
		groupBy := moduleParams.GroupBy.(string)
		if groupBy != "" {
			params.Add("group_by", groupBy)
		}
	case []interface{}:
		groupBy := moduleParams.GroupBy.([]interface{})
		for _, param := range groupBy {
			stringParam, ok := param.(string)
			if ok {
				params.Add("group_by", stringParam)
			} else {
				return "", fmt.Errorf("Couldn't stringify param: %v", param)
			}
		}
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

	moduleURL := baseURL
	if len(params) > 0 {
		moduleURL = moduleURL + params.Encode()
	}
	return moduleURL, err
}

// Returns an array of modules, including any tabbed modules, from the dashboard config.
func ListDashboardModules(dash DashboardConfig) []Module {
	var modules []Module
	for _, module := range dash.Modules {
		if module.ModuleType == "tab" {
			for _, tab := range module.Tabs {
				modules = append(modules, tab)
			}
		} else {
			modules = append(modules, module)
		}
	}
	return modules
}

func validateResponse(moduleURL string) error {
	resp, err := http.Get(moduleURL)
	defer resp.Body.Close()
	if err != nil {
		return fmt.Errorf("Error fetching module: %s (%s)", err.Error(), moduleURL)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Got status %d for module %s", resp.StatusCode, moduleURL)
	}
	return nil
}

func main() {
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
