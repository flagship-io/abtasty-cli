---
hide:
  - navigation
---
## abtasty-cli web-experimentation campaign-global-code get

Get campaign global code

### Synopsis

Get campaign global code

```
abtasty-cli web-experimentation campaign-global-code get [-i <campaign-id> | --id <campaign-id>] [flags]
```

### Options

```
      --create-file       create a file that contains campaign global code
      --create-subfiles   create a file that contains campaign and variations global code
  -h, --help              help for get
  -i, --id string         id of the campaign you want to display
      --override          override existing campaign global code file
```

### Options inherited from parent commands

```
  -f, --output-format string   output format for the get and list subcommands for AB Tasty resources. Only 3 format are possible: table, json, json-pretty (default "table")
      --user-agent string      custom user agent (default "abtasty-cli/main")
      --working-dir string     Directory where the file will be generated and pushed from (default "/Users/chadi.laoulaou/abtasty-cli")
```

### SEE ALSO

* [abtasty-cli web-experimentation campaign-global-code](abtasty-cli_web-experimentation_campaign-global-code.md)	 - Manage campaign global code

###### Auto generated by spf13/cobra on 4-Jun-2024
