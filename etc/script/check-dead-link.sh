#!/usr/bin/env bash

set -u

NC='\033[0m' # No Color
GREEN='\033[0;32m'
YELLOW='\033[0;33m'

for file in $(find . -name "*.md"); do
  markdown-link-check -c .github/dead_link_check_config.json -q "$file" >> result.txt 2>&1
done

if [ -e result.txt ] ; then
  cat result.txt
  if grep -q "ERROR:" result.txt; then
      echo -e "${YELLOW}=========================> MARKDOWN LINK CHECK RESULT<=========================${NC}"
      printf "\n"
      awk -F ' ' '/links checked/{sum+=$1}END{print "Total "sum" links checked.\n"}' result.txt
      awk -F ' ' '/ERROR/{sum+=$2}END{print "[✖] Found "sum " dead links.\n"}' result.txt 
      echo -e "${YELLOW}=========================================================================${NC}"
      exit 2
  else
      echo -e "${YELLOW}=========================> MARKDOWN LINK CHECK RESULT<=========================${NC}"
      printf "\n"
      awk -F ' ' '/links checked/{sum+=$1}END{print "Total "sum" links checked.\n"}' result.txt
      echo -e "${GREEN}[✔] All links are good!${NC}"
      echo -e "${YELLOW}=========================================================================${NC}"
  fi
else
  echo -e "${GREEN}No link need check!${NC}"
fi