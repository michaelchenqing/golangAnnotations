package rest

import (
	"fmt"
	"os"
	"testing"

	"io/ioutil"

	"github.com/MarcGrol/golangAnnotations/generator/rest/restAnnotation"
	"github.com/MarcGrol/golangAnnotations/model"
	"github.com/stretchr/testify/assert"
)

func TestGenerateForWeb(t *testing.T) {
	os.Remove("./testData/httpMyService.go")
	os.Remove("./testData/httpMyServiceHelpers_test.go")

	s := []model.Struct{
		{
			DocLines:    []string{"// @RestService( path = \"/api\")"},
			PackageName: "testData",
			Name:        "MyService",
			Operations:  []*model.Operation{},
		},
	}

	s[0].Operations = append(s[0].Operations,
		&model.Operation{
			DocLines:      []string{"// @RestOperation(path = \"/person\", method = \"GET\")"},
			Name:          "doit",
			RelatedStruct: &model.Field{TypeName: "MyService"},
			InputArgs: []model.Field{
				{Name: "uid", TypeName: "int"},
				{Name: "subuid", TypeName: "string"},
			},
			OutputArgs: []model.Field{
				{TypeName: "error"},
			},
		})

	err := Generate("testData", model.ParsedSources{Structs: s})
	assert.Nil(t, err)

	// check that generated files exisst
	_, err = os.Stat("./testData/httpMyService.go")
	assert.NoError(t, err)

	// check that generate code has 4 helper functions for MyStruct
	data, err := ioutil.ReadFile("./testData/httpMyService.go")
	assert.NoError(t, err)
	assert.Contains(t, string(data), "func (ts *MyService) HttpHandler() http.Handler {")
	assert.Contains(t, string(data), "func doit( service *MyService ) http.HandlerFunc {")

	// check that generated files exisst
	_, err = os.Stat("./testData/httpMyService.go")
	assert.NoError(t, err)

	// check that generate code has 4 helper functions for MyStruct
	data, err = ioutil.ReadFile("./testData/httpMyServiceHelpers_test.go")
	assert.NoError(t, err)
	assert.Contains(t, string(data), "func doitTestHelper")

	os.Remove("./testData/httpMyService.go")
	os.Remove("./testData/httpMyServiceHelpers_test.go")

}

func TestIsRestService(t *testing.T) {
	restAnnotation.Register()
	s := model.Struct{
		DocLines: []string{
			`//@RestService( path = "/api")`},
	}
	assert.True(t, IsRestService(s))
}

func TestGetRestServicePath(t *testing.T) {
	restAnnotation.Register()
	s := model.Struct{
		DocLines: []string{
			`//@RestService( path = "/api")`},
	}
	assert.Equal(t, "/api", GetRestServicePath(s))
}

func TestIsRestOperation(t *testing.T) {
	assert.True(t, IsRestOperation(createOper("GET")))
}

func TestGetRestOperationMethod(t *testing.T) {
	assert.Equal(t, "GET", GetRestOperationMethod(createOper("GET")))
}

func TestGetRestOperationPath(t *testing.T) {
	assert.Equal(t, "/api/person", GetRestOperationPath(createOper("DONTCARE")))
}

func TestHasInputGet(t *testing.T) {
	assert.False(t, HasInput(createOper("GET")))
}

func TestHasInputDelete(t *testing.T) {
	assert.False(t, HasInput(createOper("DELETE")))
}

func TestHasInputPost(t *testing.T) {
	assert.True(t, HasInput(createOper("POST")))
}

func TestHasInputPut(t *testing.T) {
	assert.True(t, HasInput(createOper("PUT")))
}

func TestGetInputArgTypeString(t *testing.T) {
	restAnnotation.Register()
	o := model.Operation{
		InputArgs: []model.Field{
			model.Field{TypeName: "string"},
		},
	}
	assert.Equal(t, "", GetInputArgType(o))
}

func TestGetInputArgTypePerson(t *testing.T) {
	assert.Equal(t, "Person", GetInputArgType(createOper("DONTCARE")))
}

func TestGetInputArgName(t *testing.T) {
	assert.Equal(t, "person", GetInputArgName(createOper("DONTCARE")))
}

func TestGetInputParamString(t *testing.T) {
	assert.Equal(t, "uid,person", GetInputParamString(createOper("DONTCARE")))
}

func TestHasOutput(t *testing.T) {
	assert.True(t, HasOutput(createOper("DONTCARE")))
}

func TestGetOutputArgType(t *testing.T) {
	assert.Equal(t, "Person", GetOutputArgType(createOper("DONTCARE")))
}

func TestIsPrimitiveTrue(t *testing.T) {
	f := model.Field{Name: "uid", TypeName: "string"}
	assert.True(t, IsPrimitive(f))
}

func TestIsPrimitiveFalse(t *testing.T) {
	f := model.Field{Name: "person", TypeName: "Person"}
	assert.False(t, IsPrimitive(f))
}

func TestIsNumberTrue(t *testing.T) {
	f := model.Field{Name: "uid", TypeName: "int"}
	assert.True(t, IsNumber(f))
}

func TestIsNumberFalse(t *testing.T) {
	f := model.Field{Name: "uid", TypeName: "string"}
	assert.False(t, IsNumber(f))
}

func createOper(method string) model.Operation {
	restAnnotation.Register()
	o := model.Operation{
		DocLines: []string{
			fmt.Sprintf("//@RestOperation( method = \"%s\", path = \"/api/person\")", method),
		},
		InputArgs: []model.Field{
			model.Field{Name: "uid", TypeName: "string"},
			model.Field{Name: "person", TypeName: "Person"},
		},
		OutputArgs: []model.Field{
			model.Field{TypeName: "Person"},
			model.Field{TypeName: "error"},
		},
	}
	return o
}