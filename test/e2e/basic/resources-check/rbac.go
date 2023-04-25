/*
Copyright the Velero contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the Licensm.
You may obtain a copy of the License at

    http://www.apachm.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the Licensm.
*/

/*
Copyright 2021 the Velero contributors.

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

package basic

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	. "github.com/vmware-tanzu/velero/test/e2e"
	. "github.com/vmware-tanzu/velero/test/e2e/test"
	. "github.com/vmware-tanzu/velero/test/e2e/util/k8s"
)

type RBACCase struct {
	TestCase
}

func (r *RBACCase) Init() error {
	rand.Seed(time.Now().UnixNano())
	UUIDgen, _ = uuid.NewRandom()
	r.BackupName = "backup-rbac" + UUIDgen.String()
	r.RestoreName = "restore-rbac" + UUIDgen.String()
	r.NSBaseName = "rabc-" + UUIDgen.String()
	r.NamespacesTotal = 1
	r.NSIncluded = &[]string{}
	for nsNum := 0; nsNum < r.NamespacesTotal; nsNum++ {
		createNSName := fmt.Sprintf("%s-%00000d", r.NSBaseName, nsNum)
		*r.NSIncluded = append(*r.NSIncluded, createNSName)
	}
	r.TestMsg = &TestMSG{
		Desc:      "Backup/restore of Namespaced Scoped and Cluster Scoped RBAC",
		Text:      "should be successfully backed up and restored",
		FailedMSG: "Failed to successfully backup and restore RBAC",
	}
	r.BackupArgs = []string{
		"create", "--namespace", VeleroCfg.VeleroNamespace, "backup", r.BackupName,
		"--include-namespaces", strings.Join(*r.NSIncluded, ","),
		"--default-volumes-to-fs-backup", "--wait",
	}

	r.RestoreArgs = []string{
		"create", "--namespace", VeleroCfg.VeleroNamespace, "restore", r.RestoreName,
		"--from-backup", r.BackupName, "--wait",
	}
	r.VeleroCfg = VeleroCfg
	r.Client = *r.VeleroCfg.ClientToInstallVelero
	return nil
}

func (r *RBACCase) CreateResources() error {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer ctxCancel()
	for nsNum := 0; nsNum < r.NamespacesTotal; nsNum++ {
		createNSName := fmt.Sprintf("%s-%00000d", r.NSBaseName, nsNum)
		fmt.Printf("Creating namespaces ...%s\n", createNSName)
		if err := CreateNamespace(ctx, r.Client, createNSName); err != nil {
			return errors.Wrapf(err, "Failed to create namespace %s", createNSName)
		}
		serviceAccountName := fmt.Sprintf("service-account-%s-%00000d", r.NSBaseName, nsNum)
		fmt.Printf("Creating service account ...%s\n", createNSName)
		if err := CreateServiceAccount(ctx, r.Client, createNSName, serviceAccountName); err != nil {
			return errors.Wrapf(err, "Failed to create service account %s", serviceAccountName)
		}
		clusterRoleName := fmt.Sprintf("clusterrole-%s-%00000d", r.NSBaseName, nsNum)
		clusterRoleBindingName := fmt.Sprintf("clusterrolebinding-%s-%00000d", r.NSBaseName, nsNum)
		if err := CreateRBACWithBindingSA(ctx, r.Client, createNSName, serviceAccountName, clusterRoleName, clusterRoleBindingName); err != nil {
			return errors.Wrapf(err, "Failed to create cluster role %s with role binding %s", clusterRoleName, clusterRoleBindingName)
		}
	}
	return nil
}

func (r *RBACCase) Verify() error {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer ctxCancel()
	for nsNum := 0; nsNum < r.NamespacesTotal; nsNum++ {
		checkNSName := fmt.Sprintf("%s-%00000d", r.NSBaseName, nsNum)
		checkServiceAccountName := fmt.Sprintf("service-account-%s-%00000d", r.NSBaseName, nsNum)
		checkClusterRoleName := fmt.Sprintf("clusterrole-%s-%00000d", r.NSBaseName, nsNum)
		checkClusterRoleBindingName := fmt.Sprintf("clusterrolebinding-%s-%00000d", r.NSBaseName, nsNum)
		checkNS, err := GetNamespace(ctx, r.Client, checkNSName)
		if err != nil {
			return errors.Wrapf(err, "Could not retrieve test namespace %s", checkNSName)
		}
		if checkNS.Name != checkNSName {
			return errors.Errorf("Retrieved namespace for %s has name %s instead", checkNSName, checkNS.Name)
		}

		//getting service account from the restore
		checkSA, err := GetServiceAccount(ctx, r.Client, checkNSName, checkServiceAccountName)

		if err != nil {
			return errors.Wrapf(err, "Could not retrieve test service account %s", checkSA)
		}

		if checkSA.Name != checkServiceAccountName {
			return errors.Errorf("Retrieved service account for %s has name %s instead", checkServiceAccountName, checkSA.Name)
		}

		//getting cluster role from the restore
		checkClusterRole, err := GetClusterRole(ctx, r.Client, checkClusterRoleName)

		if err != nil {
			return errors.Wrapf(err, "Could not retrieve test cluster role %s", checkClusterRole)
		}

		if checkSA.Name != checkServiceAccountName {
			return errors.Errorf("Retrieved cluster role for %s has name %s instead", checkClusterRoleName, checkClusterRole.Name)
		}

		//getting cluster role binding from the restore
		checkClusterRoleBinding, err := GetClusterRoleBinding(ctx, r.Client, checkClusterRoleBindingName)

		if err != nil {
			return errors.Wrapf(err, "Could not retrieve test cluster role binding %s", checkClusterRoleBinding)
		}

		if checkClusterRoleBinding.Name != checkClusterRoleBindingName {
			return errors.Errorf("Retrieved cluster role binding for %s has name %s instead", checkClusterRoleBindingName, checkClusterRoleBinding.Name)
		}

		//check if the role binding maps to service account
		checkSubjects := checkClusterRoleBinding.Subjects[0].Name

		if checkSubjects != checkServiceAccountName {
			return errors.Errorf("Retrieved cluster role binding for %s has name %s instead", checkServiceAccountName, checkSubjects)
		}
	}
	return nil
}

func (r *RBACCase) Destroy() error {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer ctxCancel()
	//cleanup clusterrole
	err := CleanupClusterRole(ctx, r.Client, r.NSBaseName)
	if err != nil {
		return errors.Wrap(err, "Could not cleanup clusterroles")
	}

	//cleanup cluster rolebinding
	err = CleanupClusterRoleBinding(ctx, r.Client, r.NSBaseName)
	if err != nil {
		return errors.Wrap(err, "Could not cleanup clusterrolebindings")
	}

	err = CleanupNamespacesWithPoll(ctx, r.Client, r.NSBaseName)
	if err != nil {
		return errors.Wrap(err, "Could cleanup retrieve namespaces")
	}

	return nil
}

func (r *RBACCase) Clean() error {
	return r.Destroy()
}
