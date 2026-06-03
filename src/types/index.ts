export type FuelBrand = 'Pertamina' | 'Shell' | 'Vivo' | 'BP';
export type FuelType = 'RON_90' | 'RON_92' | 'RON_95' | 'RON_98' | 'DIESEL_CN_51' | 'DIESEL_CN_53';
export type Region = 'Jakarta Pusat' | 'Jakarta Selatan' | 'Jakarta Barat' | 'Jakarta Timur' | 'Jakarta Utara';

export interface FuelPrice {
  type: FuelType;
  price: number;
}

export interface FuelStation {
  id: string;
  name: string;
  brand: FuelBrand;
  region: Region;
  address: string;
  latitude: number;
  longitude: number;
  prices: FuelPrice[];
  lastUpdated: string;
}

export interface FilterState {
  region: Region | '';
  brand: FuelBrand | '';
  fuelType: FuelType | '';
}
