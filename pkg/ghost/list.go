package ghost

import (
	"bytes"
	"fmt"
	"git-ghost/pkg/util"

	log "github.com/Sirupsen/logrus"
)

type ListOptions struct {
	WorkingEnvSpec
	Prefix string
	*ListCommitsBranchSpec
	*ListDiffBranchSpec
}

type ListResult struct {
	LocalBaseBranches LocalBaseBranches
	LocalModBranches  LocalModBranches
}

func List(options ListOptions) (*ListResult, error) {
	log.WithFields(util.ToFields(options)).Debug("list command with")

	res := ListResult{}

	if options.ListCommitsBranchSpec != nil {
		resolved := options.ListCommitsBranchSpec.Resolve(options.SrcDir)
		branches, err := resolved.GetBranches(options.GhostRepo, options.Prefix)
		if err != nil {
			return nil, err
		}
		res.LocalBaseBranches = branches
	}

	if options.ListDiffBranchSpec != nil {
		resolved := options.ListDiffBranchSpec.Resolve(options.SrcDir)
		branches, err := resolved.GetBranches(options.GhostRepo, options.Prefix)
		if err != nil {
			return nil, err
		}
		res.LocalModBranches = branches
	}

	res.LocalBaseBranches.Sort()
	res.LocalModBranches.Sort()

	return &res, nil
}

func (res *ListResult) PrettyString() string {
	// TODO: Make it prettier
	var buffer bytes.Buffer
	if len(res.LocalBaseBranches) > 0 {
		buffer.WriteString("Local Base Branches:\n")
		buffer.WriteString("\n")
	}
	for _, branch := range res.LocalBaseBranches {
		buffer.WriteString(fmt.Sprintf("%s\n", branch.BranchName()))
	}
	if len(res.LocalBaseBranches) > 0 {
		buffer.WriteString("\n")
	}
	if len(res.LocalModBranches) > 0 {
		buffer.WriteString("Local Mod Branches:\n")
		buffer.WriteString("\n")
	}
	for _, branch := range res.LocalModBranches {
		buffer.WriteString(fmt.Sprintf("%s\n", branch.BranchName()))
	}
	if len(res.LocalModBranches) > 0 {
		buffer.WriteString("\n")
	}
	return buffer.String()
}
