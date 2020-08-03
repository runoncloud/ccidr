package ccidr

import (
	"fmt"
	_ "github.com/runoncloud/ccidr/pkg/statik"
	"github.com/thoas/go-funk"
	"github.com/tidwall/gjson"
	"sort"
	"strings"
)

const (
	AzureCloud = "AzureCloud"
)

type Azure struct{}

var (
	azureJSON string
)

func (a Azure) ListRegions() []string {
	var regions []string
	json := getAzureJsonString()
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
	json := getAzureJsonString()

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
	json := getAzureJsonString()
	cidrQueryResult := gjson.Get(json, fmt.Sprintf("values.#(name==%s).properties.addressPrefixes", service))
	for _, addressPrefix := range cidrQueryResult.Array() {
		addressPrefixes = append(addressPrefixes, addressPrefix.String())
	}
	sort.Strings(addressPrefixes)
	return addressPrefixes
}

func (a Azure) ListAddressPrefixesByRegion(region string) []string {
	var addressPrefixes []string
	json := getAzureJsonString()
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
	json := getAzureJsonString()
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

func getAzureJsonString() string {
	if azureJSON == "" {
		azureJSON = GetJsonString("/azure.json")
	}
	return azureJSON
}
