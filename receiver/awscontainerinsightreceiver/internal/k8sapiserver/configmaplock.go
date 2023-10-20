// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package k8sapiserver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awscontainerinsightreceiver/internal/k8sapiserver"

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
)

type ConfigMapLock struct {
	// ConfigMapMeta should contain a Name and a Namespace of a
	// ConfigMapMeta object that the LeaderElector will attempt to lead.
	ConfigMapMeta metav1.ObjectMeta
	Client        corev1client.ConfigMapsGetter
	LockConfig    resourcelock.ResourceLockConfig
	cm            *v1.ConfigMap
}

// Get returns the election record from a ConfigMap Annotation
func (cml *ConfigMapLock) Get(ctx context.Context) (*resourcelock.LeaderElectionRecord, []byte, error) {
	var record resourcelock.LeaderElectionRecord
	var err error
	cml.cm, err = cml.Client.ConfigMaps(cml.ConfigMapMeta.Namespace).Get(ctx, cml.ConfigMapMeta.Name, metav1.GetOptions{})
	if err != nil {
		return nil, nil, err
	}
	if cml.cm.Annotations == nil {
		cml.cm.Annotations = make(map[string]string)
	}
	recordStr, found := cml.cm.Annotations[resourcelock.LeaderElectionRecordAnnotationKey]
	recordBytes := []byte(recordStr)
	if found {
		if err := json.Unmarshal(recordBytes, &record); err != nil {
			return nil, nil, err
		}
	}
	return &record, recordBytes, nil
}

// Create attempts to create a LeaderElectionRecord annotation
func (cml *ConfigMapLock) Create(ctx context.Context, ler resourcelock.LeaderElectionRecord) error {
	recordBytes, err := json.Marshal(ler)
	if err != nil {
		return err
	}
	cml.cm, err = cml.Client.ConfigMaps(cml.ConfigMapMeta.Namespace).Create(ctx, &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cml.ConfigMapMeta.Name,
			Namespace: cml.ConfigMapMeta.Namespace,
			Annotations: map[string]string{
				resourcelock.LeaderElectionRecordAnnotationKey: string(recordBytes),
			},
		},
	}, metav1.CreateOptions{})
	return err
}

// Update will update an existing annotation on a given resource.
func (cml *ConfigMapLock) Update(ctx context.Context, ler resourcelock.LeaderElectionRecord) error {
	if cml.cm == nil {
		return errors.New("configmap not initialized, call get or create first")
	}
	recordBytes, err := json.Marshal(ler)
	if err != nil {
		return err
	}
	if cml.cm.Annotations == nil {
		cml.cm.Annotations = make(map[string]string)
	}
	cml.cm.Annotations[resourcelock.LeaderElectionRecordAnnotationKey] = string(recordBytes)
	cm, err := cml.Client.ConfigMaps(cml.ConfigMapMeta.Namespace).Update(ctx, cml.cm, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	cml.cm = cm
	return nil
}

// RecordEvent in leader election while adding meta-data
func (cml *ConfigMapLock) RecordEvent(s string) {
	if cml.LockConfig.EventRecorder == nil {
		return
	}
	events := fmt.Sprintf("%v %v", cml.LockConfig.Identity, s)
	subject := &v1.ConfigMap{ObjectMeta: cml.cm.ObjectMeta}
	// Populate the type meta, so we don't have to get it from the schema
	subject.Kind = "ConfigMap"
	subject.APIVersion = v1.SchemeGroupVersion.String()
	cml.LockConfig.EventRecorder.Eventf(subject, v1.EventTypeNormal, "LeaderElection", events)
}

// Describe is used to convert details on current resource lock
// into a string
func (cml *ConfigMapLock) Describe() string {
	return fmt.Sprintf("%v/%v", cml.ConfigMapMeta.Namespace, cml.ConfigMapMeta.Name)
}

// Identity returns the Identity of the lock
func (cml *ConfigMapLock) Identity() string {
	return cml.LockConfig.Identity
}
