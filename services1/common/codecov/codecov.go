package codecov

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gocarina/gocsv"
	"golang.org/x/tools/cover"
)

const (
	testingFlag = "test.v"
)

type Coverage struct {
	Packages []Package
}

type priorityReport struct {
	Name        string `csv:"package_name"`
	NumUntested int    `csv:"untested_statements"`
}

func (c *Coverage) TestPriority(csvPath string) (err error) {
	type data struct {
		name                    string
		total, tested, untested int
	}
	var d []data
	for _, p := range c.Packages {
		cov := data{
			name: p.Name,
		}
		for _, f := range p.Functions {
			for _, s := range f.Statements {
				cov.total++
				if s.Reached > 0 {
					cov.tested++
				} else {
					cov.untested++
				}
			}
		}
		if cov.total != cov.tested {
			d = append(d, cov)
		}
	}
	sort.Slice(d, func(i, j int) bool {
		return d[i].untested > d[j].untested
	})

	var report []priorityReport
	for _, d := range d {
		report = append(report, priorityReport{Name: d.name, NumUntested: d.untested})
		if flag.Lookup(testingFlag) == nil { // Only print this to stdout when not testing
			fmt.Printf("%s, %d statement(s) of untested code\n", d.name, d.untested)
		}
	}

	var csvData []byte
	if csvData, err = gocsv.MarshalBytes(report); err != nil { // Use this to save the CSV back to the file
		return
	} else if err = ioutil.WriteFile(csvPath, csvData, os.ModePerm); err != nil {
		return
	}

	return err
}

type Package struct {
	Name      string
	Functions []Function
}

type Function struct {
	Name       string
	File       string
	Start, End int
	Statements []Statement
}

type Statement struct {
	Start, End, Reached int
}

// LoadFromFile loads a cover profile that has been converted to JSON format.
func LoadFromFile(name string) (*Coverage, error) {
	coverage := Coverage{}
	if b, err := ioutil.ReadFile(name); err != nil {
		return nil, err
	} else if err := json.Unmarshal(b, &coverage); err != nil {
		return nil, err
	} else {
		return &coverage, nil
	}
}

// CompleteCoverageFile adds information to a cover profile about untested packages
// to reflect they have no code coverage.
func CompleteCoverageFile(filename string) error {
	if profiles, err := cover.ParseProfiles(filename); err != nil {
		return err
	} else {
		// collect the files that are referenced in the coverage file
		files := map[string][]string{}
		for _, p := range profiles {
			if _, ok := files[p.FileName]; !ok {
				files[p.FileName] = strings.Split(p.FileName, `/`)
			}
		}

		// load all packages that are referenced in the coverage file
		packages := map[string]*build.Package{}
		for f := range files {
			d := path.Dir(f)
			if _, ok := packages[d]; ok {
				continue
			}
			if p, err := build.Import(d, ".", build.IgnoreVendor); err != nil {
				return err
			} else {
				packages[d] = p
			}
		}

		// find the longest common substring (this is our root package import path)
		var lcs []string
	outer:
		for _, o := range files {
			for _, i := range files {
				if i[len(lcs)] != o[len(lcs)] {
					break outer
				}
			}
			lcs = append(lcs, o[len(lcs)])
		}

		// find the root package
		root := path.Join(lcs...)
		if rp, err := build.Import(root, ".", build.FindOnly); err != nil {
			return err
		} else {
			// walk every directory under the root package and import it
			if err := filepath.Walk(rp.Dir, func(name string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() {
					dn := info.Name()
					if dn == "vendor" || dn == "generated" || strings.HasPrefix(dn, ".") {
						return filepath.SkipDir
					}
					if lp, err := filepath.Rel(rp.Dir, name); err != nil {
						return err
					} else {
						ip := path.Join(rp.ImportPath, path.Clean(filepath.ToSlash(lp)))
						if _, ok := packages[ip]; !ok {
							if p, err := build.Import(ip, ".", build.IgnoreVendor); err == nil {
								packages[ip] = p
							}
						}
					}
				}
				return nil
			}); err != nil {
				return err
			}
		}

		// collect the files that are in the packages but not in the coverage file
		var missingFiles []string
		for _, p := range packages {
			for _, f := range p.GoFiles {
				fn := path.Join(p.ImportPath, f)
				if _, ok := files[fn]; !ok {
					missingFiles = append(missingFiles, fn)
				}
			}
		}

		if f, err := os.OpenFile(filename, os.O_APPEND, os.ModePerm); err != nil {
			return err
		} else {
			w := bufio.NewWriter(f)
			for _, mf := range missingFiles {
				if _, err := w.WriteString(fmt.Sprintf("%s:0.0,0.0 0 0\n", mf)); err != nil {
					return err
				}
			}
			w.Flush()
			f.Close()
		}

		return nil
	}
}
