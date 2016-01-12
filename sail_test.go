package sail

import "testing"

func TestParse(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"Hello", "Hello"},
		{"Thank", "Than"},
	}

	for _, c := range cases {
		got := Play(c.in)
		if got == c.want {
			t.Logf("Success: %s", c.in)
		} else {
			t.Errorf("Failure: in:%s got:%s want:%s", c.in, got, c.want)
		}
	}
}
