package new

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/markbates/pkger"
	"github.com/urfave/cli/v2"
	_ "github.com/x-punch/micro-gen/cmd/new/pkger"
)

type service struct {
	Name      string
	Namespace string
	Path      string
}

func Run(cli *cli.Context) (err error) {
	srv := &service{Name: cli.String("name"), Namespace: cli.String("namespace")}
	path := cli.String("path")
	if len(path) == 0 {
		path, err = os.Getwd()
	} else {
		path, err = filepath.Abs(path)
	}
	if err != nil {
		return err
	}
	srv.Path = filepath.Join(path, srv.Name)
	if _, err := os.Stat(srv.Path); err == nil {
		fmt.Printf("Path '%s' already exist\n", srv.Path)
		return nil
	}
	if err := srv.gernate(); err != nil {
		return err
	}
	fmt.Printf("Generate %s on %s\n", srv.Name, srv.Path)
	if err := srv.modInit(); err != nil {
		return err
	}
	if err := srv.protoc(); err != nil {
		return err
	}
	fmt.Println("Generate protobuf files")
	return nil
}

func (s *service) gernate() (err error) {
	pkger.Walk("/cmd/new/templates/", func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		file, err := pkger.Open(p)
		if err != nil {
			return err
		}
		name := filepath.Join(s.Path, strings.TrimPrefix(file.Path().Name, "/cmd/new/templates"))
		i := strings.LastIndex(name, string(os.PathSeparator))
		if i > 0 {
			dir := name[:i]
			if err = os.MkdirAll(dir, 0755); err != nil {
				return err
			}
		}
		tmpl, err := ioutil.ReadAll(file)
		if err != nil {
			return err
		}
		if strings.HasSuffix(name, ".tmpl") {
			name = strings.TrimSuffix(name, ".tmpl")
			tmpl, err = s.parseTemplate(string(tmpl))
			if err != nil {
				return err
			}
		}
		if err := ioutil.WriteFile(name, tmpl, 0755); err != nil {
			return err
		}
		return nil
	})
	return nil
}

func (s *service) modInit() error {
	cmd := exec.Command("go", "mod", "init", s.Namespace+"/"+s.Name)
	cmd.Dir = s.Path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (s *service) protoc() error {
	cmd := exec.Command("protoc", "--micro_out=proto", "--go_out=proto", "proto/greeter.proto")
	cmd.Dir = s.Path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (s *service) parseTemplate(tmpl string) ([]byte, error) {
	t, err := template.New("").Parse(tmpl)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err = t.Execute(&buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
