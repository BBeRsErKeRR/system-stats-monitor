//go:build linux
// +build linux

package bpstalkers

import (
	"context"
	"strconv"
	"strings"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
	command "github.com/BBeRsErKeRR/system-stats-monitor/pkg/command"
)

func parseTCPDumpOut(line string) (interface{}, error) {
	fields := strings.Fields(line)
	length := len(fields)
	if length < 8 {
		return nil, nil
	}

	if fields[2] != "IP" {
		return nil, nil
	}

	nBytes, err := strconv.ParseFloat(fields[length-1], 64)
	if err != nil {
		return nil, err
	}
	if nBytes == 0 {
		return nil, nil
	}

	item := storage.BpsItem{
		Source:      strings.TrimSuffix(fields[5], ":"),
		Destination: fields[3],
		Numbers:     nBytes,
	}
	if len(fields) == 9 {
		item.Protocol = strings.ReplaceAll(fields[6], ",", "")
	} else {
		item.Protocol = fields[6]
	}

	return item, nil
}

func getBps(ctx context.Context) (<-chan storage.BpsItem, <-chan error) {
	parser := func(In <-chan string, errC <-chan error) (<-chan storage.BpsItem, <-chan error) {
		res := make(chan storage.BpsItem)
		resErr := make(chan error)
		go func() {
			defer close(res)
			defer close(resErr)
			for {
				select {
				case <-ctx.Done():
					return
				case err := <-errC:
					resErr <- err
					return
				case stOut, ok := <-In:
					if !ok {
						return
					}
					if strings.Contains(stOut, ", ack") {
						continue
					}
					content, err := parseTCPDumpOut(stOut)
					if err != nil {
						resErr <- err
					}
					if content == nil {
						continue
					}
					res <- content.(storage.BpsItem)
				}
			}
		}()
		return res, resErr
	}

	return parser(command.Stream(ctx, "sudo", "tcpdump", "-ntq", "-i", "any", "-Q", "inout", "-ttt", "-l"))
}

func checkCall(ctx context.Context) error {
	_, err := command.WithContext(ctx, "sudo", "tcpdump", "-h")
	return err
}
