// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import "text/template"

var (
	// ProviderConvertersTemplate is the template of the Terraform provider implementation for the Illumio CloudSecure Config API.
	ProviderConvertersTemplate = template.Must(template.New("providermodel").Parse(`
{{- define "modelDataType" }}
{{- if ne .NestedModel nil }}
types.ObjectType{
	AttrTypes: GetTypeAttrsFor{{.NestedModel.Name}}(),
}
{{- else if eq .CollectionElementType nil }}
types.{{.ModelTypeName}}Type
{{- else }}
types.{{.ModelTypeName}}Type{ElemType: {{ template "modelDataType" .CollectionElementType }}}
{{- end }}
{{- end }}

{{- define "convertRepeatedProtoValueToData" }}
	{{- if ne .NestedModel nil }}
	dataValue = Convert{{.NestedModel.Name}}ToObjectValueFromProto(protoValue)
	{{- else if eq .CollectionElementType nil }}
	dataValue = types.{{.ModelTypeName}}Value(protoValue)
	{{- else if eq .ModelTypeName "Map" }}
	{
		dataElementType := {{ template "modelDataType" .CollectionElementType }}
		protoElements := protoValue
		if protoElements == nil {
			dataValue = types.MapNull(dataElementType)
		} else {
			dataValues := make(map[string]attr.Value, len(protoElements))
			for k, protoElement := range protoElements {
				{{- if ne .CollectionElementType.NestedModel nil }}
				var protoValue *{{.CollectionElementType.ProtoTypeName}} = {{ if ne .UnwrapProtoValueElementExpr nil }}{{ .UnwrapProtoValueElementExpr }}{{ else }}protoElement{{ end }}
				{{- else }}
				var protoValue {{.CollectionElementType.ProtoTypeName}} = {{ if ne .UnwrapProtoValueElementExpr nil }}{{ .UnwrapProtoValueElementExpr }}{{ else }}protoElement{{ end }}
				{{- end }}
				var dataValue attr.Value
				{{ template "convertRepeatedProtoValueToData" .CollectionElementType }}
				dataValues[k] = dataValue
			}
			dataValue = types.MapValueMust(dataElementType, dataValues)
		}
	}
	{{- else }}
	{
		dataElementType := {{ template "modelDataType" .CollectionElementType }}
		protoElements := protoValue
		if protoElements == nil {
			dataValue = types.{{.ModelTypeName}}Null(dataElementType)
		} else {
			dataValues := make([]attr.Value, 0, len(protoElements))
			for _, protoElement := range protoElements {
				{{- if ne .CollectionElementType.NestedModel nil }}
				var protoValue *{{.CollectionElementType.ProtoTypeName}} = {{ if ne .UnwrapProtoValueElementExpr nil }}{{ .UnwrapProtoValueElementExpr }}{{ else }}protoElement{{ end }}
				{{- else }}
				var protoValue {{.CollectionElementType.ProtoTypeName}} = {{ if ne .UnwrapProtoValueElementExpr nil }}{{ .UnwrapProtoValueElementExpr }}{{ else }}protoElement{{ end }}
				{{- end }}
				var dataValue attr.Value
				{{ template "convertRepeatedProtoValueToData" .CollectionElementType }}
				dataValues = append(dataValues, dataValue)
			}
			dataValue = types.{{.ModelTypeName}}ValueMust(dataElementType, dataValues)
		}
	}
	{{- end }}
{{- end }}

{{- define "convertDataValueToProto" }}
	{{- if ne .NestedModel nil }}
	protoValue, newDiags := ConvertDataValueTo{{.NestedModel.Name}}Proto(ctx, dataValue)
	diags.Append(newDiags...)
	if diags.HasError() {
		return nil, diags
	}
	{{- else if eq .CollectionElementType nil }}
	protoValue = dataValue.(types.{{.ModelTypeName}}).Value{{.ModelTypeName}}()
	{{- else if eq .ModelTypeName "Map" }}
	{
		dataElements := dataValue.(types.Map).Elements()
		protoValues := make({{.ProtoTypeName}}, len(dataElements))
		for k, dataElement := range dataElements {
			var dataValue attr.Value = dataElement
			{{- if eq .CollectionElementType.NestedModel nil }}
			var protoValue {{.CollectionElementType.ProtoTypeName}}
			{{- else }}
			var protoValue *{{.CollectionElementType.ProtoTypeName}}
			{{- end }}
			{{ template "convertDataValueToProto" .CollectionElementType }}
			protoValues[k] = {{ if ne .WrapProtoValueElementExpr nil }}{{ .WrapProtoValueElementExpr }}{{ else }}protoValue{{ end }}
		}
		protoValue = protoValues
	}
	{{- else }}
	{
		dataElements := dataValue.(types.{{.ModelTypeName}}).Elements()
		protoValues := make({{.ProtoTypeName}}, 0, len(dataElements))
		for _, dataElement := range dataElements {
			var dataValue attr.Value = dataElement
			{{- if eq .CollectionElementType.NestedModel nil }}
			var protoValue {{.CollectionElementType.ProtoTypeName}}
			{{- else }}
			var protoValue *{{.CollectionElementType.ProtoTypeName}}
			{{- end }}
			{{ template "convertDataValueToProto" .CollectionElementType }}
			protoValues = append(protoValues, {{ if ne .WrapProtoValueElementExpr nil }}{{ .WrapProtoValueElementExpr }}{{ else }}protoValue{{ end }})
		}
		protoValue = protoValues
	}
	{{- end }}
{{- end }}

{{- define "convertersForModel" }}
type {{.Name}} struct {
	{{- range $field := .Fields }}
	{{$field.Name}} types.{{$field.Type.ModelTypeName}} ` + "`" + `tfsdk:"{{$field.AttributeName}}"` + "`" + `
	{{- end }}
}

func GetTypeAttrsFor{{.Name}}() map[string]attr.Type {
	return map[string]attr.Type{
		{{- range $field := .Fields }}
		{{- if ne $field.Type.CollectionElementType nil }}
		"{{$field.AttributeName}}": {{ template "modelDataType" $field.Type }},
		{{- else if ne $field.Type.NestedModel nil }}
		"{{$field.AttributeName}}": types.ObjectType{
			AttrTypes: GetTypeAttrsFor{{$field.Type.NestedModel.Name}}(),
		},
		{{- else }}
		"{{$field.AttributeName}}": types.{{$field.Type.ModelTypeName}}Type,
		{{- end }}
		{{- end }}
	}
}

func Convert{{.Name}}ToObjectValueFromProto(proto *configv1.{{.Name}}) basetypes.ObjectValue  {
	vals := map[string]attr.Value{}
	{{- range $field := .Fields }}
	{{- if ne $field.Type.NestedModel nil }}
	vals["{{$field.AttributeName}}"] = Convert{{$field.Type.NestedModel.Name}}ToObjectValueFromProto(proto.{{$field.Name}})
	{{- else if ne $field.Type.CollectionElementType nil }}
	{
		protoValue := proto.{{$field.Name}}
		var dataValue types.{{$field.Type.ModelTypeName}}
		{{ template "convertRepeatedProtoValueToData" $field.Type }}
		vals["{{$field.AttributeName}}"] = dataValue
	}
	{{- else }}
	vals["{{$field.AttributeName}}"] = types.{{$field.Type.ModelTypeName}}Value(proto.{{$field.Name}})
	{{- end }}
	{{- end }}
	return types.ObjectValueMust(
		GetTypeAttrsFor{{.Name}}(),
		vals,
	)
}

func ConvertDataValueTo{{.Name}}Proto(ctx context.Context, dataValue attr.Value) (*configv1.{{.Name}}, diag.Diagnostics) {
	pv := {{.Name}}{}
	diags := tfsdk.ValueAs(ctx, dataValue, &pv)
	if diags.HasError() {
		return nil, diags
	}
	proto := &configv1.{{.Name}}{}
	{{- range $field := .Fields }}
	{{- if ne $field.Type.CollectionElementType nil }}
	{
		var dataValue attr.Value = pv.{{$field.Name}}
		var protoValue {{$field.Type.ProtoTypeName}}
		{{ template "convertDataValueToProto" $field.Type }}
		proto.{{$field.Name}} = protoValue
	}
	{{- else if ne $field.Type.NestedModel nil }}
	pvModel, dvDiags := ConvertDataValueTo{{$field.Type.NestedModel.Name}}Proto(ctx, pv.{{$field.Name}})
	diags.Append(dvDiags...)
	if diags.HasError() {
		return nil, diags
	}
	proto.{{$field.Name}} = pvModel
	{{- else }}
	proto.{{$field.Name}} = pv.{{$field.Name}}.Value{{$field.Type.ModelTypeName}}()
	{{- end }}
	{{- end }}
	return proto, diags
}

{{- range $field := .Fields }}
	{{- if ne $field.Type.NestedModel nil }}
	{{ template "convertersForModel" $field.Type.NestedModel }}
	{{- end }}
	{{- if ne $field.Type.CollectionElementType nil }}
		{{- if ne $field.Type.CollectionElementType.NestedModel nil }}
		{{ template "convertersForModel" $field.Type.CollectionElementType.NestedModel }}
		{{- end }}
		{{- if ne $field.Type.CollectionElementType.CollectionElementType nil }}
			{{- if ne $field.Type.CollectionElementType.CollectionElementType.NestedModel nil }}
			{{ template "convertersForModel" $field.Type.CollectionElementType.CollectionElementType.NestedModel }}
			{{- end }}
		{{- end }}
	{{- end }}
{{- end }}
{{- end }}

{{- range $model := .Models }}
	{{- range $fields := $model.Fields }}
		{{- if ne $fields.Type.NestedModel nil }}
		{{ template "convertersForModel" $fields.Type.NestedModel }}
		{{- end }}
		{{- if ne $fields.Type.CollectionElementType nil }}
			{{- if ne $fields.Type.CollectionElementType.NestedModel nil }}
			{{ template "convertersForModel" $fields.Type.CollectionElementType.NestedModel }}
			{{- end }}
			{{- if ne $fields.Type.CollectionElementType.CollectionElementType nil }}
				{{- if ne $fields.Type.CollectionElementType.CollectionElementType.NestedModel nil }}
				{{ template "convertersForModel" $fields.Type.CollectionElementType.CollectionElementType.NestedModel }}
				{{- end }}
			{{- end }}
		{{- end }}
	{{- end }}
{{- end }}
`))
)
