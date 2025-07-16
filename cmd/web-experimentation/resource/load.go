/*
Copyright © 2022 Flagship Team flagship@abtasty.com

*/

package resource

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/d5/tengo/v2"
	"github.com/flagship-io/abtasty-cli/cmd/web-experimentation/audience"
	"github.com/flagship-io/abtasty-cli/cmd/web-experimentation/campaign"
	favorite_url "github.com/flagship-io/abtasty-cli/cmd/web-experimentation/favorite-url"
	models "github.com/flagship-io/abtasty-cli/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/utils"
	"github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/flagship-io/abtasty-cli/utils/http_request/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	resourceFile    string
	outputFile      string
	inputParams     string
	inputParamsFile string
)

var inputParamsMap map[string]any

type Data interface {
	getName() string
	Save(data []byte) ([]byte, error)
	Delete(id string) error
}

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

var resourceTypeMap = map[string]ResourceType{
	"audience":     Audience,
	"favorite_url": FavoriteURL,
	"modification": Modification,
	"campaign":     Campaign,
	"variation":    Variation,
}

type Resource struct {
	Name            ResourceType
	ParentID        string
	Data            Data
	Reference       string
	Action          string
	NestedResources []Resource
}

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

type CampaignData struct {
	*models.CampaignWEResourceLoader
}

// getName implements Data.
func (f *CampaignData) getName() string {
	return "Campaign"
}

func (f *CampaignData) Save(data []byte) ([]byte, error) {
	return campaign.CreateCampaign(data), nil
}

func (f *CampaignData) Delete(id string) error {
	return http_request.CampaignWERequester.HTTPDeleteCampaign(id)
}

func (f *CampaignData) Switch(id, state string) error {
	return http_request.CampaignWERequester.HTTPSwitchStateCampaign(id, state)
}

type AudienceData struct {
	*models.AudienceResourceLoader
}

// getName implements Data.
func (f *AudienceData) getName() string {
	return "Audience"
}

func (f *AudienceData) Save(data []byte) ([]byte, error) {
	return audience.CreateAudience(data), nil
}

// Can't delete audience for this moment
func (f *AudienceData) Delete(id string) error {
	return nil
}

type FavoriteUrlData struct {
	*models.FavoriteURL
}

// Can't create favorite URL for this moment
func (f *FavoriteUrlData) Save(data []byte) ([]byte, error) {
	return favorite_url.CreateFavoriteURL(data), nil
}

// Can't delete favorite URL for this moment
func (f *FavoriteUrlData) Delete(id string) error {
	return nil
}

/* type ModificationData struct {
	*models.Modification
}

// Delete implements Data.
func (f *ModificationData) Delete(id string) error {
	panic("unimplemented")
}

// Save implements Data.
func (f *ModificationData) Save(data string) ([]byte, error) {
	panic("unimplemented")
}

func (f *ModificationData) SaveWithParent(campaignId, data string) ([]byte, error) {
	var modelData models.ModificationCodeCreateStruct
	campaignIdInt, err := strconv.Atoi(campaignId)
	if err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	err = json.Unmarshal([]byte(data), &modelData)

	if err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	return http_request.ModificationRequester.HTTPCreateModification(campaignIdInt, modelData)
}

func (f *ModificationData) DeleteWithParent(campaignId, id string) error {
	campaignIdInt, err := strconv.Atoi(campaignId)
	if err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.Fatalf("error occurred: %v", err)
	}
	return http_request.ModificationRequester.HTTPDeleteModification(campaignIdInt, idInt)
} */

type VariationData struct {
	*models.VariationWE
}

func (f *VariationData) Save(campaignId int, data models.VariationWE) ([]byte, error) {
	return http_request.VariationWERequester.HTTPCreateVariation(campaignId, data)
}

func (f *VariationData) Delete(campaignId, id int) error {
	return http_request.VariationWERequester.HTTPDeleteVariation(campaignId, id)
}

func makeResourcesFromJSON(jsonResources JsonResources) ([]Resource, error) {
	var resources []Resource
	for _, r := range jsonResources {
		name, ok := resourceTypeMap[r.Name]
		if !ok {
			return nil, fmt.Errorf("invalid resource name: %s", r.Name)
		}

		var data Data = nil
		var err error = nil
		var nestedResource []Resource = nil

		switch name {

		case Audience:
			audienceData := AudienceData{}
			err = json.Unmarshal(r.Data, &audienceData)
			data = &audienceData

		/* case Modification:
		modificationData := ModificationData{}
		err = json.Unmarshal(r.Data, &modificationData)
		data = &modificationData */

		case Campaign:
			campaignData := CampaignData{}
			err = json.Unmarshal(r.Data, &campaignData)
			data = &campaignData
			nestedResource, err = makeResourcesFromJSON(r.NestedResource)
		}

		if err != nil {
			return nil, err
		}

		resources = append(resources, Resource{Name: name, Data: data, Reference: r.Reference, Action: r.Action, NestedResources: nestedResource})
	}

	return resources, nil
}

func resolveVariables(data any, resourceReferences map[string]any) any {
	switch val := data.(type) {
	case map[string]any:
		for k, v := range val {
			val[k] = resolveVariables(v, resourceReferences)
		}

	case []any:
		for i, v := range val {
			val[i] = resolveVariables(v, resourceReferences)
		}

	case string:
		if strings.Contains(val, "$") {
			vTrim := strings.Trim(val, "$")
			for k_, variable := range resourceReferences {
				script, _ := tengo.Eval(context.Background(), vTrim, map[string]any{
					k_: variable,
				})
				if script == nil {
					continue
				}
				// Update the string value with the result
				if resultStr, ok := script.(string); ok {
					return resultStr
				}
			}
		}
	}
	return data
}

var cred common.RequestConfig

func Init(credL common.RequestConfig) {
	cred = credL
}

func ExtractResourcesFromFile(filePath string) ([]Resource, error) {
	var config struct {
		Resources JsonResources
	}

	bytes, err := os.ReadFile(resourceFile)

	if err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	if err := json.Unmarshal(bytes, &config); err != nil {
		return nil, err
	}

	return makeResourcesFromJSON(config.Resources)
}

var gResources []Resource

// LoadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:   "load [--file=<file>]",
	Short: "Load your resources",
	Long:  `Load your resources`,
	Run: func(cmd *cobra.Command, args []string) {
		var params []byte

		if resourceFile == "" && !utils.CheckSingleFlag(inputParams != "", inputParamsFile != "") {
			log.Fatalf("error occurred: %s", "1 flag is required. (input-params, input-params-file)")
		}

		if inputParams != "" {
			params = []byte(inputParams)

		}

		if inputParamsFile != "" {
			fileContent, err := os.ReadFile(inputParamsFile)
			if err != nil {
				log.Fatalf("error occurred: %s", err)
			}

			params = fileContent
		}

		if params != nil {
			err := json.Unmarshal(params, &inputParamsMap)
			if err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "Error: %s", err)
				return
			}
		}

		jsonBytes := ScriptResource(cmd, gResources, inputParamsMap)
		if outputFile != "" {
			os.WriteFile(outputFile, jsonBytes, os.ModePerm)
			fmt.Fprintf(cmd.OutOrStdout(), "File created at %s\n", outputFile)
			return
		}

		if viper.GetString("output_format") == "json" {
			fmt.Fprintf(cmd.OutOrStdout(), "%s", string(jsonBytes))
		}

	},
}

func init() {
	cobra.OnInitialize(initResource)

	loadCmd.Flags().StringVarP(&resourceFile, "file", "", "", "resource file that contains your resource")

	if err := loadCmd.MarkFlagRequired("file"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	loadCmd.Flags().StringVarP(&outputFile, "output-file", "", "", "result of the command that contains all resource information")

	loadCmd.Flags().StringVarP(&inputParams, "input-params", "", "", "params to replace resource loader file")
	loadCmd.Flags().StringVarP(&inputParamsFile, "input-params-file", "", "", "file that contains params to replace resource loader file")

	ResourceCmd.AddCommand(loadCmd)
}

func initResource() {

	// Use config file from the flag.
	var err error
	if resourceFile != "" {
		gResources, err = ExtractResourcesFromFile(resourceFile)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
	}
}

/*
func ScriptResource(cmd *cobra.Command, resources []Resource, inputParamsMap map[string]any) []byte {

	resourceReferences := make(map[string]any)
	var loadResultJSON []string
	var loadResultOutputFile []ResourceCmdStruct

	for _, resource := range resources {
		var response []byte
		var resultOutputFile ResourceCmdStruct
		var resourceData map[string]any
		var responseData any

		var resourceName = resource.Data.getName()
		const color = "\033[0;33m"
		const colorNone = "\033[0m"

		data, err := json.Marshal(resource.Data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error occurred marshal data: %v\n", err)
		}

		var httpMethod string = "POST"

		if resource.Action == "delete" {
			httpMethod = "DELETE"
		}

		if resource.Action == "switch" {
			httpMethod = "PATCH"
		}

		err = json.Unmarshal(data, &resourceData)

		if err != nil {
			fmt.Fprintf(os.Stderr, "error occurred unmarshal resourceData: %v\n", err)
		}

		if inputParamsMap != nil {
			for k, vInterface := range resourceData {
				v, ok := vInterface.(string)
				if ok {
					if strings.Contains(v, "$") {
						vTrim := strings.Trim(v, "$")
						vTrimL := strings.Split(vTrim, ".")
						value, err := getNestedValue(inputParamsMap, vTrimL)
						if err != nil {
							fmt.Fprintf(os.Stderr, "Error: %s\n", err)
						}

						if value != nil {
							resourceData[k] = value
						}
					}
				}
			}
		}

		resourceData = resolveVariables(resourceData, resourceReferences).(map[string]any)

		dataResource, err := json.Marshal(resourceData)
		if err != nil {
			log.Fatalf("error occurred http call: %v\n", err)
		}

		if resource.Action == "switch" {
			if resource.Name == Campaign {
				fmt.Println("--id="+fmt.Sprintf("%v", resourceData["id"]), "--state="+fmt.Sprintf("%v", resourceData["state"]))
				if !(fmt.Sprintf("%v", resourceData["state"]) == "active" || fmt.Sprintf("%v", resourceData["state"]) == "paused" || fmt.Sprintf("%v", resourceData["state"]) == "interrupted") {
					fmt.Fprintln(cmd.OutOrStdout(), "Status can only have 3 values: active or paused or interrupted")
				} else {
					err := http_request.CampaignFERequester.HTTPSwitchStateCampaign(fmt.Sprintf("%v", resourceData["id"]), fmt.Sprintf("%v", resourceData["state"]))
					if err != nil {
						log.Fatalf("error occurred: %v", err)
					}
					fmt.Fprintf(cmd.OutOrStdout(), "campaign status set to %s\n", fmt.Sprintf("%v", resourceData["state"]))
				}
			}
		}

		if httpMethod == "POST" {
			response, err = resource.Data.Save(dataResource)

			resultOutputFile = ResourceCmdStruct{
				Name:      resourceName,
				Response:  string(response),
				Reference: resource.Reference,
				Action:    httpMethod,
			}

			if err != nil {
				resultOutputFile.Error = err.Error()
			}

			loadResultOutputFile = append(loadResultOutputFile, resultOutputFile)
		}

		if httpMethod == "DELETE" {
			err = resource.Data.Delete(fmt.Sprintf("%s", resourceData["id"]))

			if err == nil && viper.GetString("output_format") != "json" {
				response = []byte("The id: " + fmt.Sprintf("%v", resourceData["id"]) + " deleted successfully")
			}
		}

		if err != nil {
			fmt.Fprintf(cmd.OutOrStdout(), "%s - %s: %s %s\n", color, resourceName, colorNone, err.Error())
			continue
		}

		if viper.GetString("output_format") != "json" {
			fmt.Fprintf(cmd.OutOrStdout(), "%s - %s: %s %s\n", color, resourceName, colorNone, string(response))
		}

		if httpMethod != "DELETE" && httpMethod != "PATCH" {
			err = json.Unmarshal(response, &responseData)

			if err != nil {
				fmt.Fprintf(os.Stderr, "error occurred unmarshal responseData: %v\n", err)
			}

			if responseData == nil {
				fmt.Fprintf(os.Stderr, "error occurred not response data: %s\n", string(response))
				continue
			}

			resourceReferences[resource.Reference] = responseData
		}

		loadResultJSON = append(loadResultJSON, string(response))
	}

	var jsonBytes []byte
	var jsonString any

	if outputFile != "" {
		jsonString = loadResultOutputFile
	} else {
		jsonString = loadResultJSON
	}

	jsonBytes, err := json.Marshal(jsonString)

	if err != nil {
		log.Fatalf("Error marshaling struct: %v", err)
	}

	return jsonBytes
} */

func httpMethodFor(action string) string {
	switch action {
	case "delete":
		return "DELETE"
	case "switch":
		return "PATCH"
	default:
		return "POST"
	}
}

func substituteVars(data map[string]any, inputParams map[string]any) {
	if inputParams == nil {
		return
	}
	for k, v := range data {
		s, ok := v.(string)
		if !ok || !strings.Contains(s, "$") {
			continue
		}
		path := strings.Split(strings.Trim(s, "$"), ".")
		if val, err := getNestedValue(inputParams, path); err == nil && val != nil {
			data[k] = val
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "param substitution error: %v\n", err)
		}
	}
}

func dispatch(
	cmd *cobra.Command,
	res Resource,
	name, method string,
	payload map[string]any,
	body []byte,
) ([]byte, *ResourceCmdStruct) {
	const (
		cyan  = "\033[0;36m"
		reset = "\033[0m"
	)
	switch method {
	case "PATCH":
		if res.Name == Campaign {
			return nil, handleSwitch(cmd, payload)
		}
		return nil, nil

	case "DELETE":
		id := fmt.Sprintf("%v", payload["id"])
		if err := res.Data.Delete(id); err != nil {
			fmt.Fprintf(cmd.OutOrStdout(), "%s - %s: %v\n", cyan, name, err)
			return nil, nil
		}
		resp := []byte(fmt.Sprintf("Deleted id %s", id))
		if viper.GetString("output_format") != "json" {
			fmt.Fprintf(cmd.OutOrStdout(), "%s - %s: %s\n", cyan, name, reset, resp)
		}
		return resp, nil

	default: // POST
		resp, err := res.Data.Save(body)
		var resourceResp ResourceResp

		unmarshalErr := json.Unmarshal(resp, &resourceResp)
		if unmarshalErr != nil {
			fmt.Fprintf(os.Stderr, "unmarshal error: %v\n", err)
		}

		rec := ResourceCmdStruct{
			Name:      name,
			Response:  string(resp),
			Reference: res.Reference,
			Action:    method,
			ParentID:  resourceResp.Id,
		}

		if err != nil {
			rec.Error = err.Error()
		}

		if viper.GetString("output_format") != "json" {
			fmt.Fprintf(cmd.OutOrStdout(), "%s - %s: %s%s\n", cyan, name, reset, recResponse(rec))
		}

		return resp, &rec
	}
}

func handleSwitch(cmd *cobra.Command, payload map[string]any) *ResourceCmdStruct {
	id := fmt.Sprintf("%v", payload["id"])
	state := fmt.Sprintf("%v", payload["state"])
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

func recResponse(r ResourceCmdStruct) string {
	if r.Error != "" {
		return r.Error
	}
	return r.Response
}

func ScriptResource(cmd *cobra.Command, resources []Resource, inputParams map[string]any) []byte {
	refs := make(map[string]any, len(resources))
	var results []string
	var outputRecords []ResourceCmdStruct

	for _, res := range resources {
		name := res.Data.getName()
		action := strings.ToLower(res.Action)
		method := httpMethodFor(action)

		// 1) Prepare the payload
		raw, err := json.Marshal(res.Data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[%s] marshal error: %v\n", name, err)
			continue
		}

		var payload map[string]any
		if err := json.Unmarshal(raw, &payload); err != nil {
			fmt.Fprintf(os.Stderr, "[%s] unmarshal error: %v\n", name, err)
			continue
		}

		// 2) Substitute input params like "$foo.bar"
		substituteVars(payload, inputParams)

		// 3) Resolve references from previous resources
		payload = resolveVariables(payload, refs).(map[string]any)

		// 4) Marshal the final payload
		finalBody, err := json.Marshal(payload)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[%s] final marshal error: %v\n", name, err)
			continue
		}

		// 5) Dispatch by method
		respBytes, record := dispatch(cmd, res, name, method, payload, finalBody)
		if record != nil {
			outputRecords = append(outputRecords, *record)
		}
		if respBytes != nil {
			if method != "DELETE" && method != "PATCH" {
				var respData any
				if err := json.Unmarshal(respBytes, &respData); err == nil && respData != nil {
					refs[res.Reference] = respData
				}
				results = append(results, string(respBytes))
			}
		}
	}

	// 6) Choose output format
	var out any
	if file := viper.GetString("output_file"); file != "" {
		out = outputRecords
	} else {
		out = results
	}

	b, err := json.Marshal(out)
	if err != nil {
		fmt.Fprintf(os.Stderr, "final marshal error: %v\n", err)
		os.Exit(1)
	}
	return b
}

func getNestedValue(data map[string]any, path []string) (any, error) {
	if len(path) == 0 {
		return nil, fmt.Errorf("empty path")
	}

	current := data
	for i, key := range path {
		value, ok := current[key]
		if !ok {
			continue
		}
		if i == len(path)-1 {
			return value, nil
		}
		next, ok := value.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("value at key '%s' is not an object", key)
		}
		current = next
	}

	return nil, nil
}
