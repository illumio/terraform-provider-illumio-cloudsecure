// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import "text/template"

var (
	// ProviderConvertersTemplate is the template of the Terraform provider implementation for the Illumio CloudSecure Config API.
	ProviderConvertersTemplate = template.Must(template.New("providermodel").Parse(`
{{- define "convertersForModel"}}
type {{.Name}} struct {
	{{- range $field := .Fields}}
	{{$field.Name}} types.{{$field.Type.ModelTypeName}} ` + "`" + `tfsdk:"{{$field.AttributeName}}"` + "`" + `
	{{- end}}
}

func GetTypeAttrsFor{{.Name}}() map[string]attr.Type {
	return map[string]attr.Type{
		{{- range $field := .Fields}}
		{{- if ne $field.Type.NestedModel nil}}
		"{{$field.AttributeName}}": types.ObjectType{
			AttrTypes: GetTypeAttrsFor{{$field.Type.NestedModel.Name}}(),
		},
		{{- else if and (ne $field.Type.CollectionElementType nil) (ne $field.Type.CollectionElementType.NestedModel nil)}}
		"{{$field.AttributeName}}": types.{{$field.Type.ModelTypeName}}Type{ElemType: types.ObjectType{
			AttrTypes: GetTypeAttrsFor{{$field.Type.CollectionElementType.NestedModel.Name}}(),
		}},
		{{- else if ne $field.Type.CollectionElementType nil}}
		"{{$field.AttributeName}}": types.{{$field.Type.ModelTypeName}}Type{ElemType: types.{{$field.Type.CollectionElementType.ModelTypeName}}Type},
		{{- else}}
		"{{$field.AttributeName}}": types.{{$field.Type.ModelTypeName}}Type,
		{{- end}}
		{{- end}}
	}
}

func Convert{{.Name}}ToObjectValueFromProto(proto *configv1.{{.Name}}) basetypes.ObjectValue  {
	{{- range $field := .Fields}}
	{{- if and (ne $field.Type.CollectionElementType nil) (ne $field.Type.CollectionElementType.NestedModel nil)}}
	{{- if eq $field.Type.ModelTypeName "Map"}}
	elementsIn{{$field.Name}} := make(map[string]attr.Value, len(proto.{{$field.Name}}))
	for key, item := range proto.{{$field.Name}} {
		elementsIn{{$field.Name}}[key] = Convert{{$field.Type.CollectionElementType.NestedModel.Name}}ToObjectValueFromProto(item)
	}
	{{- else}}{{- /* "List" or "Set" */}}
	elementsIn{{$field.Name}} := make([]attr.Value, 0, len(proto.{{$field.Name}}))
	for _, item := range proto.{{$field.Name}} {
		elementsIn{{$field.Name}} = append(elementsIn{{$field.Name}}, Convert{{$field.Type.CollectionElementType.NestedModel.Name}}ToObjectValueFromProto(item))
	}
	{{- end}}
	{{- else if ne $field.Type.CollectionElementType nil}}
	{{- if eq $field.Type.ModelTypeName "Map"}}
	elementsIn{{$field.Name}} := make(map[string]attr.Value, len(proto.{{$field.Name}}))
	for key, item := range proto.{{$field.Name}} {
		elementsIn{{$field.Name}}[key] = types.{{$field.Type.CollectionElementType.ModelTypeName}}Value(item)
	}
	{{- else}}{{- /* "List" or "Set" */}}
	elementsIn{{$field.Name}} := make([]attr.Value, 0, len(proto.{{$field.Name}}))
	for _, item := range proto.{{$field.Name}} {
		elementsIn{{$field.Name}} = append(elementsIn{{$field.Name}}, types.{{$field.Type.CollectionElementType.ModelTypeName}}Value(item))
	}
	{{- end}}
	{{- end}}
	{{- end}}
	return types.ObjectValueMust(
		GetTypeAttrsFor{{.Name}}(),
		map[string]attr.Value{
			{{- range $field := .Fields}}
			{{- if ne $field.Type.NestedModel nil}}
			"{{$field.AttributeName}}": Convert{{$field.Type.NestedModel.Name}}ToObjectValueFromProto(proto.{{$field.Name}}),
			{{- else if and (ne $field.Type.CollectionElementType nil) (ne $field.Type.CollectionElementType.NestedModel nil)}}
			"{{$field.AttributeName}}": types.{{$field.Type.ModelTypeName}}ValueMust(types.ObjectType{AttrTypes: GetTypeAttrsFor{{$field.Type.CollectionElementType.NestedModel.Name}}()}, elementsIn{{$field.Name}}),
			{{- else if ne $field.Type.CollectionElementType nil}}
			"{{$field.AttributeName}}": types.{{$field.Type.ModelTypeName}}ValueMust(types.{{$field.Type.CollectionElementType.ModelTypeName}}Type, elementsIn{{$field.Name}}),
			{{- else}}
			"{{$field.AttributeName}}": types.{{$field.Type.ModelTypeName}}Value(proto.{{$field.Name}}),
			{{- end}}
			{{- end}}
		},
	)
}

func ConvertDataValueTo{{.Name}}Proto(ctx context.Context, dataValue attr.Value) (*configv1.{{.Name}}, diag.Diagnostics) {
	pv := {{.Name}}{}
	diags := tfsdk.ValueAs(ctx, dataValue, &pv)
	if diags.HasError() {
		return nil, diags
	}
	proto := &configv1.{{.Name}}{}
	{{- range $field := .Fields}}
	{{- if ne $field.Type.NestedModel nil}}
	pvModel{{$field.Name}}, dvDiags{{$field.Name}} := ConvertDataValueTo{{$field.Type.NestedModel.Name}}Proto(ctx, pv.{{$field.Name}})
	diags.Append(dvDiags{{$field.Name}}...)
	if diags.HasError() {
		return nil, diags
	}
	proto.{{$field.Name}} = pvModel{{$field.Name}}
	{{- else if and (ne $field.Type.CollectionElementType nil) (ne $field.Type.CollectionElementType.NestedModel nil)}}
	{{- if eq $field.Type.ModelTypeName "Map"}}
	pvElemModel{{$field.Name}} := pv.{{$field.Name}}.Elements()
	proto.{{$field.Name}} = make(map[string]*configv1.{{$field.Type.CollectionElementType.NestedModel.Name}}, len(pvElemModel{{$field.Name}}))
	for key, elem := range pvElemModel{{$field.Name}} {
		pvModel, dvDiags := ConvertDataValueTo{{$field.Type.CollectionElementType.NestedModel.Name}}Proto(ctx, elem)
		diags.Append(dvDiags...)
		if diags.HasError() {
			return nil, diags
		}
		proto.{{$field.Name}}[key] = pvModel
	}
	{{- else}}{{- /* "List" or "Set" */}}
	pvElemModel{{$field.Name}} := pv.{{$field.Name}}.Elements()
	proto.{{$field.Name}} = make([]*configv1.{{$field.Type.CollectionElementType.NestedModel.Name}}, 0, len(pvElemModel{{$field.Name}}))
	for _, elem := range pvElemModel{{$field.Name}} {
		pvModel, dvDiags := ConvertDataValueTo{{$field.Type.CollectionElementType.NestedModel.Name}}Proto(ctx, elem)
		diags.Append(dvDiags...)
		if diags.HasError() {
			return nil, diags
		}
		proto.{{$field.Name}} = append(proto.{{$field.Name}}, pvModel)
	}
	{{- end}}
	{{- else if ne $field.Type.CollectionElementType nil}}
	{{- if eq $field.Type.ModelTypeName "Map"}}
	var pvModel{{$field.Name}} map[string]{{$field.Type.CollectionElementType.ProtoTypeName}}
	dvDiags{{$field.Name}} := pv.{{$field.Name}}.ElementsAs(ctx, &pvModel{{$field.Name}}, false)
	diags.Append(dvDiags{{$field.Name}}...)
	if diags.HasError() {
		return nil, diags
	}
	proto.{{$field.Name}} = pvModel{{$field.Name}}
	{{- else}}{{- /* "List" or "Set" */}}
	var pvModel{{$field.Name}} []{{$field.Type.CollectionElementType.ProtoTypeName}}
	dvDiags{{$field.Name}} := pv.{{$field.Name}}.ElementsAs(ctx, &pvModel{{$field.Name}}, false)
	diags.Append(dvDiags{{$field.Name}}...)
	if diags.HasError() {
		return nil, diags
	}
	proto.{{$field.Name}} = pvModel{{$field.Name}}
	{{- end}}
	{{- else}}
	proto.{{$field.Name}} = pv.{{$field.Name}}.Value{{$field.Type.ModelTypeName}}()
	{{- end}}
	{{- end}}
	return proto, diags
}

{{- range $fields := .Fields}}
{{- if ne $fields.Type.NestedModel nil}}
{{- template "convertersForModel" $fields.Type.NestedModel}}
{{- end}}
{{- if ne $fields.Type.CollectionElementType nil}}
{{- if ne $fields.Type.CollectionElementType.NestedModel nil}}
{{- template "convertersForModel" $fields.Type.CollectionElementType.NestedModel}}
{{- end}}
{{- end}}
{{- end}}

{{- end}}

{{- range $model := .Models}}
	{{- range $fields := $model.Fields}}
		{{- if ne $fields.Type.NestedModel nil}}
		{{- template "convertersForModel" $fields.Type.NestedModel}}
		{{- end}}
		{{- if ne $fields.Type.CollectionElementType nil}}
		{{- if ne $fields.Type.CollectionElementType.NestedModel nil}}
		{{- template "convertersForModel" $fields.Type.CollectionElementType.NestedModel}}
		{{- end}}
		{{- end}}
	{{- end}}
{{- end}}
`))
)
