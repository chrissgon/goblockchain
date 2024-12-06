#!/bin/sh
while read p || [ -n "$p" ] 
do  
sed -i '' "/${p//\//\\/}/d" ./coverage.out 
done < ./coverage-ignore.txt