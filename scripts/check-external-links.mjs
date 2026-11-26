import { readFile, readdir } from 'node:fs/promises';
import { join } from 'node:path';
async function walk(dir) {
  const files = await readdir(dir, { withFileTypes: true });
  return (await Promise.all(files.map(f => f.isDirectory() ? walk(join(dir,f.name)) : [join(dir,f.name)]))).flat();
}
const urls = new Map();
for (const file of await walk('src')) {
  if (!/\.(astro|ts|md|mdx)$/.test(file)) continue;
  const source = await readFile(file,'utf8');
  for (const match of source.matchAll(/["'](https:\/\/[^"'\s<>]+)["']/g)) {
    const url = match[1];
    if (url.includes('${') || url.includes('{')) continue;
    if (!urls.has(url)) urls.set(url,new Set());
    urls.get(url).add(file);
  }
}
let dead = 0, uncertain = 0;
for (const [url, files] of urls) {
  try {
    const response = await fetch(url,{redirect:'follow',signal:AbortSignal.timeout(15000),headers:{'User-Agent':'Tenkdegom-LinkAudit/1.0'}});
    await response.body?.cancel();
    const status = response.status;
    const label = status === 404 || status === 410 ? 'DEAD' : response.ok ? 'OK' : 'REVIEW';
    if (label === 'DEAD') dead++;
    if (label === 'REVIEW') uncertain++;
    console.log(label,status,url,response.url !== url ? ' -> '+response.url : '',[...files].join(', '));
  } catch (error) {
    uncertain++;
    console.log('REVIEW',url,String(error),[...files].join(', '));
  }
}
console.log(JSON.stringify({total:urls.size,dead,uncertain}));
if (dead) process.exitCode = 1;
