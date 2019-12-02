package test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"

	"github.com/jslyzt/jsoniter"
	"github.com/stretchr/testify/require"
)

func Test_disallowUnknownFields(t *testing.T) {
	should := require.New(t)
	type TestObject struct{}
	var obj TestObject
	decoder := jsoniter.NewDecoder(bytes.NewBufferString(`{"field1":100}`))
	decoder.DisallowUnknownFields()
	should.Error(decoder.Decode(&obj))
}

func Test_new_decoder(t *testing.T) {
	should := require.New(t)
	decoder1 := json.NewDecoder(bytes.NewBufferString(`[1][2]`))
	decoder2 := jsoniter.NewDecoder(bytes.NewBufferString(`[1][2]`))
	arr1 := []int{}
	should.Nil(decoder1.Decode(&arr1))
	should.Equal([]int{1}, arr1)
	arr2 := []int{}
	should.True(decoder1.More())
	buffered, _ := ioutil.ReadAll(decoder1.Buffered())
	should.Equal("[2]", string(buffered))
	should.Nil(decoder2.Decode(&arr2))
	should.Equal([]int{1}, arr2)
	should.True(decoder2.More())
	buffered, _ = ioutil.ReadAll(decoder2.Buffered())
	should.Equal("[2]", string(buffered))

	should.Nil(decoder1.Decode(&arr1))
	should.Equal([]int{2}, arr1)
	should.False(decoder1.More())
	should.Nil(decoder2.Decode(&arr2))
	should.Equal([]int{2}, arr2)
	should.False(decoder2.More())
}

func Test_use_number(t *testing.T) {
	should := require.New(t)
	decoder1 := json.NewDecoder(bytes.NewBufferString(`123`))
	decoder1.UseNumber()
	decoder2 := jsoniter.NewDecoder(bytes.NewBufferString(`123`))
	decoder2.UseNumber()
	var obj1 interface{}
	should.Nil(decoder1.Decode(&obj1))
	should.Equal(json.Number("123"), obj1)
	var obj2 interface{}
	should.Nil(decoder2.Decode(&obj2))
	should.Equal(json.Number("123"), obj2)
}

func Test_decoder_more(t *testing.T) {
	should := require.New(t)
	decoder := jsoniter.NewDecoder(bytes.NewBufferString("abcde"))
	should.True(decoder.More())
}

func Test_decoder_join(t *testing.T) {
	var (
		should                 = require.New(t)
		buf1, buf2, buf3       = &bytes.Buffer{}, &bytes.Buffer{}, &bytes.Buffer{}
		encod1, encod2, encod3 = json.NewEncoder(buf1), json.NewEncoder(buf2), json.NewEncoder(buf3)
	)

	encod1.SetEscapeHTML(false)
	encod2.SetEscapeHTML(false)
	encod3.SetEscapeHTML(false)

	encod1.Encode(map[string]interface{}{
		"key":  "encode1",
		"key1": 1,
		"arr":  []float64{1},
	})
	encod2.Encode(map[string]interface{}{
		"key":  "encode2",
		"key2": 2,
		"arr":  []int{2},
	})
	str1, str2 := buf1.String(), buf2.String()

	cfg := jsoniter.ConfigGinUse
	value := make(map[string]interface{})

	err := cfg.UnmarshalFromString(str1+str2, &value)
	should.Nil(err)

	encod3.Encode(value)
	buff := buf3.String()
	log.Println(buff)
}
