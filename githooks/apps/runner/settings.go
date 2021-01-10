package main

import (
	cm "rycus86/githooks/common"
	"rycus86/githooks/git"
	strs "rycus86/githooks/strings"
)

// HookSettings defines hooks related settings for this run.
type HookSettings struct {
	Args               []string       // Rest arguments.
	ExecX              cm.ExecContext // Execution context for executables (working dir is this repository).
	GitX               *git.Context   // Git context to execute commands (working dir is this repository).
	RepositoryDir      string         // Repository directory (bare, non-bare).
	RepositoryHooksDir string         // Directory with hooks for this repository.
	GitDirWorktree     string         // Git directory. (for worktrees this points to the worktree Git dir).
	InstallDir         string         // Install directory.

	HookPath string // Absolute path of the hook executing this runner.
	HookName string // Name of the hook.
	HookDir  string // Directory of the hook.

	IsRepoTrusted                bool // If the repository is a trusted repository.
	FailOnNonExistingSharedHooks bool // If Githooks should fail if there are shared hooks demanded which are not existing.
}

func (s HookSettings) toString() string {
	return strs.Fmt(
		"\n- Args: '%q'\n"+
			" • Repo Path: '%s'\n"+
			" • Repo Hooks: '%s'\n"+
			" • Git Dir Worktree: '%s'\n"+
			" • Install Dir: '%s'\n"+
			" • Hook Path: '%s'\n"+
			" • Hook Name: '%s'\n"+
			" • Trusted: '%v'",
		s.Args, s.RepositoryDir,
		s.RepositoryHooksDir, s.GitDirWorktree,
		s.InstallDir, s.HookPath, s.HookName, s.IsRepoTrusted)
}
