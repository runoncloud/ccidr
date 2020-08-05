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
	Amazon        = "AMAZON"
	AmazonJsonUrl = "https://ip-ranges.amazonaws.com/ip-ranges.json"
)

type AWS struct {
	isRemote bool
}

var (
	awsJSON string
)

func (a AWS) ListRegions() []string {
	var regions []string
	json := getAWSJsonString(a.isRemote)
	queryResult := gjson.Get(json, "prefixes.#.region")
	for _, region := range queryResult.Array() {
		regions = append(regions, region.String())
	}
	sort.Strings(regions)
	return funk.UniqString(regions)
}

func (a AWS) ListServices() []string {
	var services []string
	json := getAWSJsonString(a.isRemote)

	queryResult := gjson.Get(json, "prefixes.#.service")
	for _, service := range queryResult.Array() {
		services = append(services, service.String())
	}
	sort.Strings(services)
	return funk.UniqString(services)
}

func (a AWS) ListServicesByRegion(region string) []string {
	var services []string
	json := getAWSJsonString(a.isRemote)

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
	json := getAWSJsonString(a.isRemote)
	cidrQueryResult := gjson.Get(json, fmt.Sprintf("prefixes.#(region==%s)#.ip_prefix", region))
	for _, addressPrefix := range cidrQueryResult.Array() {
		addressPrefixes = append(addressPrefixes, addressPrefix.String())
	}
	sort.Strings(addressPrefixes)
	return funk.UniqString(addressPrefixes)
}

func (a AWS) ListAddressPrefixesByService(service string) []string {
	var addressPrefixes []string
	json := getAWSJsonString(a.isRemote)
	cidrQueryResult := gjson.Get(json, fmt.Sprintf("prefixes.#(service==%s)#.ip_prefix", service))
	for _, addressPrefix := range cidrQueryResult.Array() {
		addressPrefixes = append(addressPrefixes, addressPrefix.String())
	}
	sort.Strings(addressPrefixes)
	return funk.UniqString(addressPrefixes)
}

func (a AWS) ListAddressPrefixesByServiceAndRegion(service string, region string) []string {
	var addressPrefixes []string
	json := getAWSJsonString(a.isRemote)
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

func getAWSJsonString(isRemote bool) string {
	if awsJSON == "" {
		if isRemote {
			awsJSON = getAWSRemoteJsonString()
		} else {
			awsJSON = GetJsonString("/aws.json")
		}
	}
	return awsJSON
}

func getAWSRemoteJsonString() string {
	if awsJSON == "" {
		client := http.Client{Timeout: time.Second * 10}

		req, err := http.NewRequest(http.MethodGet, AmazonJsonUrl, nil)
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
		awsJSON = string(body)
	}
	return awsJSON
}
