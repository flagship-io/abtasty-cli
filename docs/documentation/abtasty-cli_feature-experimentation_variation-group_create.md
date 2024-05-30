---
hide:
  - navigation
---
## abtasty-cli feature-experimentation variation-group create

Create a variation group

### Synopsis

Create a variation group

```
abtasty-cli feature-experimentation variation-group create [--campaign-id=<campaign-id>] [-d <data-raw> | --data-raw <data-raw>] [flags]
```

### Options

```
  -d, --data-raw string   raw data contains all the info to create your variation group, check the doc for details
  -h, --help              help for create
```

### Options inherited from parent commands

```
      --campaign-id string     the campaign id of your variation group
  -f, --output-format string   output format for the get and list subcommands for AB Tasty resources. Only 3 format are possible: table, json, json-pretty (default "table")
      --user-agent string      custom user agent (default "abtasty-cli/main")
```

### SEE ALSO

* [abtasty-cli feature-experimentation variation-group](abtasty-cli_feature-experimentation_variation-group.md)	 - Manage your variation groups

###### Auto generated by spf13/cobra on 30-May-2024