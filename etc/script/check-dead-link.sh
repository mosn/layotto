for file in $(find . -name "*.md"); do
  markdown-link-check -c .dlc.json -q "$file"
done