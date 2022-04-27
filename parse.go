package pagination

import "reflect"

type Option func(*Page) error

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

	if !(v.FieldByName("PageNum").IsValid() && v.FieldByName("PageSize").IsValid()) {
		keywords := []string{"Page", "Pagination", "PageRequest", "PaginationRequest"}
		for _, word := range keywords {
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
			case "PageNum":
				if val, ok := v.Field(i).Interface().(int); ok {
					q.Num = val
					continue
				}
				if val, ok := v.Field(i).Interface().(int32); ok {
					q.Num = int(val)
					continue
				}
				if val, ok := v.Field(i).Interface().(int64); ok {
					q.Num = int(val)
					continue
				}
				return q, ErrInvalidPageNum
			case "PageSize":
				if val, ok := v.Field(i).Interface().(int); ok {
					q.Size = val
					continue
				}
				if val, ok := v.Field(i).Interface().(int32); ok {
					q.Size = int(val)
					continue
				}
				if val, ok := v.Field(i).Interface().(int64); ok {
					q.Size = int(val)
					continue
				}
				return q, ErrInvalidPageSize
			case "OrderBy":
				if val, ok := v.Field(i).Interface().(string); ok {
					q.OrderBy = val
				} else {
					return q, ErrInvalidOrderBy
				}
			case "IsDescending":
				if val, ok := v.Field(i).Interface().(bool); ok {
					q.IsDescending = val
				} else {
					return q, ErrInvalidIsDescending
				}
			case "SearchKey":
				if val, ok := v.Field(i).Interface().(string); ok {
					q.SearchKey = val
				} else {
					return q, ErrInvalidSearchKey
				}
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
