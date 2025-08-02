cmake_minimum_required(VERSION 3.13)
project(s3-stream-processor)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_BUILD_TYPE Release)

find_package(AWSSDK REQUIRED COMPONENTS s3 transfer)

add_executable(s3_stream_processor src/main.cpp)

target_link_libraries(s3_stream_processor ${AWSSDK_LINK_LIBRARIES})
target_compile_definitions(s3_stream_processor PRIVATE AWS_DISABLE_DEPRECATION)
