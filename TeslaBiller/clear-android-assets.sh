#!/bin/sh
resPath='./android/app/src/main/res/'

ls $resPath | grep "drawable-" | while read -r folder ; do
path="$resPath$folder"
echo "Processing $folder"
ls $path | grep -E "node_modules_|assets_" | xargs -d "\n" -I % sh -c "rm $path/%;"
done
