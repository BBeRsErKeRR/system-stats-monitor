package command

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"sync"
)

func WithContext(ctx context.Context, name string, arg ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, arg...)

	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	if err := cmd.Start(); err != nil {
		return buf.Bytes(), err
	}

	if err := cmd.Wait(); err != nil {
		return buf.Bytes(), err
	}

	return buf.Bytes(), nil
}

func Stream(ctx context.Context, name string, arg ...string) (chan string, chan error) {
	cmd := exec.CommandContext(ctx, name, arg...)
	scannerErrors := make(chan error)
	scannerRes := make(chan string)
	var bufErr bytes.Buffer
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		scannerErrors <- err
		return scannerRes, scannerErrors
	}
	cmd.Stderr = &bufErr
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			scannerRes <- scanner.Text()
		}
	}()

	go func() {
		defer stdout.Close()
		defer wg.Done()
		if err := cmd.Start(); err != nil {
			scannerErrors <- err
			return
		}

		if err := cmd.Wait(); err != nil {
			scannerErrors <- fmt.Errorf("%s %w", bufErr.String(), err)
		}
	}()

	go func() {
		defer close(scannerRes)
		wg.Wait()
	}()

	return scannerRes, scannerErrors
}
