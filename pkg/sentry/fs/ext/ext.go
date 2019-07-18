// Copyright 2019 The gVisor Authors.
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

// Package ext implements readonly ext(2/3/4) filesystems.
package ext

import (
	"io"
	"sync"

	"gvisor.dev/gvisor/pkg/abi/linux"
	"gvisor.dev/gvisor/pkg/sentry/fs/ext/disklayout"
	"gvisor.dev/gvisor/pkg/syserror"
)

// Filesystem implements vfs.FilesystemImpl.
type Filesystem struct {
	// mu serializes changes to the Dentry tree.
	mu sync.Mutex

	// dev is the ReadSeeker for the underlying fs device. It is protected by mu
	// because it is only used when we are making additions to the dentry tree.
	//
	// The ext filesystems aim to maximize locality, i.e. place all the data
	// blocks of a file close together. On a spinning disk, locality reduces the
	// amount of movement of the head hence speeding up IO operations. On an SSD
	// there are no moving parts but locality increases the size of each transer
	// request. Hence, having mutual exclusion on the read seeker while reading a
	// file *should* help in achieving the intended performance gains.
	//
	// Note: This synchronization was not coupled with the ReadSeeker itself
	// because we want to synchronize across read/seek operations for the
	// performance gains mentioned above. Helps enforcing one-file-at-a-time IO.
	dev io.ReadSeeker

	// inodeCache maps absolute inode numbers to the corresponding Inode struct.
	// Inodes should be removed from this once their reference count hits 0.
	//
	// Protected by mu because every addtition and removal from this corresponds to
	// a change in the dentry tree.
	inodeCache map[uint32]*Inode

	// sb represents the filesystem superblock. Immutable after initialization.
	sb disklayout.SuperBlock

	// bgs represents all the block group descriptors for the filesystem.
	// Immutable after initialization.
	bgs []disklayout.BlockGroup
}

// newFilesystem is the Filesystem constructor.
func newFilesystem(dev io.ReadSeeker) (*Filesystem, error) {
	fs := Filesystem{dev: dev}
	var err error

	fs.sb, err = readSuperBlock(dev)
	if err != nil {
		return nil, err
	}

	if fs.sb.Magic() != linux.EXT_SUPER_MAGIC {
		// mount(2) specifies that EINVAL should be returned if the superblock is
		// invalid.
		return nil, syserror.EINVAL
	}

	fs.bgs, err = readBlockGroups(dev, fs.sb)
	if err != nil {
		return nil, err
	}

	return &fs, nil
}

// getInode gets the inode corresponding to the inode number passed in.
//
// Preconditions: must be holding fs.mu.
func (fs *Filesystem) getInode(inodeNum uint32) (*Inode, error) {
	if in, ok := fs.inodeCache[inodeNum]; ok {
		return in, nil
	}

	in, err := newInode(fs.dev, fs.sb, fs.bgs, inodeNum)
	if err != nil {
		return nil, err
	}

	fs.inodeCache[inodeNum] = in
	return in, nil
}
