#!/usr/bin/env node
'use strict';

const fs = require('fs');
const https = require('https');
const path = require('path');
const os = require('os');
const { execSync } = require('child_process');

// Try to get version from package.json
let version;
try {
    const pkg = require('../package.json');
    version = pkg.version;
} catch (e) {
    console.error('Could not read package.json:', e.message);
    process.exit(1);
}

const platform = process.platform;
const arch = process.arch;

// Map Node OS to GOOS
let goos = '';
let ext = '';
if (platform === 'win32') {
    goos = 'windows';
    ext = '.exe';
} else if (platform === 'darwin') {
    goos = 'darwin';
} else if (platform === 'linux') {
    goos = 'linux';
} else {
    console.error(`Unsupported platform: ${platform}`);
    console.error('Please download the binary manually from:');
    console.error(`  https://github.com/loophole-ai/loophole-cli/releases`);
    process.exit(1);
}

// Map Node arch to GOARCH
let goarch = '';
if (arch === 'x64') {
    goarch = 'amd64';
} else if (arch === 'arm64') {
    goarch = 'arm64';
} else {
    console.error(`Unsupported architecture: ${arch}`);
    console.error('Please download the binary manually from:');
    console.error(`  https://github.com/loophole-ai/loophole-cli/releases`);
    process.exit(1);
}

// Binary name matches GoReleaser output: loophole-linux-amd64, loophole-windows-amd64.exe, etc.
const binName = `loophole-${goos}-${goarch}${ext}`;
const downloadUrl = `https://github.com/loophole-ai/loophole-cli/releases/download/v${version}/${binName}`;

const destDir = path.join(__dirname, '..', 'bin');
if (!fs.existsSync(destDir)) {
    fs.mkdirSync(destDir, { recursive: true });
}

const destFile = path.join(destDir, `loophole${ext}`);

console.log(`\nDownloading Loophole CLI v${version} for ${goos}/${goarch}...`);
console.log(`URL: ${downloadUrl}\n`);

/**
 * Follows redirects and returns the final URL.
 */
function resolveRedirect(url, maxRedirects = 10) {
    return new Promise((resolve, reject) => {
        if (maxRedirects === 0) return reject(new Error('Too many redirects'));
        https.get(url, (res) => {
            res.resume(); // Drain the response body
            if (res.statusCode === 301 || res.statusCode === 302) {
                return resolve(resolveRedirect(res.headers.location, maxRedirects - 1));
            }
            if (res.statusCode !== 200) {
                return reject(
                    new Error(
                        `HTTP ${res.statusCode} — binary not found.\n` +
                        `Make sure release v${version} exists at:\n` +
                        `  https://github.com/loophole-ai/loophole-cli/releases/tag/v${version}`
                    )
                );
            }
            resolve(url);
        }).on('error', reject);
    });
}

/**
 * Downloads a URL to a local file (no redirects — call resolveRedirect first).
 */
function downloadDirect(url, dest) {
    return new Promise((resolve, reject) => {
        const file = fs.createWriteStream(dest);
        https.get(url, (res) => {
            res.pipe(file);
            file.on('finish', () => {
                file.close();
                const stats = fs.statSync(dest);
                if (stats.size === 0) {
                    fs.unlinkSync(dest);
                    return reject(new Error('Downloaded binary is empty — the release asset may be missing.'));
                }
                resolve();
            });
        }).on('error', (err) => {
            file.close();
            if (fs.existsSync(dest)) fs.unlinkSync(dest);
            reject(err);
        });
        file.on('error', (err) => {
            if (fs.existsSync(dest)) fs.unlinkSync(dest);
            reject(err);
        });
    });
}


resolveRedirect(downloadUrl)
    .then((finalUrl) => downloadDirect(finalUrl, destFile))
    .then(() => {
        // Make executable on Unix
        if (platform !== 'win32') {
            fs.chmodSync(destFile, 0o755);
        }
        console.log(`✅  Loophole CLI installed to: ${destFile}`);
        console.log('    Run "loophole --version" to verify the installation.\n');
    })
    .catch((err) => {
        console.error(`\n❌  Installation failed: ${err.message}\n`);
        console.error('Manual install options:');
        console.error('  • Download from: https://github.com/loophole-ai/loophole-cli/releases');
        console.error(`  • Homebrew (macOS/Linux): brew install loophole-ai/tap/loophole`);
        process.exit(1);
    });

