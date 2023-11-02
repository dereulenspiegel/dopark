package geoweb

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeowebSimple(t *testing.T) {
	s, err := NewScraper(slog.Default())
	require.NoError(t, err)
	spaces, err := s.Scrape()
	require.NoError(t, err)
	assert.NotEmpty(t, spaces)
	assert.Len(t, spaces, 16)
}
