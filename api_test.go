package jmespath

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/jmespath/go-jmespath/internal/testify/assert"
)

func TestValidUncompiledExpressionSearches(t *testing.T) {
	assert := assert.New(t)
	var j = []byte(`{"foo": {"bar": {"baz": [0, 1, 2, 3, 4]}}}`)
	var d interface{}
	err := json.Unmarshal(j, &d)
	assert.Nil(err)
	result, err := Search("foo.bar.baz[2]", d)
	assert.Nil(err)
	assert.Equal(2.0, result)
}

func TestJSONNumber(t *testing.T) {
	assert := assert.New(t)
	var d interface{}
	dec := json.NewDecoder(strings.NewReader(`{"foo": [{"baz":0}, {"baz":1}, {"baz":-2}, {"baz":3}, {"baz":-4}]}`))
	dec.UseNumber()
	err := dec.Decode(&d)
	assert.Nil(err)
	r, err := Search("sort_by(foo, &baz)", d)
	assert.Nil(err)
	result, ok := r.([]interface{})
	assert.True(ok)
	assert.Equal("-4", result[0].(map[string]interface{})["baz"].(json.Number).String())
	assert.Equal("-2", result[1].(map[string]interface{})["baz"].(json.Number).String())
	assert.Equal("0", result[2].(map[string]interface{})["baz"].(json.Number).String())
	assert.Equal("1", result[3].(map[string]interface{})["baz"].(json.Number).String())
	assert.Equal("3", result[4].(map[string]interface{})["baz"].(json.Number).String())
}

func TestValidPrecompiledExpressionSearches(t *testing.T) {
	assert := assert.New(t)
	data := make(map[string]interface{})
	data["foo"] = "bar"
	precompiled, err := Compile("foo")
	assert.Nil(err)
	result, err := precompiled.Search(data)
	assert.Nil(err)
	assert.Equal("bar", result)
}

func TestInvalidPrecompileErrors(t *testing.T) {
	assert := assert.New(t)
	_, err := Compile("not a valid expression")
	assert.NotNil(err)
}

func TestInvalidMustCompilePanics(t *testing.T) {
	defer func() {
		r := recover()
		assert.NotNil(t, r)
	}()
	MustCompile("not a valid expression")
}
