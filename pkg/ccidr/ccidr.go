package ccidr

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
)

const (
	RegionsArg  = "regions"
	ServicesArg = "services"
	IpsArg      = "ips"
	AzureArg    = "azure"
	AWSArgs     = "aws"
)

func RunCommand(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		printsCloudCommands()
		return nil
	}

	regionParam, _ := getFlagString(cmd, "region")
	serviceParam, _ := getFlagString(cmd, "service")
	remoteParam, _ := getFlagBool(cmd, "remote")

	cloud, err := getCloud(args, remoteParam)
	if err != nil {
		return err
	}

	if len(args) == 1 {
		printsResourceCommands()
		return nil
	}

	resource, err := getResource(args)
	if err != nil {
		return err
	}

	if resource == RegionsArg {
		regions := cloud.ListRegions()

		var data [][]string
		table := tablewriter.NewWriter(os.Stdout)
		for _, region := range regions {
			data = append(data, []string{region})
		}
		table.SetBorder(false)
		table.AppendBulk(data)
		table.Render()

	} else if resource == ServicesArg {
		var services []string
		if regionParam == "" {
			services = cloud.ListServices()
		} else {
			services = cloud.ListServicesByRegion(regionParam)
		}

		var data [][]string
		table := tablewriter.NewWriter(os.Stdout)
		for _, service := range services {
			data = append(data, []string{service})
		}
		table.SetBorder(false)
		table.AppendBulk(data)
		table.Render()

	} else if resource == IpsArg {
		var addresses []string
		if regionParam != "" && serviceParam != "" {
			addresses = cloud.ListAddressPrefixesByServiceAndRegion(serviceParam, regionParam)
		} else if regionParam != "" {
			addresses = cloud.ListAddressPrefixesByRegion(regionParam)
		} else if serviceParam != "" {
			addresses = cloud.ListAddressPrefixesByService(serviceParam)
		} else {
			addresses = cloud.ListAddressPrefixes()
		}

		var data [][]string
		table := tablewriter.NewWriter(os.Stdout)
		for _, adress := range addresses {
			data = append(data, []string{adress})
		}
		table.SetBorder(false)
		table.AppendBulk(data)
		table.Render()
	}

	return nil
}

func getResource(args []string) (resource string, err error) {
	resource = args[1]
	if resource != "regions" && resource != "services" && resource != "ips" {
		err = errors.New(fmt.Sprintf("Argument : %s is not supported", resource))
	}
	return
}

func getCloud(args []string, isRemote bool) (cloud Cloud, err error) {
	if args[0] == AWSArgs {
		cloud = AWS{isRemote: isRemote}
	} else if args[0] == AzureArg {
		cloud = Azure{isRemote: isRemote}
	} else {
		err = errors.New("Cloud is not valid or supported")
	}
	return
}

func getFlagString(cmd *cobra.Command, flag string) (flagValue string, err error) {
	flagValue, err = cmd.Flags().GetString(flag)
	return
}

func getFlagBool(cmd *cobra.Command, flag string) (flagValue bool, err error) {
	flagValue, err = cmd.Flags().GetBool(flag)
	return
}

func printsCloudCommands() {
	println("Cloud Commands :")
	println("  azure \t Retrieves Azure resources")
	println("  aws \t\t Retrieves AWS resources ")
}

func printsResourceCommands() {
	println("Resource supported :")
	println("  regions \t Retrieves the list of regions")
	println("  services \t Retrieves the list of services")
	println("  ips \t\t Retrieves the list of IP address ranges")
}
