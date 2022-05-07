# Pagination

简单的分页工具，只要结构体符合下面的要求，就可以使用该工具解析并填充。

1. 被解析的数据必须是结构体
    1. 结构体必须含有可导出的:
       1. int 类型： `PageNum` or `Num`
       2. int 类型： `PageSize` or `Size`
       3. string 类型： `OrderBy`
       4. bool 类型 `IsDescending` or `Descending` 
       5. string 类型: `Qeury` or `SearchKey` 字段
    
       比如说这样
       ```go
        type testRequest struct {
             PageNum      int
             PageSize     int
             OrderBy      string
             IsDescending bool
             KeyWords     string
             SearchKey    string
             CustomField  string
       }
       ```
    2. 这些必须字段可以成为一个新的结构体被名在 `Page` 或 `Pagination` 的字段下
    比如说这样：
      ```go
      type SearchDialogCasesRequest struct {
	          Page *PaginationRequest
      }
   
      // pagination.pb.go
       type PaginationRequest struct {
	          state         protoimpl.MessageState
	          sizeCache     protoimpl.SizeCache
	          unknownFields protoimpl.UnknownFields
			  
              PageNum      int64  `protobuf:"varint,1,opt,name=page_num,json=pageNum,proto3" json:"page_num,omitempty"`
	          PageSize     int64  `protobuf:"varint,2,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	          OrderBy      string `protobuf:"bytes,3,opt,name=order_by,json=orderBy,proto3" json:"order_by,omitempty"`
	          IsDescending bool   `protobuf:"varint,4,opt,name=is_descending,json=isDescending,proto3" json:"is_descending,omitempty"`
	          SearchKey    string `protobuf:"bytes,5,opt,name=search_key,json=searchKey,proto3" json:"search_key,omitempty"`
     }
      ```
3. 被导出的数据必须是结构体
    1. 被导出的结构体中必须含有 `PageNum`， `PageSize` 和 `Total` 字段

## Usage example

```go
import "github.com/uptutu/pagination"

// 自定义请求结构体
type ListSubscribeEntitiesRequest {
   PageNum      int64
   PageSize     int64
   OrderBy      string
   IsDescending bool
   Query    string
   MyData       interface{}
}

type ListSubscribeEntitiesResponse {
   Total    int64
   PageNum  int64
   LastPage int64
   PageSize int64
   // 自己的数据
   Data     string
}

func (s *SubscribeService) ListSubscribeEntities(ctx context.Context, req *pb.ListSubscribeEntitiesRequest) (*pb.ListSubscribeEntitiesResponse, error) {
   page, err := pagination.Parse(req)
   
   //输出类似于 {Num:10 Size:50 OrderBy:test IsDescending:true SearchKey:"search"}
   fmt.Println(page)
   
   // returned Bool 判读是否需要分页
   // return page_num > 0 && page_size > 0 
   page.Required()
   
   // 返回限制个数（page_size）
   page.Limit()
   
   // 返回分页请求下标（page_num-1 * page_size）
   page.Offset()
   
   // 是否倒序
   page.IsDescending
   
   // 搜索关键字
   page.Query
   
   // 排序字段
   page.OrderBy
   
   // 获取总数
   count := db.Find(data).Count()
   
   // 请务必手动填充总数
   page.SetTotal(count)
   
   resp := &ListSubscribeEntitiesResponse{}
   resp.Data = "一些数据"
   
   // 此后 resp 的分页信息将会自动被填充
   page.FillResponse(resp)
}

```

## func Required

Used to determine if the paging request passed meets the paging needs

Judgement conditions：

```go
func (p Page) Required() bool {
    return p.Num > 0 && p.Size > 0
}
```

## func Limit

if no limit set this will return the default value.

```go
func (p Page) Limit() int32 {
   if p.Size != 0 {
    return uint32(p.Size)
   }
   
   return uint32(p.defaultSize)
}
```

## func Offset

count the offset of the current page

```go
func (p Page) Offset() int32 {
   if p.Num <= 0 {
    return 0
   }
   return uint32((p.Num - 1) * p.Size)
}
```

## func FillResponse

Automatic padding of paginated data to match paginated responsive design.
