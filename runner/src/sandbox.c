#define _GNU_SOURCE

#include <sys/ptrace.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <sys/user.h>
#include <errno.h>
#include <string.h>
#include <syscall.h>
#include <unistd.h>
#include <stdio.h>
#include <stdlib.h>
#include <seccomp.h>
#include <stdbool.h>
#include <stdarg.h>
#include <linux/prctl.h>

void die(const char *format, ...) {
    va_list args;
    va_start(args, format);
    vfprintf(stderr, format, args);
    va_end(args);

    exit(1);
}

static void enter_sandbox(uint32_t action) {
    int rc;
    
    rc = clearenv();
    if (rc != 0) {
        die("clearenv failed: %s\n", strerror(errno));
    }

    scmp_filter_ctx ctx;
    ctx = seccomp_init(action);
    if (ctx == NULL) {
        die("seccomp_init failed: %s\n", strerror(errno));
    }

    // File descriptor operations
    // These only operate on already opened file descriptors, so they're safe.
    seccomp_rule_add(ctx, SCMP_ACT_ALLOW, SCMP_SYS(read), 0);
    seccomp_rule_add(ctx, SCMP_ACT_ALLOW, SCMP_SYS(write), 0);

    rc = seccomp_load(ctx);
    if (rc < 0) {
        die("seccomp_load failed: %s\n", strerror(errno));
    }

    seccomp_release(ctx);
}

static void tracer(pid_t child) {
    int status;

    // Wait for the initial SIGSTOP, so we can configure the ptrace
    waitpid(child, &status, 0);
    if (WIFEXITED(status) || WIFSIGNALED(status)) {
        return;
    }

    ptrace(PTRACE_SETOPTIONS, child, 0, PTRACE_O_TRACESECCOMP);

    ptrace(PTRACE_CONT, child, NULL, NULL);

    while (1) {
        // Wait for the next ptrace event
        waitpid(child, &status, 0);

        // On seccomp events, log the offending system call
        if (WIFSTOPPED(status) &&
                WSTOPSIG(status) == SIGTRAP &&
                status >> 16 == PTRACE_EVENT_SECCOMP) {

            struct user_regs_struct regs;

            ptrace(PTRACE_GETREGS, child, NULL, &regs);

            fprintf(stderr,
                    "unhandled system call pid=%d rip=0x%llx %lld (0x%llx, 0x%llx, 0x%llx, 0x%llx, 0x%llx, 0x%llx)\n", child,
                    regs.rip,
                    regs.orig_rax, regs.rdi, regs.rsi, regs.rdx, regs.r10, regs.r8, regs.r9);
        }

        // Stop there if the child died
        if (WIFEXITED(status) || WIFSIGNALED(status)) {
            break;
        }

        if (WIFSTOPPED(status) && WSTOPSIG(status) == SIGSEGV) {
            struct user_regs_struct regs;

            ptrace(PTRACE_GETREGS, child, NULL, &regs);
            printf("segfault rip=0x%llx\n", regs.rip);
            break;
        }

        // Keep running the child until the next event
        ptrace(PTRACE_CONT, child, NULL, NULL);
    }
}

static void enter_sandbox_with_trace(uint32_t action) {
    pid_t child = fork();
    if (child < 0) {
        die("fork failed: %s\n", strerror(errno));
    } else if (child == 0) {
        // Die immediately if the tracer gets killed
        prctl(PR_SET_PDEATHSIG, SIGKILL);

        // Start ptracing our selves, and wait for the tracer to wake us up.
        ptrace(PTRACE_TRACEME, 0, NULL, NULL);
        raise(SIGSTOP);

        enter_sandbox(action);
    } else {
        tracer(child);
        exit(0);
    }
}

int sandbox(const char *mode) {
    uint32_t action;
    bool trace = false;

    if (strcmp(mode, "trace") == 0) {
        action = SCMP_ACT_TRACE(0);
        trace = true;
    } else if (strcmp(mode, "error") == 0) {
        action = SCMP_ACT_ERRNO(EPERM);
    } else if (strcmp(mode, "kill") == 0) {
        action = SCMP_ACT_KILL;
    } else {
        fprintf(stderr, "invalid sandbox mode: %s\n", mode);
        exit(1);
    }

    if (trace) {
        enter_sandbox_with_trace(action);
    } else {
        enter_sandbox(action);
    }

    return 1;
}
