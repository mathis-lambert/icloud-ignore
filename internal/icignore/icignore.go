package icignore

import (
    "errors"
    "fmt"
    "os"
    "path/filepath"
    "strings"
)

const nosyncSuffix = ".nosync"

// Options controls behavior of operations.
type Options struct {
    DryRun  bool
    Verbose bool
}

// StatusInfo describes the current state of a path.
type StatusInfo struct {
    UnsuffixedPath string
    SuffixedPath   string
    HasSymlink     bool
    SymlinkTarget  string
    RealIsSuffixed bool
}

// ExpandPath resolves ~ and makes an absolute, cleaned path.
func ExpandPath(p string) (string, error) {
    if p == "" {
        return "", errors.New("empty path")
    }
    if strings.HasPrefix(p, "~") {
        home, err := os.UserHomeDir()
        if err != nil {
            return "", err
        }
        if p == "~" {
            p = home
        } else if strings.HasPrefix(p, "~/") {
            p = filepath.Join(home, p[2:])
        }
    }
    abs, err := filepath.Abs(p)
    if err != nil {
        return "", err
    }
    return abs, nil
}

// Status returns a StatusInfo for a path.
func Status(inputPath string) (*StatusInfo, error) {
    parent := filepath.Dir(inputPath)
    base := filepath.Base(inputPath)

    suffixed := withNoSync(base)
    unsuffixed := withoutNoSync(base)

    suffixedPath := filepath.Join(parent, suffixed)
    unsuffixedPath := filepath.Join(parent, unsuffixed)

    info := &StatusInfo{
        UnsuffixedPath: unsuffixedPath,
        SuffixedPath:   suffixedPath,
    }

    if fi, err := os.Lstat(unsuffixedPath); err == nil && fi.Mode()&os.ModeSymlink != 0 {
        info.HasSymlink = true
        target, _ := os.Readlink(unsuffixedPath)
        info.SymlinkTarget = target
    }

    if fi, err := os.Stat(suffixedPath); err == nil && fi.IsDir() {
        info.RealIsSuffixed = true
        return info, nil
    }
    if fi, err := os.Stat(unsuffixedPath); err == nil && fi.IsDir() {
        info.RealIsSuffixed = false
        return info, nil
    }

    if !exists(suffixedPath) && !exists(unsuffixedPath) {
        return nil, fmt.Errorf("path not found: %s", inputPath)
    }

    return info, nil
}

// Ignore renames the folder to add .nosync and creates a symlink at the original path.
func Ignore(inputPath string, opt Options) error {
    parent := filepath.Dir(inputPath)
    base := filepath.Base(inputPath)

    if strings.HasSuffix(base, nosyncSuffix) {
        suffixedPath := inputPath
        unsuffixedPath := filepath.Join(parent, withoutNoSync(base))

        if exists(unsuffixedPath) {
            if isSymlinkTo(unsuffixedPath, filepath.Base(suffixedPath)) {
                if opt.Verbose {
                    fmt.Println("already ignored: symlink exists")
                }
                return nil
            }
            return fmt.Errorf("destination exists and is not the correct symlink: %s", unsuffixedPath)
        }

        if opt.DryRun {
            fmt.Printf("DRY-RUN: ln -s %s %s\n", filepath.Base(suffixedPath), unsuffixedPath)
            return nil
        }
        return os.Symlink(filepath.Base(suffixedPath), unsuffixedPath)
    }

    unsuffixedPath := inputPath
    suffixedBase := withNoSync(base)
    suffixedPath := filepath.Join(parent, suffixedBase)

    if isSymlink(unsuffixedPath) {
        target, _ := os.Readlink(unsuffixedPath)
        if strings.HasSuffix(target, nosyncSuffix) {
            if opt.Verbose {
                fmt.Println("already ignored: symlink points to .nosync target")
            }
            return nil
        }
        return fmt.Errorf("%s is a symlink not pointing to a .nosync target (%s)", unsuffixedPath, target)
    }

    fi, err := os.Stat(unsuffixedPath)
    if err != nil {
        return err
    }
    if !fi.IsDir() {
        return fmt.Errorf("not a directory: %s", unsuffixedPath)
    }

    if exists(suffixedPath) {
        return fmt.Errorf("target already exists: %s (try 'icignore status' or 'icignore unignore')", suffixedPath)
    }

    if opt.DryRun {
        fmt.Printf("DRY-RUN: mv %s %s\n", unsuffixedPath, suffixedPath)
        fmt.Printf("DRY-RUN: ln -s %s %s\n", suffixedBase, unsuffixedPath)
        return nil
    }

    if err := os.Rename(unsuffixedPath, suffixedPath); err != nil {
        return err
    }
    if err := os.Symlink(suffixedBase, unsuffixedPath); err != nil {
        _ = os.Rename(suffixedPath, unsuffixedPath)
        return err
    }
    if opt.Verbose {
        fmt.Printf("ignored: %s -> %s (symlink created)\n", unsuffixedPath, suffixedBase)
    }
    return nil
}

// Unignore removes the symlink and renames the .nosync folder back.
func Unignore(inputPath string, opt Options) error {
    parent := filepath.Dir(inputPath)
    base := filepath.Base(inputPath)

    var unsuffixedPath, suffixedPath string
    if strings.HasSuffix(base, nosyncSuffix) {
        suffixedPath = inputPath
        unsuffixedPath = filepath.Join(parent, withoutNoSync(base))
    } else {
        unsuffixedPath = inputPath
        suffixedPath = filepath.Join(parent, withNoSync(base))
    }

    if isSymlinkTo(unsuffixedPath, filepath.Base(suffixedPath)) {
        if opt.DryRun {
            fmt.Printf("DRY-RUN: rm %s\n", unsuffixedPath)
            fmt.Printf("DRY-RUN: mv %s %s\n", suffixedPath, unsuffixedPath)
            return nil
        }
        if err := os.Remove(unsuffixedPath); err != nil {
            return err
        }
        if !exists(suffixedPath) {
            return fmt.Errorf("missing expected directory: %s", suffixedPath)
        }
        if err := os.Rename(suffixedPath, unsuffixedPath); err != nil {
            return err
        }
        if opt.Verbose {
            fmt.Printf("unignored: %s (symlink removed, directory restored)\n", unsuffixedPath)
        }
        return nil
    }

    if exists(unsuffixedPath) && !exists(suffixedPath) && !isSymlink(unsuffixedPath) {
        if opt.Verbose {
            fmt.Println("already unignored: no symlink and no .nosync directory")
        }
        return nil
    }

    if strings.HasSuffix(filepath.Base(inputPath), nosyncSuffix) {
        if isSymlink(unsuffixedPath) {
            if opt.DryRun {
                fmt.Printf("DRY-RUN: rm %s\n", unsuffixedPath)
            } else if err := os.Remove(unsuffixedPath); err != nil {
                return err
            }
        } else if exists(unsuffixedPath) {
            return fmt.Errorf("destination exists and is not a symlink: %s", unsuffixedPath)
        }
        if opt.DryRun {
            fmt.Printf("DRY-RUN: mv %s %s\n", inputPath, unsuffixedPath)
            return nil
        }
        return os.Rename(inputPath, unsuffixedPath)
    }

    return fmt.Errorf("unable to unignore: expected symlink at %s pointing to %s", unsuffixedPath, filepath.Base(suffixedPath))
}

// Helpers
func withNoSync(name string) string {
    if strings.HasSuffix(name, nosyncSuffix) {
        return name
    }
    return name + nosyncSuffix
}

func withoutNoSync(name string) string {
    if strings.HasSuffix(name, nosyncSuffix) {
        return strings.TrimSuffix(name, nosyncSuffix)
    }
    return name
}

func exists(p string) bool {
    _, err := os.Lstat(p)
    return err == nil
}

func isSymlink(p string) bool {
    fi, err := os.Lstat(p)
    if err != nil {
        return false
    }
    return fi.Mode()&os.ModeSymlink != 0
}

func isSymlinkTo(linkPath string, wantTargetBase string) bool {
    if !isSymlink(linkPath) {
        return false
    }
    target, err := os.Readlink(linkPath)
    if err != nil {
        return false
    }
    if target == wantTargetBase {
        return true
    }
    return filepath.Base(target) == wantTargetBase
}
