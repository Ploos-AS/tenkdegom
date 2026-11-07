import type { APIRoute } from 'astro';
import { articles } from '../data/articles';
export const GET: APIRoute = ({ site }) => {
  if (!site) throw new Error('Site URL must be configured');
  const paths = ['/', '/laer/', '/personvernsjekken/', '/om/', '/personvern/', ...articles.map(a => `/laer/${a.slug}/`)];
  const urls = paths.map(path => `<url><loc>${new URL(path, site).href}</loc></url>`).join('');
  return new Response(`<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">${urls}</urlset>`, { headers: { 'Content-Type': 'application/xml; charset=utf-8' } });
};
