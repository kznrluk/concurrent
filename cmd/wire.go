//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/totegamma/concurrent/x/activitypub"
	"github.com/totegamma/concurrent/x/agent"
	"github.com/totegamma/concurrent/x/association"
	"github.com/totegamma/concurrent/x/auth"
	"github.com/totegamma/concurrent/x/character"
	"github.com/totegamma/concurrent/x/entity"
	"github.com/totegamma/concurrent/x/host"
	"github.com/totegamma/concurrent/x/message"
	"github.com/totegamma/concurrent/x/socket"
	"github.com/totegamma/concurrent/x/stream"
	"github.com/totegamma/concurrent/x/userkv"
	"github.com/totegamma/concurrent/x/util"
)

var hostHandlerProvider = wire.NewSet(host.NewHandler, host.NewService, host.NewRepository)
var entityHandlerProvider = wire.NewSet(entity.NewHandler, entity.NewService, entity.NewRepository)
var streamHandlerProvider = wire.NewSet(stream.NewHandler, stream.NewService, stream.NewRepository, entity.NewService, entity.NewRepository)
var messageHandlerProvider = wire.NewSet(message.NewHandler, message.NewService, message.NewRepository)
var characterHandlerProvider = wire.NewSet(character.NewHandler, character.NewService, character.NewRepository)
var associationHandlerProvider = wire.NewSet(association.NewHandler, association.NewService, association.NewRepository, message.NewService, message.NewRepository)
var userkvHandlerProvider = wire.NewSet(userkv.NewHandler, userkv.NewService, userkv.NewRepository)

func SetupMessageHandler(db *gorm.DB, rdb *redis.Client, config util.Config) *message.Handler {
	wire.Build(messageHandlerProvider, stream.NewService, stream.NewRepository, entity.NewService, entity.NewRepository)
	return &message.Handler{}
}

func SetupCharacterHandler(db *gorm.DB, config util.Config) *character.Handler {
	wire.Build(characterHandlerProvider)
	return &character.Handler{}
}

func SetupAssociationHandler(db *gorm.DB, rdb *redis.Client, config util.Config) *association.Handler {
	wire.Build(associationHandlerProvider, stream.NewService, stream.NewRepository, entity.NewService, entity.NewRepository)
	return &association.Handler{}
}

func SetupStreamHandler(db *gorm.DB, rdb *redis.Client, config util.Config) *stream.Handler {
	wire.Build(streamHandlerProvider)
	return &stream.Handler{}
}

func SetupHostHandler(db *gorm.DB, config util.Config) *host.Handler {
	wire.Build(hostHandlerProvider)
	return &host.Handler{}
}

func SetupEntityHandler(db *gorm.DB, config util.Config) *entity.Handler {
	wire.Build(entityHandlerProvider)
	return &entity.Handler{}
}

func SetupSocketHandler(rdb *redis.Client, config util.Config) *socket.Handler {
	wire.Build(socket.NewHandler, socket.NewService)
	return &socket.Handler{}
}

func SetupAgent(db *gorm.DB, rdb *redis.Client, config util.Config) *agent.Agent {
	wire.Build(agent.NewAgent, host.NewService, host.NewRepository, entity.NewService, entity.NewRepository)
	return &agent.Agent{}
}

func SetupAuthHandler(db *gorm.DB, config util.Config) *auth.Handler {
	wire.Build(auth.NewHandler, auth.NewService, entity.NewService, entity.NewRepository, host.NewService, host.NewRepository)
	return &auth.Handler{}
}

func SetupAuthService(db *gorm.DB, config util.Config) *auth.Service {
	wire.Build(auth.NewService, entity.NewService, entity.NewRepository, host.NewService, host.NewRepository)
	return &auth.Service{}
}

func SetupUserkvHandler(db *gorm.DB, rdb *redis.Client, config util.Config) *userkv.Handler {
	wire.Build(userkvHandlerProvider, entity.NewService, entity.NewRepository)
	return &userkv.Handler{}
}

func SetupActivitypubHandler(db *gorm.DB, rdb *redis.Client, config util.Config) *activitypub.Handler {
	wire.Build(activitypub.NewHandler, activitypub.NewRepository, message.NewService, message.NewRepository, entity.NewService, entity.NewRepository, stream.NewService, stream.NewRepository)
	return &activitypub.Handler{}
}
