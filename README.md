# OpenAI Thread/Image Retrievel Tool
CLI retrieve threads or images (files) individually, programatically, etc.

For info:
```bash
python oait --help
```

## Examples
With just python:

```bash
# Read default input (input.txt), write to default output (output.txt)
python oait.py threads -f -o
```

```bash
# Read input.json, write to output.json
python oait.py threads -f input.json -o output.json
```

```bash
# Read 2 thread ids, write to standard out
python oait.py threads thread_id123456789 thread_id987654321
```

```bash
# retrieve image by file id, output to image.png
python oait.py images -r file_123456789 -o "image.png"
```

It will use your currently active python environment.

With shell file, it's the same but can just call "oait" directly without python and do it from anywhere.

## Shell Instructions
Either call with python from cli, or follow below steps.

Allows for calling anywhere in shell. Ex:
```bash
oait threads -f input.txt -o output.txt
```

Copy oait/oait to ~/bin (or wherever your user bin is)

```bash
cp oait/oait ~/bin/oait
chmod +x ~/bin/oait
```

If your user bin is not sourced, add user bin to path (substitute .zshrc with .bashrc if that's what you have)
```bash
echo 'export PATH=$PATH:~/bin' >> ~/.zshrc
source ~/.zshrc
```

You can either put in your OpenAI key each time with the -k flag, or add to environment with:
```bash
echo 'export OPENAI_API_KEY="your_key"'
```

For the python file, you may choose to edit where it's located, but it's defaulted to ~/scripts/oait/oait.py in the shell script.
