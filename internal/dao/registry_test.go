// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of K9s

package dao

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestExtractMeta(t *testing.T) {
	c := load(t, "dr")
	var crd apiextensionsv1.CustomResourceDefinition
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(c.Object, &crd)
	require.NoError(t, err)

	m, err := extractMeta(&crd)
	assert.NoError(t, err)
	assert.Equal(t, "destinationrules", m.Name)
	assert.Equal(t, "destinationrule", m.SingularName)
	assert.Equal(t, "DestinationRule", m.Kind)
	assert.Equal(t, "networking.istio.io", m.Group)
	assert.Equal(t, "v1", m.Version)
	assert.Equal(t, true, m.Namespaced)
	assert.Equal(t, []string{"dr"}, m.ShortNames)
	var vv metav1.Verbs
	assert.Equal(t, vv, m.Verbs)
}

// Helpers...

func load(t *testing.T, n string) *unstructured.Unstructured {
	raw, err := os.ReadFile(fmt.Sprintf("testdata/%s.json", n))
	assert.Nil(t, err)

	var o unstructured.Unstructured
	err = json.Unmarshal(raw, &o)
	assert.Nil(t, err)

	return &o
}
