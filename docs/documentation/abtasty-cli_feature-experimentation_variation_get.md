---
hide:
  - navigation
---
## abtasty-cli feature-experimentation variation get

Get a variation

### Synopsis

Get a variation

```
abtasty-cli feature-experimentation variation get [--campaign-id=<campaign-id>] [--variation-group-id=<variation-group-id>] [-i <variation-id> | --id=<variation-id>] [flags]
```

### Options

```
  -h, --help        help for get
  -i, --id string   id of the variation you want to display
```

### Options inherited from parent commands

```
      --campaign-id string          the campaign id of your variation
  -f, --output-format string        output format for the get and list subcommands for AB Tasty resources. Only 3 format are possible: table, json, json-pretty (default "table")
      --user-agent string           custom user agent (default "abtasty-cli/main")
      --variation-group-id string   the variation group id of your variation
```

### SEE ALSO

* [abtasty-cli feature-experimentation variation](abtasty-cli_feature-experimentation_variation.md)	 - Manage your variations

###### Auto generated by spf13/cobra on 30-May-2024