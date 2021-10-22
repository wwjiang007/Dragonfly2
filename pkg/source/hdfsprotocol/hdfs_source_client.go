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

package hdfsprotocol

import (
	"bytes"
	"io"
	"net/url"
	"os/user"
	"strconv"
	"strings"
	"sync"
	"time"

	"d7y.io/dragonfly/v2/cdn/types"
	"github.com/go-http-utils/headers"
	"github.com/pkg/errors"

	"d7y.io/dragonfly/v2/pkg/source"
	"github.com/colinmarc/hdfs/v2"
)

const (
	HDFSClient = "hdfs"
)
const (
	layout = "2006-01-02 15:04:05"
	// hdfsUseDataNodeHostName set hdfs client whether user hostname connect to datanode
	hdfsUseDataNodeHostName = "dfs.client.use.datanode.hostname"
	// hdfsUseDataNodeHostNameValue set value is true
	hdfsUseDataNodeHostNameValue = "true"
)

func init() {
	source.Register(HDFSClient, NewHDFSSourceClient())
}

// hdfsSourceClient is an implementation of the interface of SourceClient.
type hdfsSourceClient struct {
	sync.RWMutex
	clientMap map[string]*hdfs.Client
}

// hdfsFileReaderClose is a combination object of the  io.LimitedReader and io.Closer
type hdfsFileReaderClose struct {
	limited io.Reader
	c       io.Closer
	buf     *bytes.Buffer
}

func newHdfsFileReaderClose(r io.Reader, n int64, c io.Closer) io.ReadCloser {
	return &hdfsFileReaderClose{
		limited: io.LimitReader(r, n),
		c:       c,
		buf:     bytes.NewBuffer(make([]byte, 512)),
	}
}

type HDFSSourceClientOption func(p *hdfsSourceClient)

func (h *hdfsSourceClient) GetContentLength(request *source.Request) (int64, error) {
	hdfsClient, path, err := h.getHDFSClientAndPath(request.URL)
	if err != nil {
		return types.UnKnownSourceFileLen, err
	}
	info, err := hdfsClient.Stat(path)
	if err != nil {
		return types.UnKnownSourceFileLen, err
	}
	if request.Header.Get(source.Range) != "" {
		rangeIndex := strings.Split(request.Header.Get(source.Range), "=")
		if strconv.Atoi(rangeIndex[1]) <= info.Size() {
			return int64(rang.EndIndex - rang.StartIndex), nil
		}
		return info.Size() - int64(rang.StartIndex), nil
	}
	return info.Size(), nil
}

func (h *hdfsSourceClient) IsSupportRange(request *source.Request) (bool, error) {
	hdfsClient, path, err := h.getHDFSClientAndPath(request.URL)
	if err != nil {
		return false, err
	}
	_, err = hdfsClient.Stat(path)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (h *hdfsSourceClient) IsExpired(request *source.Request) (bool, error) {
	lastModified := expireInfo[headers.LastModified]
	//eTag := expireInfo[headers.ETag]
	if lastModified == "" {
		return true, nil
	}

	hdfsClient, path, err := h.getHDFSClientAndPath(url)
	if err != nil {
		return false, err
	}

	info, err := hdfsClient.Stat(path)
	if err != nil {
		return false, err
	}
	t, err := time.ParseInLocation(layout, lastModified, time.Local)
	if err != nil {
		return false, err
	}
	return info.ModTime().Format(layout) != t.Format(layout), nil
}

func (h *hdfsSourceClient) Download(request *source.Request) (io.ReadCloser, error) {
	hdfsClient, path, err := h.getHDFSClientAndPath(request.URL)
	if err != nil {
		return nil, err
	}
	hdfsFile, err := hdfsClient.Open(path)
	if err != nil {
		return nil, err
	}

	// default read all data when rang is nil
	var limitReadN int64 = hdfsFile.Stat().Size()

	if rang != nil {
		_, err = hdfsFile.Seek(int64(rang.StartIndex), 0)
		if err != nil {
			hdfsFile.Close()
			return nil, err
		}
		limitReadN = int64(rang.EndIndex - rang.StartIndex)
	}

	return newHdfsFileReaderClose(hdfsFile, limitReadN, hdfsFile), nil
}

func (h *hdfsSourceClient) DownloadWithResponseHeader(request *source.Request) (*source.Response, error) {

	hdfsClient, path, err := h.getHDFSClientAndPath(request.URL)
	if err != nil {
		return nil, err
	}

	hdfsFile, err := hdfsClient.Open(path)
	if err != nil {
		return nil, err
	}

	fileInfo := hdfsFile.Stat()

	// default read all data when rang is nil
	var limitReadN int64 = fileInfo.Size()

	if rang != nil {
		_, err = hdfsFile.Seek(int64(rang.StartIndex), 0)
		if err != nil {
			hdfsFile.Close()
			return nil, nil, err
		}
		limitReadN = int64(rang.EndIndex - rang.StartIndex)
	}

	return newHdfsFileReaderClose(hdfsFile, limitReadN, hdfsFile), source.ResponseHeader{
		source.LastModified: fileInfo.ModTime().Format(layout),
	}, nil
}

func (h *hdfsSourceClient) Transform(header source.Header) source.Header {
	panic("implement me")
}

func (h *hdfsSourceClient) GetLastModifiedMillis(request *source.Request) (int64, error) {

	hdfsClient, path, err := h.getHDFSClientAndPath(request.URL)
	if err != nil {
		return -1, err
	}

	info, err := hdfsClient.Stat(path)
	if err != nil {
		return -1, err
	}

	return info.ModTime().UnixNano() / time.Millisecond.Nanoseconds(), nil
}

// getHDFSClient return hdfs client
func (h *hdfsSourceClient) getHDFSClient(url *url.URL) (*hdfs.Client, error) {
	// get client for map
	h.RWMutex.RLock()
	if client, ok := h.clientMap[url.Host]; ok {
		h.RWMutex.RUnlock()
		return client, nil
	}
	h.RWMutex.RUnlock()

	// create client option
	options := hdfs.ClientOptionsFromConf(map[string]string{
		hdfsUseDataNodeHostName: hdfsUseDataNodeHostNameValue,
	})
	options.Addresses = strings.Split(url.Host, ",")
	u, err := user.Current()
	if err != nil {
		return nil, err
	}
	options.User = u.Username

	// create hdfs client and put map
	h.RWMutex.Lock()
	client, err := hdfs.NewClient(options)
	if err != nil {
		h.RWMutex.Unlock()
		return nil, err
	}
	h.clientMap[url.Host] = client
	h.RWMutex.Unlock()
	return client, err
}

// getHDFSClientAndPath return client and path
func (h *hdfsSourceClient) getHDFSClientAndPath(url *url.URL) (*hdfs.Client, string, error) {
	client, err := h.getHDFSClient(url)
	if err != nil {
		return nil, "", errors.Errorf("hdfs create client failed, url is %s", url)
	}
	return client, url.Path, nil
}

func NewHDFSSourceClient(opts ...HDFSSourceClientOption) source.ResourceClient {
	sourceClient := &hdfsSourceClient{
		clientMap: make(map[string]*hdfs.Client),
	}
	for i := range opts {
		opts[i](sourceClient)
	}
	return sourceClient
}

var _ source.ResourceClient = (*hdfsSourceClient)(nil)

func (rc *hdfsFileReaderClose) Read(p []byte) (n int, err error) {
	return rc.limited.Read(p)
}

func (rc *hdfsFileReaderClose) Close() error {
	return rc.c.Close()
}

func (rc *hdfsFileReaderClose) WriteTo(w io.Writer) (n int64, err error) {
	_, err = rc.limited.Read(rc.buf.Bytes())
	if err != nil {
		return -1, err
	}
	return rc.buf.WriteTo(w)
}
