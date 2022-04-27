package pagination

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testRequest struct {
	PageNum      int
	PageSize     int
	OrderBy      string
	IsDescending bool
	KeyWords     string
	SearchKey    string
	CustomField  string
}

type SearchDialogCasesRequest struct {
	Page *PaginationRequest
}

type SearchDialogCasesResponse struct {
	Page *PaginationResponse
	data string
}

type uintFieldData struct {
	PageNum      uint
	PageSize     uint
	OrderBy      string
	IsDescending bool
	KeyWords     string
	SearchKey    string
	CustomField  string
}

var (
	pbData = SearchDialogCasesRequest{
		Page: &PaginationRequest{
			PageNum:      10,
			PageSize:     50,
			OrderBy:      "id",
			IsDescending: true,
			SearchKey:    "search",
		},
	}
	customData = testRequest{
		PageNum:      10,
		PageSize:     50,
		OrderBy:      "id",
		IsDescending: true,
		KeyWords:     "key",
		SearchKey:    "search",
		CustomField:  "my data",
	}

	uintData = uintFieldData{
		PageNum:      10,
		PageSize:     50,
		OrderBy:      "id",
		IsDescending: true,
		KeyWords:     "key",
		SearchKey:    "search",
		CustomField:  "my data",
	}

	targetPage = Page{
		Num:          10,
		Size:         50,
		OrderBy:      "id",
		IsDescending: true,
		SearchKey:    "search",
		defaultSize:  15,
	}
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		data     interface{}
		options  []Option
		excepted Page
	}{
		{
			name:     "test pb struct",
			data:     pbData,
			excepted: targetPage,
		},
		{
			name:     "test custom data",
			data:     customData,
			excepted: targetPage,
		},
		{
			name:     "test pb ptr",
			data:     &pbData,
			excepted: targetPage,
		},
		{
			name:     "test custom ptr",
			data:     &customData,
			excepted: targetPage,
		},
		{
			name:     "test uint data",
			data:     uintData,
			excepted: targetPage,
		},
		{
			name:     "test uint ptr",
			data:     &uintData,
			excepted: targetPage,
		},
		{
			name:     "test custom ptr with default options",
			data:     &customData,
			options:  []Option{WithDefaultSize(20)},
			excepted: targetPage,
		},
	}
	for _, test := range tests {
		page, err := Parse(test.data, test.options...)
		assert.NoError(t, err)
		if test.name == "test custom ptr with default options" {
			copyOne := targetPage
			copyOne.defaultSize = 20
			assert.Equal(t, copyOne, page)
		} else {
			assert.Equal(t, test.excepted, page)
		}
	}
}

func TestInvalidDataTypeParse(t *testing.T) {
	tests := []struct {
		name     string
		data     interface{}
		excepted error
	}{
		{
			name:     "map type",
			data:     make(map[string]interface{}),
			excepted: ErrInvalidParseData,
		},
		{
			name:     "string",
			data:     "string",
			excepted: ErrInvalidParseData,
		},
		{
			name:     "int",
			data:     32,
			excepted: ErrInvalidParseData,
		},
		{
			name:     "float",
			data:     float32(32.0),
			excepted: ErrInvalidParseData,
		},
		{
			name:     "bool",
			data:     true,
			excepted: ErrInvalidParseData,
		},
		{
			name:     "struct",
			data:     struct{ invalidField string }{"test"},
			excepted: nil,
		},
		{
			name:     "mismatch struct type",
			data:     struct{ PageNum string }{"PageNum"},
			excepted: ErrInvalidPageNum,
		},
	}
	for _, test := range tests {
		_, err := Parse(test.data)
		assert.Equal(t, test.excepted, err)
	}
}

func TestPage_Offset(t *testing.T) {
	tests := []struct {
		name     string
		req      testRequest
		excepted int32
	}{
		{
			name:     "test offset",
			req:      testRequest{PageNum: 2, PageSize: 450},
			excepted: 450,
		},
		{
			name:     "test offset no page",
			req:      testRequest{},
			excepted: 0,
		},
		{
			name:     "test offset page 1 per page 30",
			req:      testRequest{PageNum: 1, PageSize: 30},
			excepted: 0,
		},
		{
			name:     "test offset page 2 per page 30",
			req:      testRequest{PageNum: 2, PageSize: 30},
			excepted: 30,
		},
	}

	for _, test := range tests {
		page, err := Parse(&test.req)
		assert.NoError(t, err)
		assert.Equal(t, test.excepted, page.Offset())
	}
}

func TestPage_Limit(t *testing.T) {
	tests := []struct {
		name     string
		req      testRequest
		excepted int32
	}{
		{
			name:     "test limit",
			req:      customData,
			excepted: 50,
		},
		{
			name:     "test limit with default when page query",
			req:      testRequest{PageNum: 1},
			excepted: 15,
		},
		{
			name:     "test limit with query all data",
			req:      testRequest{},
			excepted: 0,
		},
	}

	for _, test := range tests {
		page, err := Parse(&test.req)
		assert.NoError(t, err)
		assert.Equal(t, test.excepted, page.Limit())
	}
}

func TestPage_Required(t *testing.T) {
	tests := []struct {
		name     string
		req      testRequest
		excepted bool
	}{
		{
			name:     "test Required",
			req:      customData,
			excepted: true,
		},
		{
			name:     "test no required",
			req:      testRequest{},
			excepted: false,
		},
	}

	for _, test := range tests {
		page, err := Parse(&test.req)
		assert.NoError(t, err)
		assert.Equal(t, test.excepted, page.Required())
	}
}

type ListResponse struct {
	Total    int64  `json:"total"`
	PageNum  int64  `json:"page_num"`
	PageSize int64  `json:"page_size"`
	data     string `json:"data"`
}

func TestPage_FillResponse(t *testing.T) {
	tests := []struct {
		name     string
		req      testRequest
		count    int
		excepted ListResponse
	}{
		{
			name:  "test Required",
			req:   customData,
			count: 500,
			excepted: ListResponse{
				Total:    500,
				PageSize: int64(customData.PageSize),
				PageNum:  int64(customData.PageNum),
			},
		}, {
			name:  "test Required and more data",
			req:   customData,
			count: 501,
			excepted: ListResponse{
				Total:    501,
				PageSize: int64(customData.PageSize),
				PageNum:  int64(customData.PageNum),
			},
		},
		{
			name:  "test no required",
			req:   testRequest{},
			count: 500,
			excepted: ListResponse{
				Total:    500,
				PageSize: 500,
				PageNum:  0,
			},
		},
	}

	for _, test := range tests {
		resp := &ListResponse{}
		page, err := Parse(&test.req)
		page.SetTotal(test.count)
		err = page.FillResponse(resp)
		assert.NoError(t, err)
		assert.NoError(t, err)
		assert.True(t, reflect.DeepEqual(test.excepted, *resp))
	}

	// pb.Response test
	pbResp := &SearchDialogCasesResponse{
		Page: &PaginationResponse{},
	}
	pbReq := &SearchDialogCasesRequest{
		Page: &PaginationRequest{
			PageNum:  1,
			PageSize: 10,
		},
	}
	page, err := Parse(pbReq)
	assert.NoError(t, err)
	page.SetTotal(500)
	err = page.FillResponse(pbResp)
	assert.NoError(t, err)
	assert.Equal(t, int64(10), pbResp.Page.PageSize)
	assert.Equal(t, int64(1), pbResp.Page.PageNum)
	assert.Equal(t, int64(500), pbResp.Page.Total)

}

func TestSetNumber(t *testing.T) {
	var oneInt int
	var oneInt8 int8
	var oneInt16 int16
	var oneInt32 int32
	var oneInt64 int64

	err := SetNumber(reflect.ValueOf(&oneInt).Elem(), 10)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, oneInt)

	err = SetNumber(reflect.ValueOf(&oneInt8).Elem(), 10)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, oneInt8)

	err = SetNumber(reflect.ValueOf(&oneInt16).Elem(), 10)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, oneInt16)

	err = SetNumber(reflect.ValueOf(&oneInt32).Elem(), 10)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, oneInt32)

	err = SetNumber(reflect.ValueOf(&oneInt64).Elem(), 10)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, oneInt64)
}
