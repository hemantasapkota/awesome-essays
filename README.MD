# awesome-essays

A simple CLI tool for listening to curated collection of essays.

**IMPORTANT**: Works on OSX only. Other platforms will be supported in the future.

![screenshot](ss.png)

# Install

```
go get -u github.com/hemantasapkota/awesome-essays
```

# Usage

## List the titles by the author
```
awesome-essays -author paul_graham | less
```

## Listen to a title

```
awesome-essays -author paul_graham -title "How To Raise Money"
```

## Listen to title by specifying index
```
awesome-essays -author paul_graham -index 150
```

# Controls
```
q: Quit
p: Pause
r: Resume
o: Open link to the essay
d: Donate to the maintainer
```

# Contributing

We'd love for you to contribute to the project. If you'd like to add new essays to the collection, please follow these guidelines:

1. Start with a template ( cool_author )
2. Format text to have 80 columns - ```fmt -w 80 file.txt```
3. Link up the essay in index.yaml
4. Recreate bin data. ```go-bindata <author1> <author2> ... <authorX>```
5. Submit a pull a request.

# License

The MIT License (MIT)

Copyright (c) 2017 Hemanta Sapkota

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
