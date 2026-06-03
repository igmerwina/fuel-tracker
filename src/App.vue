<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import Header from './components/Header.vue';
import HeroSummary from './components/HeroSummary.vue';
import FilterBar from './components/FilterBar.vue';
import StationList from './components/StationList.vue';
import Map from './components/Map.vue';
import type {
  FilterState,
  FuelBrand,
  FuelStation,
  FuelType,
  SortMode,
  StationResult,
} from './types/index';
import { fuelStations, fuelTypes, JAKARTA_CENTER } from './data/stations';
import { fetchRealtimePrices } from './services/prices';

const stations = ref<FuelStation[]>(fuelStations);
const favoriteIds = ref<string[]>([]);
const sortMode = ref<SortMode>('cheapest');
const filterState = ref<FilterState>({
  query: '',
  region: '',
  brand: '',
  fuelType: 'RON_92',
  openNow: false,
});

const mapRef = ref<InstanceType<typeof Map>>();

const fuelAliases: Record<string, string[]> = {
  RON_90: ['ron 90', 'pertalite', 'revvo 90'],
  RON_92: ['ron 92', 'pertamax', 'shell super', 'bp 92', 'revvo 92'],
  RON_95: ['ron 95', 'pertamax green', 'v-power', 'bp ultimate', 'revvo 95'],
  RON_98: ['ron 98', 'pertamax turbo', 'nitro'],
  DIESEL_CN_51: ['diesel cn 51', 'dexlite', 'v-power diesel'],
  DIESEL_CN_53: ['diesel cn 53', 'pertamina dex'],
};

const fuelDisplayName = computed(() => {
  if (filterState.value.fuelType === 'RON_92') return 'Pertamax';
  return fuelTypes.find((fuel) => fuel.id === filterState.value.fuelType)?.label ?? 'BBM';
});

const distanceKm = (station: FuelStation) => {
  const toRad = (value: number) => (value * Math.PI) / 180;
  const earthKm = 6371;
  const dLat = toRad(station.latitude - JAKARTA_CENTER.lat);
  const dLng = toRad(station.longitude - JAKARTA_CENTER.lng);
  const lat1 = toRad(JAKARTA_CENTER.lat);
  const lat2 = toRad(station.latitude);
  const a =
    Math.sin(dLat / 2) ** 2 + Math.cos(lat1) * Math.cos(lat2) * Math.sin(dLng / 2) ** 2;
  return earthKm * 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
};

const priceForFuel = (station: FuelStation, fuelType: FuelType) => {
  return station.prices.find((price) => price.type === fuelType)?.price ?? station.prices[0]?.price ?? 0;
};

const averageSelectedPrice = computed(() => {
  const prices = stations.value
    .map((station) => priceForFuel(station, filterState.value.fuelType || 'RON_92'))
    .filter(Boolean);
  return prices.length
    ? Math.round(prices.reduce((sum, price) => sum + price, 0) / prices.length)
    : 0;
});

const cheapestSelectedPrice = computed(() => {
  const prices = stations.value
    .map((station) => priceForFuel(station, filterState.value.fuelType || 'RON_92'))
    .filter(Boolean);
  return prices.length ? Math.min(...prices) : 0;
});

const stationResults = computed<StationResult[]>(() => {
  const activeFuel = filterState.value.fuelType || 'RON_92';
  const average = averageSelectedPrice.value;
  const cheapest = cheapestSelectedPrice.value;

  return stations.value
    .filter((station) => {
      const query = filterState.value.query.trim().toLowerCase();
      if (query) {
        const haystack = [
          station.name,
          station.brand,
          station.region,
          station.address,
          ...station.prices.flatMap((price) => [
            price.type.replaceAll('_', ' '),
            fuelTypes.find((fuel) => fuel.id === price.type)?.label ?? '',
            ...(fuelAliases[price.type] ?? []),
          ]),
        ]
          .join(' ')
          .toLowerCase();
        if (!haystack.includes(query)) return false;
      }
      if (filterState.value.region && station.region !== filterState.value.region) return false;
      if (filterState.value.brand && station.brand !== filterState.value.brand) return false;
      if (!station.prices.some((price) => price.type === activeFuel)) return false;
      return true;
    })
    .map((station) => {
      const distance = distanceKm(station);
      const selectedPrice = priceForFuel(station, activeFuel);
      return {
        ...station,
        distanceKm: distance,
        travelMinutes: Math.max(4, Math.round((distance / 24) * 60)),
        selectedPrice,
        averagePrice: average,
        savingsPerLiter: Math.max(0, average - selectedPrice),
        isCheapest: selectedPrice === cheapest,
        isFavorite: favoriteIds.value.includes(station.id),
      };
    })
    .sort((a, b) => {
      if (sortMode.value === 'nearest') return a.distanceKm - b.distanceKm;
      return a.selectedPrice - b.selectedPrice || a.distanceKm - b.distanceKm;
    });
});

const cheapestStation = computed(() => stationResults.value.find((station) => station.isCheapest));
const nearestStation = computed(() => [...stationResults.value].sort((a, b) => a.distanceKm - b.distanceKm)[0]);

const lastUpdated = computed(() => {
  const latest = stations.value
    .map((station) => new Date(station.lastUpdated).getTime())
    .sort((a, b) => b - a)[0];
  return latest ? new Date(latest) : null;
});

const relativeUpdated = computed(() => {
  if (!lastUpdated.value) return '-';
  const minutes = Math.max(1, Math.round((Date.now() - lastUpdated.value.getTime()) / 60000));
  if (minutes < 60) return `${minutes} menit lalu`;
  const hours = Math.round(minutes / 60);
  return `${hours} jam lalu`;
});

const updateFilters = (next: Partial<FilterState>) => {
  filterState.value = { ...filterState.value, ...next };
};

const updateQuery = (query: string) => {
  updateFilters({ query });
};

const handleFlyTo = (station: StationResult) => {
  mapRef.value?.flyToStation(station);
};

const toggleFavorite = (stationId: string) => {
  favoriteIds.value = favoriteIds.value.includes(stationId)
    ? favoriteIds.value.filter((id) => id !== stationId)
    : [...favoriteIds.value, stationId];
  localStorage.setItem('favoriteStations', JSON.stringify(favoriteIds.value));
};

const applyRealtimePrices = async () => {
  const cache = await fetchRealtimePrices();
  if (!cache?.prices.length) return;

  const priceByBrandAndType = new globalThis.Map<string, number>();
  cache.prices.forEach((price) => {
    const key = `${price.brand}:${price.fuel_type}`;
    if (!priceByBrandAndType.has(key)) {
      priceByBrandAndType.set(key, price.price);
    }
  });

  stations.value = fuelStations.map((station) => ({
    ...station,
    prices: station.prices.map((price) => ({
      ...price,
      price: priceByBrandAndType.get(`${station.brand as FuelBrand}:${price.type as FuelType}`) ?? price.price,
    })),
    lastUpdated: cache.generated_at || station.lastUpdated,
  }));
};

onMounted(() => {
  favoriteIds.value = JSON.parse(localStorage.getItem('favoriteStations') || '[]') as string[];
  void applyRealtimePrices();
});
</script>

<template>
  <div class="min-h-screen bg-slate-100 text-slate-950">
    <Header
      :query="filterState.query"
      :station-count="stationResults.length"
      :last-updated="relativeUpdated"
      @update-query="updateQuery"
    />

    <main class="mx-auto flex h-[calc(100vh-72px)] max-w-[1800px] flex-col gap-4 p-3 md:p-4">
      <HeroSummary
        :fuel-name="fuelDisplayName"
        :cheapest-station="cheapestStation"
        :nearest-station="nearestStation"
        :last-updated="relativeUpdated"
      />

      <FilterBar
        :filters="filterState"
        :sort-mode="sortMode"
        @update-filters="updateFilters"
        @update-sort="sortMode = $event"
      />

      <section class="grid min-h-0 flex-1 grid-cols-1 gap-4 lg:grid-cols-[minmax(0,7fr)_minmax(360px,3fr)]">
        <div class="min-h-[520px] overflow-hidden rounded-[12px] bg-white shadow-xl shadow-slate-200/80 lg:min-h-0">
          <Map
            ref="mapRef"
            :stations="stationResults"
            :active-fuel="filterState.fuelType || 'RON_92'"
            :average-price="averageSelectedPrice"
          />
        </div>

        <StationList
          :stations="stationResults"
          :active-fuel="filterState.fuelType || 'RON_92'"
          :sort-mode="sortMode"
          @flyto="handleFlyTo"
          @toggle-favorite="toggleFavorite"
        />
      </section>
    </main>
  </div>
</template>
