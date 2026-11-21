cmake_minimum_required(VERSION 3.5.0)
project(HMDemo)

set(NATIVERENDER_ROOT_PATH ${CMAKE_CURRENT_SOURCE_DIR})
include_directories(${NATIVERENDER_ROOT_PATH}/include)

# Go 后端 c-shared 库：IMPORTED + IMPORTED_NO_SONAME，
# 避免 hvigor 把构建机的 Windows 绝对路径写进 libentry.so 的 DT_NEEDED / RUNPATH。
set(CMAKE_SKIP_RPATH TRUE)
add_library(homecast_go SHARED IMPORTED GLOBAL)
set_target_properties(homecast_go PROPERTIES
    IMPORTED_LOCATION ${NATIVERENDER_ROOT_PATH}/libs/${OHOS_ARCH}/libhomecast_go.so
    IMPORTED_NO_SONAME TRUE)

add_library(entry SHARED napi_init.cpp)
target_link_libraries(entry PUBLIC
    libace_napi.z.so
    libhilog_ndk.z.so
    homecast_go
)

# hvigor 只打包 CMake 真正构建出来的 target，IMPORTED 输入不会自动进 HAP。
# 用 POST_BUILD 把 libhomecast_go.so 复制到 CMake 输出目录，
# hvigor 打包 nativeLib 阶段会把它一起塞进 HAP 的 libs/<abi>/。
add_custom_command(TARGET entry POST_BUILD
    COMMAND ${CMAKE_COMMAND} -E copy
        ${NATIVERENDER_ROOT_PATH}/libs/${OHOS_ARCH}/libhomecast_go.so
        $<TARGET_FILE_DIR:entry>/libhomecast_go.so
    VERBATIM)
