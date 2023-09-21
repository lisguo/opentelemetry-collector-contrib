// Code generated by mdatagen. DO NOT EDIT.

package metadata

import (
	"go.opentelemetry.io/collector/pdata/pcommon"
)

// ResourceBuilder is a helper struct to build resources predefined in metadata.yaml.
// The ResourceBuilder is not thread-safe and must not to be used in multiple goroutines.
type ResourceBuilder struct {
	config ResourceAttributesConfig
	res    pcommon.Resource
}

// NewResourceBuilder creates a new ResourceBuilder. This method should be called on the start of the application.
func NewResourceBuilder(rac ResourceAttributesConfig) *ResourceBuilder {
	return &ResourceBuilder{
		config: rac,
		res:    pcommon.NewResource(),
	}
}

// SetHostArch sets provided value as "host.arch" attribute.
func (rb *ResourceBuilder) SetHostArch(val string) {
	if rb.config.HostArch.Enabled {
		rb.res.Attributes().PutStr("host.arch", val)
	}
}

// SetHostID sets provided value as "host.id" attribute.
func (rb *ResourceBuilder) SetHostID(val string) {
	if rb.config.HostID.Enabled {
		rb.res.Attributes().PutStr("host.id", val)
	}
}

// SetHostName sets provided value as "host.name" attribute.
func (rb *ResourceBuilder) SetHostName(val string) {
	if rb.config.HostName.Enabled {
		rb.res.Attributes().PutStr("host.name", val)
	}
}

// SetOsDescription sets provided value as "os.description" attribute.
func (rb *ResourceBuilder) SetOsDescription(val string) {
	if rb.config.OsDescription.Enabled {
		rb.res.Attributes().PutStr("os.description", val)
	}
}

// SetOsType sets provided value as "os.type" attribute.
func (rb *ResourceBuilder) SetOsType(val string) {
	if rb.config.OsType.Enabled {
		rb.res.Attributes().PutStr("os.type", val)
	}
}

// Emit returns the built resource and resets the internal builder state.
func (rb *ResourceBuilder) Emit() pcommon.Resource {
	r := rb.res
	rb.res = pcommon.NewResource()
	return r
}
