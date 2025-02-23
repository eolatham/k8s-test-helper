package chainsaw

import (
	"context"
	"fmt"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/chainsaw/pkg/engine/bindings"
	"github.com/kyverno/chainsaw/pkg/engine/templating"
	"github.com/kyverno/chainsaw/pkg/loaders/resource"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// ValidateTemplate ensures the template and bindings represent a valid single resource.
func ValidateTemplate(
	ctx context.Context,
	templatePath string,
	bindingsMap map[string]any,
) error {
	// Check loading errors
	resource, err := loadTemplateResource(templatePath)
	if err != nil {
		return err
	}
	// Check parsing errors
	bindings := bindingsFromMap(ctx, bindingsMap)
	if err := templating.ResourceRef(ctx, compilers, &resource, bindings); err != nil {
		return err
	}
	return nil
}

// loadTemplateResource loads the template file and returns its unstructured contents.
// Expects the template file to contain a single resource.
func loadTemplateResource(templatePath string) (unstructured.Unstructured, error) {
	var r unstructured.Unstructured
	resources, err := resource.Load(templatePath, true)
	if err != nil {
		return r, fmt.Errorf("failed to load template file %s: %w", templatePath, err)
	}
	if len(resources) != 1 {
		return r, fmt.Errorf("expected template file %s to contain a single resource; found %d", templatePath, len(resources))
	}
	r = resources[0]
	return r, nil
}

// bindingsFromMap converts the map into an object that can applied to Chainsaw templates.
func bindingsFromMap(ctx context.Context, bindingsMap map[string]any) apis.Bindings {
	bindingsObj := apis.NewBindings()
	for k, v := range bindingsMap {
		bindingsObj = bindings.RegisterBinding(ctx, bindingsObj, k, v)
	}
	return bindingsObj
}
