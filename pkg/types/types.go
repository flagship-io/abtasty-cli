package types

import "github.com/flagship-io/abtasty-cli/internal/utils/common"

type Resource struct {
	Type           string                `json:"type"`
	Ref            string                `json:"$_ref"`
	ParentID       string                `json:"$_parent_id"`
	Action         common.ResourceAction `json:"action"`
	Payload        map[string]any        `json:"payload"`
	Resources      []Resource            `json:"resources"`
	ParentResource *Resource
}

type ResourceResult struct {
	Ref      string      `json:"$_ref,omitempty"`
	Status   string      `json:"status"`
	Response interface{} `json:"response,omitempty"`
}

type LoaderResults struct {
	Results []ResourceResult `json:"results"`
}

type LoadResFile struct {
	Version   int        `json:"version"`
	Resources []Resource `json:"resources"`
}
