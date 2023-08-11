package main

import "github.com/go-git/go-git/v5/plumbing/format/diff"

type UnifiedChunk struct {
	Content string
	Type    string
}

func ToUnifiedChunk(chunks []diff.Chunk) []UnifiedChunk {
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

func main() {

}
