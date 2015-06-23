package file

import (
	"os"
	"sync"
)

type fileInfoCacheImpl struct {
	sync.RWMutex
	m map[string]os.FileInfo
}

// FileInfoCache stores the cache of FileInfo.
var FileInfoCache = &fileInfoCacheImpl{}

// Stat returns the FileInfo of the named file. If Stat has been called
// with the named file, it returns the FileInfo from the cache.
//
// To clear the FileInfo cache of the named file, use Clear function.
func (c *fileInfoCacheImpl) Stat(name string) (os.FileInfo, error) {
	c.RLock()
	if c.m != nil && c.m[name] != nil {
		defer c.RUnlock()
		return c.m[name], nil
	}

	c.RUnlock()
	c.Lock()
	defer c.Unlock()

	if c.m == nil {
		c.m = make(map[string]os.FileInfo)
	}

	info, err := os.Stat(name)
	if err != nil {
		return nil, err
	}
	c.m[name] = info
	return info, nil
}

func (c *fileInfoCacheImpl) ClearAndUnlock(name string) {
	defer c.Unlock()
	if c.m != nil && c.m[name] != nil {
		c.m[name] = nil
	}
}

func (c *fileInfoCacheImpl) ClearWithLock(name string) {
	c.Lock()
	c.ClearAndUnlock(name)
}
