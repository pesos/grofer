import os
import sys
rootdir = '.'

GO_LICENSE = """/*
Copyright © 2020 The PES Open Source Team pesos@pes.edu

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
"""

BASH_LICENSE = """# Copyright © 2020 The PES Open Source Team pesos@pes.edu
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
"""

for subdir, dirs, files in os.walk(rootdir):
    if len(subdir.split("/")) >= 2 and (subdir.split("/")[1]!=".github" and subdir.split("/")[1]!="images" and subdir.split("/")[1]!=".git"):
        for file in files:
            f = open(os.path.join(subdir, file))
            data = f.readlines()
            if file.split(".")[1] == "go":
                if "".join(data[:15]) != GO_LICENSE:
                    print("License not verified in file {}".format(file))
                    sys.exit(1)
            if file.split(".")[1] == "sh":
                if "".join(data[2:15]) != BASH_LICENSE:
                    print("License not verified in file {}".format(file))
                    sys.exit(1)
            f.close()