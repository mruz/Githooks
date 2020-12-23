package hooks

import (
	"os"
	"path"
	"regexp"
	"rycus86/githooks/build"
	cm "rycus86/githooks/common"
)

var runWrapperDetectionRegex = regexp.MustCompile("https://github.com/rycus86/githooks")

// IsRunWrapper answers the question if `filePath`
// is a Githooks hook template file.
func IsRunWrapper(filePath string) (bool, error) {
	return cm.MatchLineRegexInFile(filePath, runWrapperDetectionRegex)
}

// getHookReplacementName returns the file name of a replaced custom Git hook.
func getHookReplacementName(fileName string) string {
	return path.Join(path.Dir(fileName), path.Base(fileName)+".replaced.githooks")
}

// GetRunWrapperContent gets the bytes of the hook template.
func getRunWrapperContent() ([]byte, error) {
	return build.Asset("base-template-wrapper.sh")
}

// WriteRunWrapper writes the run wrapper to the file `filePath`.
func WriteRunWrapper(filePath string) (err error) {
	runWrapperContent, err := getRunWrapperContent()
	cm.AssertNoErrorPanic(err, "Could not get embedded run wrapper content.")

	file, err := os.Create(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = file.Write(runWrapperContent)
	if err != nil {
		return
	}
	err = file.Sync()
	if err != nil {
		return
	}

	// Make executable
	_ = file.Close()
	err = cm.MakeExecutbale(filePath)

	return
}

var lfsDetectionRe = regexp.MustCompile(`(git\s+lfs|git-lfs)`)

// HookDisableOption are the options
// how to disable a hook.
type HookDisableOption int

const (
	// BackupHook defines that a hook file gets backuped.
	BackupHook HookDisableOption = 1
	// DeleteHook defines that a hook file gets deleted.
	DeleteHook HookDisableOption = 2
)

func disableHookIfLFSDetected(
	filePath string,
	tempDir string,
	disableCallBack func(file string) HookDisableOption) (disabled bool, deleted bool, err error) {

	found, err := cm.MatchLineRegexInFile(filePath, lfsDetectionRe)
	if err != nil {
		return
	}

	if found {
		disableOption := disableCallBack(filePath)

		switch disableOption {
		default:
			fallthrough
		case BackupHook:
			err = os.Rename(filePath, filePath+".disabled.githooks")
			disabled = true
		case DeleteHook:
			// Don't delete the hook yet, move it to the temporary directory.
			// The file could potentially be opened/read.
			moveFile := cm.GetTempFileName(tempDir, "-"+path.Base(filePath))
			err = os.Rename(filePath, moveFile)
			disabled = true
			deleted = true
		}

		if err != nil {
			return
		}
	}

	return
}

func moveExistingHooks(
	dest string,
	tempDir string,
	disableHookIfLFS func(file string) HookDisableOption,
	log cm.ILogContext) error {

	// Check there is already a Git hook in place and replace it.
	if cm.IsFile(dest) {

		isTemplate, err := IsRunWrapper(dest)

		if err != nil {
			return cm.CombineErrors(err,
				cm.ErrorF("Could not detect if '%s' is a Githooks run wrapper.", dest))
		}

		if !isTemplate {

			// Try to detect a potential LFS statements and
			// disable the hook (backup or delete).
			if disableHookIfLFS != nil {
				disabled, deleted, err := disableHookIfLFSDetected(dest, tempDir, disableHookIfLFS)
				if err != nil {
					return err
				}

				if log != nil {
					if disabled && deleted {
						log.WarnF("Hook '%s' is now disabled (deleted)", dest)
					} else if disabled && !deleted {
						log.WarnF("Hook '%s' is now disabled (backuped)", dest)
					}
				}
			}

			// Replace the file normally if it is still existing.
			if cm.IsFile(dest) {
				if log != nil {
					log.InfoF("Saving existing Git hook '%s'.", dest)
				}

				newDest := getHookReplacementName(dest)
				err = os.Rename(dest, newDest)
				if err != nil {
					return cm.CombineErrors(err,
						cm.ErrorF("Could not rename file '%s' to '%s'.", dest, newDest))
				}
			}
		}
	}

	return nil
}

// InstallRunWrappers installs run wrappers for the given `hookNames` in `dir`.
// Existing custom hooks get renamed.
// All deleted hooks by this function get moved to the `tempDir` directory, because
// we should not delete them yet.
func InstallRunWrappers(
	dir string,
	hookNames []string,
	tempDir string,
	disableHookIfLFS func(file string) HookDisableOption,
	log cm.ILogContext) error {

	for _, hookName := range hookNames {

		dest := path.Join(dir, hookName)

		err := moveExistingHooks(dest, tempDir, disableHookIfLFS, log)
		if err != nil {
			return err
		}

		if log != nil {
			log.InfoF("Saving Githooks run wrapper to '%s'.", dest)
		}

		if cm.IsFile(dest) {
			// If still existing:
			// The file `dest` could be currently running,
			// therefore we move it to the temporary directly.
			// On Unix we could simply remove the file.
			// But on Windows, an opened file (mostly) cannot be deleted.
			// it might work, but is ugly.
			moveDest := cm.GetTempFileName(tempDir, "-"+path.Base(dest))
			err = os.Rename(dest, moveDest)
			if err != nil {
				return cm.CombineErrors(err,
					cm.ErrorF("Could not move file '%s' to '%s'.", dest, moveDest))
			}
		}

		err = WriteRunWrapper(dest)
		if err != nil {
			return cm.CombineErrors(err,
				cm.ErrorF("Could not write Githooks run wrapper to '%s'.", dest))
		}
	}

	return nil
}