package mysql

import (
	"bytes"
	"fmt"

	"github.com/ideadawn/dbvm/manager"
)

// print sql blocks
func (p *sqlParser) print() {
	noTrans := string(manager.MagicNoTrans)
	ignore := string(manager.MagicIgnore)

	for _, block := range p.blocks {
		for _, cmmt := range block.comments {
			fmt.Println(string(cmmt))
		}
		if block.noTrans {
			fmt.Println(string(myCnf.commentBegin), noTrans)
		}
		if len(block.ignores) > 0 {
			for idx, val := range block.ignores {
				if idx == 0 {
					fmt.Printf("%s %s %d", myCnf.commentBegin, ignore, val)
				} else {
					fmt.Printf(",%d", val)
				}
			}
			fmt.Printf(string(myCnf.newLine))
		}
		delimiter := len(block.delimiter) > 0 && bytes.Compare(block.delimiter, myCnf.defaultEnd) != 0
		if delimiter {
			fmt.Println(string(myCnf.delimiter), string(block.delimiter))
		}
		for _, sql := range block.sqlArr {
			fmt.Println(string(sql))
		}
		if delimiter {
			fmt.Println(string(block.delimiter))
			fmt.Println(string(myCnf.delimiter), string(myCnf.defaultEnd))
		}
	}
	fmt.Println("")
}

// Print 打印解析后的脚本
func Print(path string) error {
	parser := &sqlParser{
		file: path,
	}
	parser.parseSqlBlocks()
	if parser.err != nil {
		return parser
	}

	parser.print()
	return nil
}
