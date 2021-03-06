// Copyright (C) 2020 Red Hat, Inc.
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program; if not, write to the Free Software Foundation, Inc.,
// 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.

/*
Package config provides test-network-function configuration along with a config pool for aggregating configuration.
Configurations registered with the pool are automatically included in the claim.  Go structs used for configuration
should each be defined in their own files, such as `cnf.go` and `generic.go`.  The corresponding configuration yaml/json
files should be prefixed with `<filename>_test_configuration`, such as `generic_test_configuration.yaml`.
*/
package config
