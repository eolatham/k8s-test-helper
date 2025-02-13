package chainsaw

import (
	"context"
	"fmt"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/engine/templating"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ParseResource parses the resource in the template file and returns it as a structured object.
func ParseResource(c client.Client, ctx context.Context, templatePath string, bindingsMap map[string]any) (client.Object, error) {
	// Load resource
	r, err := loadTemplateResource(templatePath)
	if err != nil {
		return nil, err
	}

	// Convert bindings
	bindings := bindingsFromMap(ctx, bindingsMap)

	// Parse and merge templated fields into unstructured object
	compilers := apis.DefaultCompilers
	template := v1alpha1.NewProjection(r.UnstructuredContent())
	merged, err := templating.TemplateAndMerge(ctx, compilers, r, bindings, template)
	if err != nil {
		return nil, err
	}

	// Get GVK from unstructured object
	gvk := merged.GroupVersionKind()

	// Create new instance of the correct type
	scheme := c.Scheme()
	if scheme == nil {
		return nil, fmt.Errorf("client scheme is not set")
	}
	typed, err := scheme.New(gvk)
	if err != nil {
		return nil, fmt.Errorf("failed to create object for GVK %v: %w", gvk, err)
	}

	// Convert unstructured object to typed one
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(merged.Object, typed); err != nil {
		return nil, fmt.Errorf("failed to convert unstructured to typed object: %w", err)
	}

	// Convert to client.Object (which all K8s types implement)
	obj, ok := typed.(client.Object)
	if !ok {
		return nil, fmt.Errorf("object of type %T does not implement client.Object", typed)
	}

	// Return structured object
	return obj, nil
}
