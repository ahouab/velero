/*
Copyright 2017 the Heptio Ark contributors.

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

package backup

import (
	"bytes"
	"net/url"
	"time"

	api "github.com/heptio/ark/pkg/apis/ark/v1"
	"github.com/heptio/ark/pkg/util/errcheck"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	kapiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	kscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

// podCommandExecutor is capable of executing a command in a container in a pod.
type podCommandExecutor interface {
	// executePodCommand executes a command in a container in a pod. If the command takes longer than
	// the specified timeout, an error is returned.
	executePodCommand(log logrus.FieldLogger, item map[string]interface{}, namespace, name, hookName string, hook *api.ExecHook) error
}

type poster interface {
	Post() *rest.Request
}

type defaultPodCommandExecutor struct {
	restClientConfig *rest.Config
	restClient       poster

	streamExecutorFactory streamExecutorFactory
}

// NewPodCommandExecutor creates a new podCommandExecutor.
func NewPodCommandExecutor(restClientConfig *rest.Config, restClient poster) podCommandExecutor {
	return &defaultPodCommandExecutor{
		restClientConfig: restClientConfig,
		restClient:       restClient,

		streamExecutorFactory: &defaultStreamExecutorFactory{},
	}
}

// executePodCommand uses the pod exec API to execute a command in a container in a pod. If the
// command takes longer than the specified timeout, an error is returned (NOTE: it is not currently
// possible to ensure the command is terminated when the timeout occurs, so it may continue to run
// in the background).
func (e *defaultPodCommandExecutor) executePodCommand(log logrus.FieldLogger, item map[string]interface{}, namespace, name, hookName string, hook *api.ExecHook) error {
	if item == nil {
		return errors.New("item is required")
	}
	if namespace == "" {
		return errors.New("namespace is required")
	}
	if name == "" {
		return errors.New("name is required")
	}
	if hookName == "" {
		return errors.New("hookName is required")
	}
	if hook == nil {
		return errors.New("hook is required")
	}

	if hook.Container == "" {
		if err := setDefaultHookContainer(item, hook); err != nil {
			return err
		}
	} else if err := ensureContainerExists(item, hook.Container); err != nil {
		return err
	}

	if len(hook.Command) == 0 {
		return errors.New("command is required")
	}

	switch hook.OnError {
	case api.HookErrorModeFail, api.HookErrorModeContinue:
		// use the specified value
	default:
		// default to fail
		hook.OnError = api.HookErrorModeFail
	}

	if hook.Timeout.Duration == 0 {
		hook.Timeout.Duration = defaultHookTimeout
	}

	hookLog := log.WithFields(
		logrus.Fields{
			"hookName":      hookName,
			"hookContainer": hook.Container,
			"hookCommand":   hook.Command,
			"hookOnError":   hook.OnError,
			"hookTimeout":   hook.Timeout,
		},
	)
	hookLog.Info("running exec hook")

	req := e.restClient.Post().
		Resource("pods").
		Namespace(namespace).
		Name(name).
		SubResource("exec")

	req.VersionedParams(&kapiv1.PodExecOptions{
		Container: hook.Container,
		Command:   hook.Command,
		Stdout:    true,
		Stderr:    true,
	}, kscheme.ParameterCodec)

	executor, err := e.streamExecutorFactory.NewSPDYExecutor(e.restClientConfig, "POST", req.URL())
	if err != nil {
		return err
	}

	var stdout, stderr bytes.Buffer

	streamOptions := remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
	}

	errCh := make(chan error)

	go func() {
		err = executor.Stream(streamOptions)
		errCh <- err
	}()

	var timeoutCh <-chan time.Time
	if hook.Timeout.Duration > 0 {
		timer := time.NewTimer(hook.Timeout.Duration)
		defer timer.Stop()
		timeoutCh = timer.C
	}

	select {
	case err = <-errCh:
	case <-timeoutCh:
		return errors.Errorf("timed out after %v", hook.Timeout.Duration)
	}

	hookLog.Infof("stdout: %s", stdout.String())
	hookLog.Infof("stderr: %s", stderr.String())

	return err
}

func ensureContainerExists(pod map[string]interface{}, container string) error {
	containers, found, err := unstructured.NestedSlice(pod, "spec", "containers")
	if err := errcheck.ErrOrNotFound(found, err, "unable to get spec.containers"); err != nil {
		return errors.WithStack(err)
	}

	for _, obj := range containers {
		c, ok := obj.(map[string]interface{})
		if !ok {
			return errors.Errorf("unexpected type for container %T", obj)
		}
		name, ok := c["name"].(string)
		if !ok {
			return errors.Errorf("unexpected type for container name %T", c["name"])
		}
		if name == container {
			return nil
		}
	}

	return errors.Errorf("no such container: %q", container)
}

func setDefaultHookContainer(pod map[string]interface{}, hook *api.ExecHook) error {
	containers, found, err := unstructured.NestedSlice(pod, "spec", "containers")
	if err := errcheck.ErrOrNotFound(found, err, "unable to get spec.containers"); err != nil {
		return errors.WithStack(err)
	}

	if len(containers) < 1 {
		return errors.New("need at least 1 container")
	}

	container, ok := containers[0].(map[string]interface{})
	if !ok {
		return errors.Errorf("unexpected type for container %T", pod)
	}

	name, ok := container["name"].(string)
	if !ok {
		return errors.Errorf("unexpected type for container name %T", container["name"])
	}
	hook.Container = name

	return nil
}

type streamExecutorFactory interface {
	NewSPDYExecutor(config *rest.Config, method string, url *url.URL) (remotecommand.Executor, error)
}

type defaultStreamExecutorFactory struct{}

func (f *defaultStreamExecutorFactory) NewSPDYExecutor(config *rest.Config, method string, url *url.URL) (remotecommand.Executor, error) {
	return remotecommand.NewSPDYExecutor(config, method, url)
}
