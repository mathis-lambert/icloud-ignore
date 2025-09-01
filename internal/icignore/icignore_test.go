package icignore

import (
    "os"
    "path/filepath"
    "runtime"
    "testing"
)

func TestWithWithoutNoSync(t *testing.T) {
    if got := withNoSync("Projects"); got != "Projects.nosync" {
        t.Fatalf("withNoSync mismatch: %q", got)
    }
    if got := withoutNoSync("Projects.nosync"); got != "Projects" {
        t.Fatalf("withoutNoSync mismatch: %q", got)
    }
}

func TestExpandPath(t *testing.T) {
    if runtime.GOOS == "windows" { t.Skip("not relevant on windows") }
    home, err := os.UserHomeDir()
    if err != nil { t.Skip("no home dir") }
    got, err := ExpandPath("~")
    if err != nil { t.Fatal(err) }
    if got != home { t.Fatalf("want %q got %q", home, got) }
}

func TestSymlinkHelpers(t *testing.T) {
    dir := t.TempDir()
    real := filepath.Join(dir, "data.nosync")
    link := filepath.Join(dir, "data")
    if err := os.Mkdir(real, 0o755); err != nil { t.Fatal(err) }
    if err := os.Symlink("data.nosync", link); err != nil { t.Fatal(err) }

    if !isSymlink(link) { t.Fatalf("expected symlink at %s", link) }
    if !isSymlinkTo(link, "data.nosync") { t.Fatalf("symlink target mismatch") }
}

func TestStatusUnsuffixedReal(t *testing.T) {
    dir := t.TempDir()
    real := filepath.Join(dir, "proj")
    if err := os.Mkdir(real, 0o755); err != nil { t.Fatal(err) }
    s, err := Status(real)
    if err != nil { t.Fatal(err) }
    if s.RealIsSuffixed { t.Fatalf("expected unsuffixed real directory") }
}

