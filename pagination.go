package pagination

import (
	"reflect"

	"github.com/pkg/errors"
)

var (
	ErrInvalidPageNum         = errors.New("invalid page number")
	ErrInvalidPageSize        = errors.New("invalid page size")
	ErrInvalidOrderBy         = errors.New("invalid order")
	ErrInvalidSearchKey       = errors.New("invalid search key")
	ErrInvalidIsDescending    = errors.New("invalid is descending")
	ErrInvalidParseData       = errors.New("invalid data type parsing")
	ErrInvalidResponse        = errors.New("invalid response")
	ErrResponseFieldType      = errors.New("response filed type")
	ErrResponseFieldUnsetable = errors.New("response field unsetable")
	ErrTryToSetinvalidNumber  = errors.New("try to set invalid number to field")
)

var _defaultResponseSearchingField = []string{"Page", "Pagination"}

type Page struct {
	Num          int
	Size         int
	OrderBy      string
	IsDescending bool
	Query        string
	Total        int
	defaultSize  int
}

func (p Page) Offset() int32 {
	if p.Num <= 0 {
		return 0
	}
	return int32((p.Num - 1) * p.Size)
}

func (p Page) Limit() int32 {
	if p.Size != 0 {
		return int32(p.Size)
	}

	if p.Num != 0 {
		return int32(p.defaultSize)
	}

	return 0
}

func (p Page) Required() bool {
	return p.Num > 0 && p.Size > 0
}

func (p *Page) SetTotal(total int) {
	p.Total = total
}

func (p Page) FillResponse(resp interface{}, fields ...string) error {
	// check response type and field
	keywords := _defaultResponseSearchingField
	if len(fields) != 0 {
		keywords = fields
	}
	v := reflect.ValueOf(resp)
	for v.Type().Kind() != reflect.Struct {
		switch v.Type().Kind() {
		case reflect.Ptr:
			if v.Elem().IsValid() && !v.IsNil() {
				v = v.Elem()
				break
			}
			return ErrInvalidResponse
		default:
			return ErrInvalidResponse
		}
	}

traverse:
	for {
		if v.FieldByName("Total").IsValid() &&
			(v.FieldByName("PageNum").IsValid() || v.FieldByName("Num").IsValid() || v.FieldByName("CurrentPage").IsValid() || v.FieldByName("CurrentPageNum").IsValid()) &&
			v.FieldByName("PageSize").IsValid() || v.FieldByName("Size").IsValid() {
			break
		}

		for _, word := range keywords {
			if v.FieldByName(word).IsValid() {
				v = v.FieldByName(word)
				for v.Type().Kind() == reflect.Ptr {
					if !v.Elem().IsValid() {
						return ErrInvalidResponse
					}
					v = v.Elem()
				}
				goto traverse
			}
		}
		return ErrInvalidResponse
	}

	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).CanInterface() {
			continue
		}
		f := v.Field(i)
		switch v.Type().Field(i).Name {
		case "Total":
			if err := SetNumber(f, p.Total); err != nil {
				return err
			}
		case "PageNum", "CurrentPage", "CurrentPageNum", "Num":
			if err := SetNumber(f, p.Num); err != nil {
				return err
			}
		case "LastPage":
			if p.Size == 0 {
				if err := SetNumber(f, 0); err != nil {
					return err
				}
				continue
			}
			lastPage := p.Total / p.Size
			if p.Total%p.Size == 0 {
				if err := SetNumber(f, lastPage); err != nil {
					return err
				}
				continue
			}
			if err := SetNumber(f, lastPage+1); err != nil {
				return err
			}

		case "PageSize", "Size":
			if p.Size == 0 {
				if err := SetNumber(f, p.Total); err != nil {
					return err
				}
				continue
			}
			if err := SetNumber(f, p.Size); err != nil {
				return err
			}
		}

	}
	return nil
}

func SetNumber(f reflect.Value, number interface{}) error {
	iv := reflect.ValueOf(number)
	if !f.CanSet() {
		return ErrResponseFieldUnsetable
	}
	if f.CanInt() {
		if iv.CanInt() {
			f.SetInt(reflect.ValueOf(number).Int())
			return nil
		}
		if iv.CanConvert(reflect.TypeOf(int(0))) {
			f.SetInt(iv.Convert(reflect.TypeOf(int(0))).Int())
			return nil
		}
		return ErrTryToSetinvalidNumber
	}
	if f.CanUint() {
		if iv.CanUint() {
			f.SetUint(iv.Uint())
			return nil
		}
		if iv.CanConvert(reflect.TypeOf(uint(0))) {
			f.SetUint(iv.Convert(reflect.TypeOf(uint(0))).Uint())
			return nil
		}
		return ErrTryToSetinvalidNumber
	}

	return ErrResponseFieldType
}
