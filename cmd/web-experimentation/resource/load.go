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
	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
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
	ActionEdit   ResourceAction = "edit"
	ActionList   ResourceAction = "list"
	ActionGet    ResourceAction = "get"
	ActionDelete ResourceAction = "delete"
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

type funcModifType func(int, []byte) ([]byte, error)

func createOrEditModification(id int, res Resource, payloadBytes []byte, createOrEditModif funcModifType) ([]byte, error) {
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

	payloadBytes, err = json.Marshal(modificationResourceLoader)
	if err != nil {
		return nil, err
	}

	respBytes, err := createOrEditModif(id, payloadBytes)
	if err != nil {
		return nil, err
	}

	return respBytes, nil
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

func resolveRefs(val any, rc *RefContext) any {
	switch v := val.(type) {
	case map[string]any:
		for k, vv := range v {
			if s, ok := vv.(string); ok && strings.HasPrefix(s, "$") {
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

func LoadResources(cmd *cobra.Command, filePath, inputRefFile, inputRefRaw, outputFile string) error {

	var results []ResourceResult

	recordResult := func(ref string, status string, resp interface{}) {
		if ref != "" {
			results = append(results, ResourceResult{
				Ref:      ref,
				Status:   status,
				Response: resp,
			})
		}
	}

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

	for k, v := range inputRef {
		refCtx.Set(k, v)
	}

	if err := ValidateResources(&loadFile, refCtx); err != nil {
		return fmt.Errorf("Validation failed: %v\n", err)
	}

	var mutating, read []Resource
	for _, res := range loadFile.Resources {
		switch res.Action {
		case ActionGet, ActionList:
			read = append(read, res)
		default:
			mutating = append(mutating, res)
		}
	}

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

	loaderResults := LoaderResults{Results: results}
	if outputFile != "" {
		b, err := json.MarshalIndent(loaderResults, "", "  ")
		if err != nil {
			return err
		}

		err = os.WriteFile(outputFile, b, 0644)
		if err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Results written to %s\n", outputFile)
	} else {
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
	case ActionEdit:
		resp, err = handleEdit(res)
	case ActionList:
		resp, err = handleList(res)
	case ActionDelete:
		resp, err = handleDelete(res)
	default:
		err = fmt.Errorf("unsupported action: %s", res.Action)
	}

	if err != nil {
		return nil, err
	}

	if res.Ref != "" && resp != nil {
		rc.Set(res.Ref, resp)
	}

	if res.Action == ActionCreate || res.Action == ActionEdit {
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

func handleCreate(res Resource) (map[string]any, error) {
	payloadBytes, err := json.Marshal(res.Payload)
	if err != nil {
		return nil, err
	}

	var respBytes []byte

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

		respBytes, err = createOrEditModification(variationID, res, payloadBytes, modification.CreateModification)
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unknown resource type: %s", res.Type)
	}

	var resp map[string]any
	err = json.Unmarshal(respBytes, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func handleEdit(res Resource) (map[string]any, error) {
	var respBytes []byte

	payloadBytes, err := json.Marshal(res.Payload)
	if err != nil {
		return nil, err
	}

	var id = res.Payload["id"].(float64)

	if id == 0 {
		return nil, fmt.Errorf("error occurred: missing property %s", "id")
	}

	switch res.Type {
	case "campaign":
		var err error
		respBytes, err = campaign.EditCampaign(int(id), payloadBytes)
		if err != nil {
			return nil, err
		}
	case "variation":
		parentID, err := strconv.Atoi(res.ParentID)
		if err != nil {
			return nil, err
		}

		respBytes, err = variation.EditVariation(int(id), parentID, payloadBytes)
		if err != nil {
			return nil, err
		}
	case "modification":
		respBytes, err = createOrEditModification(int(id), res, payloadBytes, modification.EditModification)
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unknown resource type: %s", res.Type)
	}

	var resp map[string]any
	err = json.Unmarshal(respBytes, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func handleList(res Resource) (any, error) {
	var respBytes []byte
	var err error
	payloadBytes, err := json.Marshal(res.Payload)
	if err != nil {
		return nil, err
	}

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
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown resource type: %s", res.Type)
	}

	var resp any
	err = json.Unmarshal(respBytes, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func handleDelete(res Resource) (any, error) {
	var respBytes []byte

	payloadBytes, err := json.Marshal(res.Payload)
	if err != nil {
		return nil, err
	}

	var id = res.Payload["id"].(float64)

	if id == 0 {
		return nil, fmt.Errorf("error occurred: missing property %s", "id")
	}

	switch res.Type {
	case "campaign":
		resp, err := httprequest.CampaignWERequester.HTTPDeleteCampaign(int(id))
		if err != nil {
			return nil, err
		}

		respBytes, err = json.Marshal(resp)
		if err != nil {
			return nil, err
		}
	case "variation":
		parentID, err := strconv.Atoi(res.ParentID)
		if err != nil {
			return nil, err
		}

		resp, err := httprequest.VariationWERequester.HTTPDeleteVariation(parentID, int(id))
		if err != nil {
			return nil, err
		}

		respBytes, err = json.Marshal(resp)
		if err != nil {
			return nil, err
		}
	case "modification":
		var modifResourceLoader web_experimentation.ModificationResourceLoader
		err := json.Unmarshal(payloadBytes, &modifResourceLoader)
		if err != nil {
			return nil, fmt.Errorf("error occurred: %v", err)
		}

		if modifResourceLoader.CampaignID == 0 {
			return nil, fmt.Errorf("error occurred: missing property %s", "campaign_id")
		}

		resp, err := httprequest.ModificationRequester.HTTPDeleteModification(modifResourceLoader.CampaignID, int(id))
		if err != nil {
			return nil, fmt.Errorf("error occurred: %v", err)
		}

		respBytes, err = json.Marshal(resp)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown resource type: %s", res.Type)
	}

	var resp any
	err = json.Unmarshal(respBytes, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func ValidateResources(loadFile *LoadResFile, refCtx *RefContext) error {
	for _, res := range loadFile.Resources {

		if res.Type == "" {
			return fmt.Errorf("resource with $_ref=%s is missing 'type'", res.Ref)
		}
		if res.Action == "" {
			return fmt.Errorf("resource with $_ref=%s is missing 'action'", res.Ref)
		}

		/* 		if strings.HasPrefix(res.ParentID, "$") {
			parts := strings.Split(strings.TrimPrefix(res.ParentID, "$"), ".")
			if len(parts) < 2 {
				return fmt.Errorf("invalid reference format in $_parent_id for $_ref=%s", res.Ref)
			}

			if _, ok := refCtx.Get(parts[0]); !ok {
				return fmt.Errorf("reference %s not found for $_ref=%s", parts[0], res.Ref)
			}
		} */

		if len(res.Resources) > 0 {
			childFile := LoadResFile{Resources: res.Resources}
			if err := ValidateResources(&childFile, refCtx); err != nil {
				return err
			}
		}
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

	loadCmd.Flags().StringVarP(&resourceFile, "file", "", "", "resource file that contains your resource")

	if err := loadCmd.MarkFlagRequired("file"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	loadCmd.Flags().StringVarP(&outputFile, "output-file", "", "", "result of the command that contains all resource information")

	loadCmd.Flags().StringVarP(&inputRefRaw, "input-ref", "", "", "params to replace resource loader file")
	loadCmd.Flags().StringVarP(&inputRefFile, "input-ref-file", "", "", "file that contains params to replace resource loader file")

	ResourceCmd.AddCommand(loadCmd)
}
