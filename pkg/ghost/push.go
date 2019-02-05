package ghost

import (
	"git-ghost/pkg/ghost/git"
	"git-ghost/pkg/util"

	log "github.com/Sirupsen/logrus"
)

type PushOptions struct {
	WorkingEnvSpec
	*LocalBaseBranchSpec
	*LocalModBranchSpec
}

type PushResult struct {
	*LocalBaseBranch
	*LocalModBranch
}

func Push(options PushOptions) (*PushResult, error) {
	log.WithFields(util.ToFields(options)).Debug("push command with")

	var result PushResult
	if options.LocalBaseBranchSpec != nil {
		branch, err := pushGhostBranch(options.LocalBaseBranchSpec, options.WorkingEnvSpec)
		if err != nil {
			return nil, err
		}
		localBaseBranch, _ := branch.(*LocalBaseBranch)
		result.LocalBaseBranch = localBaseBranch
	}

	if options.LocalModBranchSpec != nil {
		branch, err := pushGhostBranch(options.LocalModBranchSpec, options.WorkingEnvSpec)
		if err != nil {
			return nil, err
		}
		localModBranch, _ := branch.(*LocalModBranch)
		result.LocalModBranch = localModBranch
	}

	return &result, nil
}

func pushGhostBranch(branchSpec GhostBranchSpec, workingEnvSpec WorkingEnvSpec) (GhostBranch, error) {
	workingEnv, err := workingEnvSpec.initialize()
	if err != nil {
		return nil, err
	}
	defer workingEnv.clean()
	dstDir := workingEnv.GhostDir
	branch, err := branchSpec.CreateBranch(*workingEnv)
	if err != nil {
		return nil, err
	}
	if branch == nil {
		return nil, nil
	}
	existence, err := git.ValidateRemoteBranchExistence(
		workingEnv.GhostRepo,
		branch.BranchName(),
	)
	if err != nil {
		return nil, err
	}
	if existence {
		log.WithFields(log.Fields{
			"branch":    branch.BranchName(),
			"ghostRepo": workingEnv.GhostRepo,
		}).Info("skipped pushing existing branch")
		return nil, nil
	}

	log.WithFields(log.Fields{
		"branch":    branch.BranchName(),
		"ghostRepo": workingEnv.GhostRepo,
	}).Info("pushing branch")
	err = git.Push(dstDir, branch.BranchName())
	if err != nil {
		return nil, err
	}
	return branch, nil
}
