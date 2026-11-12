// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestApplicationAwsResourcesUpdateKeepsStateIdentities(t *testing.T) {
	for _, computedPlan := range []struct {
		name string
		ids  types.List
	}{
		{"unchanged", types.ListValueMust(types.StringType, []attr.Value{types.StringValue("cloud-id-1")})},
		{"unknown", types.ListUnknown(types.StringType)},
	} {
		t.Run(computedPlan.name, func(t *testing.T) {
			before := ApplicationAwsResourcesResourceModel{
				Id: types.StringValue("association-id"), AccountId: types.StringValue("account-id"),
				ApplicationId:                     types.StringValue("application-id"),
				ApplicationResourceIds:            types.ListValueMust(types.StringType, []attr.Value{types.StringValue("cloud-id-1")}),
				Arns:                              types.ListValueMust(types.StringType, []attr.Value{types.StringValue("arn:aws:old")}),
				AwsEc2TransitGatewayAttachmentIds: types.ListNull(types.StringType),
				AwsInstancesIds:                   types.ListNull(types.StringType),
				AwsNatGatewayIds:                  types.ListNull(types.StringType),
				AwsNetworkAclIds:                  types.ListNull(types.StringType),
				AwsNetworkInterfaceIds:            types.ListNull(types.StringType),
				AwsRdsClusterIds:                  types.ListNull(types.StringType),
				AwsSecurityGroupIds:               types.ListNull(types.StringType),
				AwsSubnetIds:                      types.ListNull(types.StringType),
				AwsVpcEndpointIds:                 types.ListNull(types.StringType),
				AwsVpcIds:                         types.ListNull(types.StringType),
				AwsVpcPeeringConnectionIds:        types.ListNull(types.StringType),
			}
			after := before
			after.ApplicationResourceIds = computedPlan.ids
			after.Arns = types.ListValueMust(types.StringType, []attr.Value{types.StringValue("arn:aws:old"), types.StringValue("arn:aws:new")})
			req, diags := NewUpdateApplicationAwsResourcesRequest(context.Background(), &before, &after)
			if diags.HasError() {
				t.Fatalf("request diagnostics: %v", diags)
			}
			if req.Id != "association-id" || req.AccountId != "account-id" || req.ApplicationId != "application-id" {
				t.Fatalf("missing stable identities: %+v", req)
			}
			if !reflect.DeepEqual(req.ApplicationResourceIds, []string{"cloud-id-1"}) {
				t.Fatalf("previous tracking IDs = %v", req.ApplicationResourceIds)
			}
			if !reflect.DeepEqual(req.Arns, []string{"arn:aws:old", "arn:aws:new"}) || !reflect.DeepEqual(req.UpdateMask.Paths, []string{"arns"}) {
				t.Fatalf("mutable fields or mask: arns=%v mask=%v", req.Arns, req.UpdateMask.Paths)
			}
		})
	}
}

func TestApplicationAzureResourcesUpdateKeepsStateIdentities(t *testing.T) {
	for _, computedPlan := range []struct {
		name string
		ids  types.List
	}{
		{"unchanged", types.ListValueMust(types.StringType, []attr.Value{types.StringValue("cloud-id-1")})},
		{"unknown", types.ListUnknown(types.StringType)},
	} {
		t.Run(computedPlan.name, func(t *testing.T) {
			before := ApplicationAzureResourcesResourceModel{
				Id: types.StringValue("association-id"), ApplicationId: types.StringValue("application-id"),
				SubscriptionId:         types.StringValue("subscription-id"),
				ApplicationResourceIds: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("cloud-id-1")}),
				ResourceIds:            types.ListValueMust(types.StringType, []attr.Value{types.StringValue("azure-old")}),
			}
			after := before
			after.ApplicationResourceIds = computedPlan.ids
			after.ResourceIds = types.ListValueMust(types.StringType, []attr.Value{types.StringValue("azure-old"), types.StringValue("azure-new")})
			req, diags := NewUpdateApplicationAzureResourcesRequest(context.Background(), &before, &after)
			if diags.HasError() {
				t.Fatalf("request diagnostics: %v", diags)
			}
			if req.Id != "association-id" || req.ApplicationId != "application-id" || req.SubscriptionId != "subscription-id" {
				t.Fatalf("missing stable identities: %+v", req)
			}
			if !reflect.DeepEqual(req.ApplicationResourceIds, []string{"cloud-id-1"}) {
				t.Fatalf("previous tracking IDs = %v", req.ApplicationResourceIds)
			}
			if !reflect.DeepEqual(req.ResourceIds, []string{"azure-old", "azure-new"}) || !reflect.DeepEqual(req.UpdateMask.Paths, []string{"resource_ids"}) {
				t.Fatalf("mutable fields or mask: resource_ids=%v mask=%v", req.ResourceIds, req.UpdateMask.Paths)
			}
		})
	}
}
