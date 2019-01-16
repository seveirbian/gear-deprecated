package types

import (
    "os"
    "fmt"

    "bazil.org/fuse"
    "bazil.org/fuse/fs"

    "golang.org/x/net/context"
)

type ExtendFileInfo struct {
    Hash string
    FileInfo os.FileInfo
}

type FS struct {
    Files map[string]ExtendFileInfo
    LowerDir string
    UpperDir string
    WorkDir string
    MergedDir string
    PublicDir string
}

type Dir struct {
    Files map[string]ExtendFileInfo
    LowerDir string
    UpperDir string
    WorkDir string
    MergedDir string
    PublicDir string

    DirPath string
    DirAttr ExtendFileInfo
}

type File struct {
    FileInfo ExtendFileInfo
    LowerDir string
    UpperDir string
    WorkDir string
    MergedDir string
    PublicDir string

    FilePath string
    // File *
}

var _ fs.FS = (*FS)(nil)

func (f *FS) Root() (fs.Node, error) {
    d := &Dir {
        Files: f.Files, 
        LowerDir: f.LowerDir, 
        UpperDir: f.UpperDir, 
        WorkDir: f.WorkDir, 
        MergedDir: f.WorkDir, 
        PublicDir: f.PublicDir, 
    }

    return d, nil
}

var _ fs.Node = (*Dir)(nil)

func (d *Dir) Attr(c context.Context, a *fuse.Attr) error {
    if d.DirPath == "" {
        // root directory
        a = &fuse.Attr{Mode: os.ModeDir | 0755}
        return nil
    }

    a = dirAttr(d.DirAttr)

    return nil
}

func dirAttr(e ExtendFileInfo) *fuse.Attr {
    return &fuse.Attr{
        Size:   uint64(e.FileInfo.Size()),
        Mode:   e.FileInfo.Mode(),
        Mtime:  e.FileInfo.ModTime(),
        Ctime:  e.FileInfo.ModTime(),
        Crtime: e.FileInfo.ModTime(),
    }
}

var _ = fs.NodeRequestLookuper(&Dir{})

func (d *Dir) Lookup(c context.Context, req *fuse.LookupRequest, res *fuse.LookupResponse) (fs.Node, error) {
    path := req.Name
    fmt.Println(path)

    return nil, fuse.ENOENT
}






























