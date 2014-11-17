package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"time"
)

var InvalidFlattenError error

func init() {
	InvalidFlattenError = errors.New("URL has no group_by params so cannot be flattened")
}

const baseURL string = "https://stagecraft.staging.performance.service.gov.uk/public/dashboards"

type ResponseTimes []time.Duration

func (d ResponseTimes) Len() int           { return len(d) }
func (d ResponseTimes) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d ResponseTimes) Less(i, j int) bool { return d[i] < d[j] }

type DashboardConfigs struct {
	Items []DashboardConfigSummary
}

// DashboardConfigSummary represents a dashboard config as returned in a JSON array.
type DashboardConfigSummary struct {
	Slug  string
	Title string
}

type DashboardConfigResponse struct {
	DashboardConfig DashboardConfig
	Error           error
}

type DashboardConfig struct {
	DashboardConfigSummary
	Modules []Module
}

// Report defines how we capture information about a particular Module's performance.
type Report struct {
	URL      string
	Start    time.Time
	Elapsed  time.Duration
	BodySize int
	Error    error
}

// ModuleReport is a union of the Timing for the old version and the new flatten=true version.
type ModuleReport struct {
	Module  Report
	Flatten Report
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

func constructModuleURL(module Module, flatten bool) (string, error) {

	// TODO(mrc): use net/url to create the URL.
	baseURL := fmt.Sprintf(
		"https://www.staging.performance.service.gov.uk/data/%s/%s",
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

	// Optionally add "flatten=true" if group_by is present.
	if flatten {
		if _, ok := params["group_by"]; ok {
			params.Add("flatten", "true")
		} else {
			return "", InvalidFlattenError
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
		moduleURL = moduleURL + "?" + params.Encode()
	}
	return moduleURL, nil
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

func fetchResponse(moduleURL string) (*http.Response, error) {
	resp, err := http.Get(moduleURL)
	if err != nil {
		return resp, fmt.Errorf("Error fetching module: %s (%s)", err.Error(), moduleURL)
	}
	if resp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("Got status %d for module %s", resp.StatusCode, moduleURL)
	}
	return resp, nil
}

func main() {
	if os.Getenv("GOMAXPROCS") == "" {
		// Use all available cores if not otherwise specified
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	dashConfs, err := FetchDashboardConfigs(baseURL)
	if err != nil {
		log.Fatal(err.Error())
	}

	configs := produceConfigs(dashConfs)
	modules := produceModules(configs)

	// Arbitrary choice of workers to consume the modules
	reports := make([]chan ModuleReport, 4)
	for i, _ := range reports {
		reports[i] = produceReports(modules)
	}

	var errors []error
	var flattenErrors []error

	// Sort response times and sizes and report the worst.
	var responseTimes ResponseTimes
	responseTimesMap := map[time.Duration]string{}
	var responseDiffs ResponseTimes
	responseDiffsMap := map[time.Duration]string{}
	var responseSizes sort.IntSlice
	responseSizesMap := map[int]string{}

	for r := range merge(reports...) {

		moduleReport := r.Module
		flattenReport := r.Flatten

		if moduleReport.Error != nil {
			errors = append(errors, moduleReport.Error)
		} else {
			responseTimes = append(responseTimes, moduleReport.Elapsed)
			responseTimesMap[moduleReport.Elapsed] = moduleReport.URL
			responseSizes = append(responseSizes, moduleReport.BodySize)
			responseSizesMap[moduleReport.BodySize] = moduleReport.URL
		}

		if flattenReport.Error != nil {
			if flattenReport.Error != InvalidFlattenError {
				flattenErrors = append(flattenErrors, flattenReport.Error)
			}
		} else {
			responseTimes = append(responseTimes, flattenReport.Elapsed)
			responseTimesMap[flattenReport.Elapsed] = flattenReport.URL
			responseSizes = append(responseSizes, flattenReport.BodySize)
			responseSizesMap[flattenReport.BodySize] = flattenReport.URL
		}

		if moduleReport.Error == nil && flattenReport.Error == nil {
			responseDiffs = append(responseDiffs, flattenReport.Start.Sub(moduleReport.Start))
			responseDiffsMap[flattenReport.Start.Sub(moduleReport.Start)] = moduleReport.URL
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
	sort.Sort(sort.Reverse(responseTimes))
	if len(responseTimes) > 10 {
		responseTimes = responseTimes[0:10]
	}
	log.Print("Slowest responses...")
	for _, v := range responseTimes {
		log.Printf("%s took %s", responseTimesMap[v], v)
	}

	sort.Sort(sort.Reverse(responseDiffs))
	if len(responseDiffs) > 10 {
		responseDiffs = responseDiffs[0:10]
	}
	log.Print("Worst flatten regressions...")
	for _, v := range responseDiffs {
		log.Printf("%s took %s longer", responseDiffsMap[v], v)
	}

	sort.Sort(sort.Reverse(responseSizes))
	if len(responseSizes) > 10 {
		responseSizes = responseSizes[0:10]
	}
	log.Print("Largest response sizes...")
	for _, v := range responseSizes {
		log.Printf("%s was %d bytes", responseSizesMap[v], v)
	}
}

func produceConfigs(dashConfs DashboardConfigs) chan DashboardConfigResponse {
	out := make(chan DashboardConfigResponse)
	go func() {
		defer close(out)
		for _, dashConf := range dashConfs.Items {
			dash, err := FetchDashboardConfig(fmt.Sprintf("%s?slug=%s", baseURL, dashConf.Slug))
			out <- DashboardConfigResponse{dash, err}
		}
	}()
	return out
}

func produceModules(configs <-chan DashboardConfigResponse) <-chan []Module {
	out := make(chan []Module)
	go func() {
		defer close(out)
		for res := range configs {
			if res.Error != nil {
				log.Print(res.Error.Error())
			} else {
				out <- ListDashboardModules(res.DashboardConfig)
			}
		}
	}()
	return out
}

func newReport(URL string) Report {
	start := time.Now()
	resp, err := fetchResponse(URL)
	defer resp.Body.Close()
	elapsed := time.Since(start)
	report := Report{
		URL:     URL,
		Start:   start,
		Elapsed: elapsed,
		Error:   err,
	}

	// Only calculate response size if there is no error.
	if err == nil {
		bytes, _ := ioutil.ReadAll(resp.Body)
		report.BodySize = len(bytes)
	}
	return report
}

func produceReports(modules <-chan []Module) chan ModuleReport {
	out := make(chan ModuleReport)
	go func() {
		defer close(out)

		for dashModules := range modules {
			for _, module := range dashModules {

				var moduleReport Report
				moduleURL, err := constructModuleURL(module, false)
				if err != nil {
					moduleReport = Report{Error: err}
				} else {
					moduleReport = newReport(moduleURL)
				}

				var flattenReport Report
				flattenURL, err := constructModuleURL(module, true)
				if err != nil {
					flattenReport = Report{Error: err}
				} else {
					flattenReport = newReport(flattenURL)
				}

				out <- ModuleReport{
					moduleReport,
					flattenReport}
			}
		}
	}()
	return out
}

func merge(reports ...chan ModuleReport) <-chan ModuleReport {
	var wg sync.WaitGroup
	out := make(chan ModuleReport)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan ModuleReport) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(reports))
	for _, c := range reports {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
