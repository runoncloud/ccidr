package ccidr

import (
	"fmt"
	_ "github.com/runoncloud/ccidr/pkg/statik"
	"github.com/thoas/go-funk"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

const (
	AzureCloud   = "AzureCloud"
	AzureJsonUrl = "https://download.microsoft.com/download/7/1/D/71D86715-5596-4529-9B13-DA13A5DE5B63/ServiceTags_Public_%s.json"
)

type Azure struct {
	isRemote bool
}

var (
	azureJSON string
)

func (a Azure) ListRegions() []string {
	var regions []string
	json := getAzureJsonString(a.isRemote)
	queryResult := gjson.Get(json, "values.#.properties.region")
	for _, region := range queryResult.Array() {
		regions = append(regions, region.String())
	}
	sort.Strings(regions)
	return funk.UniqString(regions)
}

func (a Azure) ListServices() []string {
	return a.ListServicesByRegion("")
}

func (a Azure) ListServicesByRegion(region string) []string {
	var services []string
	json := getAzureJsonString(a.isRemote)

	queryResult := gjson.Get(json,
		fmt.Sprintf("values.#(properties.region==%s)#.name", region))
	for _, service := range queryResult.Array() {
		if service.String() != AzureCloud {
			services = append(services, strings.Split(service.String(), ".")[0])
		}
	}
	sort.Strings(services)
	return services
}

func (a Azure) ListAddressPrefixes() []string {
	return a.ListAddressPrefixesByService(AzureCloud)
}

func (a Azure) ListAddressPrefixesByService(service string) []string {
	var addressPrefixes []string
	json := getAzureJsonString(a.isRemote)
	cidrQueryResult := gjson.Get(json, fmt.Sprintf("values.#(name==%s).properties.addressPrefixes", service))
	for _, addressPrefix := range cidrQueryResult.Array() {
		addressPrefixes = append(addressPrefixes, addressPrefix.String())
	}
	sort.Strings(addressPrefixes)
	return addressPrefixes
}

func (a Azure) ListAddressPrefixesByRegion(region string) []string {
	var addressPrefixes []string
	json := getAzureJsonString(a.isRemote)
	regionStringQueryResult := gjson.Get(json, fmt.Sprintf("values.#(properties.region==%s)#.properties.addressPrefixes", region))
	for _, addressPrefixList := range regionStringQueryResult.Array() {
		for _, addressPrefix := range addressPrefixList.Array() {
			addressPrefixes = append(addressPrefixes, addressPrefix.String())
		}
	}
	sort.Strings(addressPrefixes)
	return addressPrefixes
}

func (a Azure) ListAddressPrefixesByServiceAndRegion(service string, region string) []string {
	var addressPrefixes []string
	json := getAzureJsonString(a.isRemote)
	regionStringQueryResult := gjson.Get(json, fmt.Sprintf("values.#(properties.region==%s).name", region))
	regionString := strings.Split(regionStringQueryResult.String(), ".")[1]
	cidrQueryResult := gjson.Get(json, fmt.Sprintf("values.#(name==%s.%s).properties.addressPrefixes",
		service, regionString))
	for _, addressPrefix := range cidrQueryResult.Array() {
		addressPrefixes = append(addressPrefixes, addressPrefix.String())
	}
	sort.Strings(addressPrefixes)
	return addressPrefixes
}

func getAzureJsonString(isRemote bool) string {
	if azureJSON == "" {
		if isRemote {
			awsJSON = getAzureRemoteJsonString()
		} else {
			azureJSON = GetJsonString("/azure.json")
		}
	}
	return azureJSON
}

func getAzureRemoteJsonString() string {
	lastMonday := getLastMondayDate()

	if azureJSON == "" {
		client := http.Client{Timeout: time.Second * 10}

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(AzureJsonUrl, lastMonday.Format("20060102")), nil)
		if err != nil {
			log.Fatal(err)
		}

		res, getErr := client.Do(req)
		if getErr != nil {
			log.Fatal(getErr)
		}

		if res.Body != nil {
			defer res.Body.Close()
		}

		body, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			log.Fatal(readErr)
		}
		azureJSON = string(body)
	}
	return azureJSON
}

func getLastMondayDate() time.Time {
	date := time.Now()
	for date.Weekday() != time.Monday {
		date = date.AddDate(0, 0, -1)
	}
	return date
}
