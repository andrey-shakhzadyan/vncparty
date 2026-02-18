#!/bin/bash 

rm -rf ~/projects/vncparty/public/
npm run build
cp -r ./dist/ ~/projects/vncparty/public
