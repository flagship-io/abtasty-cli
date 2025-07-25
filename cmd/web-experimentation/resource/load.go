/*
Copyright © 2022 Flagship Team flagship@abtasty.com

*/

package resource

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/flagship-io/abtasty-cli/cmd/web-experimentation/campaign"
	"github.com/flagship-io/abtasty-cli/cmd/web-experimentation/modification"
	"github.com/flagship-io/abtasty-cli/cmd/web-experimentation/variation"
	"github.com/flagship-io/abtasty-cli/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/spf13/cobra"
)

var (
	resourceFile string
	outputFile   string
	inputRefRaw  string
	inputRefFile string
)

type ResourceAction string

const (
	ActionCreate ResourceAction = "create"
	ActionUpdate ResourceAction = "update"
	ActionList   ResourceAction = "list"
	ActionGet    ResourceAction = "get"
	ActionDelete ResourceAction = "delete"
	ActionSwitch ResourceAction = "switch"
)

type Resource struct {
	Type           string         `json:"type"`
	Ref            string         `json:"$_ref"`
	ParentID       string         `json:"$_parent_id"`
	Action         ResourceAction `json:"action"`
	Payload        map[string]any `json:"payload"`
	Resources      []Resource     `json:"resources"`
	ParentResource *Resource
}

type LoadResFile struct {
	Version   int        `json:"version"`
	Resources []Resource `json:"resources"`
}

// Reference context for resolving $ref and $parent_id
type RefContext struct {
	refs map[string]any
}

type ResourceResult struct {
	Ref      string      `json:"$_ref,omitempty"`
	Status   string      `json:"status"`
	Response interface{} `json:"response,omitempty"`
}

type LoaderResults struct {
	Results []ResourceResult `json:"results"`
}

func NewRefContext() *RefContext {
	return &RefContext{refs: make(map[string]any)}
}

func (rc *RefContext) Set(ref string, val any) {
	if ref != "" {
		rc.refs[ref] = val
	}
}

func (rc *RefContext) Get(ref string) (any, bool) {
	val, ok := rc.refs[ref]
	return val, ok
}

// Recursively resolve $-prefixed references in payload and parent_id
func resolveRefs(val any, rc *RefContext) any {
	switch v := val.(type) {
	case map[string]any:
		for k, vv := range v {
			if s, ok := vv.(string); ok && strings.HasPrefix(s, "$") {
				// Reference: $ref.field
				parts := strings.Split(strings.TrimPrefix(s, "$"), ".")
				if len(parts) > 1 {
					if refVal, ok := rc.Get(parts[0]); ok {
						if m, ok := refVal.(map[string]any); ok {
							if field, ok := m[parts[1]]; ok {
								v[k] = field
							}
						}
					}
				}
			} else {
				v[k] = resolveRefs(vv, rc)
			}
		}
	case []any:
		for i, vv := range v {
			v[i] = resolveRefs(vv, rc)
		}
	}
	return val
}

// Main loader for new resource format
func LoadResources(cmd *cobra.Command, filePath, inputRefFile, inputRefRaw, outputFile string) error {

	var results []ResourceResult

	// Helper to record each result
	recordResult := func(ref string, status string, resp interface{}) {
		if ref != "" {
			results = append(results, ResourceResult{
				Ref:      ref,
				Status:   status,
				Response: resp,
			})
		}
	}

	// Wrap processResource to collect results
	processAndRecord := func(cmd *cobra.Command, res Resource, rc *RefContext) error {
		resp, err := processResourceWithResponse(cmd, res, rc)
		status := "success"
		if err != nil {
			status = "error"
			resp = err
		}
		recordResult(res.Ref, status, resp)
		return err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read resource file: %w", err)
	}

	var loadFile LoadResFile
	if err := json.Unmarshal(data, &loadFile); err != nil {
		return fmt.Errorf("failed to parse resource file: %w", err)
	}

	refCtx := NewRefContext()

	// Load inputRef from file or raw
	var inputRef map[string]any
	if inputRefFile != "" {
		b, err := os.ReadFile(inputRefFile)
		if err != nil {
			return fmt.Errorf("failed to read input ref file: %w", err)
		}

		err = json.Unmarshal(b, &inputRef)
		if err != nil {
			return fmt.Errorf("failed to read input ref file: %w", err)
		}
	} else if inputRefRaw != "" {
		err = json.Unmarshal([]byte(inputRefRaw), &inputRef)
		if err != nil {
			return fmt.Errorf("failed to read input ref file: %w", err)
		}
	}
	// Merge inputParams into refCtx
	for k, v := range inputRef {
		refCtx.Set(k, v)
	}

	// Separate mutating and read actions, preserving file order
	var mutating, read []Resource
	for _, res := range loadFile.Resources {
		switch res.Action {
		case ActionGet, ActionList:
			read = append(read, res)
		default:
			mutating = append(mutating, res)
		}
	}

	// Group resources by type
	var campaigns, variations, modifications, others []Resource
	for _, res := range mutating {
		switch res.Type {
		case "campaign":
			campaigns = append(campaigns, res)
		case "variation":
			variations = append(variations, res)
		case "modification":
			modifications = append(modifications, res)
		default:
			others = append(others, res)
		}
	}

	// Process in order: campaigns → variations → modifications → others
	for _, res := range campaigns {
		_ = processAndRecord(cmd, res, refCtx)
	}

	for _, res := range variations {
		_ = processAndRecord(cmd, res, refCtx)
	}

	for _, res := range modifications {
		_ = processAndRecord(cmd, res, refCtx)
	}

	for _, res := range others {
		_ = processAndRecord(cmd, res, refCtx)
	}

	for _, res := range read {
		_ = processAndRecord(cmd, res, refCtx)
	}

	// Output results
	loaderResults := LoaderResults{Results: results}
	if outputFile != "" {
		b, _ := json.MarshalIndent(loaderResults, "", "  ")
		_ = os.WriteFile(outputFile, b, 0644)
		fmt.Fprintf(cmd.OutOrStdout(), "Results written to %s\n", outputFile)
	} else {
		// Print as table
		fmt.Fprintf(cmd.OutOrStdout(), "%-10s %-10s %s\n", "$_ref", "status", "response")
		for _, r := range results {
			respStr, _ := json.Marshal(r.Response)
			fmt.Fprintf(cmd.OutOrStdout(), "%-10s %-10s %s\n", r.Ref, r.Status, string(respStr))
		}
	}

	return nil
}

func processResourceWithResponse(cmd *cobra.Command, res Resource, rc *RefContext) (interface{}, error) {
	if res.ParentID != "" && strings.HasPrefix(res.ParentID, "$") {
		parts := strings.Split(strings.TrimPrefix(res.ParentID, "$"), ".")
		if len(parts) > 1 {
			if refVal, ok := rc.Get(parts[0]); ok {
				if m, ok := refVal.(map[string]any); ok {
					if field, ok := m[parts[1]].(float64); ok {
						res.ParentID = fmt.Sprintf("%v", int(field))
					}
				}
			}
		}
	}
	res.Payload = resolveRefs(res.Payload, rc).(map[string]any)

	var resp any
	var err error
	switch res.Action {
	case ActionCreate:
		resp, err = handleCreate(res)
		//	case ActionUpdate:
		//		resp, err = handleUpdate(cmd, res)
	case ActionList:
		resp, err = handleList(res)
	case ActionDelete:
		err = handleDelete(cmd, res)
	case ActionSwitch:
		err = handleSwitch(cmd, res)
	default:
		err = fmt.Errorf("unsupported action: %s", res.Action)
	}
	if err != nil {
		return nil, err
	}

	// Store response for $ref
	if res.Ref != "" && resp != nil {
		rc.Set(res.Ref, resp)
	}

	// Recursively process children (collect their results too)
	if res.Action == ActionCreate || res.Action == ActionUpdate {
		for _, child := range res.Resources {
			if child.ParentID == "" && res.Ref != "" && resp != nil {
				if id, ok := resp.(map[string]any)["id"].(float64); ok {
					child.ParentID = fmt.Sprintf("%v", int(id))
					child.ParentResource = &res
				}
			}
			_, err = processResourceWithResponse(cmd, child, rc)
			if err != nil {
				return nil, err
			}
		}
	}

	return resp, nil
}

// Example handlers (implement as needed)
func handleCreate(res Resource) (map[string]any, error) {
	// Map type to actual creation logic
	payloadBytes, _ := json.Marshal(res.Payload)
	var respBytes []byte
	//var err error

	switch res.Type {
	case "campaign":
		var err error
		respBytes, err = campaign.CreateCampaign(payloadBytes)
		if err != nil {
			return nil, err
		}
	case "variation":
		parentID, err := strconv.Atoi(res.ParentID)
		if err != nil {
			return nil, err
		}

		respBytes, err = variation.CreateVariation(parentID, payloadBytes)
		if err != nil {
			return nil, err
		}
	case "modification":
		variationID, err := strconv.Atoi(res.ParentID)
		if err != nil {
			return nil, err
		}

		var modificationResourceLoader web_experimentation.ModificationResourceLoader
		err = json.Unmarshal(payloadBytes, &modificationResourceLoader)
		if err != nil {
			return nil, err
		}

		if res.ParentResource != nil {
			campaignIDString := res.ParentResource.ParentID
			campaignID, err := strconv.Atoi(campaignIDString)
			if err != nil {
				return nil, err
			}

			modificationResourceLoader.CampaignID = campaignID
		}

		respBytes, err = modification.CreateModification(variationID, modificationResourceLoader)
		if err != nil {
			return nil, err
		}

		fmt.Println(string(respBytes))

	default:
		return nil, fmt.Errorf("unknown resource type: %s", res.Type)
	}

	var resp map[string]any
	_ = json.Unmarshal(respBytes, &resp)
	return resp, nil
}

// Example handlers (implement as needed)
/* func handleUpdate(cmd *cobra.Command, res Resource) (map[string]any, error) {
	// Map type to actual creation logic
	payloadBytes, _ := json.Marshal(res.Payload)
	var respBytes []byte
	//var err error

	switch res.Type {
	case "campaign":
		respBytes = campaign.CreateCampaign(payloadBytes)
		fmt.Fprintf(cmd.OutOrStdout(), "Create action for type: %s\n", res.Type)
		return nil, nil
	case "variation":
		// parentID required
		//parentID, _ := strconv.Atoi(res.ParentID)
		var v models.VariationWE
		_ = json.Unmarshal(payloadBytes, &v)
		//respBytes, err = http_request.VariationWERequester.HTTPCreateVariation(parentID, v)
		fmt.Fprintf(cmd.OutOrStdout(), "Create action for type: %s\n", res.Type)
	case "modification":
		// parentID required
		//parentID, _ := strconv.Atoi(res.ParentID)
		var m models.ModificationCodeCreateStruct
		_ = json.Unmarshal(payloadBytes, &m)
		//respBytes, err = http_request.ModificationRequester.HTTPCreateModification(parentID, m)
		fmt.Fprintf(cmd.OutOrStdout(), "Create action for type: %s\n", res.Type)
	default:
		return nil, fmt.Errorf("unknown resource type: %s", res.Type)
	}
	/* 	if err != nil {
		return nil, err
	}
	var resp map[string]any
	_ = json.Unmarshal(respBytes, &resp)
	return resp, nil
} */

func handleList(res Resource) (any, error) {
	var respBytes []byte
	payloadBytes, _ := json.Marshal(res.Payload)

	switch res.Type {
	case "campaign":
		campaignList, err := campaign.ListCampaigns()
		if err != nil {
			return nil, err
		}

		respBytes, err = json.Marshal(campaignList)
		if err != nil {
			return nil, err
		}
	case "modification":
		var modificationResourceLoader web_experimentation.ModificationResourceLoader
		err := json.Unmarshal(payloadBytes, &modificationResourceLoader)
		if err != nil {
			return nil, err
		}
		if res.ParentResource != nil {
			campaignIDString := res.ParentResource.ParentID
			campaignID, err := strconv.Atoi(campaignIDString)
			if err != nil {
				return nil, err
			}

			modificationResourceLoader.CampaignID = campaignID
		}

		modificationList, err := modification.ListModifications(modificationResourceLoader.CampaignID)
		if err != nil {
			return nil, err
		}

		respBytes, err = json.Marshal(modificationList)
	default:
		return nil, fmt.Errorf("unknown resource type: %s", res.Type)
	}

	var resp any
	err := json.Unmarshal(respBytes, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func handleDelete(cmd *cobra.Command, res Resource) error {
	// Implement as needed
	fmt.Fprintf(cmd.OutOrStdout(), "Delete action for type: %s\n", res.Type)
	return nil
}

func handleSwitch(cmd *cobra.Command, res Resource) error {
	id := fmt.Sprintf("%v", res.Payload["id"])
	state := fmt.Sprintf("%v", res.Payload["state"])
	fmt.Fprintf(cmd.OutOrStdout(), "--id=%s --state=%s\n", id, state)

	valid := map[string]bool{"active": true, "paused": true, "interrupted": true}
	if !valid[state] {
		fmt.Fprintln(cmd.OutOrStdout(), "Status must be one of: active, paused, interrupted")
		return nil
	}
	if err := http_request.CampaignFERequester.HTTPSwitchStateCampaign(id, state); err != nil {
		fmt.Fprintf(cmd.OutOrStdout(), "Switch error: %v\n", err)
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "Campaign status set to %s\n", state)
	}
	return nil
}

// LoadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:   "load [--file=<file>]",
	Short: "Load your resources",
	RunE: func(cmd *cobra.Command, args []string) error {
		if resourceFile == "" {
			return fmt.Errorf("missing --file flag")
		}
		return LoadResources(cmd, resourceFile, inputRefFile, inputRefRaw, outputFile)
	},
}

func init() {
	//cobra.OnInitialize(initResource)

	loadCmd.Flags().StringVarP(&resourceFile, "file", "", "", "resource file that contains your resource")

	if err := loadCmd.MarkFlagRequired("file"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	loadCmd.Flags().StringVarP(&outputFile, "output-file", "", "", "result of the command that contains all resource information")

	loadCmd.Flags().StringVarP(&inputRefRaw, "input-ref", "", "", "params to replace resource loader file")
	loadCmd.Flags().StringVarP(&inputRefFile, "input-ref-file", "", "", "file that contains params to replace resource loader file")

	ResourceCmd.AddCommand(loadCmd)
}
