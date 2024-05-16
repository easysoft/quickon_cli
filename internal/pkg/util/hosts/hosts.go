// Copyright © 2021 Sealos Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hosts

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/emirpasic/gods/maps/linkedhashmap"
	"github.com/ergoapi/util/exstr"
	"github.com/ergoapi/util/file"

	"github.com/easysoft/qcadmin/internal/pkg/util/log"
)

type HostFile struct {
	Path string
}

type hostname struct {
	Comment string
	Domain  string
	IP      string
}

func newHostname(comment string, domain string, ip string) *hostname {
	return &hostname{comment, domain, ip}
}

func (h *hostname) toString() string {
	return h.Comment + h.IP + " " + h.Domain + "\n"
}

func appendToFile(filePath string, hostname *hostname) {
	log := log.GetInstance()
	fp, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Warnf("failed opening file %s : %s", filePath, err)
		return
	}
	defer fp.Close()

	_, err = fp.WriteString(hostname.toString())
	if err != nil {
		log.Warnf("failed append string: %s: %s", filePath, err)
		return
	}
}

func (h *HostFile) ParseHostFile(path string) (*linkedhashmap.Map, error) {
	hlog := log.GetInstance()
	if !file.CheckFileExists(path) {
		hlog.Warnf("path %s is not exists", path)
		return nil, errors.New("path %s is not exists")
	}

	fp, fpErr := os.Open(path)
	if fpErr != nil {
		hlog.Warnf("open file '%s' failed", path)
		return nil, errors.Errorf("open file '%s' failed ", path)
	}
	defer fp.Close()

	br := bufio.NewReader(fp)
	lm := linkedhashmap.New()
	curComment := ""
	for {
		str, rErr := br.ReadString('\n')
		if rErr == io.EOF {
			break
		}
		if len(str) == 0 || str == "\r\n" || exstr.IsEmptyLine(str) {
			continue
		}

		if str[0] == '#' {
			// 处理注释
			curComment += str
			continue
		}
		tmpHostnameArr := strings.Fields(str)
		curDomain := strings.Join(tmpHostnameArr[1:], " ")
		curIP := exstr.TrimSpaceWS(tmpHostnameArr[0])

		checkIP := net.ParseIP(curIP)
		if checkIP == nil {
			continue
		}
		tmpHostname := newHostname(curComment, curDomain, curIP)
		lm.Put(tmpHostname.Domain, tmpHostname)
		curComment = ""
	}

	return lm, nil
}

func (h *HostFile) AppendHost(domain string, ip string) {
	if domain == "" || ip == "" {
		return
	}

	hostname := newHostname("", domain, ip)
	appendToFile(h.Path, hostname)
}

func (h *HostFile) writeToFile(hostnameMap *linkedhashmap.Map, path string) {
	hlog := log.GetInstance()
	if !file.CheckFileExists(path) {
		hlog.Warnf("path %s is not exists", path)
		return
	}

	fp, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		hlog.Warnf("open file '%s' failed: %v", path, err)
		return
	}
	defer fp.Close()

	hostnameMap.Each(func(key interface{}, value interface{}) {
		if v, ok := value.(*hostname); ok {
			_, writeErr := fp.WriteString(v.toString())
			if writeErr != nil {
				hlog.Warn(writeErr)
				return
			}
		}
	})
}

func (h *HostFile) DeleteDomain(domain string) {
	if domain == "" {
		return
	}

	hlog := log.GetInstance()

	currHostsMap, parseErr := h.ParseHostFile(h.Path)
	if parseErr != nil {
		hlog.Warnf("parse file failed" + parseErr.Error())
		return
	}
	_, found := currHostsMap.Get(domain)
	if currHostsMap == nil || !found {
		return
	}
	currHostsMap.Remove(domain)
	h.writeToFile(currHostsMap, h.Path)
}

func (h *HostFile) HasDomain(domain string) bool {
	if domain == "" {
		return false
	}
	hlog := log.GetInstance()
	currHostsMap, parseErr := h.ParseHostFile(h.Path)
	if parseErr != nil {
		hlog.Warnf("parse file failed" + parseErr.Error())
		return false
	}
	_, found := currHostsMap.Get(domain)
	if currHostsMap == nil || !found {
		return false
	}
	return true
}

func (h *HostFile) ListCurrentHosts() {
	hlog := log.GetInstance()
	currHostsMap, parseErr := h.ParseHostFile(h.Path)
	if parseErr != nil {
		hlog.Warnf("parse file failed" + parseErr.Error())
		return
	}
	if currHostsMap == nil {
		return
	}
	currHostsMap.Each(func(key interface{}, value interface{}) {
		if v, ok := value.(*hostname); ok {
			fmt.Print(v.toString())
		}
	})
}
