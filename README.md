<span align="center">

# API Mesh Source Registry

</span>

API Mesh Source Registry allows developers to get sources reviewed and endorsed by Adobe team.

## Requirements
* The source's checks must pass.
* The source's metadata must pass

## How to submit a new source
If you would like your connector reviewed and added to the registry, please open a pull request in this repository and fill in the Pull Request template.

### Submit new source steps Instruction
* Create new file with the name that represent source in the `/connectors` folder
* Provide required metadata (`name`, `version`, `description`, `author` fields)
* Provide source content as a value of `provider` field
* Create Pull Request
* Verify that all checks are passed

### Update source steps Instruction
* Open the source file in the `/connectors` folder
* Update source content
* Create Pull Request
* Verify that all checks are passed

### How it works
After Pull Request with the new or updated source will be merged the automation start running. The automation collects metadata and made index of sources. When automation is finished you will be able to see created/updated source marked by specific version in `archive` folder and corresponded metadata in `/connectors-metadata.json`.

### How to use
To use sources please install [AIO CLI API mesh plugin](https://github.com/adobe/aio-cli-plugin-api-mesh)
