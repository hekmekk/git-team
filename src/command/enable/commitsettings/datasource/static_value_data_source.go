package datasource

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
	"fmt"
	"os"

	"github.com/hekmekk/git-team/v2/src/command/enable/commitsettings/entity"
)

type dependencies struct {
	getEnv func(string) string
}

// StaticValueDataSource reads configuration from constant values
type StaticValueDataSource struct {
	deps dependencies
}

// NewStaticValueDataSource constructs new StaticValueDataSource
func NewStaticValueDataSource() StaticValueDataSource {
	return newStaticValueDataSource(dependencies{getEnv: os.Getenv})
}

// for tests
func newStaticValueDataSource(deps dependencies) StaticValueDataSource {
	return StaticValueDataSource{deps: deps}
}

func (ds StaticValueDataSource) Read() entity.CommitSettings {
	homeDir := ds.deps.getEnv("HOME")

	cfg := entity.CommitSettings{
		TemplatesBaseDir: fmt.Sprintf("%s/.git-team/commit-templates", homeDir),
		HooksDir:         fmt.Sprintf("%s/.git-team/hooks", homeDir),
	}
	return cfg
}
