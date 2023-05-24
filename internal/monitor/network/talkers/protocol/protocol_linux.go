//go:build linux
// +build linux

package protocoltalkers

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
	if length < 5 {
		return nil, nil
	}
	if fields[0] != "IP" {
		return nil, nil
	}
	nBytes, err := strconv.ParseFloat(fields[length-1], 64)
	if err != nil {
		return nil, err
	}

	if nBytes == 0 {
		return nil, nil
	}
	item := storage.ProtocolTalkerItem{
		Protocol:  strings.ReplaceAll(fields[4], ",", ""),
		SendBytes: nBytes,
	}
	if length == 7 {
		item.Protocol = strings.ReplaceAll(fields[4], ",", "")
	} else {
		item.Protocol = fields[4]
	}

	return item, nil
}

func getTalkers(ctx context.Context) (<-chan interface{}, <-chan error) {
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
					content, err := parseTCPDumpOut(stOut)
					if err != nil {
						resErr <- err
					}
					if content == nil {
						continue
					}
					res <- content.(storage.ProtocolTalkerItem)
				}
			}
		}()
		return res, resErr
	}

	return parser(command.Stream(ctx, "sudo", "tcpdump", "-ntq", "-i", "any", "-Q", "inout", "-l"))
}

func checkCall(ctx context.Context) error {
	_, err := command.WithContext(ctx, "sudo", "tcpdump", "-h")
	return err
}
