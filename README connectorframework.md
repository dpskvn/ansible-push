# Sample Push Connector



## Getting started

### Building and running the connector

``` sh
go build
./sample-connector
```

That should start your connector on port 8080. 

API Calls are as follows:

``` sh
curl --location --request POST 'http://localhost:8080/v1/testconnection' \
--header 'Content-Type: application/json' \
--data-raw '{
"connection": {
"address": "192.168.3.7",
"password": "",
"username": "admin"
}
}'
```

### Automating the build and deploy
1. Add a `.gitlab-ci.yml` file to the root of your project:

``` yaml
workflow:
 rules:
   - if: $CI_PIPELINE_SOURCE == "merge_request_event"
   - if: $CI_COMMIT_BRANCH == $CI_DEFAULT_BRANCH

variables:
  CONNECTOR_ID:

include:
  - project: venafi/vaas-connectors/shared
    file:
      - /ci/common.gitlab-ci.yml
      - /ci/connector.gitlab-ci.yml
```

This will run the pipeline for building, testing, linting and deploying the
connector.

If you don't specify the `CONNECTOR_ID`, it won't be deployed to production, be
sure to contact an admin when you are ready for this process.

2. Create a Makefile on the root of your project:

``` makefile
include ./build/Makefile
```

You don't need to add a build folder or Makefile inside it, this `include`
references a pipeline file including all required commands.

## Naming your connector
The name of your connector should match the go module name, and it is assumed to
be the url of your gitlab repository, example:

For the Citrix ADC Connector, the url is:

`gitlab.com/venafi/vaas-connectors/citrix-adc-connector`


The package name inside `go.mod` file should contain the full path to the gitlab
repository, example:

``` go
module gitlab.com/venafi/vaas-connectors/citrix-adc-connector

go 1.18

require (
....
)
```

This will compile a binary called `citrix-adc-connector` when running `go
build.`

To run you would execute `./citrix-adc-connector`.

## Manifest template
Manifest file name should be named `manifest.json` in your connector repository,
this is a mandatory step, since it will be automatically deployed using the
pipeline.

Required fields are:
- name
- domainSchema.certificateBundle
- domainSchema.connection
- domainSchema.keystore
- domainSchema.binding
- localizationResources
- hooks

Optional fields are
- description


Use the manifest editor to render and verify your manifest.

<https://ui.venafi.cloud/manifest-simulator>
