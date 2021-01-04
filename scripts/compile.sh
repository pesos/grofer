#! /bin/bash

cd ..

archs=("amd64" "386" "arm" "arm64")

mkdir -p bin

for arch in ${archs[@]}; do
  env GOOS=linux GOARCH=${arch} go build ./grofer.go
  cp grofer bin/grofer_${arch}
  echo "Compiled grofer_${arch}"
done
