/*
Copyright © 2022 Flagship Team flagship@abtasty.com

*/

package resource

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/flagship-io/abtasty-cli/cmd/feature-experimentation/campaign"
	"github.com/flagship-io/abtasty-cli/cmd/feature-experimentation/flag"
	"github.com/flagship-io/abtasty-cli/cmd/feature-experimentation/goal"
	"github.com/flagship-io/abtasty-cli/cmd/feature-experimentation/project"
	targeting_key "github.com/flagship-io/abtasty-cli/cmd/feature-experimentation/targeting-key"
	"github.com/flagship-io/abtasty-cli/cmd/feature-experimentation/variation"
	variation_group "github.com/flagship-io/abtasty-cli/cmd/feature-experimentation/variation-group"

	"github.com/flagship-io/abtasty-cli/models/feature_experimentation"
	"github.com/flagship-io/abtasty-cli/utils"
	"github.com/spf13/cobra"
)

var (
	resourceFile string
	outputFile   string
	inputRefRaw  string
	inputRefFile string
	dryRun       bool
)

type ResourceAction string

const (
	ActionCreate ResourceAction = "create"
	ActionEdit   ResourceAction = "edit"
	ActionList   ResourceAction = "list"
	ActionGet    ResourceAction = "get"
	ActionDelete ResourceAction = "delete"
)

const (
	Project        string = "project"
	Campaign       string = "campaign"
	VariationGroup string = "variation-group"
	Variation      string = "variation"
	Flag           string = "flag"
	Goal           string = "goal"
	TargetingKey   string = "targeting-key"
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

type ResourceData struct {
	Id string `json:"id"`
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

func resolveRefs(val any, rc *RefContext) any {
	switch v := val.(type) {
	case string:
		if strings.HasPrefix(v, "$") {
			parts := strings.Split(strings.TrimPrefix(v, "$"), ".")
			if len(parts) > 1 {
				if refVal, ok := rc.Get(parts[0]); ok {
					if m, ok := refVal.(map[string]interface{}); ok {
						if field, ok := m[parts[1]]; ok {
							return field
						}
					}
				}
			}
		}
		return v

	case []interface{}:
		for i, item := range v {
			v[i] = resolveRefs(item, rc)
		}
		return v

	case map[string]interface{}:
		for k, mapVal := range v {
			v[k] = resolveRefs(mapVal, rc)
		}
		return v

	default:
		return val
	}
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

	processAndRecord := func(cmd *cobra.Command, res Resource, rc *RefContext) {
		resp, err := processResourceWithResponse(cmd, res, rc)
		status := "success"
		if err != nil {
			status = "error"
			resp = err.Error()
		}
		recordResult(res.Ref, status, resp)

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

	if dryRun {
		fmt.Fprintf(cmd.OutOrStdout(), "Dry-run mode: resources validated, no changes applied.\n")
		return nil
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

	var projects, campaigns, variationGroups, variations, goals, targetingKeys, flags, others []Resource
	for _, res := range mutating {
		switch res.Type {
		case Project:
			projects = append(projects, res)
		case Campaign:
			campaigns = append(campaigns, res)
		case VariationGroup:
			variationGroups = append(variationGroups, res)
		case Variation:
			variations = append(variations, res)
		case Goal:
			goals = append(goals, res)
		case TargetingKey:
			targetingKeys = append(targetingKeys, res)
		case Flag:
			flags = append(flags, res)
		default:
			others = append(others, res)
		}
	}

	for _, res := range projects {
		processAndRecord(cmd, res, refCtx)
	}

	for _, res := range variationGroups {
		processAndRecord(cmd, res, refCtx)
	}

	for _, res := range campaigns {
		processAndRecord(cmd, res, refCtx)
	}

	for _, res := range variations {
		processAndRecord(cmd, res, refCtx)
	}

	for _, res := range goals {
		processAndRecord(cmd, res, refCtx)
	}

	for _, res := range flags {
		processAndRecord(cmd, res, refCtx)
	}

	for _, res := range targetingKeys {
		processAndRecord(cmd, res, refCtx)
	}

	for _, res := range others {
		processAndRecord(cmd, res, refCtx)
	}

	for _, res := range read {
		processAndRecord(cmd, res, refCtx)
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

					if field, ok := m[parts[1]].(string); ok {
						res.ParentID = fmt.Sprintf("%v", field)
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
	case ActionGet:
		resp, err = handleGet(res)
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

				if id, ok := resp.(map[string]any)["id"].(string); ok {
					child.ParentID = fmt.Sprintf("%v", id)
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

func handleCreate(res Resource) (resp map[string]any, err error) {
	payloadBytes, err := json.Marshal(res.Payload)
	if err != nil {
		return nil, err
	}

	var respBytes []byte

	switch res.Type {
	case Project:
		respBytes, err = project.CreateProject(payloadBytes)
		if err != nil {
			return nil, err
		}

	case Campaign:
		respBytes, err = campaign.CreateCampaign(payloadBytes)
		if err != nil {
			return nil, err
		}

	case VariationGroup:
		respBytes, err = variation_group.CreateVariationGroup(res.ParentID, payloadBytes)
		if err != nil {
			return nil, err
		}

	case Variation:
		respBytes, err = variation.CreateVariation(res.ParentResource.ParentID, res.ParentID, payloadBytes)
		if err != nil {
			return nil, err
		}

	case Goal:
		respBytes, err = goal.CreateGoal(payloadBytes)
		if err != nil {
			return nil, err
		}

	case TargetingKey:
		respBytes, err = targeting_key.CreateTargetingKey(payloadBytes)
		if err != nil {
			return nil, err
		}

	case Flag:
		respBytes, err = flag.CreateFlag(payloadBytes)
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unknown resource type: %s", res.Type)
	}

	err = json.Unmarshal(respBytes, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func handleEdit(res Resource) (resp map[string]any, err error) {
	var respBytes []byte

	payloadBytes, err := json.Marshal(res.Payload)
	if err != nil {
		return nil, err
	}

	var id = res.Payload["id"].(string)

	if id == "" {
		return nil, fmt.Errorf("error occurred: missing property %s", "id")
	}

	switch res.Type {
	case Project:
		respBytes, err = project.EditProject(id, payloadBytes)
		if err != nil {
			return nil, err
		}

	case Campaign:
		respBytes, err = campaign.EditCampaign(id, payloadBytes)
		if err != nil {
			return nil, err
		}

	case VariationGroup:
		respBytes, err = variation_group.EditVariationGroup(res.ParentID, id, payloadBytes)
		if err != nil {
			return nil, err
		}

	case Variation:
		respBytes, err = variation.EditVariation(res.ParentResource.ParentID, res.ParentID, id, payloadBytes)
		if err != nil {
			return nil, err
		}

	case Goal:
		respBytes, err = goal.EditGoal(id, payloadBytes)
		if err != nil {
			return nil, err
		}

	case TargetingKey:
		respBytes, err = targeting_key.EditTargetingKey(id, payloadBytes)
		if err != nil {
			return nil, err
		}

	case Flag:
		respBytes, err = flag.EditFlag(id, payloadBytes)
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unknown resource type: %s", res.Type)
	}

	err = json.Unmarshal(respBytes, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func handleList(res Resource) (any, error) {
	var respBytes []byte
	var err error

	switch res.Type {
	case Project:
		projects, err := project.ListProjects()
		if err != nil {
			return nil, err
		}

		respBytes, err = json.Marshal(projects)
		if err != nil {
			return nil, err
		}

	case Campaign:
		campaigns, err := campaign.ListCampaigns()
		if err != nil {
			return nil, err
		}

		respBytes, err = json.Marshal(campaigns)
		if err != nil {
			return nil, err
		}

	case VariationGroup:
		variationGroups, err := variation_group.ListVariationGroups(res.ParentID)
		if err != nil {
			return nil, err
		}

		respBytes, err = json.Marshal(variationGroups)
		if err != nil {
			return nil, err
		}

	case Variation:
		variations, err := variation.ListVariations(res.ParentResource.ParentID, res.ParentID)
		if err != nil {
			return nil, err
		}

		respBytes, err = json.Marshal(variations)
		if err != nil {
			return nil, err
		}

	case Goal:
		goals, err := goal.ListGoals()
		if err != nil {
			return nil, err
		}

		respBytes, err = json.Marshal(goals)
		if err != nil {
			return nil, err
		}

	case TargetingKey:
		targetingKeys, err := targeting_key.ListTargetingKeys()
		if err != nil {
			return nil, err
		}

		respBytes, err = json.Marshal(targetingKeys)
		if err != nil {
			return nil, err
		}

	case Flag:
		flags, err := flag.ListFlags()
		if err != nil {
			return nil, err
		}

		respBytes, err = json.Marshal(flags)
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

func handleGet(res Resource) (resp map[string]any, err error) {
	var respBytes []byte

	var id = res.Payload["id"].(string)

	if id == "" {
		return nil, fmt.Errorf("error occurred: missing property %s", "id")
	}

	switch res.Type {
	case Project:
		project, err := project.GetProject(id)
		if err != nil {
			return nil, err
		}

		respBytes, err = json.Marshal(project)
		if err != nil {
			return nil, err
		}

	case Campaign:
		campaign, err := campaign.GetCampaign(id)
		if err != nil {
			return nil, err
		}

		respBytes, err = json.Marshal(campaign)
		if err != nil {
			return nil, err
		}

	case VariationGroup:
		variationGroup, err := variation_group.GetVariationGroup(res.ParentID, id)
		if err != nil {
			return nil, err
		}

		respBytes, err = json.Marshal(variationGroup)
		if err != nil {
			return nil, err
		}

	case Variation:
		variations, err := variation.GetVariation(res.ParentResource.ParentID, res.ParentID, id)
		if err != nil {
			return nil, err
		}

		respBytes, err = json.Marshal(variations)
		if err != nil {
			return nil, err
		}

	case Goal:
		goal, err := goal.GetGoal(id)
		if err != nil {
			return nil, err
		}

		respBytes, err = json.Marshal(goal)
		if err != nil {
			return nil, err
		}

	case TargetingKey:
		targetingKey, err := targeting_key.GetTargetingKey(id)
		if err != nil {
			return nil, err
		}

		respBytes, err = json.Marshal(targetingKey)
		if err != nil {
			return nil, err
		}

	case Flag:
		flag, err := flag.GetFlag(id)
		if err != nil {
			return nil, err
		}

		respBytes, err = json.Marshal(flag)
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unknown resource type: %s", res.Type)
	}

	err = json.Unmarshal(respBytes, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func handleDelete(res Resource) (any, error) {
	var respBytes []byte
	var id = res.Payload["id"].(string)

	if id == "" {
		return nil, fmt.Errorf("error occurred: missing property %s", "id")
	}

	switch res.Type {
	case Project:
		resp, err := project.DeleteProject(id)
		if err != nil {
			return nil, err
		}

		respBytes, err = json.Marshal(resp)
		if err != nil {
			return nil, err
		}

	case Campaign:
		resp, err := campaign.DeleteCampaign(id)
		if err != nil {
			return nil, err
		}

		respBytes, err = json.Marshal(resp)
		if err != nil {
			return nil, err
		}

	case VariationGroup:
		resp, err := variation_group.DeleteVariationGroup(res.ParentID, id)
		if err != nil {
			return nil, err
		}

		respBytes, err = json.Marshal(resp)
		if err != nil {
			return nil, err
		}

	case Variation:
		resp, err := variation.DeleteVariation(res.ParentResource.ParentID, res.ParentID, id)
		if err != nil {
			return nil, err
		}

		respBytes, err = json.Marshal(resp)
		if err != nil {
			return nil, err
		}

	case Goal:
		resp, err := goal.DeleteGoal(id)
		if err != nil {
			return nil, err
		}

		respBytes, err = json.Marshal(resp)
		if err != nil {
			return nil, err
		}

	case TargetingKey:
		resp, err := targeting_key.DeleteTargetingKey(id)
		if err != nil {
			return nil, err
		}

		respBytes, err = json.Marshal(resp)
		if err != nil {
			return nil, err
		}

	case Flag:
		resp, err := flag.DeleteFlag(id)
		if err != nil {
			return nil, err
		}

		respBytes, err = json.Marshal(resp)
		if err != nil {
			return nil, err
		}

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

func ValidateResources(loadFile *LoadResFile, refCtx *RefContext) error {
	for _, res := range loadFile.Resources {
		if res.Ref == "" && res.Type == "" {
			b, err := json.Marshal(res)
			if err != nil {
				return fmt.Errorf("error occurred: %v", err)
			}

			return fmt.Errorf("resource: %s is missing '$_ref' and 'type'", string(b))
		}

		if res.Ref == "" {
			return fmt.Errorf("resource with $type: %s is missing '$_ref'", res.Type)
		}

		if res.Type == "" {
			return fmt.Errorf("resource with $_ref: %s is missing 'type'", res.Ref)
		}

		if res.Type != Project && res.Type != Campaign && res.Type != Variation && res.Type != VariationGroup && res.Type != Flag && res.Type != TargetingKey && res.Type != Goal {
			return fmt.Errorf("resource with $_ref: %s has unknown type: %s, only %s, %s, %s, %s, %s, %s, %s are allowed", res.Ref, res.Type, Project, Campaign, VariationGroup, Variation, TargetingKey, Flag, Goal)
		}

		if res.Action == "" {
			return fmt.Errorf("resource with $_ref: %s is missing 'action'", res.Ref)
		}

		if res.Action != ActionCreate && res.Action != ActionEdit && res.Action != ActionGet && res.Action != ActionList && res.Action != ActionDelete {
			return fmt.Errorf("resource with $_ref: %s has unknown action: %s, only %s, %s, %s, %s and %s are allowed ", res.Ref, res.Action, ActionCreate, ActionEdit, ActionGet, ActionList, ActionDelete)
		}

		/* 		if len(res.Resources) != 0 {
			for _, subRes := range res.Resources {
				if res.Type == Campaign && subRes.Type != Variation {
					return fmt.Errorf("resource %s with $_ref: %s can only accept sub resource type %s", res.Type, res.Ref, Variation)
				}

				if res.Type == Variation && subRes.Type != Modification {
					return fmt.Errorf("resource %s with $_ref: %s can only accept sub resource type %s", res.Type, res.Ref, Modification)
				}
			}
		} */

		refCtx.Set(res.Ref, res.Payload)
		var resPayloadDeepCopy map[string]any
		err := utils.DeepCopyMap(res.Payload, &resPayloadDeepCopy)
		if err != nil {
			return err
		}

		payloadToValidate := preprocessPayloadForValidation(resPayloadDeepCopy, res.Type)
		payloadToValidateBytes, err := json.Marshal(payloadToValidate)
		if err != nil {
			return err
		}

		dec := json.NewDecoder(bytes.NewReader(payloadToValidateBytes))
		dec.DisallowUnknownFields()

		switch res.Type {
		case Project:
			var projectModel feature_experimentation.Project
			if err := dec.Decode(&projectModel); err != nil {
				return fmt.Errorf("%v in %s", err, res.Type)
			}

		case Campaign:
			var campaignModel feature_experimentation.CampaignFE
			if err := dec.Decode(&campaignModel); err != nil {
				return fmt.Errorf("%v in %s", err, res.Type)
			}

		case VariationGroup:
			var variationGroupModel feature_experimentation.VariationGroup
			if err := dec.Decode(&variationGroupModel); err != nil {
				return fmt.Errorf("%v in %s", err, res.Type)
			}

		case Variation:
			var variationModel feature_experimentation.VariationFE
			if err := dec.Decode(&variationModel); err != nil {
				return fmt.Errorf("%v in %s", err, res.Type)
			}

		case Flag:
			var featureFlagModel feature_experimentation.Flag
			if err := dec.Decode(&featureFlagModel); err != nil {
				return fmt.Errorf("%v in %s", err, res.Type)
			}

		case TargetingKey:
			var targetingKeyModel feature_experimentation.TargetingKey
			if err := dec.Decode(&targetingKeyModel); err != nil {
				return fmt.Errorf("%v in %s", err, res.Type)
			}

		case Goal:
			var goalModel feature_experimentation.Goal
			if err := dec.Decode(&goalModel); err != nil {
				return fmt.Errorf("%v in %s", err, res.Type)
			}

		default:
			return fmt.Errorf("unknown resource type: %s", res.Type)
		}

		if res.Action == ActionDelete || res.Action == ActionEdit || res.Action == ActionGet {
			var id = res.Payload["id"].(string)

			if id == "" {
				return fmt.Errorf("error occurred: missing property %s", "id")
			}
		}
	}

	for _, res := range loadFile.Resources {
		if strings.HasPrefix(res.ParentID, "$") {
			parts := strings.Split(strings.TrimPrefix(res.ParentID, "$"), ".")
			if len(parts) < 2 {
				return fmt.Errorf("invalid reference format in $_parent_id for $_ref=%s", res.Ref)
			}

			if _, ok := refCtx.Get(parts[0]); !ok {
				return fmt.Errorf("reference %s not found for $_ref=%s", parts[0], res.Ref)
			}
		}

		for k, v := range res.Payload {
			path := k
			if err := validateReferences(v, refCtx, path, res.Ref); err != nil {
				return err
			}
		}

		if len(res.Resources) > 0 {
			childFile := LoadResFile{Resources: res.Resources}
			if err := ValidateResources(&childFile, refCtx); err != nil {
				return err
			}
		}
	}

	return nil
}

func preprocessPayloadForValidation(payload map[string]any, structType string) map[string]any {
	/* for k, v := range payload {
		switch vv := v.(type) {
		case string:
			if strings.HasPrefix(vv, "$") {
				switch structType {
				case Modification:
					if k == "campaign_id" {
						payload[k] = 0
					}
				case Campaign:
					if k == "folder_id" {
						payload[k] = 0
					}
				case Variation:
					if k == "traffic" {
						payload[k] = 0
					}
				default:
					payload[k] = ""
				}
			}
		case map[string]any:
			payload[k] = preprocessPayloadForValidation(vv, structType)
		case []any:
			for i, item := range vv {
				if m, ok := item.(map[string]any); ok {
					vv[i] = preprocessPayloadForValidation(m, structType)
				}
			}
			payload[k] = vv
		}
	} */

	return payload
}

// LoadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:   "load [--file=<file>]",
	Short: "Load your resources",
	Long:  `Load your resources`,
	Run: func(cmd *cobra.Command, args []string) {
		err := LoadResources(cmd, resourceFile, inputRefFile, inputRefRaw, outputFile)
		if err != nil {
			log.Fatalf("%v", err)
		}
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
	loadCmd.Flags().BoolVarP(&dryRun, "dry-run", "", false, "perform dry run to validate resources without load resource to the API")

	ResourceCmd.AddCommand(loadCmd)
}

func validateReferences(value interface{}, refCtx *RefContext, path string, resRef string) error {
	switch v := value.(type) {
	case string:
		if strings.HasPrefix(v, "$") {
			parts := strings.Split(strings.TrimPrefix(v, "$"), ".")
			if len(parts) < 2 {
				return fmt.Errorf("invalid $ reference format at %s in $_ref=%s", path, resRef)
			}

			if _, ok := refCtx.Get(parts[0]); !ok {
				return fmt.Errorf("reference %s not found at %s in $_ref=%s", parts[0], path, resRef)
			}
		}

	case []interface{}:
		for i, item := range v {
			itemPath := fmt.Sprintf("%s[%d]", path, i)
			if err := validateReferences(item, refCtx, itemPath, resRef); err != nil {
				return err
			}
		}

	case map[string]interface{}:
		for k, val := range v {
			nestedPath := fmt.Sprintf("%s.%s", path, k)
			if err := validateReferences(val, refCtx, nestedPath, resRef); err != nil {
				return err
			}
		}
	}

	return nil
}
