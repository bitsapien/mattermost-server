// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package model

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClusterInfoJSON(t *testing.T) {
	cluster := ClusterInfo{IPAddress: NewID(), Hostname: NewID()}
	json := cluster.ToJSON()
	result := ClusterInfoFromJSON(strings.NewReader(json))

	assert.Equal(t, cluster.IPAddress, result.IPAddress, "Ids do not match")
}

func TestClusterInfosJSON(t *testing.T) {
	cluster := ClusterInfo{IPAddress: NewID(), Hostname: NewID()}
	clusterInfos := make([]*ClusterInfo, 1)
	clusterInfos[0] = &cluster
	json := ClusterInfosToJSON(clusterInfos)
	result := ClusterInfosFromJSON(strings.NewReader(json))

	assert.Equal(t, clusterInfos[0].IPAddress, result[0].IPAddress, "Ids do not match")
}
