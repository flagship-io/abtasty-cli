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
	"strconv"
	"strings"

	"github.com/flagship-io/abtasty-cli/cmd/web-experimentation/audience"
	"github.com/flagship-io/abtasty-cli/cmd/web-experimentation/campaign"
	"github.com/flagship-io/abtasty-cli/cmd/web-experimentation/folder"
	"github.com/flagship-io/abtasty-cli/cmd/web-experimentation/modification"
	"github.com/flagship-io/abtasty-cli/cmd/web-experimentation/variation"
	"github.com/flagship-io/abtasty-cli/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/utils"
	"github.com/flagship-io/abtasty-cli/utils/common"
	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/spf13/cobra"
)

const (
	Folder       string = "folder"
	Campaign     string = "campaign"
	Variation    string = "variation"
	Modification string = "modification"
	Audience     string = "audience"
)

type funcModifType func(int, []byte) ([]byte, error)

func LoadResources(cmd *cobra.Command, filePath, inputRefFile, inputRefRaw, outputFile string) error {

	var results []common.ResourceResult

	recordResult := func(ref string, status string, resp interface{}) {
		if ref != "" {
			results = append(results, common.ResourceResult{
				Ref:      ref,
				Status:   status,
				Response: resp,
			})
		}
	}

	processAndRecord := func(cmd *cobra.Command, res common.Resource, rc *common.RefContext) {
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

	var loadFile common.LoadResFile
	if err := json.Unmarshal(data, &loadFile); err != nil {
		return fmt.Errorf("failed to parse resource file: %w", err)
	}

	refCtx := common.NewRefContext()

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

	if common.DryRun {
		fmt.Fprintf(cmd.OutOrStdout(), "Dry-run mode: resources validated, no changes applied.\n")
		return nil
	}

	var mutating, read []common.Resource
	for _, res := range loadFile.Resources {
		switch res.Action {
		case common.ActionGet, common.ActionList:
			read = append(read, res)
		default:
			mutating = append(mutating, res)
		}
	}

	var audiences, folders, campaigns, variations, modifications, others []common.Resource
	for _, res := range mutating {
		switch res.Type {
		case Audience:
			audiences = append(audiences, res)
		case Folder:
			folders = append(folders, res)
		case Campaign:
			campaigns = append(campaigns, res)
		case Variation:
			variations = append(variations, res)
		case Modification:
			modifications = append(modifications, res)
		default:
			others = append(others, res)
		}
	}

	for _, res := range folders {
		processAndRecord(cmd, res, refCtx)
	}

	for _, res := range audiences {
		processAndRecord(cmd, res, refCtx)
	}

	for _, res := range campaigns {
		processAndRecord(cmd, res, refCtx)
	}

	for _, res := range variations {
		processAndRecord(cmd, res, refCtx)
	}

	for _, res := range modifications {
		processAndRecord(cmd, res, refCtx)
	}

	for _, res := range others {
		processAndRecord(cmd, res, refCtx)
	}

	for _, res := range read {
		processAndRecord(cmd, res, refCtx)
	}

	loaderResults := common.LoaderResults{Results: results}
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

func processResourceWithResponse(cmd *cobra.Command, res common.Resource, rc *common.RefContext) (interface{}, error) {
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
	res.Payload = common.ResolveRefs(res.Payload, rc).(map[string]any)

	var resp any
	var err error
	switch res.Action {
	case common.ActionCreate:
		resp, err = handleCreate(res)
	case common.ActionEdit:
		resp, err = handleEdit(res)
	case common.ActionList:
		resp, err = handleList(res)
	case common.ActionGet:
		resp, err = handleGet(res)
	case common.ActionDelete:
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

	if res.Action == common.ActionCreate || res.Action == common.ActionEdit {
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

func handleCreate(res common.Resource) (resp map[string]any, err error) {
	payloadBytes, err := json.Marshal(res.Payload)
	if err != nil {
		return nil, err
	}

	var respBytes []byte

	switch res.Type {
	case Folder:
		respBytes, err = folder.CreateFolder(payloadBytes)
		if err != nil {
			return nil, err
		}
	case Campaign:
		respBytes, err = campaign.CreateCampaign(payloadBytes)
		if err != nil {
			return nil, err
		}
	case Audience:
		respBytes, err = audience.CreateAudience(payloadBytes)
		if err != nil {
			return nil, err
		}
	case Variation:
		parentID, err := strconv.Atoi(res.ParentID)
		if err != nil {
			return nil, err
		}

		respBytes, err = variation.CreateVariation(parentID, payloadBytes)
		if err != nil {
			return nil, err
		}
	case Modification:
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

	err = json.Unmarshal(respBytes, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func handleEdit(res common.Resource) (resp map[string]any, err error) {
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
	case Folder:
		respBytes, err = folder.EditFolder(int(id), payloadBytes)
		if err != nil {
			return nil, err
		}
	case Campaign:
		respBytes, err = campaign.EditCampaign(int(id), payloadBytes)
		if err != nil {
			return nil, err
		}
	case Variation:
		parentID, err := strconv.Atoi(res.ParentID)
		if err != nil {
			return nil, err
		}

		respBytes, err = variation.EditVariation(int(id), parentID, payloadBytes)
		if err != nil {
			return nil, err
		}
	case Modification:
		respBytes, err = createOrEditModification(int(id), res, payloadBytes, modification.EditModification)
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

func handleList(res common.Resource) (any, error) {
	var respBytes []byte
	var err error

	switch res.Type {
	case Folder:
		folderList, err := folder.ListFolder()
		if err != nil {
			return nil, err
		}

		respBytes, err = json.Marshal(folderList)
		if err != nil {
			return nil, err
		}
	case Campaign:
		campaignList, err := campaign.ListCampaigns()
		if err != nil {
			return nil, err
		}

		respBytes, err = json.Marshal(campaignList)
		if err != nil {
			return nil, err
		}
	case Variation:
		var campaignID int
		if res.ParentResource != nil {
			campaignIDString := res.ParentResource.ParentID
			campaignIDInt, err := strconv.Atoi(campaignIDString)
			if err != nil {
				return nil, err
			}
			campaignID = campaignIDInt
		}

		variationList, err := variation.ListVariations(campaignID)
		if err != nil {
			return nil, err
		}

		respBytes, err = json.Marshal(variationList)
		if err != nil {
			return nil, err
		}
	case Modification:
		payloadBytes, err := json.Marshal(res.Payload)
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

func handleGet(res common.Resource) (resp map[string]any, err error) {
	var respBytes []byte

	var id = res.Payload["id"].(float64)

	if id == 0 {
		return nil, fmt.Errorf("error occurred: missing property %s", "id")
	}

	switch res.Type {
	case Folder:
		folder, err := folder.GetFolder(int(id))
		if err != nil {
			return nil, err
		}

		respBytes, err = json.Marshal(folder)
		if err != nil {
			return nil, err
		}
	case Campaign:
		campaign, err := campaign.GetCampaign(int(id))
		if err != nil {
			return nil, err
		}

		respBytes, err = json.Marshal(campaign)
		if err != nil {
			return nil, err
		}
	case Variation:
		parentID, err := strconv.Atoi(res.ParentID)
		if err != nil {
			return nil, err
		}

		variation, err := variation.GetVariation(parentID, int(id))
		if err != nil {
			return nil, err
		}

		respBytes, err = json.Marshal(variation)
		if err != nil {
			return nil, err
		}
	case Modification:
		payloadBytes, err := json.Marshal(res.Payload)
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

		modification, err := modification.GetModification(modificationResourceLoader.CampaignID, int(id))
		if err != nil {
			return nil, err
		}

		respBytes, err = json.Marshal(modification)
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

func handleDelete(res common.Resource) (any, error) {
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
	case Folder:
		resp, err := httprequest.FolderRequester.HTTPDeleteFolder(int(id))
		if err != nil {
			return nil, err
		}

		respBytes, err = json.Marshal(resp)
		if err != nil {
			return nil, err
		}
	case Campaign:
		resp, err := httprequest.CampaignWERequester.HTTPDeleteCampaign(int(id))
		if err != nil {
			return nil, err
		}

		respBytes, err = json.Marshal(resp)
		if err != nil {
			return nil, err
		}
	case Variation:
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
	case Modification:
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

func ValidateResources(loadFile *common.LoadResFile, refCtx *common.RefContext) error {
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

		if res.Type != Campaign && res.Type != Variation && res.Type != Modification && res.Type != Folder && res.Type != Audience {
			return fmt.Errorf("resource with $_ref: %s has unknown type: %s, only %s, %s, %s, %s, %s are allowed", res.Ref, res.Type, Folder, Campaign, Variation, Modification, Audience)
		}

		if res.Action == "" {
			return fmt.Errorf("resource with $_ref: %s is missing 'action'", res.Ref)
		}

		if res.Action != common.ActionCreate && res.Action != common.ActionEdit && res.Action != common.ActionGet && res.Action != common.ActionList && res.Action != common.ActionDelete {
			return fmt.Errorf("resource with $_ref: %s has unknown action: %s, only %s, %s, %s, %s and %s are allowed ", res.Ref, res.Action, common.ActionCreate, common.ActionEdit, common.ActionGet, common.ActionList, common.ActionDelete)
		}

		if len(res.Resources) != 0 {
			for _, subRes := range res.Resources {
				if res.Type == Campaign && subRes.Type != Variation {
					return fmt.Errorf("resource %s with $_ref: %s can only accept sub resource type %s", res.Type, res.Ref, Variation)
				}

				if res.Type == Variation && subRes.Type != Modification {
					return fmt.Errorf("resource %s with $_ref: %s can only accept sub resource type %s", res.Type, res.Ref, Modification)
				}
			}
		}

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
		case Audience:
			var audienceModel web_experimentation.AudiencePayload
			if err := dec.Decode(&audienceModel); err != nil {
				return fmt.Errorf("%v in %s", err, res.Type)
			}
		case Folder:
			var folderModel web_experimentation.Folder
			if err := dec.Decode(&folderModel); err != nil {
				return fmt.Errorf("%v in %s", err, res.Type)
			}
		case Campaign:
			var campaignModel web_experimentation.CampaignWEResourceLoader
			if err := dec.Decode(&campaignModel); err != nil {
				return fmt.Errorf("%v in %s", err, res.Type)
			}
		case Variation:
			var variationModel web_experimentation.VariationResourceLoader
			if err := dec.Decode(&variationModel); err != nil {
				return fmt.Errorf("%v in %s", err, res.Type)
			}
		case Modification:
			var modificationModel web_experimentation.ModificationResourceLoader
			if err := dec.Decode(&modificationModel); err != nil {
				return fmt.Errorf("%v in %s", err, res.Type)
			}
		default:
			return fmt.Errorf("unknown resource type: %s", res.Type)
		}

		if res.Action == common.ActionDelete || res.Action == common.ActionEdit || res.Action == common.ActionGet {
			var id = res.Payload["id"].(float64)

			if id == 0 {
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
			if err := common.ValidateReferences(v, refCtx, path, res.Ref); err != nil {
				return err
			}
		}

		if len(res.Resources) > 0 {
			childFile := common.LoadResFile{Resources: res.Resources}
			if err := ValidateResources(&childFile, refCtx); err != nil {
				return err
			}
		}
	}

	return nil
}

func preprocessPayloadForValidation(payload map[string]any, structType string) map[string]any {
	for k, v := range payload {
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
	}

	return payload
}

// LoadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:   "load [--file=<file>]",
	Short: "Load your resources",
	Run: func(cmd *cobra.Command, args []string) {
		err := LoadResources(cmd, common.ResourceFile, common.InputRefFile, common.InputRefRaw, common.OutputFile)
		if err != nil {
			log.Fatalf("%v", err)
		}
	},
}

func init() {
	loadCmd.Flags().StringVarP(&common.ResourceFile, "file", "", "", "resource file that contains your resource")

	if err := loadCmd.MarkFlagRequired("file"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	loadCmd.Flags().StringVarP(&common.OutputFile, "output-file", "", "", "result of the command that contains all resource information")
	loadCmd.Flags().StringVarP(&common.InputRefRaw, "input-ref", "", "", "params to replace resource loader file")
	loadCmd.Flags().StringVarP(&common.InputRefFile, "input-ref-file", "", "", "file that contains params to replace resource loader file")
	loadCmd.Flags().BoolVarP(&common.DryRun, "dry-run", "", false, "perform dry run to validate resources without load resource to the API")

	ResourceCmd.AddCommand(loadCmd)
}

func createOrEditModification(id int, res common.Resource, payloadBytes []byte, createOrEditModif funcModifType) ([]byte, error) {
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
