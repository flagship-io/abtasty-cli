---
hide:
  - navigation
---
## abtasty-cli feature-experimentation variation edit

Edit a variation

### Synopsis

Edit a variation

```
abtasty-cli feature-experimentation variation edit [--campaign-id=<campaign-id>] [--variation-group-id=<variation-group-id>] [-i <variation-id> | --id=<variation-id>] [-d <data-raw> | --data-raw=<data-raw>] [flags]
```

### Options

```
  -d, --data-raw string   raw data contains all the info to edit your variation, check the doc for details
  -h, --help              help for edit
  -i, --id string         id of the variation you want to edit
```

### Options inherited from parent commands

```
      --campaign-id string          the campaign id of your variation
      --output-format string        output format for the get and list subcommands for AB Tasty resources. Only 3 format are possible: table, json, json-pretty (default "table")
      --user-agent string           custom user agent (default "abtasty-cli/main")
      --variation-group-id string   the variation group id of your variation
```

### SEE ALSO

* [abtasty-cli feature-experimentation variation](abtasty-cli_feature-experimentation_variation.md)	 - Manage your variations

###### Auto generated by spf13/cobra on 24-Feb-2025
