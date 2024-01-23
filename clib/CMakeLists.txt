cmake_minimum_required(VERSION 3.17)
project(xdag_runtime)

set(CMAKE_C_STANDARD 11)
set(CMAKE_C_STANDARD_REQUIRED ON)
set(CMAKE_C_FLAGS "-std=gnu11 -O3 -g -Wall -fpic -Wmissing-prototypes -Wno-unused-result")

set(LIBRARY_OUTPUT_PATH ${CMAKE_SOURCE_DIR})

#set(ENV{PKG_CONFIG_PATH} secp256k1)
if ( NOT DEFINED ENV{SECP256K1_ROOT_DIR})
    find_package(PkgConfig)
    pkg_search_module(SECP256K1 REQUIRED libsecp256k1)
endif()

MESSAGE(STATUS ${SECP256K1_INCLUDE_DIRS})
MESSAGE(STATUS ${SECP256K1_LIBRARIES})
MESSAGE(STATUS ${SECP256K1_LIBDIR})

if (APPLE AND NOT DEFINED ENV{OPENSSL_ROOT_DIR})
    set(OPENSSL_ROOT_DIR "/usr/local/opt/openssl" "/usr/local/ssl/*" "/usr/local/Cellar/openssl/*")
    set(OPENSSL_LIBRARIES "/usr/local/opt/openssl/lib" "/usr/local/ssl/*/lib" "/usr/local/Cellar/openssl/*/lib")
endif()

file(GLOB_RECURSE DAGGER_SOURCES
        xDagWallet/src/client/*.c
        xDagWallet/src/client/utils/*.c
        xDagWallet/src/dus/*.c
        )


file(GLOB_RECURSE DAGGER_HEADERS
        xDagWallet/src/client/*.h
        xDagWallet/src/client/utils/*.h
        xDagWallet/src/dus/*.h
        )


OPTION(DEBUG "Build the project using debugging code" OFF)
if(DEBUG)
    MESSAGE("Adding Debug flag...")
    SET(CMAKE_BUILD_TYPE Debug)
    MESSAGE("Build type is " ${CMAKE_BUILD_TYPE})
else()
    add_definitions(-DNDEBUG)
endif(DEBUG)

add_definitions(-DDFSTOOLS)
add_definitions(-DSHA256_USE_OPENSSL_TXFM)
add_definitions(-DSHA256_OPENSSL_MBLOCK)

find_package(Threads QUIET)
if(NOT Threads_FOUND)
    message(FATAL_ERROR "Failed to find pthreads dependency!")
endif()

find_package(OpenSSL QUIET)
if(NOT OpenSSL_FOUND)
    message(FATAL_ERROR "Failed to find openssl dependency!")
endif()

include_directories(
        ${OPENSSL_INCLUDE_DIR}
        ${SECP256K1_INCLUDE_DIRS}
)

if(APPLE)
    add_library(
            ${PROJECT_NAME} STATIC
            ./xdag_runtime.cpp
            ./xdag_runtime.h
            ${DAGGER_HEADERS}
            ${DAGGER_SOURCES}
    )
else(UNIX)
    add_library(
            ${PROJECT_NAME} STATIC
            ./xdag_runtime.cpp
            ./xdag_runtime.h
            ${DAGGER_HEADERS}
            ${DAGGER_SOURCES}
    )
endif()

set_target_properties(${PROJECT_NAME} PROPERTIES SUFFIX _${CMAKE_SYSTEM_NAME}.a)

target_link_libraries(
        ${PROJECT_NAME}
        PUBLIC
        ${CMAKE_THREAD_LIBS_INIT}
)

target_link_libraries(
        ${PROJECT_NAME}
        PUBLIC
        OpenSSL::SSL
        OpenSSL::Crypto
        -L${SECP256K1_LIBDIR} ${SECP256K1_LIBRARIES}
)

