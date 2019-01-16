package nonCooperativeDriver

import (
    "fmt"
    // "errors"
    "io"
    "os"
    "log"
    "path"
    // "system"
    "path/filepath"

    graphDriver "github.com/docker/docker/daemon/graphdriver"
    "github.com/docker/docker/pkg/system"
    "github.com/docker/docker/pkg/archive"
    "github.com/docker/docker/pkg/containerfs"
    "github.com/docker/docker/pkg/idtools"
    "github.com/docker/docker/pkg/locker"
    graphPlugin "github.com/docker/go-plugins-helpers/graphdriver"
)

type NonCooperativeDriver struct {
    Home string
    Options []string
    UidMaps []idtools.IDMap
    GidMaps []idtools.IDMap

    locker        *locker.Locker
}

// initlize driver
func (d *NonCooperativeDriver) Init(home string, options []string, uidMaps, gidMaps []idtools.IDMap) error {
    d.Home = home
    d.Options = options
    d.UidMaps = uidMaps
    d.GidMaps = gidMaps

    fmt.Println("This is driver: ")
    fmt.Println(d)

    // get root user's uid and gid
    rootUID, rootGID, err := idtools.GetRootUIDGID(d.UidMaps, d.GidMaps)
    if err != nil {
        return err
    }
    root := idtools.Identity{UID: rootUID, GID: rootGID}

    // create the path of id's parent
    if err := idtools.MkdirAllAndChown(filepath.Join(home, "public"), 0755, root); err != nil {
        return err
    }

    return nil
}

// 
func (d *NonCooperativeDriver) Create(id, parent, mountLabel string, storageOpt map[string]string) error {
    fmt.Println("Create")
    fmt.Println("id: " + id)
    fmt.Println("parent: " + parent)
    fmt.Println("mountLabel: " + mountLabel)
    fmt.Println("storageOpt: ")
    fmt.Println(storageOpt)

    // get absolute path of id
    dir := path.Join(d.Home, id)

    // get root user's uid and gid
    rootUID, rootGID, err := idtools.GetRootUIDGID(d.UidMaps, d.GidMaps)
    if err != nil {
        return err
    }
    root := idtools.Identity{UID: rootUID, GID: rootGID}

    // create the path of id's parent
    if err := idtools.MkdirAllAndChown(path.Dir(dir), 0700, root); err != nil {
        return err
    }
    // create id dir
    if err := idtools.MkdirAndChown(dir, 0700, root); err != nil {
        return err
    }
    // create id's lower dir
    if err := idtools.MkdirAndChown(path.Join(dir, "lower"), 0755, root); err != nil {
        return err
    }

    return nil
}

func (d *NonCooperativeDriver) CreateReadWrite(id, parent, mountLabel string, storageOpt map[string]string) error {
    fmt.Println("CreateReadWrite")
    fmt.Println("id: " + id)
    fmt.Println("parent: " + parent)
    fmt.Println("mountLabel: " + mountLabel)
    fmt.Println("storageOpt: ")
    fmt.Println(storageOpt)

    // get absolute path of id
    dir := path.Join(d.Home, id)

    // get root user's uid and gid
    rootUID, rootGID, err := idtools.GetRootUIDGID(d.UidMaps, d.GidMaps)
    if err != nil {
        return err
    }
    root := idtools.Identity{UID: rootUID, GID: rootGID}

    // create the path of id's parent
    if err := idtools.MkdirAllAndChown(path.Dir(dir), 0700, root); err != nil {
        return err
    }
    // create id dir
    if err := idtools.MkdirAndChown(dir, 0700, root); err != nil {
        return err
    }
    // create id's upper dir
    if err := idtools.MkdirAndChown(path.Join(dir, "upper"), 0755, root); err != nil {
        return err
    }
    // create id's work dir
    if err := idtools.MkdirAndChown(path.Join(dir, "work"), 0700, root); err != nil {
        return err
    }
    // create id's lower link to parent's lower
    if err := os.Symlink(path.Join("..", parent, "lower"), path.Join(d.Home, id, "lower")); err != nil {
        return err
    }
    // create id's merged dir
    if err := idtools.MkdirAndChown(path.Join(dir, "merged"), 0700, root); err != nil {
        return err
    }

    return nil
}

func (d *NonCooperativeDriver) Remove(id string) error {
    if id == "" {
        return fmt.Errorf("refusing to remove the directories: id is empty")
    }

    // lock and unlock
    d.locker.Lock(id)
    defer d.locker.Unlock(id)

    // get absolute path of id
    dir := path.Join(d.Home, id)

    if err := system.EnsureRemoveAll(dir); err != nil && !os.IsNotExist(err) {
        return err
    }
    return nil
}

// Get creates and mounts the required file system for the given id and returns the mount path
func (d *NonCooperativeDriver) Get(id, mountLabel string) (containerfs.ContainerFS, error) {
    if d == nil {
        return nil, errNotInitialized
    }
    return d.driver.Get(id, mountLabel)
}

func (d *NonCooperativeDriver) Put(id string) error {
    if d == nil {
        return errNotInitialized
    }
    return d.driver.Put(id)
}

func (d *NonCooperativeDriver) Exists(id string) bool {
    if d == nil {
        return false
    }
    return d.driver.Exists(id)
}

func (d *NonCooperativeDriver) Status() [][2]string {
    if d == nil {
        return nil
    }
    return d.driver.Status()
}

func (d *NonCooperativeDriver) GetMetadata(id string) (map[string]string, error) {
    if d == nil {
        return nil, errNotInitialized
    }
    return d.driver.GetMetadata(id)
}

func (d *NonCooperativeDriver) Cleanup() error {
    if d == nil {
        return errNotInitialized
    }
    return d.driver.Cleanup()
}

func (d *NonCooperativeDriver) Diff(id, parent string) io.ReadCloser {
    if d == nil {
        return nil
    }
    // FIXME(samoht): how do we pass the error to the driver?
    archive, err := d.driver.Diff(id, parent)
    if err != nil {
        log.Fatalf("Diff: error in stream %v", err)
    }
    return archive
}

func changeKind(c archive.ChangeType) graphPlugin.ChangeKind {
    switch c {
    case archive.ChangeModify:
        return graphPlugin.Modified
    case archive.ChangeAdd:
        return graphPlugin.Added
    case archive.ChangeDelete:
        return graphPlugin.Deleted
    }
    return 0
}

func (d *NonCooperativeDriver) Changes(id, parent string) ([]graphPlugin.Change, error) {
    if d == nil {
        return nil, errNotInitialized
    }
    cs, err := d.driver.Changes(id, parent)
    if err != nil {
        return nil, err
    }
    changes := make([]graphPlugin.Change, len(cs))
    for _, c := range cs {
        change := graphPlugin.Change{
            Path: c.Path,
            Kind: changeKind(c.Kind),
        }
        changes = append(changes, change)
    }
    return changes, nil
}

func (d *NonCooperativeDriver) ApplyDiff(id, parent string, archive io.Reader) (int64, error) {
    if d == nil {
        return 0, errNotInitialized
    }
    return d.driver.ApplyDiff(id, parent, archive)
}

func (d *NonCooperativeDriver) DiffSize(id, parent string) (int64, error) {
    if d == nil {
        return 0, errNotInitialized
    }
    return d.driver.DiffSize(id, parent)
}

func (d *NonCooperativeDriver) Capabilities() graphDriver.Capabilities {
    if d == nil {
        return graphDriver.Capabilities{}
    }
    if capDriver, ok := d.driver.(graphDriver.CapabilityDriver); ok {
        return capDriver.Capabilities()
    }
    return graphDriver.Capabilities{}
}































