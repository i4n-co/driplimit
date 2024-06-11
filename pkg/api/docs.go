package api

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/structtag"
)

// NamespacedDoc is a map of namespace to GeneratedRPCDocumentation
type NamespacedDoc map[string][]GeneratedRPCDocumentation

// Docs is the documentation for the API
type Docs struct {
	RPCs NamespacedDoc
}

// RPCDocumentation is the documentation for an RPC endpoint
type RPCDocumentation struct {
	Path        string `json:"path"`
	Description string `json:"description"`
	Parameters  any    `json:"parameters,omitempty"`
	Response    any    `json:"response,omitempty"`
	Code        int    `json:"code,omitempty"`
}

// RPCDocumentationField is a field in an GeneratedRPCDocumentation
type RPCDocumentationField struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description,omitempty"`
	Type        string                  `json:"type"`
	MapOf       string                  `json:"map_of,omitempty"`
	Required    bool                    `json:"required,omitempty"`
	SubFields   []RPCDocumentationField `json:"sub_fields,omitempty"`
}

// GeneratedRPCDocumentation is the documentation for an RPC endpoint
// with parameter fields automatically discovered
type GeneratedRPCDocumentation struct {
	RPCDocumentation
	ParamFields []RPCDocumentationField `json:"param_fields"`
}

// GenerateDocs generates the documentation for the API
func (api *Server) GenerateDocs() (*Docs, error) {
	docs := new(Docs)
	docs.RPCs = make(NamespacedDoc)
	for _, rpc := range api.rpcs {
		rpcDoc := rpc.Documentation.DocumentStruct()
		rpcDoc.Path = rpc.path()
		if rpcDoc.Code == 0 {
			rpcDoc.Code = 200
		}
		docs.RPCs[rpc.Namespace] = append(docs.RPCs[rpc.Namespace], rpcDoc)
	}
	return docs, nil
}

// DocumentStruct recursively documents a struct
func (doc RPCDocumentation) DocumentStruct() GeneratedRPCDocumentation {
	ed := GeneratedRPCDocumentation{
		RPCDocumentation: doc,
	}

	documentStruct(doc.Parameters, &ed.ParamFields)
	return ed
}

func documentStruct(strct any, docfields *[]RPCDocumentationField) {
	t := reflect.TypeOf(strct)
	if t == nil {
		return
	}
	for i := 0; i < t.NumField(); i++ {
		docfield := documentField(t.Field(i))
		if docfield.Name == "-" {
			continue
		}
		if docfield.Type == "object" {
			docfield.SubFields = make([]RPCDocumentationField, 0)
			documentStruct(reflect.New(t.Field(i).Type).Elem().Interface(), &docfield.SubFields)
		}
		if docfield.Type == "map" {
			docfield.SubFields = make([]RPCDocumentationField, 1)
			docfield.SubFields[0].Name = ""
			docfield.SubFields[0].Type = docFieldType(t.Field(i).Type.Key())
			docfield.SubFields[0].Required = true
			docfield.SubFields[0].Description = "keys of the map"
			docfield.SubFields[0].SubFields = make([]RPCDocumentationField, 0)
			documentStruct(reflect.New(t.Field(i).Type.Elem()).Elem().Interface(), &docfield.SubFields[0].SubFields)
		}
		*docfields = append(*docfields, docfield)
	}
}

func isRequired(field reflect.StructField) bool {
	tag, err := structtag.Parse(string(field.Tag))
	if err != nil {
		return false
	}
	validateTag, _ := tag.Get("validate")
	if validateTag == nil {
		return false
	}
	return validateTag.Name == "required"
}

func jsonName(field reflect.StructField) string {
	tag, err := structtag.Parse(string(field.Tag))
	if err != nil {
		return field.Name
	}
	jsonTag, _ := tag.Get("json")
	if jsonTag == nil {
		return field.Name
	}
	return jsonTag.Name
}

func description(field reflect.StructField) string {
	tag, err := structtag.Parse(string(field.Tag))
	if err != nil {
		return ""
	}
	descriptionTag, _ := tag.Get("description")
	if descriptionTag == nil {
		return ""
	}
	return descriptionTag.Name
}

func documentField(field reflect.StructField) RPCDocumentationField {
	doc := RPCDocumentationField{}

	doc.Name = jsonName(field)
	doc.Description = description(field)
	doc.Type = docFieldType(field.Type)
	doc.Required = isRequired(field)
	if doc.Type == "object" || strings.HasPrefix(doc.Type, "map") {
		doc.SubFields = make([]RPCDocumentationField, 0)
	}
	return doc
}

func docFieldType(field reflect.Type) (fieldType string) {
	if field.String() == "time.Time" {
		return "timestamp"
	}
	if strings.HasPrefix(field.String(), "int") {
		return "integer"
	}
	if field.String() == "driplimit.Milliseconds" {
		return "integer"
	}
	if field.Kind() == reflect.Struct {
		return "object"
	}
	if field.Kind() == reflect.Slice {
		return "array"
	}
	if field.Kind() == reflect.Map {
		return "map"
	}
	return fmt.Sprintf("%v", field)
}
