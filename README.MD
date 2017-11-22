# awesome-essays

A simple CLI tool for listening to curated collection of essays. 

**IMPORTANT**: Works on OSX only. Other platforms will be supported in the future.

![screenshot](ss.png)

# Install

```
go get -u github.com/hemantasapkota/awesome-essays
```

# Usage

## Listen to a title

```
awesome-essays -author paul_graham -title "How To Raise Money"
```

## List the titles by the author
```
awesome-essays -author paul_graham | less
```

# Controls
```
q: Quit
p: Pause
r: Resume
o: Open link to the essay
```

# Contributing  ( TODO )