package enableutils

/*
Copyright (C) 2026 Rea Sand

This file is part of git-team.

This program is free software: you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation, either version 3
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.
*/

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
