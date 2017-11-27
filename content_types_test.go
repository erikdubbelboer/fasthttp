package fasthttp

import (
	"strings"
	"testing"
)

func TestContentType(t *testing.T) {
	if !strings.EqualFold(
		ContentType("fast.wav"), "audio/x-wav",
	) {
		t.Fatal("error: fast.wav content type defers from \"audio/x-wav\"")
	}
	if !strings.EqualFold(
		ContentType("fast.avi"), "video/x-msvideo",
	) {
		t.Fatal("error: fast.avi content type defers from \"video/x-msvideo\"")
	}
	if !strings.EqualFold(
		ContentType("fast.epub"), "application/epub+zip",
	) {
		t.Fatal("error: fast.epub content type defers from \"application/epub+zip\"")
	}
}
