import test from 'node:test';
import assert from 'node:assert/strict';
import { readFileSync } from 'node:fs';

const articles = JSON.parse(readFileSync(new URL('../src/data/articles.ts', import.meta.url), 'utf8')
  .replace(/^export const articles = /, '').replace(/\s+as const;\s*$/, ''));
const related = JSON.parse(readFileSync(new URL('../src/data/related-articles.ts', import.meta.url), 'utf8')
  .replace(/^\/\/[^\n]*\n/, '').replace(/^export const relatedArticles: Record<string, readonly string\[\]> = /, '').replace(/;\s*$/, ''));
const slugs = new Set(articles.map(a => a.slug));
const nonempty = value => typeof value === 'string' && value.trim().length > 0;
function validDate(value) {
  if (typeof value !== 'string' || !/^\d{4}-\d{2}-\d{2}$/.test(value)) return false;
  const date = new Date(value + 'T00:00:00Z');
  return Number.isFinite(date.getTime()) && date.toISOString().slice(0,10) === value;
}
test('article identities and editorial fields', () => {
  assert.equal(slugs.size, articles.length, 'Duplicate article slug');
  for (const a of articles) {
    assert.match(a.slug, /^[a-z0-9]+(?:-[a-z0-9]+)*$/);
    assert.ok(nonempty(a.title) && nonempty(a.intro), a.slug);
    assert.ok(validDate(a.publishedAt) && validDate(a.updatedAt), a.slug);
    assert.ok(a.updatedAt >= a.publishedAt, a.slug);
    assert.ok(Array.isArray(a.sections) && a.sections.length > 0, a.slug);
    for (const section of a.sections) {
      assert.ok(Array.isArray(section) && section.length === 2 && section.every(nonempty), a.slug);
    }
  }
});
test('article links have valid destinations and labels', () => {
  for (const a of articles) {
    assert.ok(Array.isArray(a.links), a.slug);
    for (const link of a.links) {
      assert.ok(Array.isArray(link) && link.length === 2 && link.every(nonempty), a.slug);
      const [, url] = link;
      if (url.startsWith('/laer/')) {
        assert.ok(slugs.has(url.slice(6,-1)) && url.endsWith('/'), a.slug + ': ' + url);
      } else if (url.startsWith('/') && !url.startsWith('//')) {
        assert.ok(!url.includes('\\\\'), url);
      } else {
        assert.equal(new URL(url).protocol, 'https:', url);
      }
    }
  }
});
test('related articles cover every article without invalid or duplicate targets', () => {
  assert.deepEqual(new Set(Object.keys(related)), slugs);
  for (const [slug, targets] of Object.entries(related)) {
    assert.ok(Array.isArray(targets) && targets.length > 0, slug);
    assert.equal(new Set(targets).size, targets.length, slug);
    for (const target of targets) {
      assert.ok(slugs.has(target) && target !== slug, slug + ': ' + target);
    }
  }
});
