package translator

import (
	"bufio"
	"log"
	"os"
	"testing"

	"github.com/xlab/cgogen/parser"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func TestLearn(t *testing.T) {
	pCfg := parser.NewConfig("/Users/xlab/Documents/dev/ctru/ctrulib/libctru/include/3ds.h")
	pCfg.SysIncludePaths = []string{"/Users/xlab/Documents/dev/ctru/ctrulib/libctru/include", "/usr/include"}
	p, err := parser.New(pCfg)
	if err != nil {
		t.Fatal(err)
	}
	unit, macros, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	buf := bufio.NewWriter(os.Stdout)
	defer buf.Flush()
	rules := Rules{
		TargetGlobal: {
			RuleSpec{From: "(?i)VPX_", Action: ActionAccept},
			RuleSpec{Transform: TransformLower},
		},
		TargetDefine: {
			RuleSpec{From: "_INLINE$", Action: ActionIgnore},
			RuleSpec{From: "vpx_", To: "_", Action: ActionReplace},
			RuleSpec{From: "_abi", Transform: TransformUpper},
			RuleSpec{From: "_img", To: "_image", Action: ActionReplace},
			RuleSpec{From: "_fmt", To: "_format", Action: ActionReplace},
			RuleSpec{From: "_([^_]+)", To: "$1", Action: ActionReplace, Transform: TransformTitle},
		},
	}
	constRules := ConstRules{
		ConstEnum:        ConstEvalFull,
		ConstDeclaration: ConstExpand,
	}
	tl, err := New(rules, constRules, nil, buf)
	if err != nil {
		t.Fatal(err)
	}
	if err := tl.Learn(unit, macros); err != nil {
		t.Fatal(err)
	}
	tl.Report()
}
