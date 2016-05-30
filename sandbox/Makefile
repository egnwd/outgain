CC=gcc
CFLAGS=-Wall -Wextra -Werror -std=c99
LDFLAGS=-lseccomp -static
SRCS=sandbox.c
TARGET=sandbox

all: $(TARGET)
clean:
	rm -f $(TARGET)
$(TARGET): $(SRCS)
	$(CC) $(CFLAGS) -o $@ $^ $(LDFLAGS)

.PHONY: all clean
