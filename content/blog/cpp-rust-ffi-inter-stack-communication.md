+++
title = "C++ Rust FFI Inter-Stack Communication Concerns"
date = "2025-01-20"
draft = false
path = "blog/2024/09/14/cpp-rust-fii-inter-stack-communication-concerns"
lang = "en"
[extra]
category = "DEV"
tags = ["cpp", "rust", "ffi"]
comment = true
+++
In this article, I'll discuss internal communication aspects of C++ and Rust by using FFI to address two important concerns.
<!-- more -->

## Communication Concerns

* How will data be transferred between stacks?
* How does the owner release the memory?

I will focus on two key concerns: first, effiecent data transfer between C++ and Rust, and second: determing ownership of memory allocation and resource management.

### Passing data
* Pointers

When passing data between different stacks, as you all you, using pointers is the most efficent option.
Pointers allow us to reference the exact memory address, avoiding data duplication of passing by value.

### Memory management
* Each stack is responsible:
    * Release the allocated resources.
    * Providing release functionalities for use by other stacks.

Regarding the memory management, each stack is responsible for releasing its own allocated memory, either directly or by providing functionality that can be invoked by other stacks.

## Example

This is the directory structure of the example:

{% code() %}
```
 strtest
├──  src
│   └──  lib.rs
├──  Cargo.lock
├──  Cargo.toml
├──  cbindgen.toml
├──  main.cpp
└──  Makefile
```
{% end %}

In this example, we have a Rust project functioning as a library and C++ as the top-layer stack.
Here we have a `print_and_return` function that takes a char pointer, converts it into Rust's string data type, process it, and finally generates a new string as output, retriving its memory address.

{% code(filename="src/lib.rs") %}
```rust
use std::ffi::{CStr, CString};
use std::os::raw::c_char;

#[no_mangle]
pub extern "C" fn print_and_return(input_str: *const c_char) -> *mut c_char {
    // Convert the C string to a Rust string
    let c_str = unsafe { CStr::from_ptr(input_str) };
    let rust_str = c_str.to_str().unwrap();

    // Using the rust_str here...

    // Create a new CString to return
    // let output = ....
    let output_cstring = CString::new(output).unwrap();
    output_cstring.into_raw()
}

#[no_mangle]
pub extern "C" fn free_string(ptr: *mut c_char) {
    if !ptr.is_null() {
        unsafe {
            let _ = CString::from_raw(ptr);
        }
    }
}
```
{% end %}

As mentioned, the stack should either release its resources or provide a release functionality for the top-layer stack.
In this case we defined a `free_string` function, allowing the top-layer stack to release the memory allocated by Rust by calling this function, which targets the char pointer address.

This is the C++ example, that uses the Rust library to call the `print_and_return` function inside Rust, retrive its output, and perform any desired operations, afterward, it releases the memory allocated by Rust by calling the `free_string` function.

{% code(filename="main.cpp") %}
```c++
#include <iostream>
#include "bindings.h"

int main() {
    std::string data{"Hello from C++"};
    char* result = print_and_return(data.c_str());
    // Performing any desired operation on the result.
    free_string(result);
    return 0;
}
```
{% end %}

For resources allocated in C++, since we are using the standard C++ library for the string data type, the allocated resources will be automatically released when the string goes out of scope.

After checking the resources by Valgrind tool, we observed 4 allocations and 4 releases, indicating that there are no memory leaks.

{% code() %}
```
...
==685024== HEAP SUMMARY:
==685024==     in use at exit: 0 bytes in 0 blocks
==685024==   total heap usage: 4 allocs, 4 frees, 74,815 bytes allocated
==685024==
==685024== All heap blocks were freed -- no leaks are possible
...
```
{% end %}

Here is the output of the sample code, which demonstrates that the C++ string points to the same address as the Rust string, and vice versa.

{% code() %}
```
C++ input address: 0x7fffa7430cb0
Rust received address: 0x7fffa7430cb0
Rust returning address: 0x2f277ad0
C++ recevied address: 0x2f277ad0
```
{% end %}


This confirms that we have achieved an efficient way to transfer data between the two stacks.

This is also the way to compile the project, if you wish to do so:

{% code(filename="Makefile") %}
```make
all:
	cbindgen --config ./cbindgen.toml --crate strtest --output ./bindings.h
	cargo build
	g++ main.cpp -o main -L ./target/debug -lstrtest -I. -pthread
prod:
	cbindgen --config ./cbindgen.toml --crate strtest --output ./bindings.h
	cargo build --release
	g++ main.cpp -o main -L ./target/release -lstrtest -I. -pthread
check:
	valgrind --leak-check=full ./main
```
{% end %}

Also, the Cargo file:

{% code(filename="Cargo.toml") %}
```toml
[package]
name = "strtest"
version = "0.1.0"
edition = "2021"

[lib]
crate-type = ["staticlib"]

[dependencies]
```
{% end %}

## Note

You can release memory allocated by Rust in C++ using free(result), but this is not recommended.
Rust's memory management includes safety checks and allocation context that C++ lacks. Using free in C++ can lead to errors or undefined behavior in Rust.
Similarly, C++ resources should not be freed directly in Rust.
It's safest to let each language manage its own resources.
