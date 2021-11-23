package goosr

import (
	"testing"
)

func TestReplay_ReadWriteFile(t *testing.T) {
	r, err := NewReplayFromFile("data/replay_test.osr")
	if err != nil {
		t.Error(err)
	}
	if r != nil {
		t.Log(r)
	} else {
		t.Error("Couldn't parse Replay file")
	}

	de, err := r.ReplayData.Decompress()
	if err != nil {
		t.Error(err)
	}
	if len(de) > 0 {
		t.Logf("Replay data successfully decompressed. Lenght %d", len(de))
	}
	r.ReplayData = de.Compress()
	err = r.WriteToFile("data/replay_write_test.osr")
	if err != nil {
		t.Error(err)
	} else {
		t.Log("Replay was successfully saved")
	}
}
