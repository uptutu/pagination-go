package pagination

import "reflect"

type Option func(*Page) error

var (
	_requestStructSearchFields = []string{"Page", "Pagination", "PageRequest", "PaginationRequest"}
)

// Parse a struct which have defined Page fields.
func Parse(req interface{}, options ...Option) (Page, error) {
	q := Page{
		defaultSize: 15,
	}
	v := reflect.ValueOf(req)
	for v.Type().Kind() != reflect.Struct {
		switch v.Type().Kind() {
		case reflect.Ptr:
			v = v.Elem()
		default:
			return q, ErrInvalidParseData
		}
	}

	if !((v.FieldByName("PageNum").IsValid() || v.FieldByName("Num").IsValid()) &&
		(v.FieldByName("PageSize").IsValid() || v.FieldByName("Size").IsValid())) {
		for _, word := range _requestStructSearchFields {
			if v.FieldByName(word).IsValid() {
				v = v.FieldByName(word)
				for v.Type().Kind() == reflect.Ptr {
					v = v.Elem()
				}
				break
			}
		}
	}

	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).CanInterface() {
			switch v.Type().Field(i).Name {
			case "PageNum", "Num":
				if v.Field(i).CanConvert(reflect.TypeOf(q.Num)) {
					q.Num = v.Field(i).Convert(reflect.TypeOf(q.Num)).Interface().(int)
					continue
				}
				return q, ErrInvalidPageNum
			case "PageSize", "Size":
				if v.Field(i).CanConvert(reflect.TypeOf(q.Size)) {
					q.Size = v.Field(i).Convert(reflect.TypeOf(q.Size)).Interface().(int)
					continue
				}
				return q, ErrInvalidPageSize
			case "OrderBy":
				if val, ok := v.Field(i).Interface().(string); ok {
					q.OrderBy = val
					continue
				}
				return q, ErrInvalidOrderBy

			case "IsDescending", "Descending":
				if val, ok := v.Field(i).Interface().(bool); ok {
					q.IsDescending = val
					continue
				}
				return q, ErrInvalidIsDescending
			case "Query", "SearchKey":
				if val, ok := v.Field(i).Interface().(string); ok {
					q.Query = val
					continue
				}
				return q, ErrInvalidSearchKey
			}
		}
	}

	for i := range options {
		if err := options[i](&q); err != nil {
			return q, err
		}
	}

	return q, nil
}

func WithDefaultSize(size int) Option {
	return func(p *Page) error {
		p.defaultSize = size
		return nil
	}
}
