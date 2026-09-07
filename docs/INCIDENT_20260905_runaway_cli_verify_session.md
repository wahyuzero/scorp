# Incident Report: Runaway CLI Verify Session

- **Incident Date**: September 4–5, 2026
- **Documentation**: September 5, 2026, 14:05 UTC+7
- **Severity**: Medium (No data loss; purely resource waste)
- **Status**: Resolved (Processes terminated manually)

## Summary

Two `scorp --cli` processes were detected running runaway for approximately 15 hours since the night of September 4, consuming ~170% CPU each (total ~340% CPU capacity across cores) without performing meaningful work.

| PID | Command | CPU | Total CPU-Time | Process Lifetime |
|---|---|---|---|---|
| 87914 | `./scorp --cli --session=verify_all_ok` | ~170% | 25h 44m | 15h 07m |
| 88650 | `./scorp --cli --session=verify_all_final` | ~170% | 25h 49m | 15h 07m |

Total wasted compute: **±51 core-hours** in a dead loop. System load average spiked to 4.7–6.1.

## Detected Runaway Characteristics

1. **Binary `(deleted)`** — `readlink /proc/<pid>/exe` pointed to `/home/wxsys/Project/scorp/scorp (deleted)`. The binary had been rebuilt at 07:46, but the older process was still holding the inode of the previous version. Recompiling did not affect already running instances.
2. **Adopted by `systemd --user`** — Parent process had become `systemd --user`; the originating terminal tab had been closed, leaving the processes unmonitored.
3. **No Network & No Child Processes** — Zero active TCP connections and no child subprocesses, indicating a pure in-process busy loop.
4. **Misleading Session Names** — Sessions `verify_all_ok` and `verify_all_final` never reached a conclusion.

## Root Cause Analysis

A CLI verification session initiated on the evening of September 4 encountered an unbounded loop condition without an execution timeout. The process lacked an enforceable wall-clock deadline to self-terminate when work stalled, leaving it spinning until manually terminated the following afternoon.

Application logs (`~/.scorp/scorp.log`) showed regular activity up to 07:48 (model loading, cost tracking, tool execution), after which the verification loop stalled while spinning CPU cycles.

## Remediation

```bash
kill 87914 88650
```

Post-kill validation:
- CPU utilization normalized: 39.7% → 6.0% (idle 91.7%)
- Load average steadily declined from 4.13
- Zero data loss; compiled binary remained intact on disk

## Permanent Fixes & Implementation Status

1. **CLI Execution Timeouts** — [RESOLVED] The agent execution loop (`agent/loop.go`) is now bounded with `context.WithTimeout` (default 30 minutes, configurable via `SCORP_MAX_TURN_TIMEOUT`), preventing indefinite runaway loops.
2. **Heartbeat & Stall Detection** — [RESOLVED] Enforced via strict context deadlines on HTTP transports, bounded retry limits (max 4), and subagent turn caps.
3. **Process Exclusivity (Session Lock)** — [RESOLVED] Kernel advisory file locking (`cli_lock.go`) is now required for CLI sessions. If a secondary CLI process attempts to attach to an active session, it is immediately rejected with a clear message and PID indicator.
4. **Safe Re-build Cleanup** — [RESOLVED] Advisory locks prevent history file corruption and duplicate execution collisions.
5. **Deterministic Termination Contract** — [RESOLVED] Replaced heuristic regex completion checks with an explicit `complete_task` tool contract.

## Operational Takeaways

- `(deleted)` in `/proc/<pid>/exe` is a quick sign that a process is running an outdated or replaced binary.
- A process under `systemd --user` with high CPU and zero network activity is a prime candidate for an orphaned runaway.
- Persistent elevated load averages without active user workload should always be checked with `ps aux --sort=-%cpu`.
