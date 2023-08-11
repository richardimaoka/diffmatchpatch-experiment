package main

import (
	"fmt"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
)

var storage *memory.Storage
var repoCache map[string]*git.Repository

func gitOpenOrClone(repoUrl string) (*git.Repository, error) {
	if repo, ok := repoCache[repoUrl]; ok {
		return repo, nil
	}

	repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{URL: repoUrl})
	if err != nil {
		return nil, fmt.Errorf("cannot clone repo %s, %s", repoUrl, err)
	}

	repoCache[repoUrl] = repo
	return repo, nil
}

func gitCommit(repoUrl, commitHashStr string) (*object.Commit, error) {
	repo, err := gitOpenOrClone(repoUrl)
	if err != nil {
		return nil, fmt.Errorf("failed in gitCommit for, %s", err)
	}

	commitHash := plumbing.NewHash(commitHashStr)
	if commitHash.String() != commitHashStr {
		return nil, fmt.Errorf("failed in gitCommit, commit hash = %s is invalid as its re-calculated hash is mismatched = %s", commitHashStr, commitHash.String())
	}

	commit, err := repo.CommitObject(commitHash)
	if err != nil {
		return nil, fmt.Errorf("failed in gitCommit, cannot get commit = %s, %s", commitHashStr, err)
	}

	return commit, nil
}

func TestFileHighlight(t *testing.T) {
	repoCache = make(map[string]*git.Repository)

	cases := []struct {
		repoUrl       string
		prevCommit    string
		currentCommit string
		filePath      string
	}{
		{
			"https://github.com/richardimaoka/file-highlight-test.git",
			"f1df093152800852cf892d4015ff56a56427716b",
			"cf3bc8ae215607bd18d50c72a48868bc4f2b5e49",
			"1.txt",
		},
	}

	for _, c := range cases {
		prevCommit, err := gitCommit(c.repoUrl, c.prevCommit)
		if err != nil {
			t.Fatalf("failed in TestFileHighlight to get prev commit, %s", err)
		}

		currentCommit, err := gitCommit(c.repoUrl, c.currentCommit)
		if err != nil {
			t.Fatalf("failed in TestFileHighlight to get current commit, %s", err)
		}

		patch, _ := prevCommit.Patch(currentCommit)
		filePatches := patch.FilePatches()
		for _, p := range filePatches {
			_, to := p.Files()
			if to.Path() == c.filePath {
				chunks := p.Chunks()
				unifiedChunks := ToUnifiedChunk(chunks)
				fmt.Println("unifiedChunks:", unifiedChunks)
			}
		}
	}
}
