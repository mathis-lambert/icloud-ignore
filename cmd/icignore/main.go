package main

import (
    "flag"
    "fmt"
    "os"
    "strings"

    core "github.com/mathis-lambert/icloud-ignore/internal/icignore"
)

var (
    // Set via -ldflags at build time: -X main.version=...
    version = "0.1.0"
)

func usage() {
    fmt.Fprintf(os.Stderr, "icignore %s\n", version)
    fmt.Fprintf(os.Stderr, "Usage:\n")
    fmt.Fprintf(os.Stderr, "  icignore [--dry-run] [--verbose] <command> <path>\n\n")
    fmt.Fprintf(os.Stderr, "Commands:\n")
    fmt.Fprintf(os.Stderr, "  ignore   Exclude a folder from iCloud sync (.nosync + symlink)\n")
    fmt.Fprintf(os.Stderr, "  unignore Restore sync for a folder (remove symlink, rename back)\n")
    fmt.Fprintf(os.Stderr, "  status   Show status for a folder (symlink and real path)\n")
    fmt.Fprintf(os.Stderr, "  version  Print version\n\n")
    fmt.Fprintf(os.Stderr, "Examples:\n")
    fmt.Fprintf(os.Stderr, "  icignore ignore ~/Documents/Projects\n")
    fmt.Fprintf(os.Stderr, "  icignore status ~/Documents/Projects\n")
}

func main() {
    dryRun := flag.Bool("dry-run", false, "Print actions without changing anything")
    verbose := flag.Bool("verbose", false, "Verbose output")
    flag.Usage = usage
    flag.Parse()

    args := flag.Args()
    if len(args) == 0 {
        usage()
        os.Exit(2)
    }

    cmd := strings.ToLower(args[0])
    if cmd == "version" {
        fmt.Println(version)
        return
    }

    if len(args) < 2 {
        fmt.Fprintln(os.Stderr, "error: missing <path>")
        usage()
        os.Exit(2)
    }

    rawPath := args[1]
    path, err := core.ExpandPath(rawPath)
    if err != nil { fatalErr(err) }

    switch cmd {
    case "ignore":
        if err := core.Ignore(path, core.Options{DryRun: *dryRun, Verbose: *verbose}); err != nil { fatalErr(err) }
    case "unignore":
        if err := core.Unignore(path, core.Options{DryRun: *dryRun, Verbose: *verbose}); err != nil { fatalErr(err) }
    case "status":
        if s, err := core.Status(path); err != nil {
            fatalErr(err)
        } else {
            if s.HasSymlink {
                fmt.Printf("SYMLINK: %s -> %s\n", s.UnsuffixedPath, s.SymlinkTarget)
            }
            if s.RealIsSuffixed {
                fmt.Printf("REAL: %s (excluded from iCloud via .nosync suffix)\n", s.SuffixedPath)
            } else {
                fmt.Printf("REAL: %s (no .nosync suffix)\n", s.UnsuffixedPath)
            }
        }
    default:
        fmt.Fprintf(os.Stderr, "error: unknown command %q\n", cmd)
        usage()
        os.Exit(2)
    }
}

func fatalErr(err error) {
    fmt.Fprintln(os.Stderr, "error:", err)
    os.Exit(1)
}
