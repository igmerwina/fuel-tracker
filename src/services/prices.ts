import type { FuelBrand, FuelType } from '../types';

interface BackendFuelPrice {
  brand: FuelBrand;
  fuel_name: string;
  fuel_type: FuelType;
  price: number;
  location: string;
  effective_at?: string;
  source_url: string;
  source_status: string;
  scraped_at: string;
}

interface BackendPriceCache {
  generated_at: string;
  regions: string[];
  prices: BackendFuelPrice[];
  errors: string[];
}

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? 'http://127.0.0.1:8080';

export const fetchRealtimePrices = async (): Promise<BackendPriceCache | null> => {
  try {
    const response = await fetch(`${API_BASE_URL}/api/v1/prices`);
    if (!response.ok) return null;
    return (await response.json()) as BackendPriceCache;
  } catch {
    return null;
  }
};
