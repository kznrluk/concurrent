package association

import (
    "log"
    "context"
    "encoding/json"
    "github.com/redis/go-redis/v9"
    "github.com/totegamma/concurrent/x/util"
    "github.com/totegamma/concurrent/x/core"
    "github.com/totegamma/concurrent/x/stream"
    "github.com/totegamma/concurrent/x/message"
)

// Service is association service
type Service struct {
    rdb *redis.Client
    repo *Repository
    stream *stream.Service
    message *message.Service
}

// NewService is used for wire.go
func NewService(rdb *redis.Client, repo *Repository, stream *stream.Service, message *message.Service) *Service {
    return &Service{rdb, repo, stream, message}
}

// PostAssociation creates new association
func (s *Service) PostAssociation(ctx context.Context, objectStr string, signature string, streams []string, targetType string) error {
    ctx, childSpan := tracer.Start(ctx, "ServicePostAssociation")
    defer childSpan.End()

    var object signedObject
    err := json.Unmarshal([]byte(objectStr), &object)
    if err != nil {
        return err
    }

    if err := util.VerifySignature(objectStr, object.Signer, signature); err != nil {
        log.Println("verify signature err: ", err)
        return err
    }

    association := core.Association {
        Author: object.Signer,
        Schema: object.Schema,
        TargetID: object.Target,
        TargetType: targetType,
        Payload: objectStr,
        Signature: signature,
        Streams: streams,
    }

    s.repo.Create(ctx, &association)
    for _, stream := range association.Streams {
        s.stream.Post(ctx, stream, association.ID, "association", association.Author, "")
    }

    targetMessage, err := s.message.Get(ctx, association.TargetID)
    if err != nil {
        return err
    }
    for _, stream := range targetMessage.Streams {
        jsonstr, _ := json.Marshal(Event{
            Stream: stream,
            Type: "association",
            Action: "create",
            Body: Element {
                ID: association.TargetID,
            },
        })
        err := s.rdb.Publish(context.Background(), stream, jsonstr).Err()
        if err != nil {
            log.Printf("fail to publish message to Redis: %v", err)
        }
    }

    return nil
}

// Get returns an association by ID
func (s *Service) Get(ctx context.Context, id string) (core.Association, error) {
    ctx, childSpan := tracer.Start(ctx, "ServiceGet")
    defer childSpan.End()

    return s.repo.Get(ctx, id)
}

// GetOwn returns associations by author
func (s *Service) GetOwn(ctx context.Context, author string) ([]core.Association, error) {
    ctx, childSpan := tracer.Start(ctx, "ServiceGetOwn")
    defer childSpan.End()

    return s.repo.GetOwn(ctx, author)
}

// Delete deletes an association by ID
func (s *Service) Delete(ctx context.Context, id string) error {
    ctx, childSpan := tracer.Start(ctx, "ServiceDelete")
    defer childSpan.End()

    deleted, err := s.repo.Delete(ctx, id)
    if err != nil {
        return err
    }
    targetMessage, err := s.message.Get(ctx, deleted.TargetID)
    if err != nil {
        return err
    }
    for _, stream := range targetMessage.Streams {
        jsonstr, _ := json.Marshal(Event{
            Stream: stream,
            Type: "association",
            Action: "create",
            Body: Element {
                ID: deleted.TargetID,
            },
        })
        err := s.rdb.Publish(context.Background(), stream, jsonstr).Err()
        if err != nil {
            log.Printf("fail to publish message to Redis: %v", err)
            return err
        }
    }
    return nil
}

