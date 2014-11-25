package main

// TODO import vegeta
// import "github.com/tsenart/vegeta/lib"
import (
	"encoding/json"
	"flag"
	"fmt"
	vegeta "github.com/tsenart/vegeta/lib"
	"io/ioutil"
	"log"
	"net/http"
	"time"
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

	// Now start attacking module urls
	rate := uint64(100) // per second
	duration := 4 * time.Second

	var targets = []*vegeta.Target{}
	for dashboardSlug, moduleSlugs := range dashboardModulesList {
		for _, moduleSlug := range moduleSlugs {
			fmt.Println("Priming:", envURL+"/"+dashboardSlug+"/"+moduleSlug)
			targets = append(targets, &vegeta.Target{
				Method: "GET",
				URL:    envURL + "/" + dashboardSlug + "/" + moduleSlug,
			})
		}
	}
	targeter := vegeta.NewStaticTargeter(targets...)

	attacker := vegeta.NewAttacker()

	var results vegeta.Results
	for res := range attacker.Attack(targeter, rate, duration) {
		fmt.Println("attacking", targeter.URL, rate, duration)
		results = append(results, res)
	}

	metrics := vegeta.NewMetrics(results)
	// fmt.Printf("99th percentile: %s\n", metrics.Latencies.P99)
	fmt.Println("Requests\t[total]\t%d\n", metrics.Requests)
	fmt.Println("Duration\t[total, attack, wait]\t%s, %s, %s\n", metrics.Duration+metrics.Wait, metrics.Duration, metrics.Wait)
	fmt.Println("Latencies\t[mean, 50, 95, 99, max]\t%s, %s, %s, %s, %s\n",
		metrics.Latencies.Mean, metrics.Latencies.P50, metrics.Latencies.P95, metrics.Latencies.P99, metrics.Latencies.Max)
	fmt.Println("Bytes In\t[total, mean]\t%d, %.2f\n", metrics.BytesIn.Total, metrics.BytesIn.Mean)
	fmt.Println("Bytes Out\t[total, mean]\t%d, %.2f\n", metrics.BytesOut.Total, metrics.BytesOut.Mean)
	fmt.Println("Success\t[ratio]\t%.2f%%\n", metrics.Success*100)
	fmt.Println("Status Codes\t[code:count]\t")

	for code, count := range metrics.StatusCodes {
		fmt.Println("%s:%d  ", code, count)
	}
	fmt.Println("\nError Set:")
	for err := range metrics.Errors {
		fmt.Println(err)
	}

}
