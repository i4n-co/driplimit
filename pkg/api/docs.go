package api

import (
	"reflect"
	"strings"

	"github.com/fatih/structtag"
)

type RPCDocumentation struct {
	Description string `json:"description"`
	Param       any    `json:"input"`
	Return      any    `json:"output"`
}

type RPCDocumentationField struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type"`
	Required    bool   `json:"required,omitempty"`
}

type EnhancedRPCDocumentation struct {
	RPCDocumentation
	ParamFields []RPCDocumentationField `json:"param_fields"`
}

// NamespaceDoc is a map of namespace to RPCDoc
type NamespaceDoc map[string][]EnhancedRPCDocumentation

type Docs struct {
	RPCs NamespaceDoc
}

func (api *Server) GenerateDocs() (*Docs, error) {
	docs := new(Docs)
	docs.RPCs = make(NamespaceDoc)
	for _, rpc := range api.rpcs {
		docs.RPCs[rpc.Namespace] = append(docs.RPCs[rpc.Namespace], rpc.Documentation.Enhance())
	}
	return docs, nil
}

func (doc RPCDocumentation) Enhance() EnhancedRPCDocumentation {
	ed := EnhancedRPCDocumentation{
		RPCDocumentation: doc,
	}

	documentStruct(doc.Param, "", &ed.ParamFields)
	return ed
}

func documentStruct(strct any, prefix string, docfields *[]RPCDocumentationField) {
	t := reflect.TypeOf(strct)
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Type.String() == "time.Time" {
			*docfields = append(*docfields, RPCDocumentationField{
				Name:     prefix + jsonName(t.Field(i)),
				Type:     "timestamp",
				Required: isRequired(t.Field(i)),
			})
			continue
		}
		if strings.HasPrefix(t.Field(i).Type.String(), "int") {
			*docfields = append(*docfields, RPCDocumentationField{
				Name:     prefix + jsonName(t.Field(i)),
				Type:     "integer",
				Required: isRequired(t.Field(i)),
			})
			continue
		}
		if t.Field(i).Type.String() == "driplimit.Milliseconds" {
			*docfields = append(*docfields, RPCDocumentationField{
				Name:     prefix + jsonName(t.Field(i)),
				Type:     "integer (in milliseconds)",
				Required: isRequired(t.Field(i)),
			})
			continue
		}
		if t.Field(i).Type.Kind() == reflect.Struct {
			documentStruct(reflect.New(t.Field(i).Type).Elem().Interface(), jsonName(t.Field(i))+".", docfields)
			continue
		}
		*docfields = append(*docfields, RPCDocumentationField{
			Name:     prefix + jsonName(t.Field(i)),
			Type:     t.Field(i).Type.String(),
			Required: isRequired(t.Field(i)),
		})
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
