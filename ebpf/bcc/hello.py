#!/usr/bin/python
from bcc import BPF

# this may not work for 4.17 on x64, you need replace kprobe__sys_clone
# kprobe____x64_sys_clone
b = BPF(text="""
int hello(void *ctx) {
  bpf_trace_printk("Hello, World!\\n");
  return 0;
}
""")

# kprobe - dynamic tracing of a kernel function call
# attach_kprobe(event="event", fn_name="name")
# Instruments the kernel function event() using kernel dynamic
# tracing of the function entry, and attaches our C defined function
# name() to be called when the kernel fucntion is called
b.attach_kprobe(event=b.get_syscall_fnname("clone"), fn_name=b"hello")

print("Tracing... Hit Ctrl+C to end")

# output
# Read kernel's trace_pipe and print the messages
print("%-18s %-16s %-6s %s" % ("TIME(s)", "COMMAND", "PID", "MESSAGE"))

while True:
    try:
        (task, pid, cpu, flags, ts, msg) = b.trace_fields()
    except ValueError:
        continue
    except KeyboardInterrupt:
        print("Bye bye!")
        break
    except Exception as e:
        raise e
    print("%-18.9f %-16s %-6d %s" % (ts, task.decode(), pid, msg.decode()))
