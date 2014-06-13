/*
	Copyright 2014 Google Inc. All rights reserved.

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package video

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"bitbucket.org/zombiezen/cardcpx/natsort"
)

type Source interface {
	List() ([]*Clip, error)
	Open(path string) (io.ReadCloser, error)
}

type Clip struct {
	Name      string   `json:"name"`
	Paths     []string `json:"paths"`
	TotalSize int64    `json:"totalSize"`
}

type dirStructSource struct {
	root    string
	f       func(path string) (clip string, ok bool)
	descend func(path string) bool
}

func (src *dirStructSource) Open(path string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(src.root, path))
}

func (src *dirStructSource) List() ([]*Clip, error) {
	m, err := src.walk()
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	natsort.Strings(keys)
	clips := make([]*Clip, len(keys))
	for i, k := range keys {
		clips[i] = m[k]
	}
	return clips, err
}

func (src *dirStructSource) walk() (map[string]*Clip, error) {
	m := make(map[string]*Clip)
	err := filepath.Walk(src.root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// TODO(light): is there anything we can handle?
			return err
		}
		normPath := strings.TrimLeft(strings.TrimPrefix(path, src.root), string(filepath.Separator))
		if normPath != "" && info.IsDir() && !src.descend(normPath) {
			return filepath.SkipDir
		} else if !info.Mode().IsRegular() {
			return nil
		}
		clipName, ok := src.f(normPath)
		if ok {
			clip := m[clipName]
			if clip == nil {
				clip = &Clip{Name: clipName}
				m[clipName] = clip
			}
			clip.Paths = append(clip.Paths, normPath)
			clip.TotalSize += info.Size()
		}
		return nil
	})
	return m, err
}

// DirectorySource returns a Source from the given directory, automatically inferring clip layout.
func DirectorySource(root string) (Source, error) {
	src := &dirStructSource{
		root:    root,
		f:       flatFileLayout,
		descend: func(string) bool { return false },
	}
	dir, err := os.Open(root)
	if err != nil {
		return nil, err
	}
	names, err := dir.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	for _, name := range names {
		if isRDMName(name) {
			src.f, src.descend = redLayout, redDescend
			break
		}
	}
	return src, nil
}

func flatFileLayout(path string) (clip string, ok bool) {
	if strings.ContainsRune(path, filepath.Separator) || strings.HasPrefix(path, ".") {
		return "", false
	}
	return path, true
}

func redLayout(path string) (clip string, ok bool) {
	parts := strings.Split(path, string(filepath.Separator))
	if len(parts) != 3 {
		return "", false
	}
	if !isRDMName(parts[0]) || !isRDCName(parts[1]) {
		return "", false
	}
	return filepath.Join(parts[0], parts[1]), true
}

func redDescend(path string) bool {
	parts := strings.Split(path, string(filepath.Separator))
	if len(parts) > 2 {
		return false
	}
	if len(parts) >= 1 && !isRDMName(parts[0]) {
		return false
	}
	if len(parts) >= 2 && !isRDCName(parts[1]) {
		return false
	}
	return true
}

func isRDMName(name string) bool {
	return strings.HasSuffix(name, ".RDM") || strings.HasSuffix(name, ".rdm")
}

func isRDCName(name string) bool {
	return strings.HasSuffix(name, ".RDC") || strings.HasSuffix(name, ".rdc")
}
