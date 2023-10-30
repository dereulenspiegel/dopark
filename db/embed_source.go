package db

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"

	"github.com/golang-migrate/migrate/v4/source"
)

func init() {
	source.Register("embed", &EmbedSource{})
}

type EmbedSource struct {
	fs         embed.FS
	migrations *source.Migrations
	path       string
	log        *slog.Logger
}

func WithInstance(fs embed.FS, fsPath string) (source.Driver, error) {
	es := &EmbedSource{
		fs:         fs,
		migrations: source.NewMigrations(),
		path:       fsPath,
		log:        slog.Default(),
	}

	embedFiles, err := fs.ReadDir(fsPath)
	es.log.Debug("Read files from embedded dir", "fileCount", len(embedFiles), "dir", fsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded directory (%s): %s", fsPath, err)
	}
	for _, embeddedFile := range embedFiles {
		m, err := source.DefaultParse(embeddedFile.Name())
		if err != nil {
			es.log.Error("failed to parse source file", "err", err, "file", embeddedFile.Name())
			continue
		}
		es.migrations.Append(m)
	}
	return es, nil
}

func (e *EmbedSource) Open(url string) (source.Driver, error) {
	return nil, fmt.Errorf("not implemented")
}

func (e *EmbedSource) Close() error {
	return nil
}

func (e *EmbedSource) First() (version uint, err error) {
	if version, ok := e.migrations.First(); !ok {
		return 0, &os.PathError{Op: "first", Path: e.path, Err: os.ErrNotExist}
	} else {
		return version, nil
	}
}

func (e *EmbedSource) Prev(version uint) (prevVersion uint, err error) {
	if v, ok := e.migrations.Prev(version); !ok {
		return 0, &os.PathError{Op: fmt.Sprintf("prev for version %v", version), Path: e.path, Err: os.ErrNotExist}
	} else {
		return v, nil
	}
}

func (e *EmbedSource) Next(version uint) (nextVersion uint, err error) {
	if v, ok := e.migrations.Next(version); !ok {
		return 0, &os.PathError{Op: fmt.Sprintf("next for version %v", version), Path: e.path, Err: os.ErrNotExist}
	} else {
		return v, nil
	}
}

func (e *EmbedSource) ReadUp(version uint) (r io.ReadCloser, identifier string, err error) {
	if m, ok := e.migrations.Up(version); ok {
		body, err := e.read(m)
		if err != nil {
			return nil, "", err
		}
		return io.NopCloser(bytes.NewReader(body)), m.Identifier, nil
	}
	return nil, "", &os.PathError{Op: fmt.Sprintf("read version %v", version), Path: e.path, Err: os.ErrNotExist}
}

func (e *EmbedSource) ReadDown(version uint) (r io.ReadCloser, identifier string, err error) {
	if m, ok := e.migrations.Down(version); ok {
		body, err := e.read(m)
		if err != nil {
			return nil, "", err
		}
		return io.NopCloser(bytes.NewReader(body)), m.Identifier, nil
	}
	return nil, "", &os.PathError{Op: fmt.Sprintf("read version %v", version), Path: e.path, Err: os.ErrNotExist}
}

func (e *EmbedSource) read(m *source.Migration) ([]byte, error) {
	e.log.Debug("Reading migration", "path", m.Raw)
	return e.fs.ReadFile(path.Join(e.path, m.Raw))
}
