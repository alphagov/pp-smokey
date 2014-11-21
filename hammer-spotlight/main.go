package main

// TODO import vegeta
// import "github.com/tsenart/vegeta/lib"
import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// TODO just learn how to do string substition in Go instead of all this map madness
var environments = map[string]string{
	"preview":    "https://spotlight.preview.performance.service.gov.uk/performance",
	"staging":    "https://spotlight.staging.performance.service.gov.uk/performance",
	"production": "https://spotlight.production.performance.service.gov.uk/performance",
}

var stagecraftURLs = map[string]string{
	"preview":    "https://stagecraft.preview.performance.service.gov.uk/public/dashboards",
	"staging":    "https://stagecraft.staging.performance.service.gov.uk/public/dashboards",
	"production": "https://stagecraft.production.performance.service.gov.uk/public/dashboards",
}

type Dashboards struct {
	Items []DashboardSlug
}

type DashboardSlug struct {
	Slug        string
	ModuleSlugs []string
}

type DashboardConfig struct {
	Modules []struct {
		Slug string
	}
}

// type DashboardSlug struct {
// 	Slug Slug `json:"slug"`
// }

func GetDashboardSlugs(URL string) (Dashboards, error) {

	var results Dashboards
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

func getModuleURLs(stageCraftURL string, Slug string) ([]string, error) {
	var results DashboardConfig
	moduleSlugs := []string{}

	fmt.Println("Getting " + stageCraftURL + "?slug=" + Slug)
	resp, err := http.Get(stageCraftURL + "?slug=" + Slug) // https://stagecraft.com/public/dashboards?slug={slug}
	if err != nil {
		return moduleSlugs, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return moduleSlugs, err
	}

	err = json.Unmarshal(body, &results)
	if err != nil {
		return moduleSlugs, err
	}

	// iter moduleSlugs
	for _, module := range results.Modules {
		moduleSlugs = append(moduleSlugs, module.Slug)
	}

	return moduleSlugs, nil
}

func main() {

	var environment string
	flag.StringVar(&environment, "environment", "preview", "Select host environment (preview/staging/production)")
	flag.StringVar(&environment, "e", "preview", "Select host environment (preview/staging/production)")

	flag.Parse()

	envURL := environments[environment]
	stagecraftURL := stagecraftURLs[environment]
	fmt.Println(stagecraftURL)
	fmt.Println(envURL)

	// Grab all dashboard slug/items from stagecraft
	// https://stagecraft.preview.performance.service.gov.uk/public/dashboards
	dashboardURLs, err := GetDashboardSlugs(stagecraftURL)
	if err != nil {
		log.Fatal(err.Error())
	}

	dashboardModulesList := map[string][]string{}
	for _, item := range dashboardURLs.Items {
		urls, err := getModuleURLs(stagecraftURL, item.Slug)
		if err != nil {
			log.Fatalln(err)
		}
		dashboardModulesList[item.Slug] = urls
	}

	fmt.Println(dashboardModulesList)

}
