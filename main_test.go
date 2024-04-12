package main

import (
	"fmt"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// var storage *memory.Storage
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

func gitFileFromCommit(repoUrl, commitHashStr, filePath string) (*object.File, error) {
	commit, err := gitCommit(repoUrl, commitHashStr)
	if err != nil {
		return nil, fmt.Errorf("failed in gitFileFromCommit, cannot get commit = %s, %s", commitHashStr, err)
	}

	rootTree, err := commit.Tree()
	if err != nil {
		return nil, fmt.Errorf("failed in gitFileFromCommit, cannot get tree for commit = %s, %s", commitHashStr, err)

	}

	gitFile, err := rootTree.File(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed in gitFileFromCommit, cannot get file = %s in commit = %s, %s", filePath, commitHashStr, err)
	}

	return gitFile, nil
}

func gitFileContents(repoUrl, commitHashStr, filePath string) (string, error) {
	file, err := gitFileFromCommit(repoUrl, commitHashStr, filePath)
	if err != nil {
		return "", fmt.Errorf("failed in gitFileContents, cannot get file = %s in commit = %s, %s", filePath, commitHashStr, err)
	}

	contents, err := file.Contents()
	if err != nil {
		return "", fmt.Errorf("failed in gitFileContents, cannot get contents for file = %s in commit = %s, %s", filePath, commitHashStr, err)
	}

	return contents, nil
}

func TestChunk(t *testing.T) {
	repoCache = make(map[string]*git.Repository)

	cases := []struct {
		repoUrl       string
		prevCommit    string
		currentCommit string
		filePath      string
		goldenFile1   string
		goldenFile2   string
	}{
		{
			"https://github.com/richardimaoka/file-highlight-test.git",
			"f1df093152800852cf892d4015ff56a56427716b",
			"cf3bc8ae215607bd18d50c72a48868bc4f2b5e49",
			"1.txt",
			"testdata/golden1-1.json",
			"testdata/golden1-2.json",
		},
		{
			"https://github.com/richardimaoka/sign-in-with-google-experiment.git",
			"6c88860799e173a271dd916791f4ec38c6c20abd",
			"4a2ec0ce7ec9fd4a8bafda822ced46b995824570",
			"package.json",
			"testdata/golden2-1.json",
			"testdata/golden2-2.json",
		},
		{
			"https://github.com/richardimaoka/sign-in-with-google-experiment.git",
			"a60b11b2d038bb47daf305ef89d2d19f1e1cbc90",
			"5d0cc4b461623706507b527366beece48fe5e68c",
			"index.html",
			"testdata/golden3-1.json",
			"testdata/golden3-2.json",
		},
	}

	for _, c := range cases {
		prevCommit, err := gitCommit(c.repoUrl, c.prevCommit)
		if err != nil {
			t.Fatalf("failed in TestChunk to get prev commit, %s", err)
		}

		currentCommit, err := gitCommit(c.repoUrl, c.currentCommit)
		if err != nil {
			t.Fatalf("failed in TestChunk to get current commit, %s", err)
		}

		// 1. testing git go-git diff
		patch, _ := prevCommit.Patch(currentCommit)
		filePatches := patch.FilePatches()
		for _, p := range filePatches {
			_, to := p.Files()
			if to.Path() == c.filePath {
				chunks := p.Chunks()
				unifiedChunks := ChunksToUnifiedChunks(chunks)
				CompareWitGoldenFile(t, *updateFlag, c.goldenFile1, unifiedChunks)

			}
		}

		// 2. testing diff-match-patch
		prevContents, err := gitFileContents(c.repoUrl, c.prevCommit, c.filePath)
		if err != nil {
			t.Fatalf("failed in TestChunk to get prev file contents, %s", err)
		}
		currentContents, err := gitFileContents(c.repoUrl, c.currentCommit, c.filePath)
		if err != nil {
			t.Fatalf("failed in TestChunk to get current file contents, %s", err)
		}

		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(prevContents, currentContents, true)
		unifiedChunks := DiffsToUnifiedChunks(diffs)
		CompareWitGoldenFile(t, *updateFlag, c.goldenFile2, unifiedChunks)
	}
}
