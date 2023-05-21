// +build wireinject

package main

import (
    "gorm.io/gorm"
    "github.com/google/wire"
    "github.com/redis/go-redis/v9"

    "github.com/totegamma/concurrent/x/util"
    "github.com/totegamma/concurrent/x/host"
    "github.com/totegamma/concurrent/x/socket"
    "github.com/totegamma/concurrent/x/stream"
    "github.com/totegamma/concurrent/x/message"
    "github.com/totegamma/concurrent/x/character"
    "github.com/totegamma/concurrent/x/association"
)

var messageHandlerProvider = wire.NewSet(message.NewHandler, message.NewService, message.NewRepository)
var characterHandlerProvider = wire.NewSet(character.NewHandler, character.NewService, character.NewRepository)
var associationHandlerProvider = wire.NewSet(association.NewHandler, association.NewService, association.NewRepository)
var streamHandlerProvider = wire.NewSet(stream.NewHandler, stream.NewService, stream.NewRepository)
var hostHandlerProvider = wire.NewSet(host.NewHandler, host.NewService, host.NewRepository)

func SetupMessageHandler(db *gorm.DB, client *redis.Client, socket *socket.Service) message.Handler {
    wire.Build(messageHandlerProvider, stream.NewService, stream.NewRepository)
    return message.Handler{}
}

func SetupCharacterHandler(db *gorm.DB) character.Handler {
    wire.Build(characterHandlerProvider)
    return character.Handler{}
}

func SetupAssociationHandler(db *gorm.DB, client *redis.Client, socket *socket.Service) association.Handler {
    wire.Build(associationHandlerProvider, stream.NewService, stream.NewRepository)
    return association.Handler{}
}

func SetupStreamHandler(db *gorm.DB, client *redis.Client) stream.Handler {
    wire.Build(streamHandlerProvider)
    return stream.Handler{}
}

func SetupHostHandler(db *gorm.DB, config util.Config) host.Handler {
    wire.Build(hostHandlerProvider)
    return host.Handler{}
}

func SetupSocketHandler(socketService *socket.Service) *socket.Handler {
    wire.Build(socket.NewHandler)
    return &socket.Handler{}
}
