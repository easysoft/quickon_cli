// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package gops

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/google/gops/goprocess"
	"github.com/xlab/treeprint"
)

// DisplayProcessTree displays a tree of all the running Go processes.
func DisplayProcessTree() {
	ps := goprocess.FindAll()
	sort.Slice(ps, func(i, j int) bool {
		return ps[i].PPID < ps[j].PPID
	})
	pstree := make(map[int][]goprocess.P, len(ps))
	for _, p := range ps {
		pstree[p.PPID] = append(pstree[p.PPID], p)
	}
	tree := treeprint.New()
	tree.SetValue("...")
	seen := map[int]bool{}
	for _, p := range ps {
		constructProcessTree(p.PPID, p, pstree, seen, tree)
	}
	fmt.Println(tree.String())
}

// constructProcessTree constructs the process tree in a depth-first fashion.
func constructProcessTree(ppid int, process goprocess.P, pstree map[int][]goprocess.P, seen map[int]bool, tree treeprint.Tree) {
	if seen[ppid] {
		return
	}
	seen[ppid] = true
	if ppid != process.PPID {
		output := strconv.Itoa(ppid) + " (" + process.Exec + ")" + " {" + process.BuildVersion + "}"
		if process.Agent {
			tree = tree.AddMetaBranch("*", output)
		} else {
			tree = tree.AddBranch(output)
		}
	} else {
		tree = tree.AddBranch(ppid)
	}
	for index := range pstree[ppid] {
		process := pstree[ppid][index]
		constructProcessTree(process.PID, process, pstree, seen, tree)
	}
}
