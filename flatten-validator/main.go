package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"
)

const baseURL string = "https://stagecraft.staging.performance.service.gov.uk/public/dashboards"

type ResponseTimes []time.Duration

func (d ResponseTimes) Len() int           { return len(d) }
func (d ResponseTimes) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d ResponseTimes) Less(i, j int) bool { return d[i] < d[j] }

type DashboardConfigs struct {
	Items []DashboardConfigSummary
}

// Represents a dashboard config as returned in a JSON array.
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

	var errors []error
	var flattenErrors []error
	// Sort these by key and report the worst.
	var responseTimes ResponseTimes
	responseTimesMap := map[time.Duration]string{}
	var responseDiffs ResponseTimes
	responseDiffsMap := map[time.Duration]string{}

	for _, dashConf := range dashConfs.Items {
		dash, err := FetchDashboardConfig(fmt.Sprintf("%s?slug=%s", baseURL, dashConf.Slug))
		if err != nil {
			log.Print(err.Error())
		}

		dashModules := ListDashboardModules(dash)
		for _, module := range dashModules {
			moduleURL, err := ConstructModuleURL(module)
			if err != nil {
				log.Print(err.Error())
				continue
			}
			t1 := time.Now()
			if err := validateResponse(moduleURL); err != nil {
				errors = append(errors, err)
			}
			firstResp := time.Since(t1)
			t2 := time.Now()
			flattenURL := moduleURL + "&flatten=true"
			if err := validateResponse(flattenURL); err != nil {
				flattenErrors = append(flattenErrors, err)
			}
			secondResp := time.Since(t2)
			responseTimes = append(responseTimes, firstResp)
			responseTimes = append(responseTimes, secondResp)
			responseTimesMap[firstResp] = moduleURL
			responseTimesMap[secondResp] = flattenURL
			responseDiffs = append(responseDiffs, t2.Sub(t1))
			responseDiffsMap[t2.Sub(t1)] = moduleURL
		}
	}
	log.Printf("Errors: %d, Flatten errors: %d", len(errors), len(flattenErrors))
	for _, err := range errors {
		log.Print(err.Error())
	}
	for _, err := range flattenErrors {
		log.Print(err.Error())
	}

	// Sort the response times and diffs, and report the slowest and worst regressions.
	sort.Sort(responseTimes)
	if len(responseTimes) > 10 {
		responseTimes = responseTimes[len(responseTimes)-11:]
	}
	log.Print("Slowest responses...")
	for _, v := range responseTimes {
		log.Printf("%s took %s", responseTimesMap[v], v)
	}

	sort.Sort(responseDiffs)
	if len(responseDiffs) > 10 {
		responseDiffs = responseDiffs[len(responseDiffs)-11:]
	}
	log.Print("Worst flatten regressions...")
	for _, v := range responseDiffs {
		log.Printf("%s took %s longer", responseDiffsMap[v], v)
	}
}
