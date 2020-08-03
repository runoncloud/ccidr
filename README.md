# ccidr

`ccidr` is a command line tool written in Go that help you get and filter Public Cloud IP address ranges.

Public clouds supported : 

- Azure
- AWS

The IP address ranges for each cloud are included in the binaries. A new version will be release each week with update IP address ranges.

## Examples

- Prints All Azure regions
  ```bash
  ccidr azure regions
  ```
  
- Prints All Azure services
  ```bash
  ccidr azure services
  ```

- Prints All Azure IP address ranges
  ```bash
  ccidr azure ips
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
## Installation

### The Go Get way
 ```bash
 go get -u github.com/runoncloud/ccidr/cmd/ccidr
 ```
### Binaries
 
#### OSX
 ```bash
 curl -L -o ccidr.gz https://github.com/runoncloud/ccidr/releases/download/v1.0.0/ccidr_darwin_amd64.tar.gz && \
   tar zxvf ccidr.gz && chmod +x ccidr && mv ccidr $GOPATH/bin/
 ```
 
#### Linux
 ```bash
 curl -L -o ccidr.gz https://github.com/runoncloud/ccidr/releases/download/v1.0.0/ccidr_linux_amd64.tar.gz && \
   tar zxvf ccidr.gz && chmod +x ccidr && mv ccidr $GOPATH/bin/
 ```

#### Windows

 ```
 https://github.com/runoncloud/ccidr/releases/download/v1.0.0/ccidr_windows_amd64.zip
 ```

### From source

Requirements:
 - go 1.13 or newer
 - GNU make
 - git
 
 ```bash
 make bin   # binaries will be placed in bin/
 ```