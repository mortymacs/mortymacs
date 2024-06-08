+++
title = "handy command-line utilities - part 1"
date = "2023-10-21"
draft = false
path = "blog/2023/10/21/handy-cli-utilities-part-1"
lang = "en"
[extra]
category = "TOOLS"
part = 1
tags = ["cli", "tui", "utilities"]
comment = true
+++
As I heavily depend on the terminal for a range of activities such as browsing, handling files, coding, and converting videos,
I wanted to share the collection of tools I use regularly. This blog post is divided into three sections, and in the final part,
I'll walk you through my shortcuts and aliases.
<!-- more -->

#### File
Let's begin our exploration of tools designed for reading and parsing files.

##### [bat](https://github.com/sharkdp/bat)

Bat functions similarly to the `cat` command, offering additional features such as syntax highlighting,
Git integration, and more. It's used just like the `cat` command:
```cpp
$ bat main.cpp
-------+--------------------------------------------------------------------------------------------
       | File: main.cpp
-------+--------------------------------------------------------------------------------------------
   1   | #include <iostream>
   2   | using namespace std;
   3   | int main() {
   4   |     std::cout << "Hello World!" << std::endl;
   5   | }
-------+--------------------------------------------------------------------------------------------
```

You can set up an alias to substitute your regular `cat` command with `bat`:
```zsh
alias cat=bat
```

##### [jq](https://jqlang.github.io/jq/) / xq  / [yq](https://github.com/kislyuk/yq) / [htmlq](https://github.com/mgdm/htmlq) / [jless](https://jless.io/) / [fq](https://github.com/wader/fq)

Let's employ an API to retrieve our present location details, such as country, city, and more.

The API that we're going to use is: [http://ip-api.com](http://ip-api.com)

**JSON:**<br>
To get the result in JSON format, we need to call `http://ip-api.com/json` address.
So, for example, we want to show the country, city, and region name of the current machine:
```bash
$ curl -s http://ip-api.com/json | jq ".country,.city,.regionName"
"Netherlands"
"Amsterdam"
"North Holland"
```

**XML:**<br>
To the same data result in XML format, we need to call `http://ip-api.com/xml/{IP}` address.
Like the previous call but in XML format:
```bash
$ curl -s http://ip-api.com/xml/188.114.97.0 | xq ".query.country,.query.city,.query.regionName"
"Netherlands"
"Amsterdam"
"North Holland"
```

**YAML:**<br>
It's similar to the `jq` command:
```bash
$ curl -s http://ip-api.com/json | yq ".country,.city,.regionName"
"Netherlands"
"Amsterdam"
"North Holland"
```

**HTML:**<br>
In this case, we're planning to retrieve a webpage and extract both the `title` and an element identified by its ID.
```bash
$ curl -s https://0t1.me | htmlq title
<title>ZERO/TO/ONE - Home</title>

$ curl -s https://0t1.me | htmlq '#search'
<input aria-label="Search" class="form-control form-control-sm focus-ring-dark" id="search" placeholder="Search" type="search">
```

For a more interactive experience when dealing with large JSON and YAML files, it's advisable to use `jless`.
```bash
$ curl -s http://ip-api.com/json | jless
{
  "status": "success",
  "country": "Netherlands",
  "countryCode": "NL",
  "region": "NH",
  "regionName": "North Holland",
  "city": "Amsterdam",
  "zip": "1065",
  "lat": 52.3584,
  "lon": 4.8295,
  "timezone": "Europe/Amsterdam",
  "isp": "T-Mobile Thuis WBA Services",
  "org": "",
  "as": "AS50266 tmobile thuis",
  "query": "87.210.88.217"
}

$ curl -s http://ip-api.com/json | jless --yaml
{
  "status": "success",
  "country": "Netherlands",
  "countryCode": "NL",
  "region": "NH",
  "regionName": "North Holland",
  "city": "Amsterdam",
  "zip": "1065",
  "lat": 52.3584,
  "lon": 4.8295,
  "timezone": "Europe/Amsterdam",
  "isp": "T-Mobile Thuis WBA Services",
  "org": "",
  "as": "AS50266 tmobile thuis",
  "query": "87.210.88.217"
}
```

You can make it work with XML as well, by piping the result of XML to `xq`:
```bash
$ curl -s http://ip-api.com/xml/188.114.97.0 | xq . | jless
{
  "query": {
    "status": "success",
    "country": "Netherlands",
    "countryCode": "NL",
    "region": "NH",
    "regionName": "North Holland",
    "city": "Amsterdam",
    "zip": "1012",
    "lat": "52.3759",
    "lon": "4.8975",
    "timezone": "Europe/Amsterdam",
    "isp": "Cloudflare, Inc.",
    "org": "CloudFlare, Inc.",
    "as": "AS13335 Cloudflare, Inc.",
    "query": "188.114.97.0"
  }
}
```

**Binary:**<br>

##### [fq](https://github.com/wader/fq)

We can get a binary file information with the similar way in the `jq` command.
In this example, we're going to get all tags of a binary file and requesting for one of them:
```bash
$ fq . session2.mp4
          |00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f 10 11 12 13 14 15|0123456789abcdef012345|.{}: session2.mp4 (mp4)
0x00000000|00 00 00 20 66 74 79 70 69 73 6f 6d 00 00 02 00 69 73 6f 6d 69 73|... ftypisom....isomis|  boxes[0:4]:
*         |until 0x231b4e8d.7 (end) (588992142)                             |                      |
0x0000002c|            00 00 03 06 06 05 ff ff ff 02 dc 45 e9 bd e6 d9 48 b7|    ...........E....H.|  tracks[0:2]:
0x00000042|96 2c d8 20 d9 23 ee ef 78 32 36 34 20 2d 20 63 6f 72 65 20 31 36|.,. .#..x264 - core 16|
*         |until 0x231b4e8d.7 (end) (588992094)                             |                      |

$ fq .boxes session2.mp4
          |00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f 10 11 12 13 14 15|0123456789abcdef012345|.boxes[0:4]:
0x00000000|00 00 00 20 66 74 79 70 69 73 6f 6d 00 00 02 00 69 73 6f 6d 69 73|... ftypisom....isomis|  [0]{}: box
0x00000016|6f 32 61 76 63 31 6d 70 34 31                                    |o2avc1mp41            |
0x00000016|                              00 00 00 08 66 72 65 65            |          ....free    |  [1]{}: box
0x00000016|                                                      22 cd 84 48|                  "..H|  [2]{}: box
0x0000002c|6d 64 61 74 00 00 03 06 06 05 ff ff ff 02 dc 45 e9 bd e6 d9 48 b7|mdat...........E....H.|
*         |until 0x22cd846f.7 (583894088)                                   |                      |
0x22cd846c|            00 4d ca 1e 6d 6f 6f 76 00 00 00 6c 6d 76 68 64 00 00|    .M..moov...lmvhd..|  [3]{}: box
0x22cd8482|00 00 00 00 00 00 00 00 00 00 00 00 03 e8 00 6f d1 a0 00 01 00 00|...............o......|
*         |until 0x231b4e8d.7 (end) (5098014)                               |                      |
```

##### [hexyl](https://github.com/sharkdp/hexyl)

For reading a binary file, hexyl is the go-to tool:

```bash
$ hexyl session2.mp4
+--------+-------------------------+-------------------------+--------+--------+
|00000000| 00 00 00 20 66 74 79 70 | 69 73 6f 6d 00 00 02 00 |*** ftyp|isom**.*|
|00000010| 69 73 6f 6d 69 73 6f 32 | 61 76 63 31 6d 70 34 31 |isomiso2|avc1mp41|
|00000020| 00 00 00 08 66 72 65 65 | 22 cd 84 48 6d 64 61 74 |***.free|"xxHmdat|
|00000030| 00 00 03 06 06 05 ff ff | ff 02 dc 45 e9 bd e6 d9 |**....xx|x.xExxxx|
|00000040| 48 b7 96 2c d8 20 d9 23 | ee ef 78 32 36 34 20 2d |Hxx,x x#|xxx264 -|
|00000050| 20 63 6f 72 65 20 31 36 | 31 20 72 33 30 32 37 20 | core 16|1 r3027 |
|00000060| 34 31 32 31             |                         |4121    |        |
+--------+-------------------------+-------------------------+--------+--------+
```

You can change the base to binary:
```bash
$ hexyl --base binary session2.mp4
+--------+-------------------------------------------------------------------------+-------------------------------------------------------------------------+--------+--------+
|00000000| 00000000 00000000 00000000 00100000 01100110 01110100 01111001 01110000 | 01101001 01110011 01101111 01101101 00000000 00000000 00000010 00000000 |*** ftyp|isom**.*|
|00000010| 01101001 01110011 01101111 01101101 01101001 01110011 01101111 00110010 | 01100001 01110110 01100011 00110001 01101101 01110000 00110100 00110001 |isomiso2|avc1mp41|
|00000020| 00000000 00000000 00000000 00001000 01100110 01110010 01100101 01100101 | 00100010 11001101 10000100 01001000 01101101 01100100 01100001 01110100 |***.free|"xxHmdat|
|00000030| 00000000 00000000 00000011 00000110 00000110 00000101 11111111 11111111 | 11111111 00000010 11011100 01000101 11101001 10111101 11100110 11011001 |**....xx|x.xExxxx|
|00000040| 01001000 10110111 10010110 00101100 11011000 00100000 11011001 00100011 | 11101110 11101111 01111000 00110010 00110110 00110100 00100000 00101101 |Hxx,x x#|xxx264 -|
|00000050| 00100000 01100011 01101111 01110010 01100101 00100000 00110001 00110110 | 00110001 00100000 01110010 00110011 00110000 00110010 00110111 00100000 | core 16|1 r3027 |
|00000060| 00110100 00110001 00110010 00110001                                     |                                                                         |4121    |        |
+--------+-------------------------------------------------------------------------+-------------------------------------------------------------------------+--------+--------+
```

##### [tokei](https://github.com/XAMPPRocky/tokei)

In order to get an aggregated information about source codes, you need Tokei:

```bash
$ tokei
===============================================================================
 Language            Files        Lines         Code     Comments       Blanks
===============================================================================
 Go                     22         1661         1222          188          251
 Makefile                1            8            7            0            1
 TOML                    2          142           40           81           21
 YAML                    1           46           42            0            4
-------------------------------------------------------------------------------
 Markdown                1           30            0           22            8
 |- BASH                 1            2            2            0            0
 (Total)                             32            2           22            8
===============================================================================
 Total                  27         1887         1311          291          285
===============================================================================
```

You can use it with `-o` flag to define the output format.
For example, we want to get only list of Golang files:
```bash
$ tokei -o json | jq '.Go.reports[].name'
"./cmd/mark.go"
"./cmd/root.go"
"./cmd/update.go"
"./cmd/today.go"
"./internal/util/util.go"
"./internal/workspace/workspace.go"
"./internal/task/task.go"
"./cmd/add.go"
"./internal/task/status.go"
"./internal/render/color.go"
"./internal/render/text.go"
"./internal/recurring/recurring.go"
"./internal/render/table.go"
"./internal/config/config.go"
"./internal/project/project.go"
"./cmd/tui.go"
"./internal/storage/file.go"
"./internal/dto/deserialize.go"
"./main.go"
"./cmd/show.go"
"./cmd/delete.go"
"./cmd/list.go"
```

#### Process
Now, let's move to list of commands to work with processes:

##### [procs](https://github.com/dalance/procs)

I opted for `procs` over `ps` mainly because it displays the bound ports of each process.
Configuration can be done through a TOML file.

```bash
$ procs
PID:   User             | State Nice CPU MEM   VmSize    VmRSS | TCP          UDP             Read Write | Docker | Command
1535   mort             | S        0 0.0 0.0   5.645M   3.855M | []           []                 0     0 |        | /nix/store/rap2690k6sw4rd3b5zgqp4yx3lc3clqh-dbus-1.14.8/bin/dbus-daemon
1540   root             | S        0 0.0 0.1 531.027M  16.500M |                                         |        | /nix/store/il17dg6g69gnj10zy6d7vq22z39wy1ri-udisks-2.9.4/libexec/udisks2
1561   mort             | S        0 0.0 0.1 574.078M  44.125M | []           []                 0     0 |        | nm-applet
1562   mort             | S        0 7.2 0.0 612.559M  14.125M | []           []                 0     0 |        | polybar main --log=error
1567   mort             | S        0 0.0 0.0 299.582M   5.500M | []           []                 0     0 |        | xss-lock -- XSECURELOCK_FONT=sans xsecurelock
1568   mort             | S        0 0.0 0.6   3.534G 182.824M | [8080]       []                 0     0 |        | test-app                                                                    >
```

##### [btop](https://github.com/aristocratos/btop)

This is part of my morning routine—running this command every day to keep an eye on my system.

Based on the documentation:
{% quote(type="info") %}
Resource monitor that shows usage and stats for processor, memory, disks, network and processes.
{% end %}

You can configure your btop by putting the config file in the `~/.config/btop/btop.conf` path.

For example, this is my config file:
```toml
color_theme = "tokyo-night"
rounded_corners = False
theme_background = False
update_ms = 1000
```

##### [kmon](https://github.com/orhun/kmon)

Similar to `btop`, but only for monitoring the Kernel.

Based on the documentation:
{% quote(type="info") %}
<code>kmon</code> provides a text-based user interface for managing the Linux kernel modules and monitoring the kernel activities.
By managing, it means loading, unloading, blacklisting and showing the information of a module.
These updates in the kernel modules, logs about the hardware and other kernel messages can be tracked with the real-time
activity monitor in kmon. Since the usage of different tools like dmesg and kmod are required for these tasks in Linux,
kmon aims to gather them in a single terminal window and facilitate the usage as much as possible while keeping the functionality.
{% end %}

Recently, while using kmon, I noticed an error in my Linux kernel related to writing to a pipe.
This could serve as a starting point to investigate and identify the root cause of the issue.

##### [ctop](https://github.com/bcicen/ctop)

Similar to `btop` or `kmon`, but only for monitoring the containers.

Based on the documentation:
{% quote(type="info") %}
<code>ctop</code> provides a concise and condensed overview of real-time metrics for multiple containers,
as well as a single container view for inspecting a specific container. <code>ctop</code> comes with built-in support for Docker and runC;
connectors for other container and cluster systems are planned for future releases.
{% end %}

I don't use this tool often since I primarily employ Docker for testing and deploying local services.
However, when needed, it proves helpful in easily monitoring container usage.

##### [sysz](https://github.com/joehillen/sysz)

This tool combines `fzf` and systemctl, simplifying the process for users to manage daemons effortlessly.
Utilizing `fzf` allows for quickly locating the desired daemon and promptly taking action through `systemctl` commands.<br>
For instance, you can swiftly locate and restart the Docker service using this tool.

#### Conclusion

Numerous CLI tools are available on the internet to help you accomplish your objectives.<br>
If you're aware of other tools worth exploring, please share them with me.
I'm getting ready to put together the second part of this post.
