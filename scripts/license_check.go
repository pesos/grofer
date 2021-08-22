/*
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

package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

func main() {

	licenseFile, err := os.Open("scripts/licenses.yaml")
	if err != nil {
		log.Fatal(err)
	}

	licenses := make(map[string]string)
	err = yaml.NewDecoder(licenseFile).Decode(&licenses)
	if err != nil {
		log.Fatal(err)
	}

	verified := true
	err = filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				bs, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}

				stringData := string(bs)

				for suffix, license := range licenses {
					if strings.HasSuffix(path, suffix) {
						if !strings.HasPrefix(stringData, license) {
							verified = false
							log.Println("License not verified in", path)
						}
					}
				}
			}

			return nil
		})

	if err != nil {
		log.Println(err)
	}

	if !verified {
		os.Exit(1)
	} else {
		log.Println("All files have license verified.")
	}
}
