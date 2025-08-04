---
title: "handy command-line utilities - part 2"
date: "2023-11-18"
draft: false
slug: "2023/11/18/handy-cli-utilities-part-2"
categories: ["TOOLS"]
tags: ["cli", "tui", "utilities"]
comment: true
---
This marks the second installment in a series about useful CLI utilities.
In this article, I'll delve into additional command-line tools that can enhance your productivity in the terminal.
While the [first article](/blog/2023/10/21/handy-cli-utilities-part-1/) covered text and file processing tools, as well as process management,
this one will focus on topics like searching within files and directories, along with project management.
<!--more-->

## File and directory
To effectively manage files and directories, it's beneficial to utilize a well-designed CLI tool that presents information in a clear and easily understandable format.

### [exa](https://github.com/ogham/exa)
Exa is the `ls` command replacement.
It displays a list of files and directories, akin to the `ls` command, but with added features such as colorized outputs,
Git integration, icons, support for a tree view, and more.

{{< code  >}}
```bash
$ exa --git -lh --octal-permissions --color-scale
Octal Permissions Size User Date Modified Git Name
0755  drwxr-xr-x     - mort 29 Oct 20:33   -- doc
0644  .rw-r--r--  5.5k mort 18 Nov 13:04   -- flake.lock
0644  .rw-r--r--  1.6k mort 29 May 07:57   -- flake.nix
0755  drwxr-xr-x     - mort 20 Apr 14:44   -M home-manager
0644  .rw-r--r--   14k mort 23 Apr 09:11   -- LICENSE
0755  drwxr-xr-x     - mort 22 Apr 10:14   -- nixos
0644  .rw-r--r--  2.3k mort 17 Nov 12:31   -- README.md
0755  drwxr-xr-x     - mort  6 Jul 23:48   -- static
```
{{< /code >}}

To simplify usage and avoid the need to remember all the flags, you can create an alias, similar to what I've done:

{{< code  >}}
```bash
alias l="exa --git -lh --octal-permissions --color-scale --icons"
```
{{< /code >}}

### [z](https://github.com/rupa/z) / [zoxide](https://github.com/ajeetdsouza/zoxide)
If you're seeking a convenient way to navigate to a directory without specifying the entire path, consider using `z` or `zoxide`.
There have been many instances where I needed to swiftly jump into a directory with a distinctive name, and these solutions proved helpful.

While both `z` and `zoxide` can address the primary concern, if you're seeking specific features,
you'll need to examine each project individually for more information.

I personally use `zoxide`, and in this example, I was searching for my "dotfiles" directory. I simply mentioned "dot" and here is the output:
{{< code  >}}
```bash
$ z dot
mort/.../dotfiles $
```
{{< /code >}}

## Search & Lookup
Enhancing productivity in the terminal involves finding the best matches in files and directories.
In this section, we'll discuss some of the ones I find most useful.

### [fzf](https://github.com/junegunn/fzf)
You might be familiar with it, fzf is a fuzzy finder designed to assist you in swiftly and interactively locating your desired file or directory path.

{{< code  >}}
```bash
$ fzf
> go < 3/4
> go.mod
  go.sum
  main.go
```
{{< /code >}}

When I'm engaged in a project, I often need to swiftly locate files, read snippets from them, and if it's the correct file, open it in Neovim.
Fzf serves as the ideal solution. It recursively indexes your directory, allowing you to search for a file while simultaneously previewing its content.
To enhance this process, you can utilize fzf's `--preview` parameter and pair it with the [`bat`](/blog/2023/08/30/handy-cli-utilities-part-1/#bat) command.

{{< code  >}}
```bash
$ fzf --preview 'bat {} --style=numbers --color=always'
> main < 1/4
> main.go                                                            | 1 package main                                                1/73
                                                                     | 2 import "fmt"
                                                                     | 3
                                                                     | 4 func main() {
```
{{< /code >}}

So, as you see, the `bat` command proves to be quite useful when working with `fzf`.

Additionally, you have the option to modify the default command behind the fzf command.
Personally, I opt for [`fd`](https://github.com/sharkdp/fd) due to its speed and efficiency:
{{< code  >}}
```zsh
export FZF_DEFAULT_COMMAND='fd --type f'
```
{{< /code >}}

I've also configured `f` and `o` aliases in my [zsh_aliases file](https://github.com/mortymacs/dotfiles/blob/main/home-manager/common/zsh/aliases.nix)
for swiftly accomplishing what I need:
{{< code  >}}
```zsh
alias f="fzf --preview 'bat {} --style=numbers --color=always'"
alias o="nvim `f || echo '-c :quitall'`"
```
{{< /code >}}

You can leverage `fzf` to perform text searches, such as finding the PID of a process:
{{< code  >}}
```bash
$ procs | fzf
```
{{< /code >}}

### [mcfly](https://github.com/cantino/mcfly)

For those accustomed to using `ctrl-r` to search their terminal history, consider using `mcfly` for a more convenient and enhanced experience.
Personally, I opt for `fzf` to search my history.

### [ripgrep-all](https://github.com/phiresky/ripgrep-all)

Many individuals utilize `ripgrep` for recursive searches, and to extend the search capabilities to include PDFs, Ebooks, archived, compressed,
and other file types, you can make use of the `ripgrep-all` package.

{{< code  >}}
```bash
$ rg main
main.go
1:package main
59:func main() {
```
{{< /code >}}

### [bat-extras](https://github.com/eth-p/bat-extras)
`bat-extras` encompasses various commands related to `bat`.
One notable example is `batgrep`, which combines `ripgrep` and the `bat` command to enhance the output for a better experience:

{{< code  >}}
```go
$ batgrep -i main
--------------------------------------------------------------------------------------------------------------------
     File: main.go
   1 package main
   2
   3 import (
 ...-----------------------------------------------------------------8<---------------------------------------------
  57 }
  58
  59 func main() {
  60     router := gin.New()
  61
--------------------------------------------------------------------------------------------------------------------
```
{{< /code >}}

Once installed, you can utilize its commands, such as `batman` and `batdiff`.

### [ast-grep](https://ast-grep.github.io)
This is one of my favorite tools that facilitates searching in your source codes for more precise results based on
the [Abstract Syntax Tree (AST)](https://en.wikipedia.org/wiki/Abstract_syntax_tree).
For example, consider a scenario where you have a collection of code containing the word "func" (indicating a Golang function definition)
somewhere in your README file. If you were to use the `grep` or `ripgrep` commands, they might return instances of the "func" word in the README file,
whereas you're specifically searching for functions in your source code.

In this instance, I'm searching for any type of function that takes no parameters:
{{< code  >}}
```go
$ ast-grep run -p 'func $A()'
./main.go
27|func main() {
28|	a := User{}
29|
30|	tag := reflect.TypeOf(a).Field(0).Tag
31|	fmt.Println(tag)
32|
33|}
```
{{< /code >}}

Now, I'm searching for functions with only one parameter:

{{< code  >}}
```go
$ ast-grep -p 'func $A($B)'
./main.go
21|func EmailParser(email string) string {
22|	fmt.Println("Not implemented")
23|	return ""
24|}
```
{{< /code >}}

As the result indicates, `ast-grep` is a perfect choiece for these use cases.

## Projects
Now, let's delve into project management in the terminal, which requires a combination of several commands.
If you believe that simply using "mkdir" or "git clone" will suffice, I would advise against it, as it can become chaotic
when dealing with numerous projects, making maintenance a challenging task.

### [ghq](https://github.com/x-motemen/ghq)
The `ghq` tool aids in project management by simplifying the workspace path introduction.
When you need to clone a project, you can provide the URL to `ghq`, and it automatically clones the project to the appropriate path.
`ghq` allows you to view a list of projects, create a new repository, and can be seamlessly combined with other tools we've discussed to swiftly
search and navigate to your desired project.

Let's clone a project:
{{< code  >}}
```bash
$ ghq get https://github.com/neovim/neovim.git
     clone https://github.com/neovim/neovim.git -> /home/mort/Workspaces/github.com/neovim/neovim
       git clone --recursive https://github.com/neovim/neovim.git /home/mort/Workspaces/github.com/neovim/neovim
Cloning into '/home/mort/Workspaces/github.com/neovim/neovim'...
remote: Enumerating objects: 224216, done.
remote: Counting objects: 100% (43840/43840), done.
remote: Compressing objects: 100% (728/728), done.
remote: Total 224216 (delta 43285), reused 43125 (delta 43112), pack-reused 180376
Receiving objects: 100% (224216/224216), 179.78 MiB | 6.65 MiB/s, done.
Resolving deltas: 100% (179868/179868), done.
```
{{< /code >}}
As the output indicates, it stores the project in the `/home/mort/Workspaces/github.com/neovim/neovim` path without me explicitly specifying it.

Now, let's see the list of projects:
{{< code  >}}
```bash
$ ghq list
github.com/mortymacs/abcmeta
github.com/mortymacs/dotfiles
github.com/neovim/neovim
github.com/mortymacs/nvim_context_vt
```
{{< /code >}}

We can employ `fzf` to search for a project and cd to its directory:
{{< code  >}}
```bash
$ cd $GHQ_ROOT/$(ghq list | fzf -e)
> vim  < 2/4
> github.com/neovim/neovim
  github.com/mortymacs/nvim_context_vt

mort/.../neovim $ pwd
/home/mort/Workspaces/github.com/neovim/neovim
```
{{< /code >}}

{{< quote type="important" >}}
<code>GHQ_ROOT</code> is a crucial variable that the <code>ghq</code> command is attentive to.
STherefore, it's essential to set it to your workspace path. What I've done was <code>export GHQ_ROOT="$HOME/Workspaces"</code>.
{{< /quote >}}

## Conclusion
In this article, we delved into displaying files and directories with enhanced, human-readable information through colorized output.
We discussed how to locate a file or directory using fzf, and also explored a method for finding a function in our project with context,
avoiding unrelated results.
Finally, we concluded with a project management approach, combining various commands with the `ghq` tool.
