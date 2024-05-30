---
hide:
  - navigation
---
## abtasty-cli web-experimentation modification create

Create a modification

### Synopsis

Create a modification

```
abtasty-cli web-experimentation modification create [--campaign-id=<campaign-id>] [-d <data-raw> | --data-raw=<data-raw>] [flags]
```

### Options

```
  -d, --data-raw string   raw data contains all the info to create your modification, check the doc for details
  -h, --help              help for create
```

### Options inherited from parent commands

```
      --campaign-id int        the campaign id of your modifications
  -f, --output-format string   output format for the get and list subcommands for AB Tasty resources. Only 3 format are possible: table, json, json-pretty (default "table")
      --user-agent string      custom user agent (default "abtasty-cli/main")
```

### SEE ALSO

* [abtasty-cli web-experimentation modification](abtasty-cli_web-experimentation_modification.md)	 - Manage your modifications

###### Auto generated by spf13/cobra on 30-May-2024