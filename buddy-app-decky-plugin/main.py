import os
import subprocess
import threading

# The decky plugin module is located at decky-loader/plugin
# For easy intellisense checkout the decky-loader code repo
# and add the `decky-loader/plugin/imports` path to `python.analysis.extraPaths` in `.vscode/settings.json`
import decky
import asyncio

def _log_pipe(pipe, logger_fn):
    with pipe:
        for line in iter(pipe.readline, b''):
            logger_fn("[steaminputdb-buddy] %s", line.decode("utf-8", errors="replace").rstrip())

class Plugin:
    buddy_proc = None

    # Asyncio-compatible long-running code, executed in a task when the plugin is loaded
    async def _main(self):
        self.loop = asyncio.get_event_loop()
        decky.logger.info("SteamInputDB-Buddy-Decky starting...")
        bin_path = os.path.join(decky.DECKY_PLUGIN_DIR, "bin", "steaminputdb-buddy")
        os.chmod(bin_path, 0o755)
        decky.logger.info("Starting steaminputdb-buddy from %s", bin_path)
        self.buddy_proc = subprocess.Popen(
            [bin_path, "--tray-display=false", "--log-level=info", "--api-cors-origins=https://steamloopback.host,https://steaminputdb.com,https://www.steaminputdb.com"],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
        )
        threading.Thread(target=_log_pipe, args=(self.buddy_proc.stdout, decky.logger.info), daemon=True).start()
        threading.Thread(target=_log_pipe, args=(self.buddy_proc.stderr, decky.logger.error), daemon=True).start()
        decky.logger.info("steaminputdb-buddy started (PID %d)", self.buddy_proc.pid)

    # Function called first during the unload process, utilize this to handle your plugin being stopped, but not
    # completely removed
    async def _unload(self):
        decky.logger.info("SteamInputDB-Buddy-Decky stopping...")
        if self.buddy_proc is not None:
            decky.logger.info("Stopping steaminputdb-buddy (PID %d)...", self.buddy_proc.pid)
            self.buddy_proc.terminate()
            try:
                self.buddy_proc.wait(timeout=5)
            except subprocess.TimeoutExpired:
                decky.logger.warning("steaminputdb-buddy did not stop gracefully, killing...")
                self.buddy_proc.kill()
                self.buddy_proc.wait(timeout=5)
            self.buddy_proc = None

    # Function called after `_unload` during uninstall, utilize this to clean up processes and other remnants of your
    # plugin that may remain on the system
    async def _uninstall(self):
        if self.buddy_proc is not None:
            self.buddy_proc.terminate()
            try:
                self.buddy_proc.wait(timeout=5)
            except subprocess.TimeoutExpired:
                self.buddy_proc.kill()
                self.buddy_proc.wait(timeout=5)
            self.buddy_proc = None

    # Migrations that should be performed before entering `_main()`.
    async def _migration(self):
        # Nothing to be done, yet
        pass
