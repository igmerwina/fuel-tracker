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
import { fetchRealtimePrices, refreshRealtimePrices } from './services/prices';

const stations = ref<FuelStation[]>(fuelStations);
const favoriteIds = ref<string[]>([]);
const sortMode = ref<SortMode>('cheapest');
const priceDay = ref<'today' | 'tomorrow'>('today');
const selectedStationId = ref<string>('');
const navigatingStationId = ref<string>('');
const showCheapestCard = ref(true);
const userLocation = ref(JAKARTA_CENTER);
const locationStatus = ref<'idle' | 'ok' | 'outside' | 'denied'>('idle');
const isRefreshingPrices = ref(false);
const priceToast = ref<{ type: 'success' | 'error'; message: string } | null>(null);
const filterState = ref<FilterState>({
  query: '',
  region: '',
  brand: '',
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
      if (!station.prices.some((price) => price.type === activeFuel)) return false;
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
  return `${hours} jam lalu`;
});

const updateFilters = (next: Partial<FilterState>) => {
  filterState.value = { ...filterState.value, ...next };
};

const updateQuery = (query: string) => {
  updateFilters({ query });
};

const resetFilters = () => {
  filterState.value = { query: '', region: '', brand: '', fuelType: 'RON_92' };
  sortMode.value = 'cheapest';
  priceDay.value = 'today';
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
  navigatingStationId.value = station.id;
  window.setTimeout(() => {
    window.open(station.googleMapsUrl, '_blank', 'noopener,noreferrer');
    navigatingStationId.value = '';
  }, 450);
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

const applyPriceCache = (cache: Awaited<ReturnType<typeof fetchRealtimePrices>>) => {
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
  const cache = await fetchRealtimePrices();
  applyPriceCache(cache);
};

const refreshPrices = async () => {
  if (isRefreshingPrices.value) return;
  isRefreshingPrices.value = true;
  const cache = await refreshRealtimePrices();
  applyPriceCache(cache);
  isRefreshingPrices.value = false;
  showPriceToast(cache?.prices.length ? 'success' : 'error', cache?.prices.length ? 'Data harga berhasil diperbarui' : 'Gagal memperbarui data harga');
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

    <main class="mx-auto flex h-[calc(100vh-60px)] max-w-[1800px] flex-col gap-2 p-1.5 md:h-[calc(100vh-72px)] md:gap-3 md:p-3">
      <FilterBar
        :filters="filterState"
        :sort-mode="sortMode"
        @update-filters="updateFilters"
        @update-sort="sortMode = $event"
      />

      <section class="flex flex-wrap items-center gap-1.5 rounded-xl bg-white px-2.5 py-1.5 text-[11px] font-black text-slate-600 shadow-sm shadow-slate-200/60 ring-1 ring-slate-100 md:gap-2 md:rounded-2xl md:px-3 md:py-2 md:text-xs">
        <div class="inline-flex rounded-full bg-slate-100 p-0.5 md:p-1">
          <button
            type="button"
            class="h-7 rounded-full px-2.5 text-[11px] transition md:h-8 md:px-3 md:text-xs"
            :class="priceDay === 'today' ? 'bg-blue-600 text-white shadow-sm shadow-blue-600/20' : 'text-slate-600'"
            @click="priceDay = 'today'"
          >
            Hari ini
          </button>
          <button
            type="button"
            class="h-7 rounded-full px-2.5 text-[11px] transition md:h-8 md:px-3 md:text-xs"
            :class="priceDay === 'tomorrow' ? 'bg-amber-500 text-white shadow-sm shadow-amber-500/20' : 'text-slate-400'"
            title="Data harga besok belum tersedia"
            @click="priceDay = 'tomorrow'"
          >
            Besok
          </button>
        </div>
        <span class="hidden sm:inline-flex items-center gap-1 rounded-full bg-emerald-50 px-2.5 py-1.5 text-emerald-700 md:px-3 md:py-2">
          <i class="fa-solid fa-lock" aria-hidden="true"></i>
          Harga brand-level
        </span>
        <span
          v-if="locationStatus === 'ok'"
          class="hidden sm:inline-flex items-center gap-1 rounded-full bg-blue-50 px-2.5 py-1.5 text-blue-700 md:px-3 md:py-2"
        >
          <i class="fa-solid fa-location-crosshairs" aria-hidden="true"></i>
          GPS aktif
        </span>
        <span
          v-else-if="locationStatus === 'outside'"
          class="inline-flex items-center gap-1 rounded-full bg-rose-50 px-2 py-1 text-[10px] text-rose-700 md:px-3 md:py-2 md:text-xs"
        >
          <i class="fa-solid fa-triangle-exclamation" aria-hidden="true"></i>
          Luar Jakarta
        </span>
        <span v-if="priceDay === 'tomorrow'" class="rounded-full bg-amber-50 px-2 py-1 text-[10px] text-amber-700 md:px-3 md:py-2 md:text-xs">
          Besok blm tersedia
        </span>
        <button
          type="button"
          class="ml-auto inline-flex h-7 items-center gap-1.5 rounded-full bg-blue-600 px-2.5 text-[11px] text-white shadow-sm shadow-blue-600/20 transition hover:bg-blue-700 disabled:cursor-wait disabled:opacity-75 md:h-9 md:gap-2 md:px-3 md:text-xs"
          :disabled="isRefreshingPrices"
          @click="refreshPrices"
        >
          <i class="fa-solid fa-rotate-right text-[10px] md:text-xs" :class="isRefreshingPrices ? 'animate-spin' : ''" aria-hidden="true"></i>
          <span class="hidden sm:inline">{{ isRefreshingPrices ? 'Memuat...' : 'Refresh' }}</span>
        </button>
      </section>

      <div
        v-if="priceToast"
        class="fixed right-4 top-20 z-[1400] rounded-2xl px-4 py-3 text-sm font-black text-white shadow-2xl"
        :class="priceToast.type === 'success' ? 'bg-emerald-500' : 'bg-rose-500'"
      >
        {{ priceToast.message }}
      </div>

      <section class="grid min-h-0 flex-1 grid-cols-1 gap-2 lg:grid-cols-[minmax(0,3fr)_minmax(320px,1fr)] lg:gap-3">
        <div class="relative min-h-[calc(100vh-120px)] overflow-hidden rounded-[10px] bg-white shadow-lg shadow-slate-200/70 lg:min-h-0 lg:rounded-[12px]">
          <article
            v-if="cheapestStation && showCheapestCard"
            class="absolute left-2 right-2 top-2 z-[700] max-w-[320px] rounded-xl bg-white/95 p-2.5 shadow-xl shadow-slate-950/12 ring-1 ring-white/70 backdrop-blur md:left-3 md:right-auto md:top-3 md:max-w-[360px] md:rounded-2xl md:p-3"
          >
            <div class="flex items-start justify-between gap-2 md:gap-3">
              <p class="inline-flex items-center gap-1.5 rounded-full bg-emerald-50 px-2 py-0.5 text-[10px] font-black uppercase tracking-wide text-emerald-700 md:gap-2 md:px-2.5 md:py-1 md:text-[11px]">
                <i class="fa-solid fa-fire-flame-curved" aria-hidden="true"></i>
                Terbaik hari ini
              </p>
              <button
                type="button"
                class="grid h-6 w-6 shrink-0 place-items-center rounded-full bg-slate-100 text-slate-500 transition hover:bg-slate-200 hover:text-slate-900 md:h-7 md:w-7"
                aria-label="Tutup kartu harga terbaik"
                @click="showCheapestCard = false"
              >
                <i class="fa-solid fa-xmark" aria-hidden="true"></i>
              </button>
            </div>
            <div class="mt-1.5 flex items-end justify-between gap-2 md:mt-2 md:gap-3">
              <div class="min-w-0">
                <h2 class="truncate text-xs font-black text-slate-950 md:text-sm">{{ cheapestStation.name }}</h2>
                <p class="mt-0.5 text-lg font-black tracking-tight text-slate-950 md:mt-0.5 md:text-2xl">
                  Rp{{ cheapestStation.selectedPrice.toLocaleString('id-ID') }}<span class="text-xs text-slate-500 md:text-sm">/L</span>
                </p>
                <p class="mt-0.5 text-[11px] font-bold text-slate-500 md:text-xs">
                  {{ cheapestStation.distanceKm.toFixed(1) }} km · Hemat Rp{{ cheapestStation.savingsPerLiter.toLocaleString('id-ID') }}/L
                </p>
              </div>
              <button
                type="button"
                class="grid h-9 w-9 shrink-0 place-items-center rounded-full bg-blue-600 text-white shadow-lg shadow-blue-600/20 transition hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-600 focus:ring-offset-2 md:h-10 md:w-10"
                aria-label="Mulai navigasi"
                @click="startNavigation(cheapestStation)"
              >
                <i v-if="navigatingStationId === cheapestStation.id" class="fa-solid fa-spinner animate-spin" aria-hidden="true"></i>
                <i v-else class="fa-solid fa-route" aria-hidden="true"></i>
              </button>
            </div>
          </article>
          <Map
            ref="mapRef"
            :stations="stationResults"
            :active-fuel="filterState.fuelType || 'RON_92'"
            :average-price="averageSelectedPrice"
            :selected-station-id="selectedStationId"
            @select-station="selectStation"
          />
        </div>

        <StationList
          :stations="stationResults"
          :active-fuel="filterState.fuelType || 'RON_92'"
          :sort-mode="sortMode"
          :selected-station-id="selectedStationId"
          @reset-filters="resetFilters"
          @flyto="handleFlyTo"
          @toggle-favorite="toggleFavorite"
          @start-navigation="startNavigation"
        />
      </section>
    </main>
  </div>
</template>
