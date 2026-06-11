import { readFile } from 'node:fs/promises';
import { join } from 'node:path';

const cachePath = join(process.cwd(), 'backend', 'data', 'fuel_prices.json');
const bpUrl = 'https://www.bp.com/id_id/indonesia/home/produk-dan-layanan/spbu/harga.html';

const setHeaders = (res) => {
  res.setHeader('Access-Control-Allow-Origin', '*');
  res.setHeader('Access-Control-Allow-Methods', 'POST,OPTIONS');
  res.setHeader('Access-Control-Allow-Headers', 'Content-Type');
};

const parseNumber = (value) => Number.parseInt(value.replace(/[^\d]/g, ''), 10);

const fetchBpPrices = async () => {
  const response = await fetch(`${bpUrl}?_=${Date.now()}`, {
    headers: {
      'User-Agent': 'Mozilla/5.0 FuelWatchJakarta/1.0',
      Accept: 'text/html,application/xhtml+xml',
      'Accept-Language': 'id-ID,id;q=0.9,en;q=0.8',
      'Cache-Control': 'no-cache',
      Pragma: 'no-cache',
    },
  });
  if (!response.ok) throw new Error(`bp_http_${response.status}`);

  const html = await response.text();
  const effective = html.match(/Harga berlaku efektif\s+([^<\n]+)/i)?.[1]?.trim() ?? '';
  const products = [
    ['BP Ultimate Diesel', 'DIESEL_CN_51'],
    ['BP Ultimate', 'RON_95'],
    ['BP 92', 'RON_92'],
  ];

  return products.flatMap(([fuelName, fuelType]) => {
    const escapedName = fuelName.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
    const match = html.match(new RegExp(`${escapedName}\\s+IDR\\s*([\\d,.]+)\\s+IDR\\s*([\\d,.]+)`, 'i'));
    if (!match) return [];
    return [{
      brand: 'BP',
      fuel_name: fuelName,
      fuel_type: fuelType,
      price: parseNumber(match[1]),
      location: 'Jabodetabek',
      effective_at: effective,
      source_url: bpUrl,
      source_status: 'scraped',
      scraped_at: new Date().toISOString(),
    }];
  });
};

const mergePrices = (cache, freshPrices) => {
  if (!freshPrices.length) return cache;
  const freshKeys = new Set(freshPrices.map((price) => `${price.brand}:${price.fuel_type}`));
  return {
    ...cache,
    generated_at: new Date().toISOString(),
    prices: [
      ...freshPrices,
      ...(cache.prices ?? []).filter((price) => !freshKeys.has(`${price.brand}:${price.fuel_type}`)),
    ],
    errors: cache.errors ?? [],
  };
};

export default async function handler(req, res) {
  setHeaders(res);
  if (req.method === 'OPTIONS') return res.status(204).end();
  if (req.method !== 'POST') return res.status(405).json({ error: 'method_not_allowed' });

  try {
    const cache = JSON.parse(await readFile(cachePath, 'utf8'));
    const bpPrices = await fetchBpPrices();
    return res.status(200).json(mergePrices(cache, bpPrices));
  } catch {
    try {
      const cache = JSON.parse(await readFile(cachePath, 'utf8'));
      return res.status(200).json({
        ...cache,
        generated_at: new Date().toISOString(),
        errors: [...(cache.errors ?? []), 'bp_live_fetch_failed'],
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
}
