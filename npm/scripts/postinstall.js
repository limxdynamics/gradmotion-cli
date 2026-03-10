#!/usr/bin/env node

'use strict';

const fs = require('fs');
const path = require('path');

const PLATFORM_MAP = {
  darwin: 'darwin',
  linux: 'linux',
  win32: 'windows',
};

const ARCH_MAP = {
  x64: 'amd64',
  arm64: 'arm64',
};

function getPackedBinaryName() {
  const platform = PLATFORM_MAP[process.platform];
  const arch = ARCH_MAP[process.arch];

  if (!platform || !arch) {
    throw new Error(
      `Unsupported platform: ${process.platform}/${process.arch}\n` +
      `Supported: darwin/linux/win32 + x64/arm64`
    );
  }

  const ext = process.platform === 'win32' ? '.exe' : '';
  return `gm-${platform}-${arch}${ext}`;
}

function main() {
  const packedBinaryName = getPackedBinaryName();
  const vendorDir = path.join(__dirname, '..', 'vendor');
  const destName = process.platform === 'win32' ? 'gm.exe' : 'gm';
  const packedPath = path.join(vendorDir, packedBinaryName);
  const destPath = path.join(vendorDir, destName);

  try {
    if (!fs.existsSync(packedPath)) {
      throw new Error(
        `Packed binary not found: ${packedPath}\n` +
        'This npm package may be incomplete. Please reinstall or contact maintainer.'
      );
    }
    fs.copyFileSync(packedPath, destPath);

    if (process.platform !== 'win32') {
      fs.chmodSync(destPath, 0o755);
    }

    console.log(`[gm] Installed to ${destPath}`);
  } catch (err) {
    console.error(`[gm] Failed to prepare local binary: ${err.message}`);
    process.exit(1);
  }
}

main();
