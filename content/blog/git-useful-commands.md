+++
title = "Git Useful Commands"
date = "2024-09-14"
draft = false
path = "blog/2024/09/14/git-useful-commands"
lang = "en"
[extra]
category = "TOOLS"
tags = ["cli", "git", "utilities", "tools"]
comment = true
+++
As a software engineer using Git almost every day, it's essential for me to master this tool to efficiently manage my code and tasks.
In this discussion, I'll go over some Git commands that have significantly boosted my productivity during development and maintenance.
<!-- more -->

### Alias

#### [Gitignore](https://git-scm.com/docs/gitignore)

You’re probably already familiar with it, but that's not the focus here.
Instead, I’ll show you how to generate a gitignore list for your project.
The main purpose of a gitignore file is to prevent certain file types from being committed or tracked in your repository.

There's an API available that allows you to generate a custom gitignore file based on your needs: `https://www.toptal.com/developers/gitignore/api/$args`.
This can be a quick and convenient way to create a gitignore list tailored to your project.

I believe the best approach is to create a simple function that utilizes the API for generating a custom gitignore file.
This way, you can easily automate the process:

for bash:
```bash
gitignore() { curl -s "https://www.toptal.com/developers/gitignore/api/$1"; }
```

for fish:
```fish
function gitignore
    curl -s "https://www.toptal.com/developers/gitignore/api/$argv";
end
```

Now, you can use it easily:

```bash
$ gitignore go > .gitignore
```

### Commands

#### [Note](https://git-scm.com/docs/git-notes)

Sometimes you might need to add a note about a commit or mention something related to it without altering the commit itself.
For instance, you might want to note that a commit was tested by X, Y, and Z.
In such cases, you can use `git note` to add this additional information.

For example, if you want to add a note to a commit indicating that it needs a better implementation due to resource constraints:
```bash
$ git notes add 9cc33538u8 -m "it needs a better implementation due to resource constraints"
```

Output in `git log`:
```gitlog
$ git log
...
commit 9cc33537d7da32cf6c5e108df6faa0784d134ab2 (tag: v3.1.0)
Author: Morteza NourelahiAlamdari <m@0t1.me>
Date:   Sun Dec 10 16:16:27 2023 +0100

    Release v3.1.0

Notes:
    it needs a better implementation
...
```

This way, you can keep track of important details or improvements needed without modifying the original commit.

and to remove the note if you want:
```bash
$ git notes remove 9cc33538u8
```

#### [Worktree](https://git-scm.com/docs/git-worktree)

When you're refactoring, developing, or making changes in your current branch,
you might need to debug or check something in the production version of the system. With `git worktree`,
you can quickly create a separate working directory for the desired branch without stashing or altering your current code.
This allows you to test and debug in a separate environment and easily remove it when you're done.
It's a straightforward and efficient solution!

```bash
$ git worktree add debug origin/main
```

Now, go to the `debug` directory (`cd debug`) and perform whatever tasks or testing you need to do.
Once you're finished, you can easily remove the directory by:

```bash
$ git worktree remove debug
```

#### [Cherry pick](https://git-scm.com/docs/git-cherry-pick)

You’ve probably heard this before, and I don’t want to bore you with the usual details.
My goal here is to give you a practical idea of when to use cherry-picking.
For example, you might get a pull request (PR) that gets rejected or canceled due to internal decisions.
Or maybe a PR gets merged into the development branch, but you only need one specific commit from it to fix something in production.
That’s where cherry-picking can be really useful.

``` bash
$ git fetch --all
$ git checkout main
$ git cherry-pick ef003e4facea1f33b2020fadc2ea844933e176b3
[master 31c2135] Fix memory leak
 Date: Sat Sep 14 16:59:58 2024 +0200
 1 file changed, 9 insertions(+), 1 deletion(-)
```

Let's review what we have so far:
```gitlog
$ git log
commit 31c2135dbbfbf623978aa0372c336c1c033f75d4 (HEAD -> master)
Author: Morteza NourelahiAlamdari <m@0t1.me>
Date:   Sat Sep 14 16:59:58 2024 +0200

    Fix memory leak
....
```

#### [Add Patch](https://git-scm.com/docs/git-add)


Sometimes, when making changes to a file, you might only want to commit a few specific lines rather than the entire file.
In the past, I’d use `git stash`, manually apply the changes I wanted, and then add the file.
Now, though, I can use `git add -p` to selectively stage only the lines I want to commit, skipping the rest of the file.

```bash
$ git add -p
```

Now, it provides an interactive environment where you can commit only the specific parts you want:
```diff
diff --git a/a.c b/a.c
index 1261f00..0653a06 100644
--- a/a.c
+++ b/a.c
@@ -2,8 +2,16 @@
 #include <stdlib.h>
 #include <string.h>

+void free_char(char **p) {
+    if (p && *p) {
+        free(*p);
+        *p = NULL;
+    }
+}
+
 int main() {
-    char *name = (char *)malloc(sizeof(char) * 15);
+    char *name __attribute__((cleanup(free_char))) = (char *)malloc(sizeof(char) * 15);
     strcpy(name, "world");
     printf("welcome: %s", name);
+    printf("goodbye\n");
 }
(1/1) Stage this hunk [y,n,q,a,d,s,e,?]?
```

Here, we only want to add the last line, `print("goodbye\n");`, and leave the rest unchanged.

As shown in the last line, `git add -p` provides a menu that lets you choose how to proceed.
Initially, you can add all the changes, but since we only want to include a specific line, you type `s` to split the current hunk into smaller chunks.

You can type `?` to see the full menu of options:
```y - stage this hunk
n - do not stage this hunk
q - quit; do not stage this hunk or any of the remaining ones
a - stage this hunk and all later hunks in the file
d - do not stage this hunk or any of the later hunks in the file
s - split the current hunk into smaller hunks
e - manually edit the current hunk
? - print help
```

After typing `s`, it now shows us this:
```diff
Split into 3 hunks.
@@ -2,4 +2,11 @@
 #include <stdlib.h>
 #include <string.h>

+void free_char(char **p) {
+    if (p && *p) {
+        free(*p);
+        *p = NULL;
+    }
+}
+
 int main() {
(3/3) Stage this hunk [y,n,q,a,d,K,g,/,e,?]?
```

But that's not what we're looking for. Let's type `n`, which means "do not stage this hunk", to move to the next chunk.
We’ll continue doing this until we find the part we want to add.

```diff
@@ -7,3 +14,4 @@
     strcpy(name, "world");
     printf("welcome: %s", name);
+    printf("goodbye\n");
 }
(3/3) Stage this hunk [y,n,q,a,d,K,g,/,e,?]?
```

Now we type `y`, which means "stage this hunk". That's it! Next, let's run `git diff --staged` to review what we’ve added:
```diff
$ git diff --staged
diff --git a/a.c b/a.c
index 1261f00..d12075d 100644
--- a/a.c
+++ b/a.c
@@ -6,4 +6,5 @@ int main() {
     char *name = (char *)malloc(sizeof(char) * 15);
     strcpy(name, "world");
     printf("welcome: %s", name);
+    printf("goodbye\n");
 }
```

This is exactly what we wanted.
If you run `git status`, you’ll see that the file appears in both staged and non-staged sections, as we only selected a part of the file for our upcoming commit.

#### [Bisect](https://git-scm.com/docs/git-bisect)

`git bisect` is a valuable tool for debugging.
It allows you to navigate through past commits and run tests or commands to identify which commit introduced a problem.
By testing a range of commits, you can pinpoint the exact commit where the issue began.

For instance, if you notice a memory leak in the code on the main branch and want to identify which commit introduced it,
you can use `git bisect` to help pinpoint the source. In C/C++, tools like Valgrind can be used to detect the memory leak.
For example:

```c
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int main() {
    char *name = (char *)malloc(sizeof(char) * 10);
    strcpy(name, "test");
    printf("%s", name);
}
```

When I compile and run the code, it works fine, but it has a memory leak that never gets released. To track it down, we use Valgrind:

```bash
  valgrind --leak-check=full --error-exitcode=1 ./a.out
==325762== Memcheck, a memory error detector
==325762== Copyright (C) 2002-2022, and GNU GPL'd, by Julian Seward et al.
==325762== Using Valgrind-3.22.0 and LibVEX; rerun with -h for copyright info
==325762== Command: ./a.out
==325762==
test==325762==
==325762== HEAP SUMMARY:
==325762==     in use at exit: 10 bytes in 1 blocks
==325762==   total heap usage: 2 allocs, 1 frees, 1,034 bytes allocated
==325762==
==325762== LEAK SUMMARY:
==325762==    definitely lost: 10 bytes in 1 blocks
==325762==    indirectly lost: 0 bytes in 0 blocks
==325762==      possibly lost: 0 bytes in 0 blocks
==325762==    still reachable: 0 bytes in 0 blocks
==325762==         suppressed: 0 bytes in 0 blocks
==325762== Rerun with --leak-check=full to see details of leaked memory
==325762==
==325762== For lists of detected and suppressed errors, rerun with: -s
==325762== ERROR SUMMARY: 0 errors from 0 contexts (suppressed: 0 from 0)
```

As indicated, there are `2 allocs` and `1 frees`, with the `LEAK SUMMARY` section providing more details on the issue.

Now, imagine your code was working perfectly, but suddenly, when you try to push your changes, you discover a memory leak.
You need to figure out which commit introduced the issue to determine if you can fix it with a patch or if other versions are affected as well.

Let's first take a look at our commit history:
```gitlog
$ git log
commit 996b0dfc6ffbc5fddb78767f26e7790aacfaa71d (HEAD -> master)
Author: Morteza NourelahiAlamdari <m@0t1.me>
Date:   Sat Sep 14 15:23:57 2024 +0200

    Remove dead code

commit e98cd3903bbe8066fe0240acfc2cdbc33075bb50
Author: Morteza NourelahiAlamdari <m@0t1.me>
Date:   Sat Sep 14 15:23:44 2024 +0200

    Fix name value to 'world'

commit 995e8283162da6a6587877c0d60b57b869a354ab
Author: Morteza NourelahiAlamdari <m@0t1.me>
Date:   Sat Sep 14 15:23:00 2024 +0200

    Change name length

commit c97e07d10ff6906570c3170e55c22163e5b83f04
Author: Morteza NourelahiAlamdari <m@0t1.me>
Date:   Sat Sep 14 15:22:45 2024 +0200

    Update message

commit 1f223b694dac72ebbf595aafbb81132c3e0de014
Author: Morteza NourelahiAlamdari <m@0t1.me>
Date:   Sat Sep 14 15:22:28 2024 +0200

    Init repo
```

The only thing we know is that the current branch, which we’ll call "bad," has the issue.
To find the problematic commit, we need to identify a "good" commit to set our testing range.
In this case, I'll use the "Init repo" commit, `1f223b694dac72ebbf595aafbb81132c3e0de014`, as our reference for the "good" commit.

Now let's get started:

```bash
$ git bisect start
status: waiting for both good and bad commits

$ git bisect bad HEAD
status: waiting for good commit(s), bad commit known

$ git bisect good 1f223b694dac72ebbf595aafbb81132c3e0de014
Bisecting: 1 revision left to test after this (roughly 1 step)
[995e8283162da6a6587877c0d60b57b869a354ab] Change name length

$ git bisect run bash -c 'gcc a.c -o debug.out && valgrind --leak-check=full --error-exitcode=1 ./debug.out'
  git bisect run bash -c 'gcc a.c -o debug.out && valgrind --leak-check=full --error-exitcode=1 ./debug.out'
running 'bash' '-c' 'gcc a.c -o debug.out && valgrind --leak-check=full --error-exitcode=1 ./debug.out'
==405388== Memcheck, a memory error detector
==405388== Copyright (C) 2002-2022, and GNU GPL'd, by Julian Seward et al.
==405388== Using Valgrind-3.22.0 and LibVEX; rerun with -h for copyright info
==405388== Command: ./debug.out
==405388==
welcome: world==405388==
==405388== HEAP SUMMARY:
==405388==     in use at exit: 15 bytes in 1 blocks
==405388==   total heap usage: 2 allocs, 1 frees, 1,039 bytes allocated
==405388==
==405388== 15 bytes in 1 blocks are definitely lost in loss record 1 of 1
==405388==    at 0x484576B: malloc (in /libexec/valgrind/vgpreload_memcheck-amd64-linux.so)
==405388==    by 0x401193: main (in /home/mort/c/debug.out)
==405388==
==405388== LEAK SUMMARY:
==405388==    definitely lost: 15 bytes in 1 blocks
==405388==    indirectly lost: 0 bytes in 0 blocks
==405388==      possibly lost: 0 bytes in 0 blocks
==405388==    still reachable: 0 bytes in 0 blocks
==405388==         suppressed: 0 bytes in 0 blocks
==405388==
==405388== For lists of detected and suppressed errors, rerun with: -s
==405388== ERROR SUMMARY: 1 errors from 1 contexts (suppressed: 0 from 0)
e98cd3903bbe8066fe0240acfc2cdbc33075bb50 is the first bad commit
commit e98cd3903bbe8066fe0240acfc2cdbc33075bb50
Author: Morteza NourelahiAlamdari <m@0t1.me>
Date:   Sat Sep 14 15:23:44 2024 +0200

    Fix name value to 'world'

 a.c | 4 ++--
 1 file changed, 2 insertions(+), 2 deletions(-)
bisect found first bad commit
```

Look, we've identified the commit that introduced the memory leak! Let's examine the changes:
```diff
$ git show e98cd3903bbe8066fe0240acfc2cdbc33075bb50
commit e98cd3903bbe8066fe0240acfc2cdbc33075bb50
Author: Morteza NourelahiAlamdari <m@0t1.me>
Date:   Sat Sep 14 15:23:44 2024 +0200

    Fix name value to 'world'

diff --git a/a.c b/a.c
index 7ddec09..c862215 100644
--- a/a.c
+++ b/a.c
@@ -10,7 +10,7 @@ void free_char(char **p) {
 }

 int main() {
-    char *name __attribute__((cleanup(free_char))) = (char *)malloc(sizeof(char) * 15);
-    strcpy(name, "test");
+    char *name = (char *)malloc(sizeof(char) * 15);
+    strcpy(name, "world");
     printf("welcome: %s", name);
 }
```

As soon as the `__attribute__((cleanup(free_char)))` was removed, we started encountering the memory leak.

Finally, run the reset command to return to the master branch:
```bash
$ git biset reset
Previous HEAD position was e98cd39 Fix name value to 'world'
Switched to branch 'master'
```

#### [Rev-List](https://git-scm.com/docs/git-rev-list)

One of my favorite ways to recover a deleted file is using this method.
For example, if I removed a config or source code file from the project six months ago and now need to see what was in
that file to copy a single line for use in my current branch, this approach is perfect.

```bash
$ git rev-list -1 --all -- /path/to/deleted-file
1c55d557f214b4c70ab2c2a0130f080762d0a3d0
```

* `-1` refers to just the first commit.
* `--all` includes all the commits.

Now, we need to checkout that commit hash to recover the file.
```bash
$ git checkout 1c55d557f214b4c70ab2c2a0130f080762d0a3d0^ -- /path/to/deleted-file
```

Now, if you check the output of `git status`:
```bash
$ git status
On branch main
Changes to be committed:
  (use "git restore --staged <file>..." to unstage)
	new file:   /path/to/deleted-file
```

Here's an interesting thread on Stack Overflow about restoring a deleted folder in a Git repository: https://stackoverflow.com/questions/30875205/restore-a-deleted-folder-in-a-git-repo

### Conclusion

Git is an incredibly useful tool for tracking changes in complex scenarios.
For example, you can find only changed files using `git grep`, trigger pipelines by committing with an empty message, and more.
It all depends on your needs. I recommend reading the official documentation, and if you like,
you can set up a bunch of aliases in your shell or Git to streamline your workflow and make processes easier to handle.
