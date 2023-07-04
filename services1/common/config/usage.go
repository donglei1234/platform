package config

// usageFormat is a template string that extracts envconfig settings into config.Help structs.
var usageFormat = `{{range .}}{
"Name": "{{usage_key .}}",
"FieldType": "{{usage_type .}}",
"Default": "{{usage_default .}}",
"Required": {{usage_type .}},
"Description": "{{usage_description .}}"
}
{{end}}`
