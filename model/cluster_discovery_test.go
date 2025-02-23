// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package model

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClusterDiscovery(t *testing.T) {
	o := ClusterDiscovery{
		Type:        "test_type",
		ClusterName: "cluster_name",
		Hostname:    "test_hostname",
	}

	buf, err := json.Marshal(o)
	require.NoError(t, err)
	json := string(buf)
	result1 := ClusterDiscoveryFromJson(strings.NewReader(json))

	assert.NotNil(t, result1)
	assert.Equal(t, "cluster_name", result1.ClusterName)

	result2 := ClusterDiscoveryFromJson(strings.NewReader(json))
	result3 := ClusterDiscoveryFromJson(strings.NewReader(json))

	o.Id = "0"
	result1.Id = "1"
	result2.Id = "2"
	result3.Id = "3"
	result3.Hostname = "something_diff"

	assert.True(t, o.IsEqual(result1))

	list := make([]*ClusterDiscovery, 0)
	list = append(list, &o)
	list = append(list, result1)
	list = append(list, result2)
	list = append(list, result3)

	rlist := FilterClusterDiscovery(list, func(in *ClusterDiscovery) bool {
		return !o.IsEqual(in)
	})

	assert.Len(t, rlist, 1)

	o.AutoFillHostname()
	o.Hostname = ""
	o.AutoFillHostname()

	o.AutoFillIpAddress("", "")
	o.Hostname = ""
	o.AutoFillIpAddress("", "")
}
