package helper

import "sigs.k8s.io/controller-runtime/pkg/client"

type DeleteOption interface {
	ApplyToDelete(opts DeleteOptions) DeleteOptions
}

type DeleteOptions struct {
	Timeout  Timeout
	Interval Interval
	Bindings Bindings
	Template Template
}

func NewDeleteOptions(opts []DeleteOption) DeleteOptions {
	options := DeleteOptions{}
	for _, opt := range opts {
		options = opt.ApplyToDelete(options)
	}
	return options
}

func (o DeleteOptions) ApplyToDelete(opts DeleteOptions) DeleteOptions {
	opts = o.Timeout.ApplyToDelete(opts)
	opts = o.Interval.ApplyToDelete(opts)
	return opts
}

// Delete deletes the specified resource and ensures the client cache is synced within the timeout.
// Uses Chainsaw to delete the resource if given a template and optional bindings.
// Stores the state of the deleted resource in the given struct.
func (h *Helper) Delete(obj client.Object, opts ...DeleteOption) {
	// options := NewDeleteOptions(append([]DeleteOption{h.Options}, opts...))
	// TODO
}
