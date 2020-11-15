package rcache

import "testing"

func TestEnity(t *testing.T) {
	var err error
	_, err = loadEnity(nil, false, TestData)
	if err != ErrUnresolvableCacheEnity {
		t.Fatal(err)
	}
	_, err = loadEnity([]byte{enityTypecodeIrrevocable}, true, TestData)
	if err != ErrEnityTypecodeNotMatch {
		t.Fatal(err)
	}
	_, err = loadEnity([]byte{enityTypecodeRevocable}, false, TestData)
	if err != ErrEnityTypecodeNotMatch {
		t.Fatal(err)
	}
	_, err = loadEnity([]byte{255}, true, TestData)
	if err != ErrUnresolvableCacheEnity {
		t.Fatal(err)
	}

}
