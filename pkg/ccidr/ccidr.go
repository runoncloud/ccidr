package ccidr

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
)

func RunCommand(cmd *cobra.Command, args []string) error {
	regionParam, _ := getFlagString(cmd, "region")
	serviceParam, _ := getFlagString(cmd, "service")

	cloud, err := getCloud(args)
	if err != nil {
		return err
	}

	resource, err := getResource(args)
	if err != nil {
		return err
	}

	if resource == "regions" {
		regions := cloud.ListRegions()

		var data [][]string
		table := tablewriter.NewWriter(os.Stdout)
		for _, region := range regions {
			data = append(data, []string{region})
		}
		table.SetBorder(false)
		table.AppendBulk(data)
		table.Render()
	}

	if resource == "services" {
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
	}

	if resource == "ips" {
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

func getCloud(args []string) (cloud Cloud, err error) {
	if args[0] == "aws" {
		cloud = AWS{}
	} else if args[0] == "azure" {
		cloud = Azure{}
	} else {
		err = errors.New("Cloud is not valid or supported")
	}
	return
}

func getFlagString(cmd *cobra.Command, flag string) (flagValue string, err error) {
	flagValue, err = cmd.Flags().GetString(flag)
	return
}
