---
hide:
  - navigation
---
## abtasty-cli feature-experimentation analyze flag list

Analyze your codebase and list flags detected

### Synopsis

Analyze your codebase and list flags detected and check if it exist in Flagship platform

```
abtasty-cli feature-experimentation analyze flag list [flags]
```

### Options

```
      --codebase-analyzer   list codebase analyzer extract information.
  -h, --help                help for list
```

### Options inherited from parent commands

```
      --code-edge int              number of line code edges (default 1)
      --custom-regex string        regex for the pattern you want to analyze
      --custom-regex-json string   json file that contains the regex for the pattern you want to analyze
      --directory string           directory to analyze in your codebase (default ".")
      --files-exclude string       list of files to exclude in analysis (default "[\".git\", \".github\", \".vscode\", \".idea\", \".yarn\", \"node_modules\"]")
      --origin-platform string     analyze flags made with feature flag platform, we support launchdarkly, optimizely, split and vwo (latest version only)
      --output-format string       output format for the get and list subcommands for AB Tasty resources. Only 3 format are possible: table, json, json-pretty (default "table")
      --repository-branch string   branch of the code you want to analyze, and is used to track the links of the files where your flags are used (default "main")
      --repository-url string      root URL of your repository, and is used to track the links of the files where your flags are used (default "https://github.com/org/repo")
      --user-agent string          custom user agent (default "abtasty-cli/main")
```

### SEE ALSO

* [abtasty-cli feature-experimentation analyze flag](abtasty-cli_feature-experimentation_analyze_flag.md)	 - Analyze your codebase and detect the usage of Flagship or custom flags

###### Auto generated by spf13/cobra on 24-Feb-2025
