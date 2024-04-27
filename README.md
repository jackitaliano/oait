# OpenAI Tools
CLI retrieve, modify, or add threads, files, images, etc. individually, programatically, etc.

# Install
## Go
1. Install Go Lang: http://golang.org/doc/install.html
2. `go install github.com/jackitaliano/oait@latest`

## Homebrew
1. `brew tap jackitaliano/tap`
2. `brew install oait`

# Use
You can either put in your OpenAI key each time with the -k flag, or add to environment with:
```bash
echo 'export OPENAI_API_KEY="your_key"'
```

For info:
```bash
oait --help
```

## Examples
```bash
# Read input.txt, write to output.json
oait threads -f input.txt -o output.json
```

```bash
# Read 2 thread ids, write to standard out
oait threads -i "thread_id123456789 thread_id987654321"
```

```bash
# 
oait files -i "file_123456789" -o "image.png"
```

```bash
# Get all files and filter by less than or equal to 1 day old files
oait files -A -d 1
```


