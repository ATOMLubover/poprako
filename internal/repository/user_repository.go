package repository

import intf "labelplus-next-web-be/internal/domain/repository"

type UserRepository interface {
	intf.Transactor
}
