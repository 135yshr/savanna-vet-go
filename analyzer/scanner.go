package analyzer

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// Scanner はテストファイルを走査してスメルを検出する。
type Scanner struct {
	Detectors      []Detector
	EnabledSmells  map[SmellType]bool
	DisabledSmells map[SmellType]bool
}

// NewScanner はデフォルト設定のScannerを生成する。
func NewScanner() *Scanner {
	return &Scanner{
		Detectors: AllDetectors(),
	}
}

// ScanDir はディレクトリ配下の _test.go ファイルを再帰的にスキャンする。
func (s *Scanner) ScanDir(dir string) ([]TestSmell, error) {
	var allSmells []TestSmell
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			base := filepath.Base(path)
			if base == "vendor" || strings.HasPrefix(base, ".") {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, "_test.go") {
			return nil
		}
		smells, err := s.ScanFile(path)
		if err != nil {
			return nil // パースエラーは無視して続行
		}
		allSmells = append(allSmells, smells...)
		return nil
	})
	return allSmells, err
}

// ScanFile は単一のテストファイルをスキャンする。
func (s *Scanner) ScanFile(path string) ([]TestSmell, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	return s.scanAST(fset, file), nil
}

// ScanSource はソースコード文字列をスキャンする（テスト用）。
func (s *Scanner) ScanSource(src string) ([]TestSmell, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test_test.go", src, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	return s.scanAST(fset, file), nil
}

func (s *Scanner) scanAST(fset *token.FileSet, file *ast.File) []TestSmell {
	var smells []TestSmell
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		for _, detector := range s.Detectors {
			if !s.isEnabled(detector.Type()) {
				continue
			}
			smells = append(smells, detector.Detect(fset, file, fn)...)
		}
	}
	return smells
}

func (s *Scanner) isEnabled(t SmellType) bool {
	if len(s.EnabledSmells) > 0 {
		return s.EnabledSmells[t]
	}
	if len(s.DisabledSmells) > 0 {
		return !s.DisabledSmells[t]
	}
	return true
}
