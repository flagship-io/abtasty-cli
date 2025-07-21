/*
Copyright © 2022 Flagship Team flagship@abtasty.com

*/

package resource

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	models "github.com/flagship-io/abtasty-cli/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/flagship-io/abtasty-cli/utils/http_request/common"
	"github.com/spf13/cobra"
)

var (
	resourceFile    string
	outputFile      string
	inputParamsRaw  string
	inputParamsFile string
)

type ResourceLoaderModification struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Selector string `json:"selector"`
}

var inputParamsMap map[string]any

type ResourceData struct {
	Id string `json:"id"`
}

type ResourceType int

const (
	Campaign ResourceType = iota
	Audience
	FavoriteURL
	Modification
	Variation
)

type ResourceCmdStruct struct {
	Name      string `json:"name,omitempty"`
	ParentID  string `json:"parentId,omitempty"`
	Reference string `json:"ref,omitempty"`
	Response  string `json:"response,omitempty"`
	Action    string `json:"action,omitempty"`
	Error     string `json:"error,omitempty"`
}

type ResourceResp struct {
	Id string `json:"id,omitempty"`
}

// Resources is a slice of Resource.
type JsonResources []JsonResource

type JsonResource struct {
	Name           string          `json:"name"`
	ParentID       string          `json:"parentId"`
	Data           json.RawMessage `json:"data"`
	Reference      string          `json:"reference"`
	Action         string          `json:"action"`
	NestedResource JsonResources   `json:"nestedResource"`
}

type VariationData struct {
	*models.VariationWE
}

var cred common.RequestConfig

func Init(credL common.RequestConfig) {
	cred = credL
}

type ResourceAction string

const (
	ActionCreate ResourceAction = "create"
	ActionList   ResourceAction = "list"
	ActionDelete ResourceAction = "delete"
	ActionSwitch ResourceAction = "switch"
)

type Resource struct {
	Type      string         `json:"type"`
	Ref       string         `json:"$_ref,omitempty"`
	ParentID  string         `json:"$_parent_id,omitempty"`
	Action    ResourceAction `json:"action"`
	Payload   map[string]any `json:"payload,omitempty"`
	Resources []Resource     `json:"resources,omitempty"`
}

type LoadResFile struct {
	Version   int        `json:"version"`
	Resources []Resource `json:"resources"`
}

// Reference context for resolving $ref and $parent_id
type RefContext struct {
	refs map[string]any
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
func LoadResources(cmd *cobra.Command, filePath string, inputParamsFile string, inputParamsRaw string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read resource file: %w", err)
	}

	var loadFile LoadResFile
	if err := json.Unmarshal(data, &loadFile); err != nil {
		return fmt.Errorf("failed to parse resource file: %w", err)
	}

	refCtx := NewRefContext()

	// Load inputParams from file or raw
	var inputParams map[string]any
	if inputParamsFile != "" {
		b, err := os.ReadFile(inputParamsFile)
		if err != nil {
			return fmt.Errorf("failed to read input params file: %w", err)
		}
		_ = json.Unmarshal(b, &inputParams)
	} else if inputParamsRaw != "" {
		_ = json.Unmarshal([]byte(inputParamsRaw), &inputParams)
	}
	// Merge inputParams into refCtx
	for k, v := range inputParams {
		refCtx.Set(k, v)
	}

	for _, res := range loadFile.Resources {
		if err := processResource(cmd, res, refCtx); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Resource error: %v\n", err)
		}
	}
	return nil
}

// Process a resource and its children recursively
func processResource(cmd *cobra.Command, res Resource, rc *RefContext) error {
	if res.ParentID != "" && strings.HasPrefix(res.ParentID, "$") {
		parts := strings.Split(strings.TrimPrefix(res.ParentID, "$"), ".")
		if len(parts) > 1 {
			if refVal, ok := rc.Get(parts[0]); ok {
				if m, ok := refVal.(map[string]any); ok {
					if field, ok := m[parts[1]]; ok {
						res.ParentID = fmt.Sprintf("%v", field)
					}
				}
			}
		}
	}
	res.Payload = resolveRefs(res.Payload, rc).(map[string]any)

	var resp map[string]any
	var err error
	switch res.Action {
	case ActionCreate:
		resp, err = handleCreate(cmd, res)
	case ActionList:
		resp, err = handleList(cmd, res)
	case ActionDelete:
		err = handleDelete(cmd, res)
	case ActionSwitch:
		err = handleSwitch_(cmd, res)
	default:
		return fmt.Errorf("unsupported action: %s", res.Action)
	}
	if err != nil {
		return err
	}

	// Store response for $ref
	if res.Ref != "" && resp != nil {
		rc.Set(res.Ref, resp)
	}

	// Recursively process children
	for _, child := range res.Resources {
		// Set parent_id in child if not set
		if child.ParentID == "" && res.Ref != "" && resp != nil {
			if id, ok := resp["id"]; ok {
				child.ParentID = fmt.Sprintf("%v", id)
			}
		}
		if err := processResource(cmd, child, rc); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Child resource error: %v\n", err)
		}
	}
	return nil
}

// Example handlers (implement as needed)
func handleCreate(cmd *cobra.Command, res Resource) (map[string]any, error) {
	// Map type to actual creation logic
	payloadBytes, _ := json.Marshal(res.Payload)
	var respBytes []byte
	//var err error

	switch res.Type {
	case "campaign":
		//respBytes = campaign.CreateCampaign(payloadBytes)
		fmt.Fprintf(cmd.OutOrStdout(), "Create action for type: %s\n", res.Type)
		return nil, nil
	case "variation":
		// parentID required
		//parentID, _ := strconv.Atoi(res.ParentID)
		var v models.VariationWE
		_ = json.Unmarshal(payloadBytes, &v)
		//respBytes, err = http_request.VariationWERequester.HTTPCreateVariation(parentID, v)
		fmt.Fprintf(cmd.OutOrStdout(), "Create action for type: %s\n", res.Type)
		return nil, nil
	case "modification":
		// parentID required
		//parentID, _ := strconv.Atoi(res.ParentID)
		var m models.ModificationCodeCreateStruct
		_ = json.Unmarshal(payloadBytes, &m)
		//respBytes, err = http_request.ModificationRequester.HTTPCreateModification(parentID, m)
		fmt.Fprintf(cmd.OutOrStdout(), "Create action for type: %s\n", res.Type)
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown resource type: %s", res.Type)
	}
	/* 	if err != nil {
		return nil, err
	} */
	var resp map[string]any
	_ = json.Unmarshal(respBytes, &resp)
	return resp, nil
}

func handleList(cmd *cobra.Command, res Resource) (map[string]any, error) {
	// Implement as needed
	fmt.Fprintf(cmd.OutOrStdout(), "List action for type: %s\n", res.Type)
	return nil, nil
}

func handleDelete(cmd *cobra.Command, res Resource) error {
	// Implement as needed
	fmt.Fprintf(cmd.OutOrStdout(), "Delete action for type: %s\n", res.Type)
	return nil
}

func handleSwitch_(cmd *cobra.Command, res Resource) error {
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
		return LoadResources(cmd, resourceFile, inputParamsFile, inputParamsRaw)
	},
}

func init() {
	//cobra.OnInitialize(initResource)

	loadCmd.Flags().StringVarP(&resourceFile, "file", "", "", "resource file that contains your resource")

	if err := loadCmd.MarkFlagRequired("file"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	loadCmd.Flags().StringVarP(&outputFile, "output-file", "", "", "result of the command that contains all resource information")

	loadCmd.Flags().StringVarP(&inputParamsRaw, "input-params", "", "", "params to replace resource loader file")
	loadCmd.Flags().StringVarP(&inputParamsFile, "input-params-file", "", "", "file that contains params to replace resource loader file")

	ResourceCmd.AddCommand(loadCmd)
}
