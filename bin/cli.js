#!/usr/bin/env node
'use strict';

const { spawn } = require('child_process');
const path = require('path');
const fs = require('fs');

const isWindows = process.platform === 'win32';
const binaryName = isWindows ? 'loophole.exe' : 'loophole';
const binaryPath = path.join(__dirname, binaryName);

// Check if the binary was downloaded by the postinstall script
if (!fs.existsSync(binaryPath)) {
    console.error(
        '\n Loophole binary not found at: ' + binaryPath + '\n' +
        '   Try reinstalling: npm install -g @loophole-ai/loophole-cli\n' +
        '   Or download manually: https://github.com/loophole-ai/loophole-cli/releases\n'
    );
    process.exit(1);
}

const child = spawn(binaryPath, process.argv.slice(2), { stdio: 'inherit' });

child.on('exit', (code, signal) => {
    if (signal) {
        process.kill(process.pid, signal);
    } else {
        process.exit(code ?? 0);
    }
});

child.on('error', (err) => {
    if (err.code === 'EACCES') {
        console.error('\n Permission denied. Fixing permissions...');
        try {
            fs.chmodSync(binaryPath, 0o755);
            console.error('    Fixed! Please re-run your command.\n');
        } catch (_) {
            console.error(`    Could not fix automatically. Run: chmod +x ${binaryPath}\n`);
        }
    } else {
        console.error(`\n Failed to launch Loophole: ${err.message}\n`);
    }
    process.exit(1);
});
