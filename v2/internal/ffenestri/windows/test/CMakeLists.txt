cmake_minimum_required(VERSION 3.19)
project(test C)

set(CMAKE_C_STANDARD 99)
set(SOURCES ../../ffenestri_windows.cpp)

add_executable(test ${SOURCES} main.c)