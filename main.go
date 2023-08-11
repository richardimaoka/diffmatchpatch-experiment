package main

import (
	"fmt"

	"github.com/go-git/go-git/v5/plumbing/format/diff"
	"github.com/sergi/go-diff/diffmatchpatch"
)

type UnifiedChunk struct {
	Content string
	Type    string
}

func ChunksToUnifiedChunks(chunks []diff.Chunk) []UnifiedChunk {
	convert := func(chunk diff.Chunk) UnifiedChunk {
		var t string
		switch chunk.Type() {
		case diff.Add:
			t = "ADD"
		case diff.Delete:
			t = "DELETE"
		case diff.Equal:
			t = "EQUAL"
		}

		return UnifiedChunk{
			Content: chunk.Content(),
			Type:    t,
		}
	}

	var unifiedChunks []UnifiedChunk
	for _, c := range chunks {
		unifiedChunks = append(unifiedChunks, convert(c))
	}

	return unifiedChunks
}

func DiffsToUnifiedChunks(diffs []diffmatchpatch.Diff) []UnifiedChunk {
	convert := func(diff diffmatchpatch.Diff) UnifiedChunk {
		var t string
		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			t = "ADD"
		case diffmatchpatch.DiffDelete:
			t = "DELETE"
		case diffmatchpatch.DiffEqual:
			t = "EQUAL"
		}

		return UnifiedChunk{
			Content: diff.Text,
			Type:    t,
		}
	}

	var unifiedChunks []UnifiedChunk
	for _, d := range diffs {
		unifiedChunks = append(unifiedChunks, convert(d))
	}

	return unifiedChunks
}

func main() {
	fmt.Println("\\")
}
