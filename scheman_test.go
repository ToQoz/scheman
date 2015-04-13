package scheman

import (
	"testing"
)

func TestSplitStmt(t *testing.T) {
	stmts := splitStmt([]byte("A; B; C;"))
	stmts2 := splitStmt([]byte(`SELECT *
FROM A;

SELECT * FROM B;
`))

	tests := []struct {
		Got      interface{}
		Expected interface{}
	}{
		{Got: len(stmts), Expected: 3},
		{Got: stmts[0], Expected: "A;"},
		{Got: stmts[1], Expected: "B;"},
		{Got: stmts[2], Expected: "C;"},
		{Got: len(stmts2), Expected: 2},
		{Got: stmts2[0], Expected: `SELECT *
FROM A;`},
		{Got: stmts2[1], Expected: "SELECT * FROM B;"},
	}

	for _, test := range tests {
		if test.Got != test.Expected {
			t.Errorf("expected %v, but got %v", test.Expected, test.Got)
		}
	}
}
