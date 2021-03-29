package herbcache

import "testing"

func TestNop(t *testing.T) {
	c := New()

	if c.Storage().Engine != DefaultEngine {
		t.Fatal(c)
	}
	err := c.Storage().Start()
	if err != ErrStorageRequired {
		t.Fatal(err)
	}
	d, err := c.Get([]byte("testkey"))
	if len(d) != 0 || err != ErrStorageRequired {
		t.Fatal(d, err)
	}
	err = c.SetWithTTL([]byte("testkey"), []byte("testdata"), 3600)
	if err != ErrStorageRequired {
		t.Fatal(err)
	}
	err = c.Delete([]byte("testkey"))
	if err != ErrStorageRequired {
		t.Fatal(err)
	}
	err = c.Flush()
	if err != ErrStorageRequired {
		t.Fatal(err)
	}
	err = c.Storage().Stop()
	if err != ErrStorageRequired {
		t.Fatal(err)
	}
}
