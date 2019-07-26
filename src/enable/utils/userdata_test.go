package enableutils

import (
	"testing"
)

func TestPartitionNoInputData(t *testing.T) {
	expectedCoauthors := []string{}
	expectedAliases := []string{}

	coauthors, aliases := Partition([]string{})
	if len(coauthors) > 0 || len(aliases) > 0 {
		t.Errorf("unexpected coauthors: expected: %s, got: %s", expectedCoauthors, coauthors)
		t.Errorf("unexpected aliases: expected: %s, got: %s", expectedAliases, aliases)
		t.Fail()
	}
}

func TestPartitionAllCoauthors(t *testing.T) {
	expectedCoauthors := []string{"Mrs. Noujz <noujz@mrs.se>", "Mr. Noujz <noujz@mr.se>"}
	expectedAliases := []string{}

	coauthors, aliases := Partition([]string{"Mrs. Noujz <noujz@mrs.se>", "Mr. Noujz <noujz@mr.se>"})
	if len(coauthors) != 2 || len(aliases) > 0 {
		t.Errorf("unexpected coauthors: expected: %s, got: %s", expectedCoauthors, coauthors)
		t.Errorf("unexpected aliases: expected: %s, got: %s", expectedAliases, aliases)
		t.Fail()
	}
}

func TestPartitionAllAliases(t *testing.T) {
	expectedCoauthors := []string{}
	expectedAliases := []string{"alias1", "alias2"}

	coauthors, aliases := Partition([]string{"alias1", "alias2"})
	if len(coauthors) > 0 || len(aliases) != 2 {
		t.Errorf("unexpected coauthors: expected: %s, got: %s", expectedCoauthors, coauthors)
		t.Errorf("unexpected aliases: expected: %s, got: %s", expectedAliases, aliases)
		t.Fail()
	}
}

func TestPartition(t *testing.T) {
	expectedCoauthors := []string{"Mrs. Noujz <noujz@mrs.se>", "Mr. Noujz <noujz@mr.se>"}
	expectedAliases := []string{"alias1", "alias2"}

	coauthors, aliases := Partition([]string{"Mrs. Noujz <noujz@mrs.se>", "Mr. Noujz <noujz@mr.se>", "alias1", "alias2"})
	if len(coauthors) != 2 || len(aliases) != 2 {
		t.Errorf("unexpected coauthors: expected: %s, got: %s", expectedCoauthors, coauthors)
		t.Errorf("unexpected aliases: expected: %s, got: %s", expectedAliases, aliases)
		t.Fail()
	}
}
