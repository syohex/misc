cmake_minimum_required(VERSION 3.15)

project(simple_http_server01 VERSION 0.1.0)

set(CMAKE_C_STANDARD 11)
set(CMAKE_C_STANDARD_REQUIRED ON)
set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -O2 -Wall -Wextra -fsanitize=address,undefined -g")

set(LIBRARY_SOURCES
  "lib/error.c"
  "lib/server.c")

add_library(server ${LIBRARY_SOURCES})

set(MAIN_SOURCES "src/main.c")

add_executable(server01 ${MAIN_SOURCES})
target_link_libraries(server01 server)
