package common

import (
	"fmt"
	"strings"
)

var (
	ResourceFile string
	OutputFile   string
	InputRefRaw  string
	InputRefFile string
	DryRun       bool
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

func ResolveRefs(val any, rc *RefContext) any {
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
			v[i] = ResolveRefs(item, rc)
		}
		return v

	case map[string]interface{}:
		for k, mapVal := range v {
			v[k] = ResolveRefs(mapVal, rc)
		}
		return v

	default:
		return val
	}
}

func ValidateReferences(value interface{}, refCtx *RefContext, path string, resRef string) error {
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
			if err := ValidateReferences(item, refCtx, itemPath, resRef); err != nil {
				return err
			}
		}

	case map[string]interface{}:
		for k, val := range v {
			nestedPath := fmt.Sprintf("%s.%s", path, k)
			if err := ValidateReferences(val, refCtx, nestedPath, resRef); err != nil {
				return err
			}
		}
	}

	return nil
}
