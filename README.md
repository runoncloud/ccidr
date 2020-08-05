# ccidr
![build](https://github.com/runoncloud/ccidr/workflows/build/badge.svg?branch=master)
![release](https://github.com/runoncloud/ccidr/workflows/release/badge.svg)
![files-update](https://github.com/runoncloud/ccidr/workflows/files-update/badge.svg)

`ccidr` is a command line tool written in Go that help you get and filter Public Cloud IP address ranges.

Public clouds supported : 

- Azure
- AWS

The IP address ranges for each cloud are included in the binaries. A new version will be release on monday each week with update IP address ranges. If you don't want to update on a weekly basis, you can use the `--remote` flag to fetch the information directly from the source : 

- AWS : https://ip-ranges.amazonaws.com/ip-ranges.json
- Azure : https://www.microsoft.com/en-us/download/details.aspx?id=56519&WT.mc_id=rss_alldownloads_all

## Examples

- Prints All Azure regions
  ```bash
  ccidr azure regions
  ```

- Prints All AWS regions
  ```bash
  ccidr aws regions
  ```
  
- Prints All Azure services
  ```bash
  ccidr azure services
  ```

- Prints All Azure IP address ranges
  ```bash
  ccidr azure ips
  ```
  
- Prints All AWS IP address ranges
  ```bash
  ccidr aws ips
  ```

- Prints All Azure IP address ranges for a specific region
  ```bash
  ccidr azure ips -r eastus
  ```

- Prints All Azure IP address ranges for a specific service
  ```bash
  ccidr azure ips -s AppService
  ```

- Prints All Azure IP address ranges for a specific service and region
  ```bash
  ccidr azure ips -s AppService -r eastus
  ```

- Prints All Azure IP address ranges for a specific service and region using the remote source of data
  ```bash
  ccidr azure ips -s AppService -r eastus --remote
  ```
## Installation

### The Go Get way
 ```bash
 go get -u github.com/runoncloud/ccidr/cmd/ccidr
 ```
### Binaries
 
#### OSX
 ```bash
 latestVersion=$(curl --silent "https://api.github.com/repos/runoncloud/ccidr/releases/latest" | jq -r .tag_name) && \
   curl -L -o ccidr.gz https://github.com/runoncloud/ccidr/releases/download/$latestVersion/ccidr_darwin_amd64.tar.gz && \
   tar zxvf ccidr.gz && chmod +x ccidr && mv ccidr $GOPATH/bin/
 ```
 
#### Linux
 ```bash
 latestVersion=$(curl --silent "https://api.github.com/repos/runoncloud/ccidr/releases/latest" | jq -r .tag_name) &&
   curl -L -o ccidr.gz https://github.com/runoncloud/ccidr/releases/download/$latestVersion/ccidr_linux_amd64.tar.gz && \
   tar zxvf ccidr.gz && chmod +x ccidr && mv ccidr $GOPATH/bin/
 ```

### From source

Requirements:
 - go 1.13 or newer
 - GNU make
 - git
 
 ```bash
 make bin        # binaries will be placed in bin/
 ```

## Use ccidr as a library
 ```go
package main

import (
	"fmt"
	"github.com/runoncloud/ccidr/pkg/ccidr"
)

func main() {
	azure := ccidr.Azure{}
	aws := ccidr.AWS{isRemote: true}
    
	fmt.Println(azure.ListAddressPrefixes())
	fmt.Println(aws.ListAddressPrefixes())
}
 ```