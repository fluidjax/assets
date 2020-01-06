# List of examples
file(GLOB_RECURSE SRCS *.c)

# Add the binary tree directory to the search path for linking and include files
link_directories (${PROJECT_BINARY_DIR}/src
                  /usr/local/lib)

include_directories (${PROJECT_SOURCE_DIR}/include
                     /usr/local/include)

foreach(example ${SRCS})
  # Extract the filename without an extension
  get_filename_component(target ${example} NAME_WE)

  add_executable(${target} ${example})

  target_link_libraries(${target} pqnist)
endforeach(example)
