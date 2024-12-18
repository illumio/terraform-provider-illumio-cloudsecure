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
		{{- else}}
		"{{$field.AttributeName}}": types.{{$field.Type.ModelTypeName}}Type,
		{{- end}}
		{{- end}}
	}
}


func Convert{{.Name}}ToObjectValueFromProto(proto *configv1.{{.Name}}) basetypes.ObjectValue  {
	return types.ObjectValueMust(
		GetTypeAttrsFor{{.Name}}(),
		map[string]attr.Value{
			{{- range $field := .Fields}}
			{{- if ne $field.Type.NestedModel nil}}
			"{{$field.AttributeName}}": Convert{{$field.Type.NestedModel.Name}}ToObjectValueFromProto(proto.{{$field.Name}}),
			{{- else}}
			"{{$field.AttributeName}}": types.{{$field.Type.ModelTypeName}}Value(proto.{{$field.Name}}),
			{{- end}}
			{{- end}}
		},
	)
}
func ConvertDataValueTo{{.Name}}Proto(ctx context.Context, dataValue attr.Value) *configv1.{{.Name}} {
	pv := {{.Name}}{}
	diags := tfsdk.ValueAs(context.Background(), dataValue, &pv)
	if len(diags) > 0 {
		tflog.Error(ctx, "Unexpected diagnostics", map[string]any{"diags": diags})
	}
	proto := &configv1.{{.Name}}{}
	{{- range $field := .Fields}}
	{{- if ne $field.Type.NestedModel nil}}
	proto.{{$field.Name}} = ConvertDataValueTo{{$field.Type.NestedModel.Name}}Proto(ctx, pv.{{$field.Name}})
	{{- else}}
	proto.{{$field.Name}} = pv.{{$field.Name}}.Value{{$field.Type.ModelTypeName}}()
	{{- end}}
	{{- end}}
	return proto
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
