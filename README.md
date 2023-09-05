# ~~Strike~~

Did you know that `"hello".strike()` is a valid JavaScript expression? This project will blow your mind with this and other ~~useless~~ usefull facts.

# Why

> We do what we must because we can.

-- GLaDOS

This is a proof of concept that it's actually possible to generate React Server Components with backends other than Node.js. In this case, Go.

By using Go and simplified implementation of RSC, hopefully we give you a nice playground that will help you understand the principles behind RSC better.

The name `strike` doesn't have any special meaning, but it had two things going for it:
- name wasn't taken
- both html tag `<strike>` and JavaScript method `String.prototype.strike` are deprecated

# Project goals

- educate people about React Server Components (outside of Next.js/Node/Vercel ecosystem)
- get 10x performance improvement over Next.js for some synthetic use cases

# Pre-requisites

Install [air](https://github.com/cosmtrek/air) for development.

# Run

    air
    # or
    go run app.go

# Test

    air --build.bin "go test ./..." --build.exclude_regex "" --build.cmd "true"
    # or
    go test ./...

# Caution

I don't understand how HTML escaping works in Go, so I assure you there's prenty of XSS vulnerabilities in this code.

# More

This was built in public, you can watch it in the [YT Playlist](https://youtube.com/playlist?list=PLuPYpWKKQ-H12ajPoPdUO5jAhfjTeprhI&si=3ioo0SA3sP7mWuQa).
