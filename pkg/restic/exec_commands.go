/*
Copyright 2018 the Velero contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package restic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	velerov1api "github.com/heptio/velero/pkg/apis/velero/v1"
	"github.com/heptio/velero/pkg/util/exec"
)

const backupProgressCheckInterval = 10 * time.Second

type backupStatusLine struct {
	MessageType string `json:"message_type"`
	// seen in status lines
	TotalBytes int64 `json:"total_bytes"`
	BytesDone  int64 `json:"bytes_done"`
	// seen in summary line at the end
	TotalBytesProcessed int64 `json:"total_bytes_processed"`
}

// GetSnapshotID runs a 'restic snapshots' command to get the ID of the snapshot
// in the specified repo matching the set of provided tags, or an error if a
// unique snapshot cannot be identified.
func GetSnapshotID(repoIdentifier, passwordFile string, tags map[string]string, env []string) (string, error) {
	cmd := GetSnapshotCommand(repoIdentifier, passwordFile, tags)
	if len(env) > 0 {
		cmd.Env = env
	}

	stdout, stderr, err := exec.RunCommand(cmd.Cmd())
	if err != nil {
		return "", errors.Wrapf(err, "error running command, stderr=%s", stderr)
	}

	type snapshotID struct {
		ShortID string `json:"short_id"`
	}

	var snapshots []snapshotID
	if err := json.Unmarshal([]byte(stdout), &snapshots); err != nil {
		return "", errors.Wrap(err, "error unmarshalling restic snapshots result")
	}

	if len(snapshots) != 1 {
		return "", errors.Errorf("expected one matching snapshot, got %d", len(snapshots))
	}

	return snapshots[0].ShortID, nil
}

// RunBackup runs a `restic backup` command and watches the output to provide
// progress updates to the caller.
func RunBackup(backupCmd *Command, log logrus.FieldLogger, updateFunc func(velerov1api.PodVolumeOperationProgress)) (string, string, error) {
	// buffers for copying command stdout/err output into
	stdoutBuf := new(bytes.Buffer)
	stderrBuf := new(bytes.Buffer)

	// create a channel to signal when to end the goroutine scanning for progress
	// updates
	quit := make(chan struct{})
	var wg sync.WaitGroup

	cmd := backupCmd.Cmd()
	cmd.Stdout = stdoutBuf
	cmd.Stderr = stderrBuf

	cmd.Start()

	wg.Add(1)
	go func() {
		ticker := time.NewTicker(backupProgressCheckInterval)
		for {
			select {
			case <-ticker.C:
				lastLine := getLastLine(stdoutBuf.Bytes())
				stat, err := decodeBackupStatusLine(lastLine)
				if err != nil {
					log.WithError(err).Errorf("error getting restic backup progress")
				}

				// if the line contains a non-empty bytes_done field, we can update the
				// caller with the progress
				if stat.BytesDone != 0 {
					updateFunc(velerov1api.PodVolumeOperationProgress{
						TotalBytes: stat.TotalBytes,
						BytesDone:  stat.BytesDone,
					})
				}
			case <-quit:
				ticker.Stop()
				wg.Done()
				return
			}
		}
	}()

	cmd.Wait()
	close(quit)
	wg.Wait()

	summary, err := getSummaryLine(stdoutBuf.Bytes())
	if err != nil {
		return stdoutBuf.String(), stderrBuf.String(), err
	}
	stat, err := decodeBackupStatusLine(summary)
	if err != nil {
		return stdoutBuf.String(), stderrBuf.String(), err
	}
	if stat.MessageType != "summary" {
		return stdoutBuf.String(), stderrBuf.String(), errors.WithStack(fmt.Errorf("error getting restic backup summary: %s", string(summary)))
	}

	// update progress to 100%
	updateFunc(velerov1api.PodVolumeOperationProgress{
		TotalBytes: stat.TotalBytesProcessed,
		BytesDone:  stat.TotalBytesProcessed,
	})

	return string(summary), stderrBuf.String(), nil
}

func decodeBackupStatusLine(lastLine []byte) (backupStatusLine, error) {
	var stat backupStatusLine
	if err := json.Unmarshal(lastLine, &stat); err != nil {
		return stat, errors.Wrapf(err, "unable to decode backup JSON line: %s", string(lastLine))
	}
	return stat, nil
}

// getLastLine returns the last line of a byte array. The string is assumed to
// have a newline at the end of it, so this returns the substring between the
// last two newlines and the index of the penultimate line.
func getLastLine(b []byte) []byte {
	// subslice the byte array to ignore the newline at the end of the string
	lastNewLineIdx := bytes.LastIndex(b[:len(b)-1], []byte("\n"))
	return b[lastNewLineIdx+1 : len(b)-1]
}

// getSummaryLine looks for the summary JSON line
// (`{"message_type:"summary",...`) in the restic backup command output. Due to
// an issue in Restic, this might not always be the last line
// (https://github.com/restic/restic/issues/2389). It returns an error if it
// can't be found.
func getSummaryLine(b []byte) ([]byte, error) {
	summaryLineIdx := bytes.LastIndex(b, []byte(`{"message_type":"summary"`))
	if summaryLineIdx < 0 {
		return nil, errors.New("unable to find summary in restic backup command output")
	}
	// find the end of the summary line
	newLineIdx := bytes.Index(b[summaryLineIdx:], []byte("\n"))
	if newLineIdx < 0 {
		return nil, errors.New("unable to get summary line from restic backup command output")
	}
	return b[summaryLineIdx : summaryLineIdx+newLineIdx], nil
}
