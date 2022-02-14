package mysql

import (
	"regexp"
)

// MySQL配置
type mysqlConfig struct {
	newLine    []byte
	defaultEnd []byte
	commaGap   []byte
	delimiter  []byte
	empty      []byte
	space      []byte

	commentBegin []byte
	blockBegin   []byte
	blockCommit  []byte

	reAlter        *regexp.Regexp
	reAlterSub     *regexp.Regexp
	reAddColumn    *regexp.Regexp
	reAddPrimary   *regexp.Regexp
	reAddIndex     *regexp.Regexp
	reDropColumn   *regexp.Regexp
	reChangeColumn *regexp.Regexp
	reModifyColumn *regexp.Regexp

	reCreateTable    *regexp.Regexp
	reCreateTableINE *regexp.Regexp
	reDropTable      *regexp.Regexp
	reDropTableIE    *regexp.Regexp
}

// 初始化配置
var myCnf = &mysqlConfig{
	newLine:    []byte{'\n'},
	defaultEnd: []byte{';'},
	commaGap:   []byte{','},
	delimiter:  []byte(`DELIMITER`),
	empty:      []byte{},
	space:      []byte{' '},

	commentBegin: []byte(`--`),
	blockBegin:   []byte(`BEGIN`),
	blockCommit:  []byte(`COMMIT`),

	reAlter:        regexp.MustCompile("(?is)(ALTER[ \t\n]+TABLE.*?)((?:ADD|CHANGE|MODIFY|DROP).*)"),
	reAlterSub:     regexp.MustCompile("(?is)((?:ADD|CHANGE|MODIFY|DROP)[ \t\n]+(?:COLUMN|INDEX|KEY|PRIMARY|UNIQUE).*?(?:,[ \t\n]+|;|$))"),
	reAddColumn:    regexp.MustCompile("^(?is)[ \t\n]*ADD[ \t\n]+COLUMN"),
	reAddPrimary:   regexp.MustCompile("^(?is)[ \t\n]*ADD[ \t\n]+PRIMARY"),
	reAddIndex:     regexp.MustCompile("^(?is)[ \t\n]*ADD[ \t\n]+(?:UNIQUE[ \t\n]+)?(?:INDEX|KEY)"),
	reDropColumn:   regexp.MustCompile("^(?is)[ \t\n]*DROP[ \t\n]+(?:COLUMN|INDEX|KEY|PRIMARY)"),
	reChangeColumn: regexp.MustCompile("^(?is)[ \t\n]*CHANGE[ \t\n]+COLUMN"),
	reModifyColumn: regexp.MustCompile("^(?is)[ \t\n]*MODIFY[ \t\n]+COLUMN"),

	reCreateTable:    regexp.MustCompile("^(?is)[ \t\n]*CREATE[ \t\n]+TABLE"),
	reCreateTableINE: regexp.MustCompile("^(?is)[ \t\n]*CREATE[ \t\n]+TABLE[ \t\n]+IF[ \t\n]+NOT[ \t\n]+EXISTS"),
	reDropTable:      regexp.MustCompile("^(?is)[ \t\n]*DROP[ \t\n]+TABLE"),
	reDropTableIE:    regexp.MustCompile("^(?is)[ \t\n]*DROP[ \t\n]+TABLE[ \t\n]+IF[ \t\n]+EXISTS"),
}
