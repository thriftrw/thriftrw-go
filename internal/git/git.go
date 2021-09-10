package git

import (
	"fmt"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	// "go.uber.org/thriftrw/internal/compare"
)

func NewGitFS(repo *git.Repository, gitDir string, from bool) *gitFS {
	return &gitFS{
		repo: repo,
		gitDir: gitDir,
		from: from,
	}
}

type gitFS struct{
	repo *git.Repository
	gitDir string
	from bool // Whether we are looking for previous version.
}

func (fs gitFS) Read(path string) ([]byte, error) {
	// findChangedThrift(fs.gitDir)
	// thrift/v1.thrift
	r, err := git.PlainOpenWithOptions(fs.gitDir, &git.PlainOpenOptions{
		DetectDotGit:          true,
		EnableDotGitCommonDir: true,
	})
	if err != nil {
		return nil, err
	}
	refHead, err := r.Head()
	if err != nil {
		return nil, err
	}
	var commit *object.Commit
	if !fs.from {
		commit, err = fs.repo.CommitObject(refHead.Hash())
		if err != nil {
			return nil, err
		}
	} else {
		commitIter, err := r.Log(&git.LogOptions{From: refHead.Hash()})
		if err != nil {
			return nil, err
		}
		_, err = commitIter.Next()
		if err != nil {
			return nil, err
		}
		commit, _ = commitIter.Next()
		if err != nil {
			return nil, err
		}
	}

	// filename is going to be the full path. We don't want that.
	filename, err := filepath.Rel(fs.gitDir, path)
	if err != nil {
		return nil, err
	}
	f, err := commit.File(filename)
	if err != nil {
		return nil, err
	}
	s, err := f.Contents()
	if err != nil {
		return nil, err
	}
	body := []byte(s)

	return body, nil

	//
	// commitIter, err := r.Log(&git.LogOptions{From: refHead.Hash()})
	// commit, err := commitIter.Next()
	// parentCommit, _ := commit.Parent(0)
	// fmt.Println(parentCommit.Hash)
	// fmt.Println(commit.Hash)
	//
	// c, _ := commit.Tree()
	// pc, _ := parentCommit.Tree()





	return nil, nil
	// return ioutil.ReadFile(filename)
}

func (fs gitFS) Abs(p string) (string, error) {

	return filepath.Join(fs.gitDir, p), nil
}



func findChangedThrift(gitDir string) ([]string, error) {
	r, err := git.PlainOpenWithOptions(gitDir, &git.PlainOpenOptions{
		DetectDotGit:          true,
		EnableDotGitCommonDir: true,
	})
	if err != nil {
		return nil, fmt.Errorf("could not open repo: %v", err)
	}
	// Get Repo's HEAD
	refHead, err := r.Head()
	// Look at the log.
	commitIter, err := r.Log(&git.LogOptions{From: refHead.Hash()})
	commit, err := commitIter.Next()
	parentCommit, _ := commit.Parent(0)
	// Get the two commmit trees.
	c, _ := commit.Tree()
	pc, _ := parentCommit.Tree()
	// Diff the trees and find what changed.
	objects, _ := object.DiffTree(c, pc)
	changed := []string{}
	for _, o := range objects {
		a, _ := o.Action()
		if a.String() == "Modify" {
			to, _, _:= o.Files()
			if filepath.Ext(to.Name) == ".thrift" {
				changed = append(changed, to.Name)
				fmt.Printf("changed Thrift file: %s\n", to.Name)
				// TODO: compiler needs a location of the file, not its content.
				// toFile, err := to.Contents()
				if err != nil {
					return changed, err
				}
				// err = compare.CompileFiles(toFile, toFile)
			}
		} else if a.String() == "Delete" {
			// TODO: deal with deletes
		}
	}

	return changed, nil

	// fmt.Println(objects)

	// for _, f := range objects {
	//
	// 	// commit, err := repo.CommitObject(baseHash)
	// 	// if err != nil { ... }
	// 	//
	// 	// // Load and read the file.
	// 	// f, err := commit.File("path/to/file.thrift")
	// 	// if err != nil { ... }
	// 	//
	// 	// body, err := f.Contents()
	//
	//
	// 	out, _ := f.From.Tree.File(f.From.Name)
	// 	fmt.Println(out)
	// 	err := compare.CompileFiles(f.From.Name, f.To.Name)
	// 	fmt.Println(err)
	// }

	// for _, e := range treeHead.Entries {
	// 	if filepath.Ext(e.Name) == ".thrift" {
	// 		fmt.Println(e.Name)
	// 	}
	// 	if filepath.Ext(e.Name) == ".go" {
	// 		fmt.Println(e.Name)
	// 	}
	// }
	//
	// objs, _ := r.TreeObjects()
	// treeHead, _ := objs.Next()
	// treePrev, _ := objs.Next()

	// commitPrev, err := r.CommitObject(refPrev.Hash())
	// treePrev, _ := commitPrev.Tree()




	// fmt.Println(treeHead.Hash)
	// fmt.Println(treePrev.Hash)

	// commitHead, err := r.CommitObject(refHead.Hash())
	// treeHead, err := commitHead.Tree()
	// objects, _ := object.DiffTree(treeHead, treePrev)

	// fmt.Println(objects)


	// tree.Files().ForEach(func(f *object.File) error {
	// 	fmt.Println(f.Name)
	// 	return nil
	// })

	// return nil
}
