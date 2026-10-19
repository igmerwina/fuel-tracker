<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import Header from './components/Header.vue';
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
import { refreshRealtimePrices } from './services/prices';

const stations = ref<FuelStation[]>(fuelStations);
const favoriteIds = ref<string[]>([]);
const sortMode = ref<SortMode>('cheapest');
const selectedStationId = ref<string>('');
const showCheapestCard = ref(true);
const userLocation = ref(JAKARTA_CENTER);
const locationStatus = ref<'idle' | 'ok' | 'outside' | 'denied'>('idle');
const isRefreshingPrices = ref(false);
const isInitialPriceLoading = ref(true);
const priceToast = ref<{ type: 'success' | 'error'; message: string } | null>(null);
const filterState = ref<FilterState>({
  query: '',
  region: '',
  brand: 'Pertamina',
  fuelType: 'RON_92',
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

const distanceKm = (station: FuelStation) => {
  return distanceBetweenKm(userLocation.value.lat, userLocation.value.lng, station.latitude, station.longitude);
};

const distanceBetweenKm = (fromLat: number, fromLng: number, toLat: number, toLng: number) => {
  const toRad = (value: number) => (value * Math.PI) / 180;
  const earthKm = 6371;
  const dLat = toRad(toLat - fromLat);
  const dLng = toRad(toLng - fromLng);
  const lat1 = toRad(fromLat);
  const lat2 = toRad(toLat);
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

const stationResults = computed<StationResult[]>(() => {
  const activeFuel = filterState.value.fuelType || 'RON_92';
  const average = averageSelectedPrice.value;
  const candidates = stations.value
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
      if (!station.prices.some((price) => price.type === activeFuel && price.price > 0)) return false;
      return true;
    });
  const cheapest = candidates
    .map((station) => priceForFuel(station, activeFuel))
    .filter(Boolean)
    .sort((a, b) => a - b)[0] ?? 0;
  const samePriceCount = candidates.filter((station) => priceForFuel(station, activeFuel) === cheapest).length;
  const cheapestId = candidates
    .map((station) => ({ id: station.id, price: priceForFuel(station, activeFuel), distance: distanceKm(station) }))
    .sort((a, b) => a.price - b.price || a.distance - b.distance)[0]?.id;

  return candidates
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
        isCheapest: station.id === cheapestId,
        hasSamePrice: samePriceCount > 1 && selectedPrice === cheapest,
        isFavorite: favoriteIds.value.includes(station.id),
      };
    })
    .sort((a, b) => {
      if (sortMode.value === 'nearest') return a.distanceKm - b.distanceKm;
      return a.selectedPrice - b.selectedPrice || a.distanceKm - b.distanceKm;
    });
});

const cheapestStation = computed(() => stationResults.value.find((station) => station.isCheapest));

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
  if (hours < 24) return `${hours} jam lalu`;
  return lastUpdated.value.toLocaleDateString('id-ID', { day: 'numeric', month: 'short' });
});

const updateFilters = (next: Partial<FilterState>) => {
  filterState.value = { ...filterState.value, ...next };
};

const updateQuery = (query: string) => {
  updateFilters({ query });
};

const activeFuelLabel = computed(() => fuelTypes.find((fuel) => fuel.id === filterState.value.fuelType)?.label ?? 'RON 92');

const locationNote = computed(() => {
  if (locationStatus.value === 'ok') return 'Jarak dari lokasimu';
  if (locationStatus.value === 'outside') return 'Kamu di luar Jakarta, jarak dari Monas';
  return 'Jarak dari Monas';
});

const resetFilters = () => {
  filterState.value = { query: '', region: '', brand: 'Pertamina', fuelType: 'RON_92' };
  sortMode.value = 'cheapest';
  showCheapestCard.value = true;
};

const handleFlyTo = (station: StationResult) => {
  selectedStationId.value = station.id;
  mapRef.value?.flyToStation(station);
};

const selectStation = (station: StationResult | null) => {
  selectedStationId.value = station?.id ?? '';
};

const startNavigation = (station: StationResult) => {
  window.open(station.googleMapsUrl, '_blank', 'noopener,noreferrer');
};

const toggleFavorite = (stationId: string) => {
  favoriteIds.value = favoriteIds.value.includes(stationId)
    ? favoriteIds.value.filter((id) => id !== stationId)
    : [...favoriteIds.value, stationId];
  localStorage.setItem('favoriteStations', JSON.stringify(favoriteIds.value));
};

const showPriceToast = (type: 'success' | 'error', message: string) => {
  priceToast.value = { type, message };
  window.setTimeout(() => {
    priceToast.value = null;
  }, 2800);
};

const applyPriceCache = (cache: Awaited<ReturnType<typeof refreshRealtimePrices>>) => {
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

const applyRealtimePrices = async () => {
  const cache = await refreshRealtimePrices();
  applyPriceCache(cache);
};

const refreshPrices = async () => {
  if (isRefreshingPrices.value) return;
  isRefreshingPrices.value = true;
  const cache = await refreshRealtimePrices();
  applyPriceCache(cache);
  isRefreshingPrices.value = false;
  showPriceToast(cache?.prices.length ? 'success' : 'error', cache?.prices.length ? 'Harga diperbarui' : 'Harga belum bisa diperbarui. Server harga tidak merespons, coba lagi sebentar.');
};

const initUserLocation = () => {
  if (!navigator.geolocation) {
    locationStatus.value = 'denied';
    return;
  }
  navigator.geolocation.getCurrentPosition(
    ({ coords }) => {
      const distanceFromJakarta = distanceBetweenKm(coords.latitude, coords.longitude, JAKARTA_CENTER.lat, JAKARTA_CENTER.lng);
      if (distanceFromJakarta > 60) {
        locationStatus.value = 'outside';
        return;
      }
      userLocation.value = { lat: coords.latitude, lng: coords.longitude };
      locationStatus.value = 'ok';
    },
    () => {
      locationStatus.value = 'denied';
    },
    { enableHighAccuracy: true, maximumAge: 300000, timeout: 8000 },
  );
};

onMounted(() => {
  favoriteIds.value = JSON.parse(localStorage.getItem('favoriteStations') || '[]') as string[];
  document.documentElement.classList.remove('dark');
  localStorage.removeItem('darkMode');
  initUserLocation();
  void applyRealtimePrices().finally(() => {
    isInitialPriceLoading.value = false;
  });
});
</script>

<template>
  <div class="flex h-full flex-col bg-trotoar text-aspal">
    <Header
      :query="filterState.query"
      :last-updated="relativeUpdated"
      :is-refreshing="isRefreshingPrices"
      @update-query="updateQuery"
      @refresh="refreshPrices"
    />

    <FilterBar
      :filters="filterState"
      :sort-mode="sortMode"
      @update-filters="updateFilters"
      @update-sort="sortMode = $event"
    />

    <main class="grid min-h-0 flex-1 grid-cols-1 grid-rows-[minmax(0,1fr)] lg:grid-cols-[minmax(0,1fr)_420px]">
      <div class="relative min-h-0">
        <article
          v-if="cheapestStation && showCheapestCard"
          class="absolute top-3 left-3 z-[700] w-[min(340px,calc(100%-24px))] rounded-xl bg-petrol-deep p-3 md:p-4 text-white shadow-[0_12px_32px_rgb(30_35_40/0.35)] md:top-4 md:left-4"
          aria-label="SPBU termurah"
        >
          <div class="flex items-start justify-between gap-3">
            <p class="text-[13px] text-white/70">
              {{ activeFuelLabel }} termurah
              <span v-if="isInitialPriceLoading" class="text-lampu/80">, memuat harga terbaru</span>
            </p>
            <button
              type="button"
              class="-mt-1 -mr-1 grid h-7 w-7 shrink-0 place-items-center rounded-md text-white/60 transition hover:bg-white/10 hover:text-white"
              aria-label="Sembunyikan papan harga"
              @click="showCheapestCard = false"
            >
              <i class="fa-solid fa-xmark" aria-hidden="true"></i>
            </button>
          </div>
          <p class="price mt-1 text-[40px] leading-[0.9] md:text-[56px] font-extrabold text-lampu">
            {{ cheapestStation.selectedPrice.toLocaleString('id-ID') }}<span class="ml-1 text-lg font-semibold text-white/60">/liter</span>
          </p>
          <div class="mt-2 flex items-end justify-between gap-3 md:mt-3">
            <div class="min-w-0 text-sm">
              <p class="truncate font-semibold">{{ cheapestStation.name }}</p>
              <p class="text-white/60">
                {{ cheapestStation.brand }}, {{ cheapestStation.distanceKm.toFixed(1) }} km
              </p>
            </div>
            <button
              type="button"
              class="inline-flex h-9 shrink-0 items-center gap-2 rounded-lg bg-lampu px-3 text-sm font-bold text-aspal transition hover:brightness-95 focus-visible:outline-white"
              @click="startNavigation(cheapestStation)"
            >
              <i class="fa-solid fa-diamond-turn-right text-xs" aria-hidden="true"></i>
              Rute
            </button>
          </div>
        </article>

        <button
          v-else-if="cheapestStation"
          type="button"
          class="absolute top-3 left-3 z-[700] inline-flex h-9 items-center gap-2 rounded-lg bg-petrol-deep px-3 text-sm font-semibold text-white md:top-4 md:left-4"
          @click="showCheapestCard = true"
        >
          Termurah <span class="price text-base text-lampu">{{ cheapestStation.selectedPrice.toLocaleString('id-ID') }}</span>
        </button>

        <Map
          ref="mapRef"
          :stations="stationResults"
          :active-fuel="filterState.fuelType || 'RON_92'"
          :selected-station-id="selectedStationId"
          @select-station="selectStation"
        />
      </div>

      <StationList
        :stations="stationResults"
        :sort-mode="sortMode"
        :selected-station-id="selectedStationId"
        :location-note="locationNote"
        @reset-filters="resetFilters"
        @flyto="handleFlyTo"
        @toggle-favorite="toggleFavorite"
        @start-navigation="startNavigation"
      />
    </main>

    <div
      v-if="priceToast"
      class="fixed bottom-6 left-1/2 z-[1400] max-w-[calc(100%-32px)] -translate-x-1/2 rounded-lg px-4 py-3 text-sm font-semibold text-white shadow-lg lg:bottom-8"
      :class="priceToast.type === 'success' ? 'bg-teal' : 'bg-[#a3292f]'"
      role="status"
    >
      {{ priceToast.message }}
    </div>
  </div>
</template>
