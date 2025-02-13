package db

type CRUD [T any] interface  {
	CreateByObject(T) error
	FindAll() ([]T , error)
	FindAllOrdered(fieldName string, order string) ([]T, error)
	FindById(id uint) (*T, error)
	FindAllByField(fieldName string, value interface{}, orderBy string, order string) ([]T, error)
	FindOneByUniqueField(fieldName string, value interface{}) (*T, error)
	UpdateByObject(T) error
	UpdateByField(string, interface{}, T) error
	DeleteById(id uint) error
	Search(ops ...searchOption) ([]T, error)
}

const (
	OrderAsc = "ASC"
	OrderDesc = "DESC"
)

//用于Search方法的参数
type SearchParams struct {
	Keyword string
	OrderBy string
	Order string
}

type searchOption struct {
	f func(*SearchParams)
}

func WithKeyword(keyword string) searchOption {
	return searchOption{
		f: func(params *SearchParams) {
			params.Keyword = keyword
		},
	}
}

func WithOrderBy(orderBy string) searchOption {
	return searchOption{
		f: func(params *SearchParams) {
			params.OrderBy = orderBy
		},
	}
}

func WithOrder(order string) searchOption {
	return searchOption{
		f: func(params *SearchParams) {
			params.Order = order
		},
	}
}