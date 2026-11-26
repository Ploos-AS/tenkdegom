import assert from 'node:assert/strict';
import { readFileSync, existsSync } from 'node:fs';
import { test } from 'node:test';

const read = (path) => readFileSync(path, 'utf8');
test('privacy check does not transmit or persist answers', () => {
 const page = read('src/pages/personvernsjekken/index.astro');
 assert.match(page, /SiteLayout/);
 for (const forbidden of [/fetch\s*\(/, /XMLHttpRequest/, /localStorage/, /sessionStorage/, /sendBeacon\s*\(/]) {
  assert.doesNotMatch(page, forbidden);
 }
});
test('site has privacy notice and shared footer link', () => {
 assert.ok(existsSync('src/pages/personvern.astro'));
 assert.match(read('src/layouts/SiteLayout.astro'), /href="\/personvern\/"/);
});
test('all article slugs are unique and article sources are declared', () => {
 const source = read('src/data/articles.ts');
 const slugs = [...source.matchAll(/"slug":\s*"([^"]+)"/g)].map(m=>m[1]);
 assert.ok(slugs.length >= 5, "Expected at least five articles");
 assert.equal(new Set(slugs).size, slugs.length);
 assert.equal((source.match(/"links":/g) || []).length, slugs.length);
});

test('help resources are linked from navigation, homepage and sitemap', () => {
 assert.ok(existsSync('src/pages/ressurser.astro'));
 for (const path of ['src/layouts/SiteLayout.astro', 'src/pages/index.astro', 'src/pages/sitemap.xml.ts']) {
  assert.ok(read(path).includes('/ressurser/'));
 }
 const page = read('src/pages/ressurser.astro');
 for (const url of ['https://nettvett.no/', 'https://slettmeg.no/', 'https://sikkert.no/']) assert.ok(page.includes(url));
});
