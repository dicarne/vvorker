import { Command } from 'commander';
import * as fs from 'fs-extra';
import * as path from 'path';
import { createHash } from 'crypto';
import { config, getToken, getUrl, setUrl } from '../utils/config';
import { loadVVorkerConfig, saveVVorkerConfig } from '../utils/vvorker-config';
import { runCommand } from '../utils/system';
import { apiClient, requireOTP } from '../utils/api';
import pc from "picocolors"

export const deployCommand = new Command('deploy')
  .description('部署到VVorker')
  .action(async () => {
    if (!getUrl()) {
      console.error('请先配置VVorker平台的url');
      return;
    }
    if (!getToken()) {
      console.error('请先配置VVorker平台的token');
      return;
    }
    console.log(`环境：${pc.yellow(config.current_env)}`)
    await requireOTP()

    // 读取当前目录下的 vvorker.json 文件
    const vvorkerJson = loadVVorkerConfig();
    const packageJson = await fs.readJson('package.json');
    let serviceName = vvorkerJson.name;
    if (!serviceName) {
      console.error('服务名称不能为空');
      return;
    }

    // 读取 VERSION.txt 中的版本号
    const versionFile = path.join(process.cwd(), 'VERSION.txt');
    if (fs.existsSync(versionFile)) {
      const version = fs.readFileSync(versionFile, 'utf-8').trim();
      // 读取当前目录下的 vvorker.json 文件
      // 更新 version 字段
      vvorkerJson.version = version;
      // 保存回配置文件
      saveVVorkerConfig(vvorkerJson);
      console.log(pc.green(`✓ 版本号已更新: ${version}`));
    }

    const uid = vvorkerJson.project?.uid ?? vvorkerJson.uid;
    if (!uid) {
      console.error('uid不能为空');
      return;
    }

    const token = getToken();
    if (!token) {
      console.error('token不能为空');
      console.error('请先配置token');
      console.error('执行命令：vvcli config set token <token>');
      return;
    }
    // url join
    if (getUrl()!.endsWith('/')) {
      setUrl(getUrl()!.slice(0, -1));
    }
    // 运行 pnpm run build 命令
    // 如果存在pnpm-lock.yaml，则使用pnpm
    if (fs.existsSync('pnpm-lock.yaml')) {
      await runCommand('pnpm', ['run', 'build']);
    } else {
      await runCommand('npm', ['run', 'build']);
    }

    console.log(pc.gray("--------------"))

    const userinfo = await apiClient.get(`/api/user/info`)
    const vk = userinfo.data?.data?.vk ?? ""

    let jsFilePath = "";

    if (vvorkerJson.assets && vvorkerJson.assets.length > 0) {
      console.log(pc.white("Assets 文件上传开始..."))
      let wwwAssetsPath = path.join(process.cwd(), vvorkerJson.assets[0].directory)
      
      // 收集所有文件信息和哈希值
      const fileInfos: Array<{path: string, hash: string, fullPath: string}> = [];
      const collectFiles = (dir: string) => {
        const files = fs.readdirSync(dir);
        for (let file of files) {
          const filePath = path.join(dir, file);
          const stat = fs.statSync(filePath);
          if (stat.isDirectory()) {
            collectFiles(filePath);
          } else {
            const fileUrl = filePath.replace(wwwAssetsPath, '').replace(/\\/g, '/');
            const fileContent = fs.readFileSync(filePath);
            const hash = createHash('sha256').update(fileContent).digest('hex');
            fileInfos.push({
              path: fileUrl,
              hash: hash,
              fullPath: filePath
            });
          }
        }
      }
      collectFiles(wwwAssetsPath);
      
      console.log(pc.gray(`发现 ${fileInfos.length} 个文件，检查更新中...`));
      
      // 调用服务器接口检查哪些文件需要更新
      let checkResp = await apiClient.post(`/api/ext/assets/check-assets`, {
        worker_uid: uid,
        files: fileInfos.map(f => ({ path: f.path, hash: f.hash }))
      })
      
      if (checkResp.status !== 200) {
        throw new Error(`检查文件更新失败: ${checkResp.status} ${checkResp.statusText}`);
      }
      
      const { needUpload, needDelete } = checkResp.data.data;
      
      // 删除不再需要的文件
      if (needDelete.length > 0) {
        console.log(pc.gray(`正在删除 ${needDelete.length} 个已不存在的 Assets 文件...`));
        for (const filePath of needDelete) {
          let deleteResp = await apiClient.post(`/api/ext/assets/delete-assets`, {
            worker_uid: uid,
            path: filePath,
          });
          
          if (deleteResp.status !== 200) {
            console.log(pc.yellow("⚠") + pc.gray(` 删除文件失败：${filePath} ${deleteResp.status} ${deleteResp.statusText}`));
          }
        }
        console.log(pc.green("✓") + pc.gray(` 已删除 ${needDelete.length} 个 Assets 文件`));
      }
      
      if (needUpload.length === 0) {
        console.log(pc.green("✓") + pc.gray(" 所有 Assets 文件都是最新的，无需上传"));
      } else {
        console.log(pc.gray(`需要上传/更新 ${needUpload.length} 个文件...`));
        
        // 只上传需要更新的文件
        for (const fileInfo of fileInfos) {
          if (needUpload.includes(fileInfo.path)) {
            console.log(pc.gray(fileInfo.path));
            const fileContent = fs.readFileSync(fileInfo.fullPath);
            
            let uploadResp = await apiClient.post(`/api/file/upload`, {
              file: fileContent.toString('base64'),
              path: fileInfo.path,
            }, {
              headers: {
                'Content-Type': 'application/json',
                'x-encrypted-data': vk != "" ? vk : undefined
              },
            })

            if (uploadResp.status !== 200) {
              throw new Error(`上传失败：${fileInfo.path} ${uploadResp.status} ${uploadResp.statusText}`);
            }

            let fileuid = uploadResp.data.data.fileId;

            let createResp = await apiClient.post(`/api/ext/assets/create-assets`, {
              uid: fileuid,
              "worker_uid": uid,
              "path": fileInfo.path,
            })

            if (createResp.status != 200) {
              console.log(pc.red("✗") + pc.gray(`${fileInfo.path} ${createResp.status} ${createResp.statusText}`));
              throw new Error(`上传失败：${fileInfo.path} ${createResp.status} ${createResp.statusText}`);
            }
          }
        }
        
        console.log(pc.green("✓") + pc.gray(` Assets 文件上传完成，共 ${needUpload.length} 个文件`));
      }
    }

    if (vvorkerJson.project.type === "vue") {
      jsFilePath = `${process.cwd()}/dist/${packageJson.name.replaceAll("-", "_")}/index.js`;
    } else {
      jsFilePath = `${process.cwd()}/dist/index.js`;
      if (!fs.existsSync(jsFilePath)) {
        jsFilePath = `${process.cwd()}/dist/${packageJson.name.replaceAll("-", "_")}/index.js`;
      }
    }

    const distFilePath = `${process.cwd()}/dist/vvorker.json`;
    await fs.writeJson(distFilePath, vvorkerJson, { spaces: 2 });

    const jsContent = await fs.readFile(jsFilePath, 'utf-8');
    const jsBase64 = Buffer.from(jsContent).toString('base64');

    let resp = await apiClient.post(`/api/worker/v2/get-worker`, {
      uid: uid,
    })
    let prev = resp.data.data[0]

    const distVvorkerJson = await fs.readJson(`${process.cwd()}/dist/vvorker.json`);
    if (distVvorkerJson.pgsql && distVvorkerJson.pgsql.length > 0) {
      console.log("开始迁移PostgreSQL数据库...")
      for (let pgsql of distVvorkerJson.pgsql) {
        let rid = pgsql.resource_id;
        if (pgsql.migrate) {
          let migrateFiles = fs.readdirSync(pgsql.migrate);
          let allFile: {
            name: string,
            content: string
          }[] = []
          for (let migrateFile of migrateFiles) {
            if (!migrateFile.endsWith('.sql')) {
              continue;
            }
            let migrateFilePath = path.join(pgsql.migrate, migrateFile);
            let migrateFileContent = fs.readFileSync(migrateFilePath, 'utf-8');
            allFile.push({
              name: migrateFile,
              content: migrateFileContent.replaceAll("--> statement-breakpoint", "\n")
            })
          }

          let resp = await apiClient.post(`/api/ext/pgsql/migrate`, {
            resource_id: rid ?? ("worker_resource:pgsql:" + uid + ":" + pgsql.migrate),
            files: allFile,
            custom_db_name: pgsql.database,
            custom_db_user: pgsql.user,
            custom_db_host: pgsql.host,
            custom_db_port: pgsql.port,
            custom_db_password: pgsql.password,
          }, {
            headers: {
              'Content-Type': 'application/json',
              'x-encrypted-data': vk != "" ? vk : undefined
            }
          })
          if (resp.data.code !== 0) {
            console.log(`迁移失败：${resp.data}`);
            throw new Error(`迁移失败：${resp.data}`);
          }
          console.log(`迁移成功：${pgsql.migrate}`);
        }
      }
    }

    if (distVvorkerJson.mysql && distVvorkerJson.mysql.length > 0) {
      console.log("开始迁移MYSQL数据库...")
      for (let mysql of distVvorkerJson.mysql) {
        let rid = mysql.resource_id;
        if (mysql.migrate) {
          let migrateFiles = fs.readdirSync(mysql.migrate);
          let allFile: {
            name: string,
            content: string
          }[] = []
          for (let migrateFile of migrateFiles) {
            if (!migrateFile.endsWith('.sql')) {
              continue;
            }
            let migrateFilePath = path.join(mysql.migrate, migrateFile);
            let migrateFileContent = fs.readFileSync(migrateFilePath, 'utf-8');
            allFile.push({
              name: migrateFile,
              content: migrateFileContent.replaceAll("--> statement-breakpoint", "\n")
            })
          }

          let resp = await apiClient.post(`/api/ext/mysql/migrate`, {
            resource_id: rid ?? ("worker_resource:mysql:" + uid + ":" + mysql.migrate),
            files: allFile,
            custom_db_name: mysql.database,
            custom_db_user: mysql.user,
            custom_db_host: mysql.host,
            custom_db_port: mysql.port,
            custom_db_password: mysql.password,
          }, {
            headers: {
              'Content-Type': 'application/json',
              'x-encrypted-data': vk != "" ? vk : undefined
            }
          })
          if (resp.data.code !== 0) {
            console.log(`迁移失败：${resp.data}`);
            throw new Error(`迁移失败：${resp.data}`);
          }
          console.log(`迁移成功：${mysql.migrate}`);
        }
      }
    }

    prev.Code = jsBase64;
    prev.Template = JSON.stringify(distVvorkerJson);
    prev.HostName = undefined;
    prev.ExternalPath = undefined;
    prev.TunnelID = undefined;


    resp = await apiClient.post(`/api/worker/v2/update-worker`, prev, {
      headers: {
        'Content-Type': 'application/json',
        'x-encrypted-data': vk != "" ? vk : undefined
      }
    })

    if (resp.data.code === 0) {
      console.log(pc.green("✓ 部署成功！"));
    } else {
      console.log(pc.red("✗ 部署失败！"));
      throw new Error(`部署失败：${resp.data.message}`);
    }
  });
