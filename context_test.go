package sprite_sass

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func fileString(path string) string {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func fileReader(path string) io.Reader {
	reader, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	return reader
}

func init() {
	// Setup build directory
	os.MkdirAll("test/build", 0755)
}

func TestContextFile(t *testing.T) {

	ipath := "test/sass/file.scss"
	exp, err := ioutil.ReadFile("test/expected/file.css")
	if err != nil {
		panic(err)
	}

	ctx, scanned, _ := setupCtx(ipath)

	defer cleanUpSprites(ctx.Parser.Sprites)

	if string(scanned) != string(exp) {
		t.Errorf("Expected did not match returned exp:\n%s\n"+
			"was:\n%s\n", string(exp), string(scanned))
	}
}

func TestContextRun(t *testing.T) {

	ipath := "test/sass/_var.scss"
	exp, err := ioutil.ReadFile("test/expected/var.css")
	if err != nil {
		t.Error(err)
	}

	ctx, scanned, _ := setupCtx(ipath)

	defer cleanUpSprites(ctx.Parser.Sprites)

	if string(scanned) != string(exp) {
		t.Errorf("Processor file did not match exp: "+
			"\n%s\n was:\n%s", string(exp), string(scanned))
	}

}

func TestContextImport(t *testing.T) {

	ipath := "test/sass/import.scss"
	exp, err := ioutil.ReadFile("test/expected/import.css")
	if err != nil {
		t.Error(err)
	}

	ctx, scanned, _ := setupCtx(ipath)

	defer cleanUpSprites(ctx.Parser.Sprites)

	res := string(scanned)
	if e := string(exp); res != e {
		t.Errorf("Processor file did not match \nexp: "+
			"\n%s\n was:\n%s", e, res)
	}

}

func TestContextFail(t *testing.T) {
	return
	ctx := Context{
		OutputStyle:  NESTED_STYLE,
		IncludePaths: make([]string, 0),
		Out:          "",
	}

	var scanned []byte
	ipath := "test/_failimport.scss"

	r, w := io.Pipe()
	go func(ipath string, w io.WriteCloser) {

		err := ctx.Run(fileReader(ipath), w, "test")
		if err == nil {
			t.Error("Invalid SCSS was not found")
		}
		errs := strings.Split(err.Error(), "\n")
		libsassErr := errs[0]
		parsedErr := errs[1]

		e := "source string:9: error: invalid top-level expression"
		if e != libsassErr {
			t.Errorf("expected:\n%s\nwas:\n%s\n", e, libsassErr)
		}

		e = "error in fail:4"
		if e != parsedErr {
			t.Errorf("expected:\n%s\nwas:\n%s\n", e, parsedErr)
		}

	}(ipath, w)

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanBytes)

	for scanner.Scan() {
		scanned = append(scanned, scanner.Bytes()...)
	}
	defer cleanUpSprites(ctx.Parser.Sprites)

	scanned = rerandom.ReplaceAll(scanned, []byte(""))
	_ = scanned
}

func TestContextNilRun(t *testing.T) {
	ctx := Context{}
	var w io.WriteCloser
	err := ctx.Run(nil, w, "test")
	if err == nil {
		t.Errorf("No error returned: %s", err)
	}

}

func BenchmarkContextCompile(b *testing.B) {
	ctx := Context{
		OutputStyle:  NESTED_STYLE,
		IncludePaths: make([]string, 0),
		Src:          fileString("test/file1.scss"),
		Out:          "",
		ImageDir:     "/tmp",
	}
	for n := 0; n < b.N; n++ {
		ctx.Compile()
	}
}
