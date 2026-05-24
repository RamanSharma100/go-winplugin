package builder

import (
	"io"
	"os"
	"path/filepath"
)

func CopyFile(
	src string,
	dst string,
) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}

	defer source.Close()

	target, err := os.Create(dst)
	if err != nil {
		return err
	}

	defer target.Close()

	_, err = io.Copy(
		target,
		source,
	)

	return err
}

func CopyDirectory(
	src string,
	dst string,
) error {
	return filepath.Walk(
		src,

		func(
			path string,
			info os.FileInfo,
			err error,
		) error {

			if err != nil {
				return err
			}

			relative, err := filepath.Rel(
				src,
				path,
			)

			if err != nil {
				return err
			}

			target := filepath.Join(
				dst,
				relative,
			)

			if info.IsDir() {
				return os.MkdirAll(
					target,
					0755,
				)
			}

			return CopyFile(
				path,
				target,
			)
		},
	)
}
