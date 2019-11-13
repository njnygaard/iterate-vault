package iterate

import "testing"

func TestHello(t *testing.T) {
    var want error
    if got := Hello(); got != want {
        t.Errorf("Hello() = %q, want %q", got, want)
    }
}
