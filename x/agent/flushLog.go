package agent

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

func (a *agent) FlushLog(ctx context.Context) {
	ctx, span := tracer.Start(ctx, "Agent.FlushLog")
	defer span.End()

	slog.Info("flush log Start")

	allPath := filepath.Join(a.repositoryPath, "/all")

	err := os.MkdirAll(allPath, 0755)
	if err != nil {
		slog.Error("failed to create repository directory:", err)
		panic(err)
	}

	timestamp := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("%s.log", timestamp)

	alllogpath := filepath.Join(allPath, filename)
	storage, err := os.OpenFile(alllogpath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		slog.Error("failed to open repository log file:", err)
		panic(err)
	}
	defer storage.Close()

	// find last log entry

	lastID := a.rdb.Get(ctx, "agent:CommitLogLastFlushed").Val()
	if lastID == "" {
		lastID = "0"
	}

	entries, err := a.store.Since(ctx, lastID)

	var log string
	for _, entry := range entries {
		log += fmt.Sprintf("%s %s %s\n", entry.ID, entry.Owner, entry.Content)
	}

	storage.WriteString(log)

	// flush to each user log

	userlogPath := filepath.Join(a.repositoryPath, "/user")
	err = os.MkdirAll(userlogPath, 0755)
	if err != nil {
		slog.Error("failed to create repository directory:", err)
		panic(err)
	}

	bucket := make(map[string]string)
	for _, entry := range entries {
		if _, ok := bucket[entry.Owner]; !ok {
			bucket[entry.Owner] = ""
		}
		bucket[entry.Owner] += fmt.Sprintf("%s %s %s\n", entry.ID, entry.Owner, entry.Content)
	}

	for owner, log := range bucket {
		filename := fmt.Sprintf("%s.log", owner)
		logpath := filepath.Join(userlogPath, filename)
		userstore, err := os.OpenFile(logpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			slog.Error("failed to open repository log file:", err)
			continue
		}
		defer userstore.Close()

		userstore.WriteString(log)
	}

	// update last flushed log id
	if len(entries) > 0 {
		a.rdb.Set(ctx, "agent:CommitLogLastFlushed", entries[len(entries)-1].ID, 0)
	}

	slog.Info(fmt.Sprintf("%d entries flushed", len(entries)))
}