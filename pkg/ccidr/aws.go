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
	Amazon = "AMAZON"
)

type AWS struct{}

var (
	awsJSON string
)

func (a AWS) ListRegions() []string {
	var regions []string
	json := getAWSJsonString()
	queryResult := gjson.Get(json, "prefixes.#.region")
	for _, region := range queryResult.Array() {
		regions = append(regions, region.String())
	}
	sort.Strings(regions)
	return funk.UniqString(regions)
}

func (a AWS) ListServices() []string {
	var services []string
	json := getAWSJsonString()

	queryResult := gjson.Get(json, "prefixes.#.service")
	for _, service := range queryResult.Array() {
		services = append(services, service.String())
	}
	sort.Strings(services)
	return funk.UniqString(services)
}

func (a AWS) ListServicesByRegion(region string) []string {
	var services []string
	json := getAWSJsonString()

	queryResult := gjson.Get(json,
		fmt.Sprintf("prefixes.#(region==%s)#.service", region))
	for _, service := range queryResult.Array() {
		services = append(services, service.String())
	}
	sort.Strings(services)
	return funk.UniqString(services)
}

func (a AWS) ListAddressPrefixes() []string {
	return a.ListAddressPrefixesByService(Amazon)
}

func (a AWS) ListAddressPrefixesByRegion(region string) []string {
	var addressPrefixes []string
	json := getAWSJsonString()
	cidrQueryResult := gjson.Get(json, fmt.Sprintf("prefixes.#(region==%s)#.ip_prefix", region))
	for _, addressPrefix := range cidrQueryResult.Array() {
		addressPrefixes = append(addressPrefixes, addressPrefix.String())
	}
	sort.Strings(addressPrefixes)
	return funk.UniqString(addressPrefixes)
}

func (a AWS) ListAddressPrefixesByService(service string) []string {
	var addressPrefixes []string
	json := getAWSJsonString()
	cidrQueryResult := gjson.Get(json, fmt.Sprintf("prefixes.#(service==%s)#.ip_prefix", service))
	for _, addressPrefix := range cidrQueryResult.Array() {
		addressPrefixes = append(addressPrefixes, addressPrefix.String())
	}
	sort.Strings(addressPrefixes)
	return funk.UniqString(addressPrefixes)
}

func (a AWS) ListAddressPrefixesByServiceAndRegion(service string, region string) []string {
	var addressPrefixes []string
	json := getAWSJsonString()
	cidrServiceQueryResult := gjson.Get(json, fmt.Sprintf("prefixes.#(service==%s)#.ip_prefix", service))
	cidrRegionQueryResult := gjson.Get(json, fmt.Sprintf("prefixes.#(region==%s)#.ip_prefix", region))
	for _, addressPrefix := range cidrServiceQueryResult.Array() {
		if strings.Contains(cidrRegionQueryResult.String(), addressPrefix.String()) {
			addressPrefixes = append(addressPrefixes, addressPrefix.String())
		}
	}
	sort.Strings(addressPrefixes)
	return funk.UniqString(addressPrefixes)
}

func getAWSJsonString() string {
	if awsJSON == "" {
		awsJSON = GetJsonString("/aws.json")
	}
	return awsJSON
}
