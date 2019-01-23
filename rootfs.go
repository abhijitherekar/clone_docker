package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

func mountProc(newroot string) error {
	return syscall.Mount("proc", filepath.Join(newroot, "/proc"), "proc", 0, "")
}

func mountPivot(root string) error {
	//1st mount the root to root itself as the pivot root as to think
	//that the root and root/.pivot_root are on 2 different FS
	if err := syscall.Mount(root, root, "", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return err
	}
	// then create a put_hold so that the pivot_root can use this
	// to store the FS of the exisiting process
	putHold := filepath.Join(root, "/.pivot_root")
	if err := os.Mkdir(putHold, 0755); err != nil {
		return err
	}
	if err := syscall.PivotRoot(root, putHold); err != nil {
		return err
	}
	//now according to the manpage change the curr dir to /
	if err := syscall.Chdir("/"); err != nil {
		return err
	}
	// now time to detach the current FS which is at /.pivot_root
	pivotDir := filepath.Join("/", ".pivot_root")
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount pivot_root dir %v", err)
	}
	// remove temporary directory
	return os.Remove(pivotDir)

}
