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
#include <sys/prctl.h>
#include <arpa/inet.h>
#include <stdarg.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <sched.h>
#include <linux/futex.h>
#include <sys/mman.h>

void die(const char *format, ...) {
    va_list args;
    va_start(args, format);
    vfprintf(stderr, format, args);
    va_end(args);

    exit(1);
}

void setup_seccomp(uint32_t action) {
    int rc;

    scmp_filter_ctx ctx;
    ctx = seccomp_init(action);
    if (ctx == NULL) {
        die("seccomp_init failed: %s\n", strerror(errno));
    }

    // Memory allocation
    seccomp_rule_add(ctx, SCMP_ACT_ALLOW, SCMP_SYS(brk), 0);
    seccomp_rule_add(ctx, SCMP_ACT_ALLOW, SCMP_SYS(mmap), 1, SCMP_A3(SCMP_CMP_MASKED_EQ, MAP_ANONYMOUS, MAP_ANONYMOUS));
    seccomp_rule_add(ctx, SCMP_ACT_ALLOW, SCMP_SYS(munmap), 0);
    seccomp_rule_add(ctx, SCMP_ACT_ALLOW, SCMP_SYS(mprotect), 0);

    // Futexes
    seccomp_rule_add(ctx, SCMP_ACT_ALLOW, SCMP_SYS(futex), 1, SCMP_A1(SCMP_CMP_EQ, FUTEX_WAIT));
    seccomp_rule_add(ctx, SCMP_ACT_ALLOW, SCMP_SYS(futex), 1, SCMP_A1(SCMP_CMP_EQ, FUTEX_WAKE));
    seccomp_rule_add(ctx, SCMP_ACT_ALLOW, SCMP_SYS(set_robust_list), 0);

    // Signal handlers
    seccomp_rule_add(ctx, SCMP_ACT_ALLOW, SCMP_SYS(sigaction), 0);
    seccomp_rule_add(ctx, SCMP_ACT_ALLOW, SCMP_SYS(sigprocmask), 0);
    seccomp_rule_add(ctx, SCMP_ACT_ALLOW, SCMP_SYS(sigreturn), 0);
    seccomp_rule_add(ctx, SCMP_ACT_ALLOW, SCMP_SYS(rt_sigaction), 0);
    seccomp_rule_add(ctx, SCMP_ACT_ALLOW, SCMP_SYS(rt_sigprocmask), 0);
    seccomp_rule_add(ctx, SCMP_ACT_ALLOW, SCMP_SYS(rt_sigreturn), 0);
    seccomp_rule_add(ctx, SCMP_ACT_ALLOW, SCMP_SYS(sigaltstack), 0);

    // Process manipulation
    // These only modify the current procees, so they're safe.
    seccomp_rule_add(ctx, SCMP_ACT_ALLOW, SCMP_SYS(arch_prctl), 0);
    seccomp_rule_add(ctx, SCMP_ACT_ALLOW, SCMP_SYS(sched_yield), 0);
    seccomp_rule_add(ctx, SCMP_ACT_ALLOW, SCMP_SYS(sched_getaffinity), 0);
    seccomp_rule_add(ctx, SCMP_ACT_ALLOW, SCMP_SYS(set_tid_address), 0);
    seccomp_rule_add(ctx, SCMP_ACT_ALLOW, SCMP_SYS(gettid), 0);
    seccomp_rule_add(ctx, SCMP_ACT_ALLOW, SCMP_SYS(getrlimit), 0);
    seccomp_rule_add(ctx, SCMP_ACT_ALLOW, SCMP_SYS(exit_group), 0);

    // Thread creation
    // Only allow the exact flags used by pthread.
    uint64_t clone_flags =  CLONE_VM | CLONE_FS | CLONE_FILES | CLONE_SIGHAND | CLONE_THREAD |
        CLONE_SYSVSEM | CLONE_SETTLS | CLONE_PARENT_SETTID | CLONE_CHILD_CLEARTID;
    seccomp_rule_add(ctx, SCMP_ACT_ALLOW, SCMP_SYS(clone), 1, SCMP_A0(SCMP_CMP_EQ, clone_flags));

    // Time getters
    // These only get the time, they don't have any side effects.
    seccomp_rule_add(ctx, SCMP_ACT_ALLOW, SCMP_SYS(time), 0);
    seccomp_rule_add(ctx, SCMP_ACT_ALLOW, SCMP_SYS(clock_gettime), 0);

    // File descriptor operations
    // These only operate on already opened file descriptors, so they're safe.
    seccomp_rule_add(ctx, SCMP_ACT_ALLOW, SCMP_SYS(read), 0);
    seccomp_rule_add(ctx, SCMP_ACT_ALLOW, SCMP_SYS(write), 0);
    seccomp_rule_add(ctx, SCMP_ACT_ALLOW, SCMP_SYS(close), 0);
    seccomp_rule_add(ctx, SCMP_ACT_ALLOW, SCMP_SYS(select), 0);
    seccomp_rule_add(ctx, SCMP_ACT_ALLOW, SCMP_SYS(fstat), 0);
    seccomp_rule_add(ctx, SCMP_ACT_ALLOW, SCMP_SYS(setsockopt), 0);

    // Ideally we'd want to block this, but we need it to execute the process.
    // However, anything that gets executed will have the same seccomp rules applied
    // to it, and suid binaries don't raise privileges thanks to PR_SET_NO_NEW_PRIVS.
    // Thus there's not much that can be done by the executed binary.
    seccomp_rule_add(ctx, SCMP_ACT_ALLOW, SCMP_SYS(execve), 0);

    seccomp_rule_add(ctx, SCMP_ACT_ALLOW, SCMP_SYS(uname), 0);

    // Filesystem and Network access.
    // These may be issued in various context, so we shouldn't die on them,
    // merely block them by return -EPERM
    // The Go runtime probes IP support by creating dummy sockets.
    seccomp_rule_add(ctx, SCMP_ACT_ERRNO(EPERM), SCMP_SYS(open), 0);
    seccomp_rule_add(ctx, SCMP_ACT_ERRNO(EPERM), SCMP_SYS(openat), 0);
    seccomp_rule_add(ctx, SCMP_ACT_ERRNO(EPERM), SCMP_SYS(access), 0);
    seccomp_rule_add(ctx, SCMP_ACT_ERRNO(EPERM), SCMP_SYS(readlink), 0);
    seccomp_rule_add(ctx, SCMP_ACT_ERRNO(EPERM), SCMP_SYS(socket), 0);

    rc = seccomp_load(ctx);
    if (rc < 0) {
        die("seccomp_load failed: %s\n", strerror(errno));
    }

    seccomp_release(ctx);
}

void execute_sandbox(uint32_t action, const char *prog, const char* argv[]) {
    setup_seccomp(action);

    // Don't leak any environment variables to the sandboxed process
    char *env[] = { NULL };

    // glibc's execve uses a bunch of system calls which are now denied
    // Use a raw system call so only sys_execve gets called.
    syscall(SYS_execve, prog, argv, env);

    die("execve(%s) failed: %s\n", prog, strerror(errno));
}

void execute_sandbox_with_trace(uint32_t action, const char *prog, const char* argv[]) {
    pid_t child = fork();
    if (child < 0) {
        die("fork failed: %s\n", strerror(errno));
    } else if (child == 0) {
        // Die immediately if the tracer gets killed
        prctl(PR_SET_PDEATHSIG, SIGKILL);

        // Start ptracing our selves, and wait for the tracer to wake us up.
        ptrace(PTRACE_TRACEME, 0, NULL, NULL);
        raise(SIGSTOP);

        execute_sandbox(action, prog, argv);
    } else {
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

                if (regs.orig_rax == 49) {
                    size_t sock_len = regs.rdx;
                    void *sock = (void*)regs.rsi;

                    uint64_t buffer[64];
                    for (unsigned int i = 0; i * 4 < sock_len; i ++) {
                        buffer[i] = ptrace(PTRACE_PEEKDATA, child, sock + i * 4, NULL);
                    }

                    char desc[1024];
                    struct sockaddr *addr = (struct sockaddr *)buffer;
                    if (addr->sa_family == AF_INET) {
                        struct sockaddr_in *addr = (struct sockaddr_in *)buffer;
                        inet_ntop(addr->sin_family, &addr->sin_addr, desc, 1024);
                        printf("bind: %s %x\n", desc, addr->sin_port);
                    } else if (addr->sa_family == AF_INET6) {
                        struct sockaddr_in6 *addr = (struct sockaddr_in6 *)buffer;
                        inet_ntop(addr->sin6_family, &addr->sin6_addr, desc, 1024);
                        printf("bind: %s %x\n", desc, addr->sin6_port);
                    }

                }
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
        exit(0);
    }
}

int main(int argc, const char* argv[]) {
    if (argc < 3) {
        fprintf(stderr, "Usage: %s MODE PROG ARGV..\n", argv[0]);
        exit(1);
    }

    uint32_t action;
    bool trace = false;

    if (strcmp(argv[1], "trace") == 0) {
        action = SCMP_ACT_TRACE(0);
        trace = true;
    } else if (strcmp(argv[1], "error") == 0) {
        action = SCMP_ACT_ERRNO(EPERM);
    } else if (strcmp(argv[1], "kill") == 0) {
        action = SCMP_ACT_KILL;
    } else {
        fprintf(stderr, "invalid sandbox mode: %s\n", argv[1]);
        exit(1);
    }

    if (trace) {
        execute_sandbox_with_trace(action, argv[2], &argv[3]);
    } else {
        execute_sandbox(action, argv[2], &argv[3]);
    }

    return 1;
}
