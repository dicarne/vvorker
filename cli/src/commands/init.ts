import { Command } from "commander";
import inquirer from "inquirer";
import * as fs from "fs-extra";
import * as path from "path";
import json5 from "json5";
import { runCommand } from "../utils/system";
import pc from "picocolors";
import { getEnv } from "../utils/config";
import { execSync } from "child_process";
import { checkForUpdate } from "../utils/update-check";

/**
 * 公共初始化函数，处理项目创建的通用逻辑
 * @param projectName 项目名称
 * @param jsonData vvorker配置数据
 * @param gitRepo Git仓库地址
 * @param modifyJsonData 可选的回调函数，用于修改jsonData
 */
async function initializeProject(
  projectName: string,
  jsonData: object,
  gitRepo: string,
  options?: { modifyJsonData?: (data: object) => void },
) {
  // 1. 克隆Git仓库并删除.git目录
  try {
    execSync(`git clone ${gitRepo} ${projectName}`, { stdio: "inherit" });
    await fs.remove(path.join(projectName, ".git"));
  } catch (error) {
    console.log(pc.red(`克隆 Git 仓库失败: ${error}`));
    throw error;
  }

  // 2. 查找并修改wrangler配置文件
  let wconfigPath = "wrangler.jsonc";
  if (!fs.existsSync(path.join(projectName, wconfigPath))) {
    wconfigPath = "wrangler.json";
  }

  const wranglerJsonPath = path.join(projectName, wconfigPath);
  const wranglerJson = json5.parse(
    await fs.readFile(wranglerJsonPath, "utf-8"),
  );
  wranglerJson.compatibility_flags = ["nodejs_compat"];
  wranglerJson.durable_objects = undefined;
  wranglerJson.migrations = undefined;
  wranglerJson.name = projectName;
  await fs.writeJson(wranglerJsonPath, wranglerJson, { spaces: 2 });

  // 3. 可选的jsonData修改（如Vue项目需要添加assets配置）
  if (options?.modifyJsonData) {
    options?.modifyJsonData(jsonData);
  }

  // 4. 写入vvorker配置文件
  const env = getEnv();
  const jsonFilePath = `vvorker.${env}.json`;
  await fs.writeJson(path.join(projectName, jsonFilePath), jsonData, {
    spaces: 2,
  });

  // 5. 修改package.json中的项目名称
  const packageJsonPath = path.join(projectName, "package.json");
  const packageJson = json5.parse(await fs.readFile(packageJsonPath, "utf-8"));
  packageJson.name = projectName;
  await fs.writeJson(packageJsonPath, packageJson, { spaces: 2 });

  console.log(pc.green(`项目 ${projectName} 初始化完成`));

  // 6. 安装依赖
  try {
    await runCommand("pnpm", ["install"], projectName);
  } catch (error) {
    console.log(pc.red(`安装依赖失败，请手动安装`));
  }

  // 7. 初始化Git、生成类型、运行cf-typegen
  try {
    await runCommand("git", ["init"], projectName);
    await runCommand("vvcli", ["types"], projectName);
    await runCommand("pnpm", ["run", "cf-typegen"], projectName);
  } catch (error) {
    console.log(pc.red(`生成类型提示失败，请手动运行 vvcli types 生成`));
  } finally {
    await runCommand("git", ["add", "*"], projectName);
    await runCommand(
      "git",
      ["commit", "-m", '"init: Create with vvcli."'],
      projectName,
    );
  }
}

async function createWorkerProject(
  projectName: string,
  jsonData: object,
  gitRepo: string,
) {
  await initializeProject(projectName, jsonData, gitRepo);
}

async function createVueProject(
  projectName: string,
  jsonData: object,
  gitRepo: string,
) {
  await initializeProject(projectName, jsonData, gitRepo, {
    modifyJsonData: (data) => {
      (data as any)["assets"] = [
        {
          directory: "./dist/client",
          binding: "ASSETS",
        },
      ];
    },
  });
}

export const initCommand = new Command("init")
  .command("init <projectName>")
  .description("初始化VVorker项目")
  .action(async (projectName, options) => {
    await checkForUpdate();
    
    // 交互式输入uid
    const { uid, projtype } = await inquirer.prompt([
      {
        type: "input",
        name: "uid",
        message: "请输入vvorker平台worker的uid:",
      },
      {
        type: "list",
        name: "projtype",
        message: "请选择工程类型",
        choices: [
          {
            name: "Worker",
            value: "worker",
            description: "Worker后端项目，适合复杂业务",
          },
          { name: "Vue", value: "vue", description: "Vue项目，前后端工程" },
          {
            name: "Simple",
            value: "simple",
            description: "简单Worker项目，不包含数据库",
          },
        ],
        default: "worker",
      },
    ]);
    let datatype = projtype;
    if (projtype === "simple") {
      datatype = "worker";
    }
    const jsonData = {
      $schema: "vvorker-schema.json",
      name: projectName,
      project: {
        uid: uid,
        type: datatype,
      },
      version: "0.0.0",
      services: [],
      vars: {},
      ai: [],
      oss: [],
      pgsql: [],
      mysql: [],
      kv: [],
      assets: [],
      compatibility_flags: ["nodejs_compat"],
    };

    switch (projtype) {
      case "worker": {
        await createWorkerProject(
          projectName,
          jsonData,
          "https://git.cloud.zhishudali.ink/template/vv-template-worker.git",
        );
        break;
      }
      case "vue": {
        await createVueProject(
          projectName,
          jsonData,
          "https://git.cloud.zhishudali.ink/template/vv-template-vue.git",
        );
        break;
      }
      case "simple": {
        await createWorkerProject(
          projectName,
          jsonData,
          "https://git.cloud.zhishudali.ink/template/vv-template-simple.git",
        );
        break;
      }
    }
  });
