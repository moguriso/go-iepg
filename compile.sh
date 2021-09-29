#!/bin/bash

java -jar closure-compiler-v20210907.jar --js bookmarklet.js --js_output_file bookmarklet-compiled.js

# コピーしてbookmarkから実行する
echo -n "javascript:"
cat bookmarklet-compiled.js

