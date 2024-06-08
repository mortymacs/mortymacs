+++
title = "Enforce adherence of derived classes to parent signatures in Python"
date = "2023-12-21"
draft = false
path = "blog/2023/12/21/enforce-adherence-of-derived-classes-to-parent-signatures-in-python"
lang = "en"
[extra]
category = "CODE"
tags = ["python", "classes", "oop"]
comment = true
+++
When we're developing an object oriented systems in Python, there are often instances where we have to create [abstract classes](https://en.wikipedia.org/wiki/Abstract_type).
These classes define a basic structure that we expect other parts of the code to implement.
or instance, in the context of databases, we create an abstract class that contains fundamental structure for a database class.
The derived classes then need to implement the abstract methods to form a complete structure.
The challenge arises when a derived class doesn't follow the expected signature of the abstract class.
<!-- more -->

Abstract class:

```python,linenos
from abc import ABC, abstractmethod

class Database(ABC):
    def setup(self):
        print("setup step")

    @abstractmethod
    def connect(self) -> bool:
        pass

    @abstractmethod
    def query(self, query: str) -> list[dict]:
        pass
```

Derived classes:
```python,linenos
from database import Database

class PostgreSQL(Database):
    def connect(self) -> bool:
        print("PostgreSQL trying to establish a connection")
        return True

    def query(self, query: str) -> list[dict]:
        print("running PostgreSQL query")
        return [{"data": "something"}]

class MySQL(Database):
    def connect(self) -> bool:
        print("MySQL trying to establish a connection")
        return True

    def query(self, query: str) -> list[dict]:
        print("running MySQL query")
        return [{"data": "something else"}]
```

Test:
```python,linenos
from database import MySQL
MySQL().connect()
output = MySQL().query("where 1")
print(output)
```

Output:
```bashe
$ python test.py
MySQL trying to establish a connection
running MySQL query
[{'data': 'something else'}]
```

In this particular scenario, everything functions smoothly without any issues.
However, as a project expands and more individuals become involved, there is a possibility that some might not follow the
established standards exactly.
Python itself doesn't prevent this.

Consider this implementation of the MySQL class, which completely deviates from the signatures of the abstract class:

```python,linenos,hl_lines=16 18,hide_lines=3-11
from database import Database

class PostgreSQL(Database):
    def connect(self) -> bool:
        print("PostgreSQL trying to establish a connection")
        return True

    def query(self, query: str) -> list[dict]:
        print("running PostgreSQL query")
        return [{"data": "something"}]

class MySQL(Database):
    def connect(self):
        print("MySQL trying to establish a connection")

    def query(self, query, limit):
        print("running MySQL query")
        return "something else"
```

Python test:
```python,linenos
from database import MySQL
MySQL().connect()
output = MySQL().query("where 1", 100)
print(output)
```
Output:
```bash
$ python test.py
MySQL trying to establish a connection
running MySQL query
something else
```

As you can see, Python didn't check the signature at all.

Allowing this structure to merge into your codebase would violate the established standards.

Now, can we compel everyone to adhere to the abstract class signatures?
This is where the [`abcmeta`](https://github.com/mortymacs/abcmeta) project comes into play.

```bash
$ pip install abcmeta
```

The only thing needs to be changed in your code, is in your abstract class file:

from:
```python
from abc import ABC, abstractmethod
```
to:
```python
from abcmeta import ABC, abstractmethod
```

This library then examines and ensures that all signatures align with the abstract class.

Now, let's run the previous example and see what will happen:

```bash
$ python test_db.py
Traceback (most recent call last):
  File "/home/mort/project/test.py", line 1, in <module>
    from database import MySQL
  File "/home/mort/project/database.py", line 13, in <module>
    class MySQL(Database):
  File "<frozen abc>", line 106, in __new__
  File "/home/mort/project/.venv/lib/python3.12/site-packages/abcmeta/__init__.py", line 198, in __init_subclass__
    raise AttributeError("\n{}".format("\n\n".join(errors)))
AttributeError:
1: incorrect signature.
Signature of the derived method is not the same as parent class:
- connect(self) -> bool
?              --------

+ connect(self)
Derived method expected to return in '<class 'bool'>' type, but returns 'typing.Any'

2: incorrect signature.
Signature of the derived method is not the same as parent class:
- query(self, query: str) -> list[dict]
+ query(self, query, limit)
Derived method expected to return in 'list[dict]' type, but returns 'typing.Any'
```

The error message clarifies that the signature differs from that of the abstract class, and explains with all details.

{% quote(type="important") %}
Additionally, it's important to know that <code>abcmeta</code> utilizes <a href="https://docs.python.org/3/reference/datamodel.html" target="_blank">metaclasses</a>
implying that it examines the class when the class is defined!
{% end %}

Therefore, if you modify the test code to something like this:

```python
from database import MySQL
```

Then you'll get the same error result.

So, by using these kinds of libraries in Python, we can proactively avoid future errors by forcing derived classes to follow
the abstract class signatures, similar to [strong-typed](https://en.wikipedia.org/wiki/Strong_and_weak_typing) programming languags.
