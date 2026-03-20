#!/usr/bin/env node

'use strict';

const { spawn } = require('child_process');
const path = require('path');
const fs = require('fs');

const PLATFORM_MAP = {
  darwin: 'darwin',
  linux: 'linux',
  win32: 'windows',
};

const ARCH_MAP = {
  x64: 'amd64',
  arm64: 'arm64',
};

function resolveBinaryPath() {
  const binName = process.platform === 'win32' ? 'gm.exe' : 'gm';
  const primaryPath = path.join(__dirname, '..', 'vendor', binName);
  if (fs.existsSync(primaryPath)) {
    return primaryPath;
  }

  const platform = PLATFORM_MAP[process.platform];
  const arch = ARCH_MAP[process.arch];
  if (!platform || !arch) {
    return primaryPath;
  }
  const packedName = `gm-${platform}-${arch}${process.platform === 'win32' ? '.exe' : ''}`;
  return path.join(__dirname, '..', 'vendor', packedName);
}

const binPath = resolveBinaryPath();

if (!fs.existsSync(binPath)) {
  console.error(
    `[gm] Binary not found at ${binPath}\n` +
    `[gm] Try reinstalling: npm install -g @limxdynamics/gm-cli`
  );
  process.exit(1);
}

const child = spawn(binPath, process.argv.slice(2), {
  stdio: 'inherit',
  env: process.env,
});

child.on('error', (err) => {
  console.error(`[gm] Failed to start: ${err.message}`);
  process.exit(1);
});

child.on('exit', (code) => {
  process.exit(code ?? 1);
});
