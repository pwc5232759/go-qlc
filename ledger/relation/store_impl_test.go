package relation

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/qlcchain/go-qlc/common/types"
	"github.com/qlcchain/go-qlc/config"
	"github.com/qlcchain/go-qlc/test/mock"
)

func TestRelation_CreateData(t *testing.T) {
	dir := filepath.Join(config.QlcTestDataDir(), "sqlite.db")
	blk := mock.StateBlockWithoutWork()
	blk.Type = types.Send
	blk.Sender = []byte("1580000")
	blk.Receiver = []byte("1851111")
	r, err := NewRelation(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := r.Close(); err != nil {
			t.Fatal(err)
		}
		if err := os.Remove(dir); err != nil {
			t.Fatal(err)
		}
	}()
	err = r.AddBlock(blk)
	if err != nil {
		t.Fatal(err)
	}
	b, err := r.AccountBlocks(blk.Address, 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(b)
	b, err = r.Blocks(10, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(b)
	count, err := r.BlocksCount()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(count)
	g, err := r.BlocksCountByType()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(g)
	err = r.DeleteBlock(blk.GetHash())
	if err != nil {
		t.Fatal(err)
	}
}