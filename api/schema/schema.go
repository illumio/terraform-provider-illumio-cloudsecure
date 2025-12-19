// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"context"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// cloudSecureSchema is an implementation of Schema that returns the schemas of CloudSecure resources and data sources.
type cloudSecureSchema struct{}

func (c *cloudSecureSchema) Version() string {
	return "v1"
}

const (
	commonRequestTimeoutDescription = `If not specified, defaults to the provider's "request_timeout" attribute. ` +
		`Must be a string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) ` +
		`consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are ` +
		`"s" (seconds), "m" (minutes), "h" (hours).`

	createRequestTimeoutDescription = `Maximum duration of each Create API request for this resource. ` +
		commonRequestTimeoutDescription

	readRequestTimeoutDescription = `Maximum duration of each Read API request for this resource. ` +
		commonRequestTimeoutDescription

	updateRequestTimeoutDescription = `Maximum duration of each Update API request for this resource. ` +
		commonRequestTimeoutDescription

	deleteRequestTimeoutDescription = `Maximum duration of each Delete API request for this resource. ` +
		commonRequestTimeoutDescription +
		` Setting a timeout for Delete API requests is only applicable if ` +
		`changes are saved into state before the destroy operation occurs.`
)

func (c *cloudSecureSchema) Resources() Resources {
	// Keep all resources sorted by lexicographic name order.
	resources := Resources{
		applicationAwsResourcesResource,
		applicationAzureResourcesResource,
		applicationPolicyRuleResource,
		applicationResource,
		awsAccountResource,
		awsFlowLogsS3BucketResource,
		azureFlowLogsStorageAccountResource,
		azureSubscriptionResource,
		deploymentResource,
		gcpFlowLogsPubsubTopicResource,
		gcpProjectResource,
		ipListResource,
		k8sClusterOnboardingCredentialResource,
		k8sClusterResource,
		organizationPolicyResource,
		organizationPolicyRuleResource,
		tagToLabelResource,
	}

	sort.Sort(resources)

	// Add a "timeouts" block to each resource to support overriding the
	// request timeout for each Create/Read/Update/Delete API call individually.
	timeoutOpts := timeouts.Opts{
		Create:            true,
		Read:              true,
		Update:            true,
		Delete:            true,
		CreateDescription: createRequestTimeoutDescription,
		ReadDescription:   readRequestTimeoutDescription,
		UpdateDescription: updateRequestTimeoutDescription,
		DeleteDescription: deleteRequestTimeoutDescription,
	}

	for i := range resources {
		if resources[i].Schema.Blocks == nil {
			resources[i].Schema.Blocks = make(map[string]resource_schema.Block, 1)
		}

		resources[i].Schema.Blocks["timeouts"] = timeouts.Block(context.Background(), timeoutOpts)
	}

	return resources
}

func (c *cloudSecureSchema) DataSources() DataSources {
	// Keep all data sources sorted by lexicographic name order.
	dataSources := DataSources{}
	sort.Sort(dataSources)

	return dataSources
}

var _ Schema = &cloudSecureSchema{}

// CloudSecure returns the schemas of CloudSecure resources and data sources.
func CloudSecure() Schema {
	return &cloudSecureSchema{}
}
