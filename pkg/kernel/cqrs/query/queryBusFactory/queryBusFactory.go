package queryBusFactory

import (
	userFinderPkg "github.com/rcrespodev/user_manager/pkg/app/user/application/querys/userFinder"
	"github.com/rcrespodev/user_manager/pkg/app/user/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/query"
)

type NewQueryBusCommand struct {
	UserRepository domain.UserRepository
}

func NewQueryBusInstance(queryBusCommand NewQueryBusCommand) *query.Bus {
	userFinder := userFinderPkg.NewUserFinder(queryBusCommand.UserRepository)
	userFinderQueryHandler := userFinderPkg.NewQueryHandler(userFinder)

	return query.NewBus(query.HandlersMap{
		query.FindUser: userFinderQueryHandler,
	})
}
