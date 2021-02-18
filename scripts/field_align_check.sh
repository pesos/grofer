#!/bin/bash

# Copyright © 2020 The PES Open Source Team pesos@pes.edu
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Check if fields of structs in go files are ordered to take up the 
# least space possible. 

function not_aligned() {
  echo -e "\nHmmm, it looks the fields of certain structs can be re-ordered\n"
  echo -e "You can use the \`fieldalignment\` tool to apply the recommended reordering as follows:\n"
  echo -e "\tUsage: fieldalignment -fix /path/to/file/having/struct"
}

function command_not_found() {
  echo -e "The command fieldalignment doesn't seem to be installed.\n"
  echo -e "\tIntsallation: go get golang.org/x/tools/go/analysis/passes/fieldalignment\n"
}

if ! type fieldalignment &> /dev/null; then
  command_not_found
  exit 1
fi

if ! fieldalignment ./... ; then
  not_aligned
  exit 1
fi
