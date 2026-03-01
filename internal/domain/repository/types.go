package repository

import "gorm.io/gorm"

type Executor = *gorm.DB

type QueryOption func(executor Executor) Executor

type Transactor interface {
	Execute(func(executor Executor) error) error
}
