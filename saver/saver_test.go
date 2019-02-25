package saver

import (
	"github.com/stretchr/testify/assert"
	"image"
	"image/png"
	"os"
	"testing"
)

type buildDestinationPathFixtureData struct {
	Dir      string
	Filename string
}
type buildDestinationPathFixture struct {
	Data   buildDestinationPathFixtureData
	Result string
}

func (test *buildDestinationPathFixture) Do(t *testing.T) {
	assert.Equal(t, test.Result, buildDestinationPath(test.Data.Dir, test.Data.Filename))
}

func TestBuildDestinationPath(t *testing.T) {
	for name, test := range map[string]buildDestinationPathFixture{
		"simple dir with / & simple filename": {
			buildDestinationPathFixtureData{
				Dir:      "test/",
				Filename: "test.test",
			},
			"test/test.test",
		},
		"long dir with / & simple filename": {
			buildDestinationPathFixtureData{
				Dir:      "a/b/c/test_1/test_2/",
				Filename: "test.test",
			},
			"a/b/c/test_1/test_2/test.test",
		},
		"no dir & simple filename": {
			buildDestinationPathFixtureData{
				Dir:      "",
				Filename: "test.test",
			},
			"test.test",
		},
		"no dir & no filename": {
			buildDestinationPathFixtureData{
				Dir:      "",
				Filename: "",
			},
			"undefined",
		},
		"no dir & filename /": {
			buildDestinationPathFixtureData{
				Dir:      "",
				Filename: "/",
			},
			"undefined",
		},
		"no dir & filename with /": {
			buildDestinationPathFixtureData{
				Dir:      "",
				Filename: "a/a.jpg",
			},
			"aa.jpg",
		},
		"long dir with / & filename with /": {
			buildDestinationPathFixtureData{
				Dir:      "b/b/b/b",
				Filename: "a/a.jpg",
			},
			"b/b/b/b/aa.jpg",
		},
		"long dir with / and /-suffix & filename with /": {
			buildDestinationPathFixtureData{
				Dir:      "b/b/b/b/",
				Filename: "a/a.jpg",
			},
			"b/b/b/b/aa.jpg",
		},
		"long dir with / and /-prefix & filename with /": {
			buildDestinationPathFixtureData{
				Dir:      "/b/b/b/b",
				Filename: "a/a.jpg",
			},
			"/b/b/b/b/aa.jpg",
		},
	} {
		t.Run(name, test.Do)
	}
}

type createThumbnailFixtureData struct {
	ImageFileName   string
	PreviewFileName string
	Width           int
	Height          int
	CreateTestImage bool
}

type createThumbnailFixture struct {
	Data   createThumbnailFixtureData
	Result error
}

func (test *createThumbnailFixture) Do(t *testing.T) {
	if test.Data.CreateTestImage {
		img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{10, 10}})
		f, _ := os.Create(test.Data.ImageFileName)
		_ = png.Encode(f, img)
	}
	defer func() { _ = os.Remove(test.Data.ImageFileName) }()
	defer func() { _ = os.Remove(test.Data.PreviewFileName) }()
	assert.IsType(t, test.Result, createThumbnail(test.Data.ImageFileName, test.Data.PreviewFileName, test.Data.Width, test.Data.Height))
}

func TestCreateThumbnail(t *testing.T) {
	for name, test := range map[string]createThumbnailFixture{
		"100x100_true": {
			createThumbnailFixtureData{
				ImageFileName:   "1.png",
				PreviewFileName: "min_1.png",
				Width:           100,
				Height:          100,
				CreateTestImage: true,
			},
			nil,
		},
		"100x100_false": {
			createThumbnailFixtureData{
				ImageFileName:   "1.png",
				PreviewFileName: "min_1.png",
				Width:           100,
				Height:          100,
				CreateTestImage: false,
			},
			new(os.PathError),
		},
		"0x100_false": {
			createThumbnailFixtureData{
				ImageFileName:   "1.png",
				PreviewFileName: "min_1.png",
				Width:           0,
				Height:          100,
				CreateTestImage: true,
			},
			new(sizeError),
		},
		"100x0_false": {
			createThumbnailFixtureData{
				ImageFileName:   "1.png",
				PreviewFileName: "min_1.png",
				Width:           100,
				Height:          0,
				CreateTestImage: true,
			},
			new(sizeError),
		},
		"0x0_false": {
			createThumbnailFixtureData{
				ImageFileName:   "1.png",
				PreviewFileName: "min_1.png",
				Width:           0,
				Height:          0,
				CreateTestImage: true,
			},
			new(sizeError),
		},
		"badname_true": {
			createThumbnailFixtureData{
				ImageFileName:   "111",
				PreviewFileName: "min_1.png",
				Width:           100,
				Height:          100,
				CreateTestImage: true,
			},
			nil,
		},
	} {
		t.Run(name, test.Do)
	}
}

type SaveImageFromBase64FixtureData struct {
	Image     string
	Directory string
}

type SaveImageFromBase64Fixture struct {
	Data   SaveImageFromBase64FixtureData
	Result error
}

func (test *SaveImageFromBase64Fixture) Do(t *testing.T) {
	if test.Data.Directory != "" {
		_ = os.Mkdir(test.Data.Directory, 0777)
	}
	defer func() { _ = os.RemoveAll(test.Data.Directory) }()
	assert.IsType(t, test.Result, SaveImageFromBase64(test.Data.Image, test.Data.Directory))
}

func TestSaveImageFromBase64(t *testing.T) {
	for name, test := range map[string]SaveImageFromBase64Fixture{
		"empty": {
			SaveImageFromBase64FixtureData{
				Image:     "",
				Directory: "",
			},
			new(formatError),
		},
		"good": {
			SaveImageFromBase64FixtureData{
				Image:     "data:image/png;base64,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",
				Directory: "files",
			},
			nil,
		},

		"bad": {
			SaveImageFromBase64FixtureData{
				Image:     "data:image/png;base64,iVBvjL/SAniTUzynI4xQ9yrk10yB6urH8WUp3C8Ee6yiOgpyq1NgWcm8c9nb5wXcCXAqU3LiBHyS2WkBPDQokuBN9dyOMB4kusUw+1Sm+4GQfUoB80/NgqTTWX7t8gozrbjzXm+cKoplpzUpzb1MyeJRXdlQKgvSIV67zM9oINi2qrXpnfuRia/jw6bQLJdMBKmUu5ctjxqnCwfiYqkOtfJfyXgG1nzyHEw8voRqnmV8+RLIL3zO3ZbWhxfNGGyncZo4jV0WiBVuZJaLEnAa1FCKdsKIK712wSTYsyIFitalQ9H6GbD7bQVULprs/o12TlBkkls4Kw7hVH+BvajETvdlnpo0hv2PJI4QvbAP9YoJoE+bQv3tTsFAf+FI5RMr3Ag7LaEVY8GpeHvrsgJj6sZKECr1EKsywNRNNbA2VwnHnl8IjclJ1PNCgf22sTODWOBGxpUjB1lfetK/oEzyB161+dK88r6xSv3iWKWhiQ3izfjCc5Cz6jM4uOA1HLGJ/n0EAGdBHinPnVyfCkax/1iwrxt+07r6idam3NFgmhO32aSQjiBN3Hqy9Fp4jNKwKUZlaPUNxCIuughL15TBrldyQZ9Kn4I4+i4cjSzSJXe6tU0wnCmOyGaDOEsF9eyC56jit8FFXj/SxE9xSn10jFJ/lBbGQ5zyCry4/jBseuMDKFu0HD5/Jeevk7NmvnN9UnwE6aZ7gjxeXzQ906UWHHIVR9jZAG5KkLL2VdxRwf8TwzzrrOw/iJxGVaxKn4wYfEPQGKfSczsKaH6px7xzrwAcR8LbmAp+q/cWxIOjKuwHclpARZdd7DdOZViO042hBYhR6W+x9xdpzp1ulkSOwxCVzH7+Xj/jGr6N0YQ5nBq+Af/WIRvIaSE1ZusvAhyd19H4l/ipwwfDtzGqK6tiiwyPkCY9I6e2/xNMVdhZ3KYr7Vr+hBadgE9+uqd7+n+Jx/sbgLnTlrubkbkAAAAASUVORK5CYII=",
				Directory: "files",
			},
			*new(png.FormatError),
		},
	} {
		t.Run(name, test.Do)
	}
}

// TODO: tests for SaveImageFromPart

// TODO: tests for SaveImageFromUrl
