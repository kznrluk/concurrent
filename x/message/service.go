package message

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/pkg/errors"

	"github.com/totegamma/concurrent/client"
	"github.com/totegamma/concurrent/x/cdid"
	"github.com/totegamma/concurrent/x/core"
	"github.com/totegamma/concurrent/x/util"
)

type service struct {
	repo     Repository
	client   client.Client
	entity   core.EntityService
	timeline core.TimelineService
	key      core.KeyService
	config   util.Config
}

// NewService creates a new message service
func NewService(repo Repository, client client.Client, entity core.EntityService, timeline core.TimelineService, key core.KeyService, config util.Config) core.MessageService {
	return &service{repo, client, entity, timeline, key, config}
}

// Count returns the count number of messages
func (s *service) Count(ctx context.Context) (int64, error) {
	ctx, span := tracer.Start(ctx, "Message.Service.Count")
	defer span.End()

	return s.repo.Count(ctx)
}

// Get returns a message by ID
func (s *service) Get(ctx context.Context, id string, requester string) (core.Message, error) {
	ctx, span := tracer.Start(ctx, "Message.Service.Get")
	defer span.End()

	message, err := s.repo.Get(ctx, id)
	if err != nil {
		span.RecordError(err)
		return core.Message{}, err
	}

	canRead := true

	for _, timelineID := range message.Timelines {
		ok := s.timeline.HasReadAccess(ctx, timelineID, requester)
		if !ok {
			canRead = false
			break
		}
	}

	if !canRead {
		return core.Message{}, fmt.Errorf("no read access")
	}

	return message, nil
}

// GetWithOwnAssociations returns a message by ID with associations
func (s *service) GetWithOwnAssociations(ctx context.Context, id string, requester string) (core.Message, error) {
	ctx, span := tracer.Start(ctx, "Message.Service.GetWithOwnAssociations")
	defer span.End()

	message, err := s.repo.GetWithOwnAssociations(ctx, id, requester)
	if err != nil {
		span.RecordError(err)
		return core.Message{}, err
	}

	canRead := true

	for _, timelineID := range message.Timelines {
		ok := s.timeline.HasReadAccess(ctx, timelineID, requester)
		if !ok {
			canRead = false
			break
		}
	}

	if !canRead {
		return core.Message{}, fmt.Errorf("no read access")
	}

	return message, nil
}

// Create creates a new message
// It also posts the message to the timelines
func (s *service) Create(ctx context.Context, mode core.CommitMode, document string, signature string) (core.Message, error) {
	ctx, span := tracer.Start(ctx, "Message.Service.Create")
	defer span.End()

	created := core.Message{}

	var doc core.CreateMessage[any]
	err := json.Unmarshal([]byte(document), &doc)
	if err != nil {
		span.RecordError(err)
		return created, err
	}

	hash := util.GetHash([]byte(document))
	hash10 := [10]byte{}
	copy(hash10[:], hash[:10])
	signedAt := doc.SignedAt
	id := "m" + cdid.New(hash10, signedAt).String()

	signer, err := s.entity.Get(ctx, doc.Signer)
	if err != nil {
		span.RecordError(err)
		return core.Message{}, err
	}

	if signer.Domain == s.config.Concurrent.FQDN { // signerが自ドメイン管轄の場合、リソースを作成

		message := core.Message{
			ID:        id,
			Author:    doc.Signer,
			Schema:    doc.Schema,
			Document:  document,
			Signature: signature,
			Timelines: doc.Timelines,
		}

		if !doc.SignedAt.IsZero() {
			message.CDate = doc.SignedAt
		}

		created, err = s.repo.Create(ctx, message)
		if err != nil {
			span.RecordError(err)
			return message, err
		}

	}

	ispublic := true
	destinations := make(map[string][]string)
	for _, timeline := range doc.Timelines {

		ok := s.timeline.HasReadAccess(ctx, timeline, "")
		if !ok {
			ispublic = false
		}

		normalized, err := s.timeline.NormalizeTimelineID(ctx, timeline)
		if err != nil {
			span.RecordError(errors.Wrap(err, "failed to normalize timeline id"))
			continue
		}
		split := strings.Split(normalized, "@")
		if len(split) != 2 {
			span.RecordError(fmt.Errorf("invalid timeline id: %s", normalized))
			continue
		}
		domain := split[1]

		if _, ok := destinations[domain]; !ok {
			destinations[domain] = []string{}
		}
		destinations[domain] = append(destinations[domain], timeline)
	}

	sendDocument := ""
	sendSignature := ""
	if ispublic {
		sendDocument = document
		sendSignature = signature
	}

	for domain, timelines := range destinations {
		if domain == s.config.Concurrent.FQDN {
			// localなら、timelineのエントリを生成→Eventを発行
			for _, timeline := range timelines {
				posted, err := s.timeline.PostItem(ctx, timeline, core.TimelineItem{
					ResourceID: id,
					Owner:      doc.Signer,
					TimelineID: timeline,
				}, sendDocument, sendSignature)
				if err != nil {
					span.RecordError(errors.Wrap(err, "failed to post item"))
					continue
				}

				if mode == core.CommitModeExecute {
					// eventを放流
					event := core.Event{
						Timeline:  timeline,
						Item:      posted,
						Document:  document,
						Signature: signature,
					}

					err = s.timeline.PublishEvent(ctx, event)
					if err != nil {
						slog.ErrorContext(ctx, "failed to publish event", slog.String("error", err.Error()), slog.String("module", "timeline"))
						span.RecordError(errors.Wrap(err, "failed to publish event"))
						continue
					}
				}
			}
		} else if signer.Domain == s.config.Concurrent.FQDN && mode != core.CommitModeLocalOnlyExec { // ここでリソースを作成したなら、リモートにもリレー
			// remoteならdocumentをリレー
			packet := core.Commit{
				Document:  document,
				Signature: signature,
			}

			packetStr, err := json.Marshal(packet)
			if err != nil {
				span.RecordError(err)
				continue
			}
			s.client.Commit(ctx, domain, string(packetStr))
		}
	}

	// TODO: 実際に送信できたstreamの一覧を返すべき

	return created, nil
}

// Delete deletes a message by ID
// It also emits a delete event to the sockets
func (s *service) Delete(ctx context.Context, mode core.CommitMode, document, signature string) (core.Message, error) {
	ctx, span := tracer.Start(ctx, "Message.Service.Delete")
	defer span.End()

	var doc core.DeleteDocument
	err := json.Unmarshal([]byte(document), &doc)
	if err != nil {
		span.RecordError(err)
		return core.Message{}, err
	}

	deleteTarget, err := s.repo.Get(ctx, doc.Target)
	if err != nil {
		span.RecordError(err)
		return core.Message{}, err
	}

	if deleteTarget.Author != doc.Signer {
		return core.Message{}, fmt.Errorf("you are not authorized to perform this action")
	}

	deleted, err := s.repo.Delete(ctx, doc.Target)
	slog.DebugContext(ctx, fmt.Sprintf("deleted: %v", deleted), slog.String("module", "message"))

	if err != nil {
		span.RecordError(err)
		return core.Message{}, err
	}

	if mode != core.CommitModeLocalOnlyExec {
		for _, desttimeline := range deleted.Timelines {
			event := core.Event{
				Timeline:  desttimeline,
				Document:  document,
				Signature: signature,
			}
			err := s.timeline.PublishEvent(ctx, event)
			if err != nil {
				span.RecordError(err)
				return deleted, err
			}
		}
	}

	return deleted, err
}
