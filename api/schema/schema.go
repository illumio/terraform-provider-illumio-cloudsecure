// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"sort"
)

// cloudSecureSchema is an implementation of Schema that returns the schemas of CloudSecure resources and data sources.
type cloudSecureSchema struct{}

func (c *cloudSecureSchema) Version() string {
	return "v1"
}

func (c *cloudSecureSchema) Resources() Resources {
	// Keep all resources sorted by lexicographic name order.
	resources := Resources{
		applicationPolicyRuleResource,
		awsAccountResource,
		awsFlowLogsS3Bucket,
		azureFlowLogsStorageAccount,
		azureSubscriptionResource,
		deploymentResource,
		ipListResource,
		k8sClusterOnboardingCredential,
		tagToLabelResource,
	}
	sort.Sort(resources)

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
