import { readFile } from 'node:fs/promises';
import { join } from 'node:path';

const cachePath = join(process.cwd(), 'backend', 'data', 'fuel_prices.json');

const setHeaders = (res) => {
  res.setHeader('Access-Control-Allow-Origin', '*');
  res.setHeader('Access-Control-Allow-Methods', 'POST,OPTIONS');
  res.setHeader('Access-Control-Allow-Headers', 'Content-Type');
};

export default async function handler(req, res) {
  setHeaders(res);
  if (req.method === 'OPTIONS') return res.status(204).end();
  if (req.method !== 'POST') return res.status(405).json({ error: 'method_not_allowed' });

  try {
    const cache = JSON.parse(await readFile(cachePath, 'utf8'));
    return res.status(200).json({
      ...cache,
      generated_at: new Date().toISOString(),
      errors: cache.errors ?? [],
    });
  } catch {
    return res.status(500).json({
      generated_at: new Date().toISOString(),
      regions: [],
      prices: [],
      errors: ['price_cache_unavailable'],
    });
  }
}
