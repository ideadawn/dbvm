package mysql

import (
	"fmt"

	"github.com/ideadawn/dbvm/manager"
)

// print sql blocks
func (p *sqlParser) print() {
	noTrans := string(manager.MagicNoTrans)
	ignore := string(manager.MagicIgnore)

	for _, blk := range p.blocks {
		for _, cmmt := range blk.comments {
			fmt.Println(string(cmmt))
		}
		if blk.noTrans {
			fmt.Println(string(myCnf.commentBegin), noTrans)
		}
		if len(blk.ignores) > 0 {
			for idx, val := range blk.ignores {
				if idx == 0 {
					fmt.Printf("%s %s %d", myCnf.commentBegin, ignore, val)
				} else {
					fmt.Printf(",%d", val)
				}
			}
			fmt.Printf(string(myCnf.newLine))
		}
		if blk.inBlock {
			fmt.Println(string(myCnf.blockBegin))
		}
		if len(blk.delimiter) > 0 {
			fmt.Println(string(myCnf.delimiter), string(blk.delimiter))
		}
		for _, itm := range blk.items {
			for _, cmmt := range itm.comments {
				fmt.Println(string(cmmt))
			}
			for _, sql := range itm.sqlArr {
				fmt.Println(string(sql))
			}
		}
		if len(blk.delimiter) > 0 {
			fmt.Println(string(blk.delimiter))
			fmt.Println(string(myCnf.delimiter), string(myCnf.defaultEnd))
		}
		if blk.inBlock {
			fmt.Println(string(myCnf.blockCommit))
		}
	}
	fmt.Println("")
}
