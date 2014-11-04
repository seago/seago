package helper

import (
	"github.com/huichen/sego"
)

func Segmenter(segmenter sego.Segmenter, key string) []string {
	segments := segmenter.Segment([]byte(key))
	keys := sego.SegmentsToSlice(segments, true)
	return keys
}
