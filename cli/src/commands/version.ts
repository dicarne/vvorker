import { Command } from "commander";
import * as fs from "fs-extra";
import * as path from "path";
import { withWorkingDir } from "../utils/working-dir";

interface Version {
  major: number;
  minor: number;
  patch: number;
}

function parseVersion(versionStr: string): Version {
  const parts = versionStr.split(".").map(Number);
  if (parts.length !== 3 || parts.some(isNaN)) {
    throw new Error("Invalid version format");
  }
  return { major: parts[0], minor: parts[1], patch: parts[2] };
}

function formatVersion(version: Version): string {
  return `${version.major}.${version.minor}.${version.patch}`;
}

function updateVersion(
  version: Version,
  options: { major?: boolean; minor?: boolean }
): Version {
  if (options.major) {
    return { major: version.major + 1, minor: 0, patch: 0 };
  } else if (options.minor) {
    return { major: version.major, minor: version.minor + 1, patch: 0 };
  } else {
    return {
      major: version.major,
      minor: version.minor,
      patch: version.patch + 1,
    };
  }
}

export const versionCommand = new Command("ver")
  .description("用于更新版本号")
  .option("-m, --major", "主版本号+1，次版本号和修订号归零")
  .option("-s, --minor", "次版本号+1，修订号归零")
  .action(async (options) => {
    await withWorkingDir(async () => {
      try {
        const cwd = process.cwd();
        const versionFile = path.join(cwd, "VERSION.txt");

        let version: Version;

        if (fs.existsSync(versionFile)) {
          const versionStr = fs.readFileSync(versionFile, "utf-8").trim();
          version = parseVersion(versionStr);
          console.log(`当前版本: ${formatVersion(version)}`);
        } else {
          version = { major: 0, minor: 0, patch: 0 };
          console.log("VERSION.txt 不存在，初始化版本: 0.0.0");
        }

        const newVersion = updateVersion(version, options);
        const newVersionStr = formatVersion(newVersion);

        fs.writeFileSync(versionFile, newVersionStr, "utf-8");

        const updateType = options.major
          ? "主版本号"
          : options.minor
          ? "次版本号"
          : "修订号";
        console.log(
          `${updateType}已更新: ${formatVersion(version)} -> ${newVersionStr}`
        );
        console.log(`版本号已写入: ${versionFile}`);
      } catch (error: any) {
        console.error("更新版本号失败:", error.message);
      }
    });
  });
