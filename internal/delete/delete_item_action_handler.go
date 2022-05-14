/*
Copyright 2020 the Velero contributors.

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

package delete

import (
	"context"
	"encoding/json"
	"io"

	"github.com/vmware-tanzu/velero/pkg/features"
	"github.com/vmware-tanzu/velero/pkg/plugin/clientmgmt"
	isv1 "github.com/vmware-tanzu/velero/pkg/plugin/velero/item_snapshotter/v1"
	"github.com/vmware-tanzu/velero/pkg/volume"

	"github.com/vmware-tanzu/velero/pkg/plugin/framework"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"

	velerov1api "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
	"github.com/vmware-tanzu/velero/pkg/archive"
	"github.com/vmware-tanzu/velero/pkg/discovery"
	"github.com/vmware-tanzu/velero/pkg/plugin/velero"
	"github.com/vmware-tanzu/velero/pkg/util/filesystem"
)

// Context provides the necessary environment to run DeleteItemAction and ItemSnapshotter plugins
type Context struct {
	Backup                    *velerov1api.Backup
	BackupReader              io.Reader
	Filesystem                filesystem.Interface
	Log                       logrus.FieldLogger
	DiscoveryHelper           discovery.Helper
	DeleteItemResolvedActions []framework.DeleteItemResolvedAction
	ItemSnapshots             []*volume.ItemSnapshot
	ItemSnapshotters          []isv1.ItemSnapshotter
}

func InvokeDeleteActions(ctx *Context) error {
	// No actions installed and no item snapshots means we don't have to continue;
	// just do the backup deletion without worrying about plugins.
	if len(ctx.DeleteItemResolvedActions) == 0 && len(ctx.ItemSnapshots) == 0 {
		ctx.Log.Debug("No delete item actions or item snapshots present, proceeding with rest of backup deletion process")
		return nil
	}

	// get items out of backup tarball into a temp directory
	dir, err := archive.NewExtractor(ctx.Log, ctx.Filesystem).UnzipAndExtractBackup(ctx.BackupReader)
	if err != nil {
		return errors.Wrapf(err, "error extracting backup")
	}
	defer ctx.Filesystem.RemoveAll(dir)
	ctx.Log.Debugf("Downloaded and extracted the backup file to: %s", dir)

	backupResources, err := archive.NewParser(ctx.Log, ctx.Filesystem).Parse(dir)
	if existErr := errors.Is(err, archive.ErrNotExist); existErr {
		ctx.Log.Debug("ignore invoking delete item actions: ", err)
		return nil
	} else if err != nil {
		return errors.Wrapf(err, "error parsing backup %q", dir)
	}
	processdResources := sets.NewString()

	for resource := range backupResources {
		groupResource := schema.ParseGroupResource(resource)

		// We've already seen this group/resource, so don't process it again.
		if processdResources.Has(groupResource.String()) {
			continue
		}

		// Get a list of all items that exist for this resource
		resourceList := backupResources[groupResource.String()]
		if resourceList == nil {
			continue
		}

		// Iterate over all items, grouped by namespace.
		for namespace, items := range resourceList.ItemsByNamespace {
			nsLog := ctx.Log.WithField("namespace", namespace)
			nsLog.Info("Starting to check for items in namespace")

			// Filter applicable actions based on namespace only once per namespace.
			actions := ctx.getApplicableActions(groupResource, namespace)

			// Process individual items from the backup
			for _, item := range items {
				itemPath := archive.GetItemFilePath(dir, resource, namespace, item)

				// obj is the Unstructured item from the backup
				obj, err := archive.Unmarshal(ctx.Filesystem, itemPath)
				if err != nil {
					return errors.Wrapf(err, "Could not unmarshal item: %v", item)
				}

				itemLog := nsLog.WithField("item", obj.GetName())
				itemLog.Infof("invoking DeleteItemAction plugins")

				for _, action := range actions {
					if !action.Selector.Matches(labels.Set(obj.GetLabels())) {
						continue
					}
					err = action.DeleteItemAction.Execute(&velero.DeleteItemActionExecuteInput{
						Item:   obj,
						Backup: ctx.Backup,
					})
					// Since we want to keep looping even on errors, log them instead of just returning.
					if err != nil {
						itemLog.WithError(err).Error("plugin error")
					}
				}
			}
		}
	}

	if features.IsEnabled(velerov1api.UploadProgressFeatureFlag) {
		ctx.Log.Info("Deleting item snapshots")
		// Handle item snapshots
		for _, snapshot := range ctx.ItemSnapshots {
			rid := velero.ResourceIdentifier{}
			json.Unmarshal([]byte(snapshot.Spec.ResourceIdentifier), &rid)
			itemLogger := ctx.Log.WithFields(logrus.Fields{
				"namespace":  rid.Namespace,
				"resource":   rid.Resource,
				"item":       rid.Name,
				"snapshotID": snapshot.Status.ProviderSnapshotID,
			})
			itemLogger.Info("Deleting item snapshot")

			itemSnapshotter := clientmgmt.ItemSnapshotterForSnapshot(snapshot, ctx.ItemSnapshotters)

			itemPath := archive.GetItemFilePath(dir, rid.Resource, rid.Namespace, rid.Name)

			// obj is the Unstructured item from the backup
			obj, err := archive.Unmarshal(ctx.Filesystem, itemPath)
			if err != nil {
				itemLogger.WithError(err).Errorf("Could not unmarshal item: %v", rid)
				continue
			}
			dsi := isv1.DeleteSnapshotInput{
				SnapshotID:       snapshot.Status.ProviderSnapshotID,
				ItemFromBackup:   obj,
				SnapshotMetadata: snapshot.Status.Metadata,
				Params:           nil, // TBD
			}
			if err := itemSnapshotter.DeleteSnapshot(context.TODO(), &dsi); err != nil {
				itemLogger.WithError(err).Error("Error deleting snapshot")
			}
		}
	}
	return nil
}

// getApplicableActions takes resolved DeleteItemActions and filters them for a given group/resource and namespace.
func (ctx *Context) getApplicableActions(groupResource schema.GroupResource, namespace string) []framework.DeleteItemResolvedAction {
	var actions []framework.DeleteItemResolvedAction
	for _, action := range ctx.DeleteItemResolvedActions {
		if action.ShouldUse(groupResource, namespace, nil, ctx.Log) {
			actions = append(actions, action)
		}
	}
	return actions
}
