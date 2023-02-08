/*
 *     Copyright 2020 The Dragonfly Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package evaluator

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	commonv1 "d7y.io/api/pkg/apis/common/v1"
	commonv2 "d7y.io/api/pkg/apis/common/v2"

	"d7y.io/dragonfly/v2/pkg/idgen"
	"d7y.io/dragonfly/v2/pkg/types"
	"d7y.io/dragonfly/v2/scheduler/resource"
)

var (
	mockRawHost = resource.Host{
		ID:              idgen.HostIDV1("hostname", 8003),
		Type:            types.HostTypeNormal,
		Hostname:        "hostname",
		IP:              "127.0.0.1",
		Port:            8003,
		DownloadPort:    8001,
		OS:              "darwin",
		Platform:        "darwin",
		PlatformFamily:  "Standalone Workstation",
		PlatformVersion: "11.1",
		KernelVersion:   "20.2.0",
		CPU:             mockCPU,
		Memory:          mockMemory,
		Network:         mockNetwork,
		Disk:            mockDisk,
		Build:           mockBuild,
	}

	mockRawSeedHost = resource.Host{
		ID:              idgen.HostIDV1("hostname_seed", 8003),
		Type:            types.HostTypeSuperSeed,
		Hostname:        "hostname_seed",
		IP:              "127.0.0.1",
		Port:            8003,
		DownloadPort:    8001,
		OS:              "darwin",
		Platform:        "darwin",
		PlatformFamily:  "Standalone Workstation",
		PlatformVersion: "11.1",
		KernelVersion:   "20.2.0",
		CPU:             mockCPU,
		Memory:          mockMemory,
		Network:         mockNetwork,
		Disk:            mockDisk,
		Build:           mockBuild,
	}

	mockCPU = resource.CPU{
		LogicalCount:   4,
		PhysicalCount:  2,
		Percent:        1,
		ProcessPercent: 0.5,
		Times: resource.CPUTimes{
			User:      240662.2,
			System:    317950.1,
			Idle:      3393691.3,
			Nice:      0,
			Iowait:    0,
			Irq:       0,
			Softirq:   0,
			Steal:     0,
			Guest:     0,
			GuestNice: 0,
		},
	}

	mockMemory = resource.Memory{
		Total:              17179869184,
		Available:          5962813440,
		Used:               11217055744,
		UsedPercent:        65.291858,
		ProcessUsedPercent: 41.525125,
		Free:               2749598908,
	}

	mockNetwork = resource.Network{
		TCPConnectionCount:       10,
		UploadTCPConnectionCount: 1,
		SecurityDomain:           "security_domain",
		Location:                 "location",
		IDC:                      "idc",
	}

	mockDisk = resource.Disk{
		Total:             499963174912,
		Free:              37226479616,
		Used:              423809622016,
		UsedPercent:       91.92547406065952,
		InodesTotal:       4882452880,
		InodesUsed:        7835772,
		InodesFree:        4874617108,
		InodesUsedPercent: 0.1604884305611568,
	}

	mockBuild = resource.Build{
		GitVersion: "v1.0.0",
		GitCommit:  "221176b117c6d59366d68f2b34d38be50c935883",
		GoVersion:  "1.18",
		Platform:   "darwin",
	}

	mockTaskURLMeta = &commonv1.UrlMeta{
		Digest: "digest",
		Tag:    "tag",
		Range:  "range",
		Filter: "filter",
		Header: map[string]string{
			"content-length": "100",
		},
	}

	mockTaskBackToSourceLimit int32 = 200
	mockTaskURL                     = "http://example.com/foo"
	mockTaskID                      = idgen.TaskIDV1(mockTaskURL, mockTaskURLMeta)
	mockPeerID                      = idgen.PeerIDV1("127.0.0.1")
)

func TestEvaluatorBase_NewEvaluatorBase(t *testing.T) {
	tests := []struct {
		name   string
		expect func(t *testing.T, e any)
	}{
		{
			name: "new evaluator commonv1",
			expect: func(t *testing.T, e any) {
				assert := assert.New(t)
				assert.Equal(reflect.TypeOf(e).Elem().Name(), "evaluatorBase")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, NewEvaluatorBase())
		})
	}
}

func TestEvaluatorBase_Evaluate(t *testing.T) {
	tests := []struct {
		name            string
		parent          *resource.Peer
		child           *resource.Peer
		totalPieceCount int32
		mock            func(parent *resource.Peer, child *resource.Peer)
		expect          func(t *testing.T, score float64)
	}{
		{
			name: "security domain is not the same",
			parent: resource.NewPeer(idgen.PeerIDV1("127.0.0.1"),
				resource.NewTask(mockTaskID, mockTaskURL, commonv2.TaskType_DFDAEMON, mockTaskURLMeta, resource.WithBackToSourceLimit(mockTaskBackToSourceLimit)),
				resource.NewHost(
					mockRawSeedHost.ID, mockRawSeedHost.IP, mockRawSeedHost.Hostname,
					mockRawSeedHost.Port, mockRawSeedHost.DownloadPort, mockRawSeedHost.Type)),
			child: resource.NewPeer(idgen.PeerIDV1("127.0.0.1"),
				resource.NewTask(mockTaskID, mockTaskURL, commonv2.TaskType_DFDAEMON, mockTaskURLMeta, resource.WithBackToSourceLimit(mockTaskBackToSourceLimit)),
				resource.NewHost(
					mockRawHost.ID, mockRawHost.IP, mockRawHost.Hostname,
					mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.Type)),
			totalPieceCount: 1,
			mock: func(parent *resource.Peer, child *resource.Peer) {
				parent.Host.Network.SecurityDomain = "foo"
				child.Host.Network.SecurityDomain = "bar"
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0))
			},
		},
		{
			name: "security domain is same",
			parent: resource.NewPeer(idgen.PeerIDV1("127.0.0.1"),
				resource.NewTask(mockTaskID, mockTaskURL, commonv2.TaskType_DFDAEMON, mockTaskURLMeta, resource.WithBackToSourceLimit(mockTaskBackToSourceLimit)),
				resource.NewHost(
					mockRawSeedHost.ID, mockRawSeedHost.IP, mockRawSeedHost.Hostname,
					mockRawSeedHost.Port, mockRawSeedHost.DownloadPort, mockRawSeedHost.Type)),
			child: resource.NewPeer(idgen.PeerIDV1("127.0.0.1"),
				resource.NewTask(mockTaskID, mockTaskURL, commonv2.TaskType_DFDAEMON, mockTaskURLMeta, resource.WithBackToSourceLimit(mockTaskBackToSourceLimit)),
				resource.NewHost(
					mockRawHost.ID, mockRawHost.IP, mockRawHost.Hostname,
					mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.Type)),
			totalPieceCount: 1,
			mock: func(parent *resource.Peer, child *resource.Peer) {
				parent.Host.Network.SecurityDomain = "bac"
				child.Host.Network.SecurityDomain = "bac"
				parent.FinishedPieces.Set(0)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0.55))
			},
		},
		{
			name: "parent security domain is empty",
			parent: resource.NewPeer(idgen.PeerIDV1("127.0.0.1"),
				resource.NewTask(mockTaskID, mockTaskURL, commonv2.TaskType_DFDAEMON, mockTaskURLMeta, resource.WithBackToSourceLimit(mockTaskBackToSourceLimit)),
				resource.NewHost(
					mockRawSeedHost.ID, mockRawSeedHost.IP, mockRawSeedHost.Hostname,
					mockRawSeedHost.Port, mockRawSeedHost.DownloadPort, mockRawSeedHost.Type)),
			child: resource.NewPeer(idgen.PeerIDV1("127.0.0.1"),
				resource.NewTask(mockTaskID, mockTaskURL, commonv2.TaskType_DFDAEMON, mockTaskURLMeta, resource.WithBackToSourceLimit(mockTaskBackToSourceLimit)),
				resource.NewHost(
					mockRawHost.ID, mockRawHost.IP, mockRawHost.Hostname,
					mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.Type)),
			totalPieceCount: 1,
			mock: func(parent *resource.Peer, child *resource.Peer) {
				parent.Host.Network.SecurityDomain = ""
				child.Host.Network.SecurityDomain = "baz"
				parent.FinishedPieces.Set(0)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0.55))
			},
		},
		{
			name: "child security domain is empty",
			parent: resource.NewPeer(idgen.PeerIDV1("127.0.0.1"),
				resource.NewTask(mockTaskID, mockTaskURL, commonv2.TaskType_DFDAEMON, mockTaskURLMeta, resource.WithBackToSourceLimit(mockTaskBackToSourceLimit)),
				resource.NewHost(
					mockRawSeedHost.ID, mockRawSeedHost.IP, mockRawSeedHost.Hostname,
					mockRawSeedHost.Port, mockRawSeedHost.DownloadPort, mockRawSeedHost.Type)),
			child: resource.NewPeer(idgen.PeerIDV1("127.0.0.1"),
				resource.NewTask(mockTaskID, mockTaskURL, commonv2.TaskType_DFDAEMON, mockTaskURLMeta, resource.WithBackToSourceLimit(mockTaskBackToSourceLimit)),
				resource.NewHost(
					mockRawHost.ID, mockRawHost.IP, mockRawHost.Hostname,
					mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.Type)),
			totalPieceCount: 1,
			mock: func(parent *resource.Peer, child *resource.Peer) {
				parent.Host.Network.SecurityDomain = "baz"
				child.Host.Network.SecurityDomain = ""
				parent.FinishedPieces.Set(0)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0.55))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			eb := NewEvaluatorBase()
			tc.mock(tc.parent, tc.child)
			tc.expect(t, eb.Evaluate(tc.parent, tc.child, tc.totalPieceCount))
		})
	}
}

func TestEvaluatorBase_calculatePieceScore(t *testing.T) {
	mockHost := resource.NewHost(
		mockRawHost.ID, mockRawHost.IP, mockRawHost.Hostname,
		mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.Type)
	mockTask := resource.NewTask(mockTaskID, mockTaskURL, commonv2.TaskType_DFDAEMON, mockTaskURLMeta, resource.WithBackToSourceLimit(mockTaskBackToSourceLimit))

	tests := []struct {
		name            string
		parent          *resource.Peer
		child           *resource.Peer
		totalPieceCount int32
		mock            func(parent *resource.Peer, child *resource.Peer)
		expect          func(t *testing.T, score float64)
	}{
		{
			name:            "total piece count is zero and child pieces are empty",
			parent:          resource.NewPeer(idgen.PeerIDV1("127.0.0.1"), mockTask, mockHost),
			child:           resource.NewPeer(idgen.PeerIDV1("127.0.0.1"), mockTask, mockHost),
			totalPieceCount: 0,
			mock: func(parent *resource.Peer, child *resource.Peer) {
				parent.FinishedPieces.Set(0)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(1))
			},
		},
		{
			name:            "total piece count is zero and parent pieces are empty",
			parent:          resource.NewPeer(idgen.PeerIDV1("127.0.0.1"), mockTask, mockHost),
			child:           resource.NewPeer(idgen.PeerIDV1("127.0.0.1"), mockTask, mockHost),
			totalPieceCount: 0,
			mock: func(parent *resource.Peer, child *resource.Peer) {
				child.FinishedPieces.Set(0)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(-1))
			},
		},
		{
			name:            "total piece count is zero and child pieces of length greater than parent pieces",
			parent:          resource.NewPeer(idgen.PeerIDV1("127.0.0.1"), mockTask, mockHost),
			child:           resource.NewPeer(idgen.PeerIDV1("127.0.0.1"), mockTask, mockHost),
			totalPieceCount: 0,
			mock: func(parent *resource.Peer, child *resource.Peer) {
				parent.FinishedPieces.Set(0)
				child.FinishedPieces.Set(0)
				child.FinishedPieces.Set(1)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(-1))
			},
		},
		{
			name:            "total piece count is zero and child pieces of length equal than parent pieces",
			parent:          resource.NewPeer(idgen.PeerIDV1("127.0.0.1"), mockTask, mockHost),
			child:           resource.NewPeer(idgen.PeerIDV1("127.0.0.1"), mockTask, mockHost),
			totalPieceCount: 0,
			mock: func(parent *resource.Peer, child *resource.Peer) {
				parent.FinishedPieces.Set(0)
				child.FinishedPieces.Set(0)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0))
			},
		},
		{
			name:            "total piece count is zero and parent pieces of length greater than child pieces",
			parent:          resource.NewPeer(idgen.PeerIDV1("127.0.0.1"), mockTask, mockHost),
			child:           resource.NewPeer(idgen.PeerIDV1("127.0.0.1"), mockTask, mockHost),
			totalPieceCount: 0,
			mock: func(parent *resource.Peer, child *resource.Peer) {
				parent.FinishedPieces.Set(0)
				parent.FinishedPieces.Set(1)
				child.FinishedPieces.Set(0)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(1))
			},
		},
		{
			name:            "parent pieces are empty",
			parent:          resource.NewPeer(idgen.PeerIDV1("127.0.0.1"), mockTask, mockHost),
			child:           resource.NewPeer(idgen.PeerIDV1("127.0.0.1"), mockTask, mockHost),
			totalPieceCount: 10,
			mock:            func(parent *resource.Peer, child *resource.Peer) {},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0))
			},
		},
		{
			name:            "parent pieces of length greater than zero",
			parent:          resource.NewPeer(idgen.PeerIDV1("127.0.0.1"), mockTask, mockHost),
			child:           resource.NewPeer(idgen.PeerIDV1("127.0.0.1"), mockTask, mockHost),
			totalPieceCount: 10,
			mock: func(parent *resource.Peer, child *resource.Peer) {
				parent.FinishedPieces.Set(0)
				parent.FinishedPieces.Set(1)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0.2))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock(tc.parent, tc.child)
			tc.expect(t, calculatePieceScore(tc.parent, tc.child, tc.totalPieceCount))
		})
	}
}

func TestEvaluatorBase_calculatehostUploadSuccessScore(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(host *resource.Host)
		expect func(t *testing.T, score float64)
	}{
		{
			name: "UploadFailedCount is larger than UploadCount",
			mock: func(host *resource.Host) {
				host.UploadCount.Add(1)
				host.UploadFailedCount.Add(2)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0))
			},
		},
		{
			name: "UploadFailedCount and UploadCount is zero",
			mock: func(host *resource.Host) {
				host.UploadCount.Add(0)
				host.UploadFailedCount.Add(0)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(1))
			},
		},
		{
			name: "UploadCount is larger than UploadFailedCount",
			mock: func(host *resource.Host) {
				host.UploadCount.Add(2)
				host.UploadFailedCount.Add(1)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0.5))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			host := resource.NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.Type)
			mockTask := resource.NewTask(mockTaskID, mockTaskURL, commonv2.TaskType_DFDAEMON, mockTaskURLMeta, resource.WithBackToSourceLimit(mockTaskBackToSourceLimit))
			mockPeer := resource.NewPeer(mockPeerID, mockTask, host)
			tc.mock(host)
			tc.expect(t, calculateParentHostUploadSuccessScore(mockPeer))
		})
	}
}

func TestEvaluatorBase_calculateFreeUploadScore(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(host *resource.Host, mockPeer *resource.Peer)
		expect func(t *testing.T, score float64)
	}{
		{
			name: "host peers is not empty",
			mock: func(host *resource.Host, mockPeer *resource.Peer) {
				mockPeer.Host.ConcurrentUploadCount.Add(1)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0.98))
			},
		},
		{
			name: "host peers is empty",
			mock: func(host *resource.Host, mockPeer *resource.Peer) {},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(1))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			host := resource.NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.Type)
			mockTask := resource.NewTask(mockTaskID, mockTaskURL, commonv2.TaskType_DFDAEMON, mockTaskURLMeta, resource.WithBackToSourceLimit(mockTaskBackToSourceLimit))
			mockPeer := resource.NewPeer(mockPeerID, mockTask, host)
			tc.mock(host, mockPeer)
			tc.expect(t, calculateFreeUploadScore(host))
		})
	}
}

func TestEvaluatorBase_calculateHostTypeScore(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(peer *resource.Peer)
		expect func(t *testing.T, score float64)
	}{
		{
			name: "host is normal peer",
			mock: func(peer *resource.Peer) {},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0.5))
			},
		},
		{
			name: "host is seed peer but peer state is PeerStateSucceeded",
			mock: func(peer *resource.Peer) {
				peer.FSM.SetState(resource.PeerStateSucceeded)
				peer.Host.Type = types.HostTypeSuperSeed
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0))
			},
		},
		{
			name: "host is seed peer but peer state is PeerStateRunning",
			mock: func(peer *resource.Peer) {
				peer.FSM.SetState(resource.PeerStateRunning)
				peer.Host.Type = types.HostTypeSuperSeed
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(1))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockHost := resource.NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.Type)
			mockTask := resource.NewTask(mockTaskID, mockTaskURL, commonv2.TaskType_DFDAEMON, mockTaskURLMeta, resource.WithBackToSourceLimit(mockTaskBackToSourceLimit))
			peer := resource.NewPeer(mockPeerID, mockTask, mockHost)
			tc.mock(peer)
			tc.expect(t, calculateHostTypeScore(peer))
		})
	}
}

func TestEvaluatorBase_calculateIDCAffinityScore(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(dstHost *resource.Host, srcHost *resource.Host)
		expect func(t *testing.T, score float64)
	}{
		{
			name: "idc is empty",
			mock: func(dstHost *resource.Host, srcHost *resource.Host) {
				dstHost.Network.IDC = ""
				srcHost.Network.IDC = ""
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0))
			},
		},
		{
			name: "dst host idc is empty",
			mock: func(dstHost *resource.Host, srcHost *resource.Host) {
				dstHost.Network.IDC = ""
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0))
			},
		},
		{
			name: "src host idc is empty",
			mock: func(dstHost *resource.Host, srcHost *resource.Host) {
				srcHost.Network.IDC = ""
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0))
			},
		},
		{
			name: "idc is not the same",
			mock: func(dstHost *resource.Host, srcHost *resource.Host) {
				dstHost.Network.IDC = "foo"
				srcHost.Network.IDC = "bar"
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0))
			},
		},
		{
			name: "idc is the same",
			mock: func(dstHost *resource.Host, srcHost *resource.Host) {
				dstHost.Network.IDC = "example"
				srcHost.Network.IDC = "example"
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(1))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dstHost := resource.NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.Type)
			srcHost := resource.NewHost(
				mockRawSeedHost.ID, mockRawSeedHost.IP, mockRawSeedHost.Hostname,
				mockRawSeedHost.Port, mockRawSeedHost.DownloadPort, mockRawSeedHost.Type)
			tc.mock(dstHost, srcHost)
			tc.expect(t, calculateIDCAffinityScore(dstHost.Network.IDC, srcHost.Network.IDC))
		})
	}
}

func TestEvaluatorBase_calculateMultiElementAffinityScore(t *testing.T) {
	tests := []struct {
		name   string
		dst    string
		src    string
		expect func(t *testing.T, score float64)
	}{
		{
			name: "dst is empty and src is empty",
			dst:  "",
			src:  "",
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0))
			},
		},
		{
			name: "dst is empty",
			dst:  "",
			src:  "baz",
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0))
			},
		},
		{
			name: "src is empty",
			dst:  "bar",
			src:  "",
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0))
			},
		},
		{
			name: "has only one element and matches",
			dst:  "foo",
			src:  "foo",
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(1))
			},
		},
		{
			name: "has only one element and does not match",
			dst:  "foo",
			src:  "bar",
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0))
			},
		},
		{
			name: "has multi element and match",
			dst:  "foo|bar",
			src:  "foo|bar",
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(1))
			},
		},
		{
			name: "has multi element and does not match",
			dst:  "foo|bar",
			src:  "bar|foo",
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0))
			},
		},
		{
			name: "dst length is greater than src",
			dst:  "foo|bar|baz",
			src:  "foo|bar",
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0.4))
			},
		},
		{
			name: "src length is greater than dst",
			dst:  "foo|bar",
			src:  "foo|bar|baz",
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0.4))
			},
		},
		{
			name: "dst exceeds maximum length",
			dst:  "foo|bar|baz|bac|bae|baf",
			src:  "foo|bar|baz",
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0.6))
			},
		},
		{
			name: "src exceeds maximum length",
			dst:  "foo|bar|baz",
			src:  "foo|bar|baz|bac|bae|baf",
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0.6))
			},
		},
		{
			name: "dst and src both exceeds maximum length",
			dst:  "foo|bar|baz|bac|bae|baf",
			src:  "foo|bar|baz|bac|bae|baf",
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(1))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, calculateMultiElementAffinityScore(tc.dst, tc.src))
		})
	}
}

func TestEvaluatorBase_IsBadNode(t *testing.T) {
	mockHost := resource.NewHost(
		mockRawHost.ID, mockRawHost.IP, mockRawHost.Hostname,
		mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.Type)
	mockTask := resource.NewTask(mockTaskID, mockTaskURL, commonv2.TaskType_DFDAEMON, mockTaskURLMeta, resource.WithBackToSourceLimit(mockTaskBackToSourceLimit))

	tests := []struct {
		name            string
		peer            *resource.Peer
		totalPieceCount int32
		mock            func(peer *resource.Peer)
		expect          func(t *testing.T, isBadNode bool)
	}{
		{
			name:            "peer state is PeerStateFailed",
			peer:            resource.NewPeer(mockPeerID, mockTask, mockHost),
			totalPieceCount: 1,
			mock: func(peer *resource.Peer) {
				peer.FSM.SetState(resource.PeerStateFailed)
			},
			expect: func(t *testing.T, isBadNode bool) {
				assert := assert.New(t)
				assert.True(isBadNode)
			},
		},
		{
			name:            "peer state is PeerStateLeave",
			peer:            resource.NewPeer(mockPeerID, mockTask, mockHost),
			totalPieceCount: 1,
			mock: func(peer *resource.Peer) {
				peer.FSM.SetState(resource.PeerStateLeave)
			},
			expect: func(t *testing.T, isBadNode bool) {
				assert := assert.New(t)
				assert.True(isBadNode)
			},
		},
		{
			name:            "peer state is PeerStatePending",
			peer:            resource.NewPeer(mockPeerID, mockTask, mockHost),
			totalPieceCount: 1,
			mock: func(peer *resource.Peer) {
				peer.FSM.SetState(resource.PeerStatePending)
			},
			expect: func(t *testing.T, isBadNode bool) {
				assert := assert.New(t)
				assert.True(isBadNode)
			},
		},
		{
			name:            "peer state is PeerStateReceivedTiny",
			peer:            resource.NewPeer(mockPeerID, mockTask, mockHost),
			totalPieceCount: 1,
			mock: func(peer *resource.Peer) {
				peer.FSM.SetState(resource.PeerStateReceivedTiny)
			},
			expect: func(t *testing.T, isBadNode bool) {
				assert := assert.New(t)
				assert.True(isBadNode)
			},
		},
		{
			name:            "peer state is PeerStateReceivedSmall",
			peer:            resource.NewPeer(mockPeerID, mockTask, mockHost),
			totalPieceCount: 1,
			mock: func(peer *resource.Peer) {
				peer.FSM.SetState(resource.PeerStateReceivedSmall)
			},
			expect: func(t *testing.T, isBadNode bool) {
				assert := assert.New(t)
				assert.True(isBadNode)
			},
		},
		{
			name:            "peer state is PeerStateReceivedNormal",
			peer:            resource.NewPeer(mockPeerID, mockTask, mockHost),
			totalPieceCount: 1,
			mock: func(peer *resource.Peer) {
				peer.FSM.SetState(resource.PeerStateReceivedNormal)
			},
			expect: func(t *testing.T, isBadNode bool) {
				assert := assert.New(t)
				assert.True(isBadNode)
			},
		},
		{
			name:            "download costs does not meet the normal distribution and last cost is twenty times more than mean",
			peer:            resource.NewPeer(mockPeerID, mockTask, mockHost),
			totalPieceCount: 1,
			mock: func(peer *resource.Peer) {
				peer.FSM.SetState(resource.PeerStateRunning)
				peer.AppendPieceCost(10)
				peer.AppendPieceCost(201)
			},
			expect: func(t *testing.T, isBadNode bool) {
				assert := assert.New(t)
				assert.True(isBadNode)
			},
		},
		{
			name:            "download costs does not meet the normal distribution and last cost is twenty times lower than mean",
			peer:            resource.NewPeer(mockPeerID, mockTask, mockHost),
			totalPieceCount: 1,
			mock: func(peer *resource.Peer) {
				peer.FSM.SetState(resource.PeerStateRunning)
				peer.AppendPieceCost(10)
				peer.AppendPieceCost(200)
			},
			expect: func(t *testing.T, isBadNode bool) {
				assert := assert.New(t)
				assert.False(isBadNode)
			},
		},
		{
			name:            "download costs meet the normal distribution and last cost is too long",
			peer:            resource.NewPeer(mockPeerID, mockTask, mockHost),
			totalPieceCount: 1,
			mock: func(peer *resource.Peer) {
				peer.FSM.SetState(resource.PeerStateRunning)
				for i := 0; i < 30; i++ {
					peer.AppendPieceCost(int64(i))
				}
				peer.AppendPieceCost(50)
			},
			expect: func(t *testing.T, isBadNode bool) {
				assert := assert.New(t)
				assert.True(isBadNode)
			},
		},
		{
			name:            "download costs meet the normal distribution and last cost is normal",
			peer:            resource.NewPeer(mockPeerID, mockTask, mockHost),
			totalPieceCount: 1,
			mock: func(peer *resource.Peer) {
				peer.FSM.SetState(resource.PeerStateRunning)
				for i := 0; i < 30; i++ {
					peer.AppendPieceCost(int64(i))
				}
				peer.AppendPieceCost(18)
			},
			expect: func(t *testing.T, isBadNode bool) {
				assert := assert.New(t)
				assert.False(isBadNode)
			},
		},
		{
			name:            "download costs meet the normal distribution and last cost is too short",
			peer:            resource.NewPeer(mockPeerID, mockTask, mockHost),
			totalPieceCount: 1,
			mock: func(peer *resource.Peer) {
				peer.FSM.SetState(resource.PeerStateRunning)
				for i := 20; i < 50; i++ {
					peer.AppendPieceCost(int64(i))
				}
				peer.AppendPieceCost(0)
			},
			expect: func(t *testing.T, isBadNode bool) {
				assert := assert.New(t)
				assert.False(isBadNode)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			eb := NewEvaluatorBase()
			tc.mock(tc.peer)
			tc.expect(t, eb.IsBadNode(tc.peer))
		})
	}
}
