//go:build linux
// +build linux

package bpstalkers

import (
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
	command "github.com/BBeRsErKeRR/system-stats-monitor/pkg/command"
)

func parseTcpDumpOut(line string) (interface{}, error) {
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

func getBps(ctx context.Context) (<-chan interface{}, <-chan error) {
	executor := func() (chan string, chan error) {
		var out chan string
		var errC chan error

		if os.Geteuid() == 0 {
			out, errC = command.Stream(ctx, "tcpdump", "-ntq", "-i", "any", "-Q", "inout", "-ttt", "-l")
		} else {
			out, errC = command.Stream(ctx, "sudo", "tcpdump", "-ntq", "-i", "any", "-Q", "inout", "-ttt", "-l")
		}

		return out, errC
	}

	parser := func(In <-chan string, errC <-chan error) (<-chan interface{}, <-chan error) {
		res := make(chan interface{})
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
					content, err := parseTcpDumpOut(stOut)
					if err != nil {
						resErr <- err
					}
					if content == nil {
						continue
					}
					res <- content
				}
			}
		}()
		return res, resErr
	}

	return parser(executor())
}
