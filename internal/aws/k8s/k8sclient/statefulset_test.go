// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package k8sclient

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/fake"
)

var statefulSetObjects = []runtime.Object{
	&appsv1.StatefulSet{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-statefulset-1",
			Namespace: "test-namespace",
			UID:       types.UID("test-statefulset-1-uid"),
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &desired,
		},
		Status: appsv1.StatefulSetStatus{
			Replicas:          5,
			AvailableReplicas: 5,
			ReadyReplicas:     5,
		},
	},
	&appsv1.StatefulSet{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-statefulset-2",
			Namespace: "test-namespace",
			UID:       types.UID("test-statefulset-2-uid"),
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &desired,
		},
		Status: appsv1.StatefulSetStatus{
			Replicas:          10,
			AvailableReplicas: 10,
			ReadyReplicas:     10,
		},
	},
}

func TestStatefulSetClient(t *testing.T) {
	setOption := statefulSetSyncCheckerOption(&mockReflectorSyncChecker{})

	fakeClientSet := fake.NewSimpleClientset(statefulSetObjects...)
	client, _ := newStatefulSetClient(fakeClientSet, zap.NewNop(), setOption)

	statefulSets := make([]interface{}, len(statefulSetObjects))
	for i := range statefulSetObjects {
		statefulSets[i] = statefulSetObjects[i]
	}
	assert.NoError(t, client.store.Replace(statefulSets, ""))

	expected := []*StatefulSetInfo{
		{
			Name:      "test-statefulset-1",
			Namespace: "test-namespace",
			Spec: &StatefulSetSpec{
				Replicas: 20,
			},
			Status: &StatefulSetStatus{
				Replicas:          5,
				AvailableReplicas: 5,
				ReadyReplicas:     5,
			},
		},
		{
			Name:      "test-statefulset-2",
			Namespace: "test-namespace",
			Spec: &StatefulSetSpec{
				Replicas: 20,
			},
			Status: &StatefulSetStatus{
				Replicas:          10,
				AvailableReplicas: 10,
				ReadyReplicas:     10,
			},
		},
	}
	actual := client.StatefulSetInfos()
	sort.Slice(actual, func(i, j int) bool {
		return actual[i].Name < actual[j].Name
	})
	assert.Equal(t, expected, actual)
	client.shutdown()
	assert.True(t, client.stopped)
}

func TestTransformFuncStatefulSet(t *testing.T) {
	info, err := transformFuncStatefulSet(nil)
	assert.Nil(t, info)
	assert.NotNil(t, err)
}
