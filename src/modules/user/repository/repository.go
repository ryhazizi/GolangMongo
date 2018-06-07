package repository

import (
	"golangmongo/src/modules/user/model"
)

type UserRepository interface {
	Insert(*model.User) error
}
