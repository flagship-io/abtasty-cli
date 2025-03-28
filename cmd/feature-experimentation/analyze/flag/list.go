/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package flag

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/fatih/color"
	models "github.com/flagship-io/abtasty-cli/models/feature_experimentation"
	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/flagship-io/codebase-analyzer/pkg/handler"
	"github.com/kyokomi/emoji/v2"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
)

var codebaseAnalyzer bool

func summaryTableFlagListed(flagExistLen, flagNotExistLen int) {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	totalFlag := flagExistLen + flagNotExistLen

	summtbl := table.New("\nSummary")
	summtbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	summtbl.AddRow(fmt.Sprintf("Total flags: %d (%d Flag exist %s, %d Flag don't exist%s)", totalFlag, flagExistLen, emoji.Sprint(":check_mark_button:"), flagNotExistLen, emoji.Sprint(":cross_mark:")))
	summtbl.Print()
}

func flagListedTable(cmd *cobra.Command, listedFlags []models.Flag) error {

	var flagExistLen int = 0
	var flagNotExistLen int = 0
	var flagKeyNotDetected []string
	var flagLocationAddedToTable []string

	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Flag", "Type", "defaultValue", "File", fmt.Sprintf("Exists ? (%s/%s)", emoji.Sprint(":check_mark_button:"), emoji.Sprint(":cross_mark:")))
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt).WithPadding(2)

	var existedFlagKey []string

	for _, flag := range listedFlags {
		existedFlagKey = append(existedFlagKey, strings.ToLower(flag.Name))
	}

	results, err := handler.ExtractFlagsInfo(FSConfig)
	if err != nil {
		return err
	}

	if viper.GetString("output_format") == "json" {
		var filesAnalyzed []models.FileAnalyzed
		for _, r := range results {
			var fileAnalyzed = models.FileAnalyzed{
				File:    r.File,
				FileURL: r.FileURL,
				Error:   r.Error,
			}
			var flagsAnalyzed []models.FlagAnalyzed
			for _, analyzedFlag := range r.Results {
				var flagAnalyzed = models.FlagAnalyzed{
					LineNumber:       analyzedFlag.LineNumber,
					FlagKey:          analyzedFlag.FlagKey,
					FlagDefaultValue: analyzedFlag.FlagDefaultValue,
					FlagType:         analyzedFlag.FlagType,
					Exists:           false,
				}

				if analyzedFlag.FlagKey == "" {
					continue
				}

				if slices.Contains(existedFlagKey, strings.ToLower(analyzedFlag.FlagKey)) {
					flagAnalyzed.Exists = true
					flagsAnalyzed = append(flagsAnalyzed, flagAnalyzed)
					continue
				}
				flagsAnalyzed = append(flagsAnalyzed, flagAnalyzed)

			}
			fileAnalyzed.Results = flagsAnalyzed
			if len(fileAnalyzed.Results) != 0 {
				filesAnalyzed = append(filesAnalyzed, fileAnalyzed)
			}

		}

		json, _ := json.Marshal(filesAnalyzed)
		fmt.Fprintln(cmd.OutOrStdout(), string(json))
		return nil
	}

	for _, r := range results {
		pathArray := strings.Split(r.File, "/")
		for _, analyzedFlag := range r.Results {

			if analyzedFlag.FlagKey == "" {
				if !slices.Contains(flagLocationAddedToTable, fmt.Sprintf("%s:%d", r.File, analyzedFlag.LineNumber)) {
					flagKeyNotDetected = append(flagKeyNotDetected, fmt.Sprintf("%s:%d", r.File, analyzedFlag.LineNumber))
				}
				continue
			}

			flagLocationAddedToTable = append(flagLocationAddedToTable, fmt.Sprintf("%s:%d", r.File, analyzedFlag.LineNumber))

			if slices.Contains(existedFlagKey, strings.ToLower(analyzedFlag.FlagKey)) {
				flagExistLen += 1
				tbl.AddRow(analyzedFlag.FlagKey, analyzedFlag.FlagType, analyzedFlag.FlagDefaultValue, fmt.Sprintf("%s:%d", pathArray[len(pathArray)-1], analyzedFlag.LineNumber), emoji.Sprint(":check_mark_button:"))
				continue
			}

			if analyzedFlag.FlagKey == "" {
				flagKeyNotDetected = append(flagKeyNotDetected, fmt.Sprintf("%s, line: %d", r.File, analyzedFlag.LineNumber))
				continue
			}

			flagNotExistLen += 1

			tbl.AddRow(analyzedFlag.FlagKey, analyzedFlag.FlagType, analyzedFlag.FlagDefaultValue, fmt.Sprintf("%s:%d", pathArray[len(pathArray)-1], analyzedFlag.LineNumber), emoji.Sprint(":cross_mark:"))
		}
	}

	totalFlag := flagExistLen + flagNotExistLen

	if totalFlag == 0 {
		tbl.AddRow("No flag found")
	}

	tbl.Print()

	summaryTableFlagListed(flagExistLen, flagNotExistLen)

	if len(flagKeyNotDetected) != 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "\n%sWarning: feature flags functions detected in these files, but flags are unknown: \n", emoji.Sprint(":construction:"))
		for _, flag := range RemoveDuplicateStr(flagKeyNotDetected) {
			fmt.Fprintf(cmd.OutOrStdout(), "%s\n", flag)
		}
	}

	return nil
}

// ListCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Analyze your codebase and list flags detected",
	Long:  `Analyze your codebase and list flags detected and check if it exist in Flagship platform`,
	PreRun: func(cmd *cobra.Command, args []string) {
		PreRunConfiguration()
	},
	Run: func(cmd *cobra.Command, args []string) {

		if codebaseAnalyzer {
			results, err := handler.ExtractFlagsInfo(FSConfig)
			if err != nil {
				log.Fatalf("error occurred when extracting flags info: %s", err)
			}
			json, _ := json.Marshal(results)
			fmt.Fprintln(cmd.OutOrStdout(), string(json))
			return
		}

		listExistingFlags, errListFlag := httprequest.FlagRequester.HTTPListFlag()
		if errListFlag != nil {
			log.Fatalf("error occurred when listing existing flag: %s", errListFlag)
		}

		err := flagListedTable(cmd, listExistingFlags)
		if err != nil {
			log.Fatalf("error occurred in listed flag table: %s", err)
		}

		if CustomRegexJson != "" {
			fmt.Fprintf(cmd.OutOrStdout(), "\n%sContribution: If this custom regexes comes from a competitor or it's an improvement of existing regexes, we invite you to create a PR in our github repo: https://github.com/flagship-io/abtasty-cli \n", emoji.Sprint(":glowing_star:"))
		}

		if OriginPlatform != "" {
			fmt.Fprintf(cmd.OutOrStdout(), "\n%sContribution: If these regexes are outdated or contains errors, we invite you to create an issue or contribute in our github repo: https://github.com/flagship-io/abtasty-cli \n", emoji.Sprint(":glowing_star:"))
		}

	},
}

func init() {
	FlagCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVarP(&codebaseAnalyzer, "codebase-analyzer", "", false, "list codebase analyzer extract information.")
}
