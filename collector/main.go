// Command factory-telemetry-collector reads the dedicated Factorio server's
// script-output JSONL stream. Phase 1 validates and prints events; SQLite
// persistence, durable checkpoints and save rollback recovery are Phase 2.
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

type Event struct {
	SchemaVersion int             `json:"schema_version"`
	WorldID       string          `json:"world_id"`
	RunID         string          `json:"run_id"`
	Sequence      uint64          `json:"sequence"`
	Tick          uint64          `json:"tick"`
	EventType     string          `json:"event_type"`
	Payload       json.RawMessage `json:"payload"`
}

func validate(line []byte) error {
	var e Event
	if err := json.Unmarshal(line, &e); err != nil {
		return fmt.Errorf("decode: %w", err)
	}
	if e.SchemaVersion != 1 || e.Sequence == 0 || e.EventType == "" || len(e.Payload) == 0 {
		return errors.New("invalid event envelope")
	}
	if e.WorldID != "unverified" || e.RunID != "unverified" {
		// Phase 1 has no trustworthy cross-save identifiers yet. Do not infer
		// production readiness from this validator.
		return errors.New("unexpected Phase 1 session marker")
	}
	return nil
}

// consume reads only newline-terminated records. An unfinished final line
// remains uncommitted and will be retried on a later poll.
func consume(file io.ReadSeeker, offset *int64, output io.Writer) (int, error) {
	if _, err := file.Seek(*offset, io.SeekStart); err != nil {
		return 0, err
	}
	reader := bufio.NewReader(file)
	count := 0
	for {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			return count, nil
		}
		if err != nil {
			return count, err
		}
		trimmed := bytes.TrimSpace(line)
		if len(trimmed) == 0 {
			*offset += int64(len(line))
			continue
		}
		if err := validate(trimmed); err != nil {
			return count, fmt.Errorf("invalid JSONL at byte %d: %w", *offset, err)
		}
		if _, err := output.Write(line); err != nil {
			return count, err
		}
		*offset += int64(len(line))
		count++
	}
}

func poll(path string, offset *int64, output io.Writer) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	stat, err := f.Stat()
	if err != nil {
		return 0, err
	}
	if stat.Size() < *offset {
		// WARNING: a reset means stream discontinuity. Phase 2 must reconcile.
		fmt.Fprintln(os.Stderr, "warning: stream truncated or replaced; resetting offset (state untrusted)")
		*offset = 0
	}
	return consume(f, offset, output)
}

func main() {
	input := flag.String("input", "", "absolute path to Factorio server script-output/factory-telemetry/events.jsonl")
	once := flag.Bool("once", false, "consume available complete records and exit")
	interval := flag.Duration("interval", time.Second, "poll interval")
	flag.Parse()
	if *input == "" {
		fmt.Fprintln(os.Stderr, "error: -input is required")
		os.Exit(2)
	}
	if *interval <= 0 {
		fmt.Fprintln(os.Stderr, "error: -interval must be positive")
		os.Exit(2)
	}

	var offset int64
	for {
		_, err := poll(*input, &offset, os.Stdout)
		if err != nil && !os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, "collector:", err)
			os.Exit(1)
		}
		if *once {
			if err != nil {
				fmt.Fprintln(os.Stderr, "collector:", err)
				os.Exit(1)
			}
			return
		}
		time.Sleep(*interval)
	}
}
