# OpenAI Tools
CLI retrieve, modify, or add threads, files, images, etc. individually, programatically, etc.

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
# retrieve image by file id, output to image.png
oait files -i "file_123456789" -o "image.png"
```

You can either put in your OpenAI key each time with the -k flag, or add to environment with:
```bash
echo 'export OPENAI_API_KEY="your_key"'
