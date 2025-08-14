import { Command } from 'commander';
import * as fs from 'fs-extra';
import * as path from 'path';
import { config, getToken, getUrl, setUrl } from '../utils/config';
import { loadVVorkerConfig } from '../utils/vvorker-config';
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

    let up1 = await apiClient.post(`/api/ext/assets/clear-assets`, {
      worker_uid: uid,
    })
    if (up1.data.delete_count > 0) {
      console.log(pc.green("✓") + pc.gray(` 已清除 ${up1.data.delete_count} 个 Assets 文件，又为数据库腾出了空间！`))
    }

    let jsFilePath = "";

    if (vvorkerJson.assets && vvorkerJson.assets.length > 0) {
      console.log(pc.white("Assets 文件上传开始..."))
      let wwwAssetsPath = path.join(process.cwd(), vvorkerJson.assets[0].directory)
      // walk wwwAssetsPath，调用接口上传每一个文件
      const walk = async (dir: string) => {
        const files = fs.readdirSync(dir);
        for (let file of files) {
          const filePath = path.join(dir, file);
          const stat = fs.statSync(filePath);
          if (stat.isDirectory()) {
            await walk(filePath);
          } else {
            const fileContent = fs.readFileSync(filePath);
            const fileUrl = filePath.replace(wwwAssetsPath, '').replace(/\\/g, '/');
            console.log(pc.gray(fileUrl));

            let up1 = await apiClient.post(`/api/file/upload`, {
              file: fileContent.toString('base64'),
              path: fileUrl,
            }, {
              headers: {
                'Content-Type': 'application/json',
                'x-encrypted-data': vk != "" ? vk : undefined
              },
            })

            if (up1.status !== 200) {
              throw new Error(`上传失败：${fileUrl} ${up1.status} ${up1.statusText}`);
            }

            let fileuid = up1.data.data.fileId;

            let resp = await apiClient.post(`/api/ext/assets/create-assets`, {
              uid: fileuid,
              "worker_uid": uid,
              "path": fileUrl,
            })

            if (resp.status != 200) {
              console.log(pc.red("✗") + pc.gray(`${fileUrl} ${resp.status} ${resp.statusText}`));
              throw new Error(`上传失败：${fileUrl} ${resp.status} ${resp.statusText}`);
            }
          }
        }
      }
      await walk(wwwAssetsPath);
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
