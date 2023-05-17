//go:build linux
// +build linux

package protocoltalkers

import (
	"context"
	"os"
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

func getTalkers(ctx context.Context) (<-chan storage.ProtocolTalkerItem, <-chan error) {
	executor := func() (chan string, chan error) {
		var out chan string
		var errC chan error

		if os.Geteuid() == 0 {
			out, errC = command.Stream(ctx, "tcpdump", "-ntq", "-i", "any", "-Q", "inout", "-l")
		} else {
			out, errC = command.Stream(ctx, "sudo", "tcpdump", "-ntq", "-i", "any", "-Q", "inout", "-l")
		}

		return out, errC
	}

	parser := func(In <-chan string, errC <-chan error) (<-chan storage.ProtocolTalkerItem, <-chan error) {
		res := make(chan storage.ProtocolTalkerItem)
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

	return parser(executor())
}
