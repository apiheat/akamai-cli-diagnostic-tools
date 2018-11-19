# Akamai CLI for Diagnostic Tools

The Diagnostic Tools API allows you to diagnose many common problems Akamai customers experience when delivering content to their end users. It offers a programmatic alternative to many of the features available in the Luna Control Center, under the Resolve ⇒ Diagnostic Tools menu.

Should you miss something we *gladly accept patches* :)

CLI uses custom [Akamai API client](https://github.com/apiheat/go-edgegrid)

## Configuration & Installation

### Credentials

Set up your credential files as described in the [authorization](https://developer.akamai.com/introduction/Prov_Creds.html) and [credentials](https://developer.akamai.com/introduction/Conf_Client.html) sections of the getting started guide on developer.akamai.com.

Tools expect proper format of sections in edgerc file which example is shown below

*NOTE:* Default file location is *~/.edgerc*

```
[default]
client_secret = XXXXXXXXXXXX
host = XXXXXXXXXXXX
access_token = XXXXXXXXXXXX
client_token = XXXXXXXXXXXX
```

In order to change section which is being actively used you can

* change it via `--config parameter` of the tool itself
* change it via env variable `export AKAMAI_EDGERC_CONFIG=/Users/jsmitsh/.edgerc`

In order to change section which is being actively used you can

* change it via `--section parameter` of the tool itself
* change it via env variable `export AKAMAI_EDGERC_SECTION=mycustomsection`

> *NOTE:* Make sure your API client do have appropriate scopes enabled

### Installation

The tool can be used as a stand-alone binary or in conjuction with [Akamai CLI](https://developer.akamai.com/cli).

#### Akamai-cli ( recommended )

Execute the following from console

```shell
> akamai install https://github.com/apiheat/akamai-cli-diagnostic-tools
```

#### Stand-alone

As part of automated releases/builds you can download latest version from the project release page

## Usage

```shell
NAME:
   akamai-cli-diagnostic-tools - A CLI to interact with Akamai Diagnostic Tools

USAGE:
   akamai-cli-diagnostic-tools [global options] command [command options] [arguments...]

VERSION:
   X.X.X

AUTHORS:
   Petr Artamonov
   Rafal Pieniazek

COMMANDS:
     diagnostic-link     Generates a unique link to send to a user to diagnose a problem
     is-akamai-ip, i     Checks whether the specified ip address is part of the Akamai edge network
     translate-error, t  Get information about error strings produced by edge servers when a request to retrieve content fails
     help, h             Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config FILE, -c FILE   Location of the credentials FILE (default: "/Users/partamonov/.edgerc") [$AKAMAI_EDGERC_CONFIG]
   --debug value            Debug Level [$AKAMAI_EDGERC_DEBUGLEVEL]
   --section NAME, -s NAME  NAME of section to use from credentials file (default: "default") [$AKAMAI_EDGERC_SECTION]
   --help, -h               show help
   --version, -v            print the version
```

## Development

In order to develop the tool with us do the following:

1. Fork repository
1. Clone it to your folder ( within *GO* path )
1. Ensure you can restore dependencies by running

   ```shell
   dep ensure
   ```

1. Make necessary changes
1. Make sure solution builds properly ( feel free to add tests )

   ```shell
   go build -ldflags="-s -w -X main.appVer=1.2.3 -X main.appName=$(basename `pwd`)"
   ```
