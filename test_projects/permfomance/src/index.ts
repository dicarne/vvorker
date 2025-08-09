
import { Hono } from "hono";
import { EnvBinding } from "./binding";
import { init, useDebugEndpoint } from "@dicarne/vvorker-sdk";
import { env } from "cloudflare:workers";
init(env)

const app = new Hono<{ Bindings: EnvBinding }>();
useDebugEndpoint(app)

/**
 * 经典耗时算法：埃拉托色尼筛法
 * @param {number} limit 查找的上限
 * @returns {number[]}   2~limit 之间的所有素数
 */
function sieve(limit: number) {
  if (limit < 2) return [];

  // 初始化布尔数组：true 表示“可能是素数”
  const isPrime = new Array(limit + 1).fill(true);
  isPrime[0] = isPrime[1] = false;

  const sqrt = Math.floor(Math.sqrt(limit));
  for (let i = 2; i <= sqrt; i++) {
    if (isPrime[i]) {
      // 从 i*i 开始，把 i 的倍数全部标记为非素数
      for (let j = i * i; j <= limit; j += i) {
        isPrime[j] = false;
      }
    }
  }

  // 收集结果
  const primes: number[] = [];
  for (let i = 2; i <= limit; i++) {
    if (isPrime[i]) primes.push(i);
  }
  return primes;
}

app.get("/", async (c) => {
  const { limit } = c.req.query()
  return c.text(String(sieve(Number(limit ?? "10000")).length))
})

export default app;
