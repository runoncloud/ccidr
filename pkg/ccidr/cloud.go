package ccidr

type Cloud interface {
	ListRegions() []string
	ListServices() []string
	ListServicesByRegion(region string) []string
	ListAddressPrefixes() []string
	ListAddressPrefixesByService(service string) []string
	ListAddressPrefixesByRegion(region string) []string
	ListAddressPrefixesByServiceAndRegion(service string, region string) []string
}
