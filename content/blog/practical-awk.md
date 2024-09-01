+++
title = "Practical AWK"
date = "2024-09-01"
draft = false
path = "blog/2024/09/01/practical-awk"
lang = "en"
[extra]
category = "TOOLS"
part = 2
tags = ["cli", "awk", "utilities"]
comment = true
+++
Sometimes, we need to run a series of commands or API calls, gather their outputs, and store them in an organized way.
For example, imagine we have a CSV file and want to use multiple cURL commands to fetch information about the data in it.
Then, we might need to clean up and format the results from each step and use them as inputs for other cURL calls or commands.
In this article, I'll show how to handle these real-life scenarios using a powerful tool called AWK.
<!-- more -->

### What is AWK
Based of [Wikipedia](https://en.wikipedia.org/wiki/AWK):
{% quote(type="quote") %}
AWK is a domain-specific language designed for text processing and typically used as a data extraction and reporting tool.
  Like sed and grep, it is a filter, and is a standard feature of most Unix-like operating systems.
{% end %}

Let's say you have a sample file like this:
```
a 10 30
b 2  25
c 15 5
```

And you want to filter the output with this query: "Get the sum of the second and third columns for all records where the third column is greater than or equal to 10."
You can easily do this with AWK:

```bash
awk '$3 >= 10 { print $1, $2 + $3}' input.txt
```

And the output will be:

```
a 40
b 27
```

### Let's get started!

In our scenario, we have a CSV file, and we want to go through multiple steps to process the data.

Sample input CSV file:

_input.csv_

```csv
Origin
https://www.gnu.org
https://0t1.me
```

Sample output CSV file:

_output.csv_

```csv
Origin,Status,Title,ContentType,IP,Country,City,PhoneCode
https://gnu.org,OK,The GNU Operating System and the Free Software Movement,text/html,209.51.188.116,United States,Boston,+1
https://0t1.me,OK,ZeroToOne - Home,text/html,104.21.84.218,Canada,Toronto,+1
```

We need to follow these steps to achieve the output:

1. Check the website status by looking at the HTTP code.
2. Retrieve the website's title by parsing the `<title>` tag.
3. Get the website's content type by reading the `Content-Type` header.
4. Find the website's IP address.
5. Use the IP address from step 4 to determine the geo location.
6. Use the geo location from step 5 to find the phone code for that area.

```
+-------------+     +------------+     +------------------+     +---------------+     +-----------------+     +---------------+
| FetchStatus | --> | FetchTitle | --> | FetchContentType | --> | FindIPAddress | --> | FindGeoLocation | --> | FindPhoneCode |
+-------------+     +------------+     +------------------+     +---------------+     +-----------------+     +---------------+
```

To create a final clean code, we should define each step as a separate function.
This makes the code more readable and easier to maintain.

#### Prepration
First, we need to create an AWK file with a different delimiter since the input is a CSV file:

We'll set `FS` to `","` to parse based on commas, and `OFS` to an empty string to ensure the output has no extra spaces around each element.

_pipeline.awk_
```awk,linenos
#!/usr/bin/awk -f

BEGIN {
    FS = ","
    OFS = ""
}
```

Now, we need to ignore the first line as well, since it's the header line:
```awk,linenos,hl_lines=9-12
#!/usr/bin/awk -f

BEGIN {
    FS = ","
    OFS = ""
}

{
    # Ignore the header line.
    if (NR == 1) {
        next
    }
}
```

Next, we need a normalizer to clean up our input addresses and prepare them for processing:
```awk,linenos,hl_lines=8-11,hl_lines=19,hl_lines=21
#!/usr/bin/awk -f

BEGIN {
    FS = ","
    OFS = ""
}

function normalize_origin(origin) {
    sub(/www./, "", origin) # remove www from the origin.
    return origin
}

{
    # Ignore the header line.
    if (NR == 1) {
        next
    }

    origin = normalize_origin($1)

    print origin
}
```

Also, make sure to grant execute permission to the file:

```bash
$ chmox u+x ./pipeline.awk
```

Now, let's run our code by passing the input file to make sure it works as expected:
```bash
$ ./pipeline.awk input.csv
```

And the output will be just what we expected, showing only the Origins column:
```
https://gnu.org
https://0t1.me
```

{% quote(type="info") %}
Sometimes, the path to your `awk` might be different from what’s mentioned in this article.
To make the script work, you can call the `awk` command directly:

```bash
$ awk -f ./pipeline.awk input.csv
```
{% end %}

Now, let's get started with the flow.

#### Step1 (Fetch Status)

Based on the flow, we need to fetch the website's status. To do this, we can use a `curl` command to get the HTTP code of the website:
```bash
$ curl -s -w '%{http_code}' URL -o /dev/null
```

We use `-s` to make `curl` silent, and `-o` to discard the website content.
Without these options, `curl` would print the entire content, which could interfere with our other processing logic.

Now, let's implement this inside the AWK file:

_pipeline.awk_
```awk,linenos,hl_lines=12-21,hl_lines=30,hl_lines=32
#!/usr/bin/awk -f

BEGIN {
    FS = ","
    OFS = ""
}

function normalize_origin(origin) {
    ...
}

function get_website_status(origin) {
    command = "curl -s -w '%{http_code}' -L '" origin "' -o /dev/null"
    command | getline status_code
    close(command)

    if (status_code == "200") {
        return "OK"
    }
    return "ERR"
}

{
    # Ignore the header line.
    if (NR == 1) {
        next
    }

    origin = normalize_origin($1)
    status = get_website_status(origin)

    print origin, ",", status
}
```

{% quote(type="important") %}
Don’t forget to include `","` (comma) between variables in the `print` statement, as we want the output in CSV format!
{% end %}

Let's run it to see the results:

```bash
$ ./pipeline.awk input.csv
https://gnu.org,OK
https://0t1.me,OK
```

Great! Now let's move on to the next step: fetching the title of the website.

#### Step2 (Fetch Title)

To complete this step, we need a tool to process HTML content. Here, we use [htmlq](https://github.com/mgdm/htmlq) for that purpose.
```bash
curl -s URL | htmlq title -t
```

Now, let's add this inside the AWK file:

_pipeline.awk_
```awk,linenos,hl_lines=16-26,hl_lines=36,hl_lines=38
#!/usr/bin/awk -f

BEGIN {
    FS = ","
    OFS = ""
}

function normalize_origin(origin) {
    ...
}

function get_website_status(origin) {
    ...
}

function fetch_website_title(origin) {
    command = "curl -s -L '" origin "' | htmlq title -t"
    if ((command | getline title) > 0) {
        gsub(/,/, " ", title)
    } else {
        title = "couldn't fetch it"
    }
    close(command)

    return title
}

{
    # Ignore the header line.
    if (NR == 1) {
        next
    }

    origin = normalize_origin($1)
    status = get_website_status(origin)
    title = fetch_website_title(origin)

    print origin, ",", status, ",", title
}
```

To parse various types of website content, you'll need different tools.
For more information about these tools, check out my other post: [Handy Command-Line Utilities - Part 1](https://0t1.me/blog/2023/10/21/handy-cli-utilities-part-1/#jq-xq-yq-htmlq-jless-fq).

#### Step3 (Fetch Content-Type)

Now, let's fetch the website's `Content-Type` value using `curl`, `grep` and `cut`:

```bash
curl -s -L -I URL | grep -i '^Content-Type: ' | cut -d' ' -f2
```

We use the `-I` flag to fetch only the headers.

Inside the AWK file:

_pipeline.awk_
```awk,linenos,hl_lines=20-28,hl_lines=39,hl_lines=41
#!/usr/bin/awk -f

BEGIN {
    FS = ","
    OFS = ""
}

function normalize_origin(origin) {
    ...
}

function get_website_status(origin) {
    ...
}

function fetch_website_title(origin) {
    ...
}

function fetch_website_content_type(origin) {
    command = "curl -s -L -I '" origin "' | grep -i '^Content-Type: ' | cut -d' ' -f2"
    if ((command | getline content_type) > 0) {
        sub(/;/, "", content_type)
    }
    close(command)

    return content_type
}

{
    # Ignore the header line.
    if (NR == 1) {
        next
    }

    origin = normalize_origin($1)
    status = get_website_status(origin)
    title = fetch_website_title(origin)
    content_type = fetch_website_content_type(origin)

    print origin, ",", status, ",", title, ",", content_type
}
```

#### Step4 (Find IP Address)

Now, it's time to find the IP address of the website. There are several options for this, but I’ve chosen the `dig` command.
```bash
dig +short HOST
```

Since the `dig` command requires the HOST, not the full website address, we need to remove the `http` and `https` prefixes.
AWK has the `sub` function for finding and replacing strings based on regex. Let’s use it to implement our function:

_pipeline.awk_
```awk,linenos,hl_lines=24-33,hl_lines=45,hl_lines=47
#!/usr/bin/awk -f

BEGIN {
    FS = ","
    OFS = ""
}

function normalize_origin(origin) {
    ...
}

function get_website_status(origin) {
    ...
}

function fetch_website_title(origin) {
    ...
}

function fetch_website_content_type(origin) {
    ...
}

function find_website_ip(origin) {
    sub(/http(s)?:\/\//, "", origin) # remove Origin's prefix
    command = "dig +short '" origin "'"
    if ((command | getline ip) == 0) {
        ip = "-"
    }
    close(command)

    return ip
}

{
    # Ignore the header line.
    if (NR == 1) {
        next
    }

    origin = normalize_origin($1)
    status = get_website_status(origin)
    title = fetch_website_title(origin)
    content_type = fetch_website_content_type(origin)
    ip = find_website_ip(origin)

    print origin, ",", status, ",", title, ",", content_type, ",", ip
}
```


Let’s run it to see the results:

```bash
$ awk -f ./pipeline.awk input.csv
https://gnu.org,OK,The GNU Operating System and the Free Software Movement,text/html,209.51.188.116
https://0t1.me,OK,ZeroToOne - Home,text/html,188.114.97.0
```

Awesome! We’ve got some results now!

Let’s move on to the next step: finding the location of the IP address we obtained with the `dig` command.

#### Step5 (Find GEO Location)

Here, I’ll use the ip-api service: `http://ip-api.com/json/{IP}`.

The output of this URL is in JSON format. So, we can use the [`jq`](https://github.com/jqlang/jq) command to parse it:

```bash
curl -s http://ip-api.com/json/{IP} | jq -r '.country + "," + .city'
```

Now, let's add this inside the AWK file:

_pipeline.awk_
```awk,linenos,hl_lines=28-36,hl_lines=49,hl_lines=51
#!/usr/bin/awk -f

BEGIN {
    FS = ","
    OFS = ""
}

function normalize_origin(origin) {
    ...
}

function get_website_status(origin) {
    ...
}

function fetch_website_title(origin) {
    ...
}

function fetch_website_content_type(origin) {
    ...
}

function find_website_ip(origin) {
    ...
}

function find_ip_location(ip) {
    command = "curl -s http://ip-api.com/json/'" ip "' | jq -r '.country + \",\" + .city'"
    if ((command | getline location) == 0) {
        location = "unknown"
    }
    close(command)

    return location
}

{
    # Ignore the header line.
    if (NR == 1) {
        next
    }

    origin = normalize_origin($1)
    status = get_website_status(origin)
    title = fetch_website_title(origin)
    content_type = fetch_website_content_type(origin)
    ip = find_website_ip(origin)
    location = find_ip_location(ip)

    print origin, ",", status, ",", title, ",", content_type, ",", ip, ",", location
}
```

Now, let's see the output:

```bash
$ ./pipeline.awk input.csv
https://gnu.org,OK,The GNU Operating System and the Free Software Movement,text/html,209.51.188.116,United States
,Boston
https://0t1.me,OK,ZeroToOne - Home,text/html,188.114.97.0,Canada,Toronto
```

#### Step6 (Find Phone Code)

The last step in the flow is to find the phone code for that location.
I’ll get this information from the [countrycode.org](https://www.countrycode.org/) website by parsing the HTML output:
```bash
curl -s https://www.countrycode.org/{country} | htmlq 'h2.text-center' -t
```

Before implementing this logic, we need to parse the location to extract only the country name and convert it to lowercase, as the website requires.
AWK's built-in `split` and `tolower` functions can help us achieve this.

So, let's implement the step:

_pipeline.awk_
```awk,linenos,hl_lines=32-43,hl_lines=57-58,hl_lines=60
#!/usr/bin/awk -f

BEGIN {
    FS = ","
    OFS = ""
}

function normalize_origin(origin) {
    ...
}

function get_website_status(origin) {
    ...
}

function fetch_website_title(origin) {
    ...
}

function fetch_website_content_type(origin) {
    ...
}

function find_website_ip(origin) {
    ...
}

function find_ip_location(ip) {
    ...
}

function find_location_phone_code(country) {
    if (country == "United States") {
        country = "usa"
    }
    command = "curl -s 'https://www.countrycode.org/" tolower(country) "' | htmlq 'h2.text-center' -t"
    if ((command | getline phone_code) == 0) {
        phone_code = "unknown"
    }
    close(command)

    return phone_code
}

{
    # Ignore the header line.
    if (NR == 1) {
        next
    }

    origin = normalize_origin($1)
    status = get_website_status(origin)
    title = fetch_website_title(origin)
    content_type = fetch_website_content_type(origin)
    ip = find_website_ip(origin)
    location = find_ip_location(ip)
    split(location, location_parts, ",")
    phone_code = find_location_phone_code(location_parts[1])

    print origin, ",", status, ",", title, ",", content_type, ",", ip, ",", location, ",", phone_code
}
```

Let's run it:

```bash
$ ./pipeline.awk input.csv
https://gnu.org,OK,The GNU Operating System and the Free Software Movement,text/html,209.51.188.116,United States
,Boston,+1
https://0t1.me,OK,ZeroToOne - Home,text/html,188.114.97.0,Canada,Toronto,+1
```

#### Finalization

Now that we have the results we're looking for, the only remaining step is to prepare it for output.

As we wanted to achieve this output structure:

```csv
Origin,Status,Title,ContentType,IP,Country,City,PhoneCode
```

We need to include a header in our output file, so we should add a `print` statement in the `BEGIN` section:
```awk,linenos,hl_lines=4
BEGIN {
    FS = ","
    OFS = ""
    print "Origin,Status,Title,ContentType,IP,Country,City,PhoneCode"
}
```

So, all in one place:

_pipeline.awk_
```awk,linenos
#!/usr/bin/awk -f

BEGIN {
    FS = ","
    OFS = ""
    print "Origin,Status,Title,ContentType,IP,Country,City,PhoneCode"
}

function normalize_origin(origin) {
    sub(/www./, "", origin) # remove www from the origin.
    return origin
}

function get_website_status(origin) {
    command = "curl -s -w '%{http_code}' -L '" origin "' -o /dev/null"
    command | getline status_code
    close(command)

    if (status_code == "200") {
        return "OK"
    }
    return "ERR"
}

function fetch_website_title(origin) {
    command = "curl -s -L '" origin "' | htmlq title -t"
    if ((command | getline title) > 0) {
        gsub(/,/, " ", title)
    } else {
        title = "couldn't fetch it"
    }
    close(command)

    return title
}

function fetch_website_content_type(origin) {
    command = "curl -s -L -I '" origin "' | grep -i '^Content-Type: ' | cut -d' ' -f2"
    if ((command | getline content_type) > 0) {
        sub(/;/, "", content_type)
    }
    close(command)

    return content_type
}

function find_website_ip(origin) {
    sub(/http(s)?:\/\//, "", origin) # remove Origin's prefix
    command = "dig +short '" origin "'"
    if ((command | getline ip) == 0) {
        ip = "-"
    }
    close(command)

    return ip
}

function find_ip_location(ip) {
    command = "curl -s http://ip-api.com/json/'" ip "' | jq -r '.country + \",\" + .city'"
    if ((command | getline location) == 0) {
        location = "unknown"
    }
    close(command)

    return location
}

function find_location_phone_code(country) {
    if (country == "United States") {
        country = "usa"
    }
    command = "curl -s 'https://www.countrycode.org/" tolower(country) "' | htmlq 'h2.text-center' -t"
    if ((command | getline phone_code) == 0) {
        phone_code = "unknown"
    }
    close(command)

    return phone_code
}

{
    # Ignore the header line.
    if (NR == 1) {
        next
    }

    origin = normalize_origin($1)
    status = get_website_status(origin)
    title = fetch_website_title(origin)
    content_type = fetch_website_content_type(origin)
    ip = find_website_ip(origin)
    location = find_ip_location(ip)
    split(location, location_parts, ",")
    phone_code = find_location_phone_code(location_parts[1])

    print origin, ",", status, ",", title, ",", content_type, ",", ip, ",", location, ",", phone_code
}
```

Finally, run the script and store the results in _output.csv_.
You can use [csview](https://github.com/wfxr/csview) to present the CSV file.

```bash
$ ./pipeline.awk input.csv > output.csv
$ csview -s ascii2 output.csv
      Origin      | Status |                          Title                          | ContentType |       IP       |    Country    |  City   | PhoneCode
 -----------------+--------+---------------------------------------------------------+-------------+----------------+---------------+---------+-----------
  https://gnu.org | OK     | The GNU Operating System and the Free Software Movement | text/html   | 209.51.188.116 | United States | Boston  | +1
  https://0t1.me  | OK     | ZeroToOne - Home                                        | text/html   | 188.114.97.0   | Canada        | Toronto | +1
```

### Conclusion

As we’ve seen, AWK is a very handy and useful tool for working with data.
If you have a series of data processing tasks, you can use AWK to achieve the desired output.
For example, in one of my real cases, I needed to call multiple APIs to gather information, such as fetching user details and then retrieving user-registered processes, and etc.

I hope you find this article useful!
