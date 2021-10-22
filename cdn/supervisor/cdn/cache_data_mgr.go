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

package cdn

import (
	"fmt"
	"io"
	"sort"

	"d7y.io/dragonfly/v2/cdn/storedriver"
	"d7y.io/dragonfly/v2/cdn/supervisor/cdn/storage"
	"d7y.io/dragonfly/v2/cdn/types"
	logger "d7y.io/dragonfly/v2/internal/dflog"
	"d7y.io/dragonfly/v2/pkg/synclock"
	"d7y.io/dragonfly/v2/pkg/util/digestutils"
	"d7y.io/dragonfly/v2/pkg/util/stringutils"
	"d7y.io/dragonfly/v2/pkg/util/timeutils"
	"github.com/pkg/errors"
)

// metaDataManager manages the meta file and piece meta file of each TaskID.
type metaDataManager struct {
	storage     storage.Manager
	cacheLocker *synclock.LockerPool
}

func newCacheDataManager(storeMgr storage.Manager) *metaDataManager {
	return &metaDataManager{
		storeMgr,
		synclock.NewLockerPool(),
	}
}

// writeFileMetaDataByTask stores the metadata of task by task to storage.
func (mm *metaDataManager) writeFileMetaDataByTask(task *types.SeedTask) (*storage.FileMetaData, error) {
	mm.cacheLocker.Lock(task.ID, false)
	defer mm.cacheLocker.UnLock(task.ID, false)
	metaData := &storage.FileMetaData{
		TaskID:          task.ID,
		TaskURL:         task.TaskURL,
		PieceSize:       task.PieceSize,
		SourceFileLen:   task.SourceFileLength,
		AccessTime:      getCurrentTimeMillisFunc(),
		CdnFileLength:   task.CdnFileLength,
		Finish:          false,
		Success:         false,
		TotalPieceCount: task.PieceTotal,
	}

	if err := mm.storage.WriteFileMetaData(task.ID, metaData); err != nil {
		return nil, errors.Wrapf(err, "write task metadata file")
	}

	return metaData, nil
}

// updateAccessTime update access and interval
func (mm *metaDataManager) updateAccessTime(taskID string, accessTime int64) error {
	mm.cacheLocker.Lock(taskID, false)
	defer mm.cacheLocker.UnLock(taskID, false)

	originMetaData, err := mm.readFileMetaData(taskID)
	if err != nil {
		return err
	}
	// access interval
	interval := accessTime - originMetaData.AccessTime
	originMetaData.Interval = interval
	if interval <= 0 {
		logger.WithTaskID(taskID).Warnf("file hit interval: %d, accessTime: %s", interval, timeutils.MillisUnixTime(accessTime))
		originMetaData.Interval = 0
	}

	originMetaData.AccessTime = accessTime

	return mm.storage.WriteFileMetaData(taskID, originMetaData)
}

func (mm *metaDataManager) updateExpireInfo(taskID string, expireInfo map[string]string) error {
	mm.cacheLocker.Lock(taskID, false)
	defer mm.cacheLocker.UnLock(taskID, false)

	originMetaData, err := mm.readFileMetaData(taskID)
	if err != nil {
		return err
	}

	originMetaData.ExpireInfo = expireInfo

	return mm.storage.WriteFileMetaData(taskID, originMetaData)
}

func (mm *metaDataManager) updateStatusAndResult(taskID string, metaData *storage.FileMetaData) error {
	mm.cacheLocker.Lock(taskID, false)
	defer mm.cacheLocker.UnLock(taskID, false)

	originMetaData, err := mm.readFileMetaData(taskID)
	if err != nil {
		return err
	}

	originMetaData.Finish = metaData.Finish
	originMetaData.Success = metaData.Success
	if originMetaData.Success {
		originMetaData.CdnFileLength = metaData.CdnFileLength
		originMetaData.SourceFileLen = metaData.SourceFileLen
		if metaData.TotalPieceCount > 0 {
			originMetaData.TotalPieceCount = metaData.TotalPieceCount
		}
		if !stringutils.IsBlank(metaData.SourceRealDigest) {
			originMetaData.SourceRealDigest = metaData.SourceRealDigest
		}
		if !stringutils.IsBlank(metaData.PieceMd5Sign) {
			originMetaData.PieceMd5Sign = metaData.PieceMd5Sign
		}
	}
	return mm.storage.WriteFileMetaData(taskID, originMetaData)
}

// appendPieceMetaData append piece meta info to storage
func (mm *metaDataManager) appendPieceMetaData(taskID string, record *storage.PieceMetaRecord) error {
	mm.cacheLocker.Lock(taskID, false)
	defer mm.cacheLocker.UnLock(taskID, false)
	// write to the storage
	return mm.storage.AppendPieceMetaData(taskID, record)
}

// appendPieceMetaData append piece meta info to storage
func (mm *metaDataManager) writePieceMetaRecords(taskID string, records []*storage.PieceMetaRecord) error {
	mm.cacheLocker.Lock(taskID, false)
	defer mm.cacheLocker.UnLock(taskID, false)
	// write to the storage
	return mm.storage.WritePieceMetaRecords(taskID, records)
}

// readAndCheckPieceMetaRecords reads pieceMetaRecords from storage and check data integrity by the md5 file of the TaskId
func (mm *metaDataManager) readAndCheckPieceMetaRecords(taskID, pieceMd5Sign string) ([]*storage.PieceMetaRecord, error) {
	mm.cacheLocker.Lock(taskID, true)
	defer mm.cacheLocker.UnLock(taskID, true)
	md5Sign, pieceMetaRecords, err := mm.getPieceMd5Sign(taskID)
	if err != nil {
		return nil, err
	}
	if md5Sign != pieceMd5Sign {
		return nil, fmt.Errorf("check piece meta data integrity fail, expectMd5Sign: %s, actualMd5Sign: %s",
			pieceMd5Sign, md5Sign)
	}
	return pieceMetaRecords, nil
}

// readPieceMetaRecords reads pieceMetaRecords from storage and without check data integrity
func (mm *metaDataManager) readPieceMetaRecords(taskID string) ([]*storage.PieceMetaRecord, error) {
	mm.cacheLocker.Lock(taskID, true)
	defer mm.cacheLocker.UnLock(taskID, true)
	return mm.storage.ReadPieceMetaRecords(taskID)
}

func (mm *metaDataManager) getPieceMd5Sign(taskID string) (string, []*storage.PieceMetaRecord, error) {
	pieceMetaRecords, err := mm.storage.ReadPieceMetaRecords(taskID)
	if err != nil {
		return "", nil, errors.Wrapf(err, "read piece meta file")
	}
	var pieceMd5 []string
	sort.Slice(pieceMetaRecords, func(i, j int) bool {
		return pieceMetaRecords[i].PieceNum < pieceMetaRecords[j].PieceNum
	})
	for _, piece := range pieceMetaRecords {
		pieceMd5 = append(pieceMd5, piece.Md5)
	}
	return digestutils.Sha256(pieceMd5...), pieceMetaRecords, nil
}

func (mm *metaDataManager) readFileMetaData(taskID string) (*storage.FileMetaData, error) {
	return mm.storage.ReadFileMetaData(taskID)
}

func (mm *metaDataManager) statDownloadFile(taskID string) (*storedriver.StorageInfo, error) {
	return mm.storage.StatDownloadFile(taskID)
}

func (mm *metaDataManager) readDownloadFile(taskID string) (io.ReadCloser, error) {
	return mm.storage.ReadDownloadFile(taskID)
}

func (mm *metaDataManager) resetRepo(task *types.SeedTask) error {
	mm.cacheLocker.Lock(task.ID, false)
	defer mm.cacheLocker.UnLock(task.ID, false)
	return mm.storage.ResetRepo(task)
}

func (mm *metaDataManager) writeDownloadFile(taskID string, offset int64, len int64, data io.Reader) error {
	return mm.storage.WriteDownloadFile(taskID, offset, len, data)
}
