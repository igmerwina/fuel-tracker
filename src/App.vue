<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import Header from './components/Header.vue';
import SearchBar from './components/SearchBar.vue';
import FilterBar from './components/FilterBar.vue';
import StationList from './components/StationList.vue';
import Map from './components/Map.vue';
import type { FilterState, FuelBrand, FuelStation, FuelType } from './types/index';
import { fuelStations, fuelTypes } from './data/stations';
import { fetchRealtimePrices } from './services/prices';

const stations = ref<FuelStation[]>(fuelStations);

const filterState = ref<FilterState>({
  query: '',
  region: '',
  brand: '',
  fuelType: '',
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

const filteredStations = computed(() => {
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
      if (
        filterState.value.fuelType &&
        !station.prices.some((price) => price.type === filterState.value.fuelType)
      ) {
        return false;
      }
      return true;
    })
    .sort((a, b) => {
      const selectedFuel = filterState.value.fuelType;
      const priceA =
        a.prices.find((price) => price.type === selectedFuel)?.price ?? a.prices[0]?.price ?? 0;
      const priceB =
        b.prices.find((price) => price.type === selectedFuel)?.price ?? b.prices[0]?.price ?? 0;
      return priceA - priceB;
    });
});

const lowestPrice = computed(() => {
  const prices = filteredStations.value.flatMap((station) => station.prices.map((price) => price.price));
  return prices.length ? Math.min(...prices) : 0;
});

const lastUpdated = computed(() => {
  const latest = filteredStations.value
    .map((station) => new Date(station.lastUpdated).getTime())
    .sort((a, b) => b - a)[0];

  return latest
    ? new Intl.DateTimeFormat('id-ID', {
        day: '2-digit',
        month: 'short',
        hour: '2-digit',
        minute: '2-digit',
      }).format(latest)
    : '-';
});

const handleFiltersUpdate = (newFilters: Partial<FilterState>) => {
  filterState.value = { ...filterState.value, ...newFilters };
};

const handleQueryUpdate = (query: string) => {
  filterState.value = { ...filterState.value, query };
};

const handleFlyTo = (station: FuelStation) => {
  mapRef.value?.flyToStation(station);
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
  void applyRealtimePrices();
});
</script>

<template>
  <div class="app-shell">
    <Header
      :station-count="filteredStations.length"
      :lowest-price="lowestPrice"
      :last-updated="lastUpdated"
    />

    <SearchBar
      :query="filterState.query"
      :result-count="filteredStations.length"
      @update-query="handleQueryUpdate"
    />

    <main class="workspace">
      <section class="map-section" aria-label="Peta SPBU Jakarta">
        <Map ref="mapRef" :stations="filteredStations" />
      </section>

      <aside class="sidebar" aria-label="Filter dan daftar SPBU">
        <FilterBar :result-count="filteredStations.length" @update-filters="handleFiltersUpdate" />
        <StationList :stations="filteredStations" :filters="filterState" @flyto="handleFlyTo" />
      </aside>
    </main>
  </div>
</template>

<style scoped>
.app-shell {
  width: 100%;
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: var(--page);
}

.workspace {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 440px;
  gap: 18px;
  min-height: 0;
  flex: 1;
  padding: 18px;
}

.map-section,
.sidebar {
  min-height: 0;
}

.map-section {
  overflow: hidden;
  border: 1px solid var(--line);
  border-radius: 8px;
  background: var(--panel);
  box-shadow: var(--shadow);
}

.sidebar {
  display: flex;
  min-width: 0;
  flex-direction: column;
  gap: 14px;
}

@media (max-width: 1180px) {
  .workspace {
    grid-template-columns: minmax(0, 1fr) 390px;
  }
}

@media (max-width: 920px) {
  .workspace {
    grid-template-columns: 1fr;
    grid-template-rows: minmax(380px, 48vh) minmax(420px, 1fr);
    overflow: auto;
  }
}

@media (max-width: 760px) {
  .app-shell {
    height: auto;
    min-height: 100vh;
  }

  .workspace {
    grid-template-rows: 420px auto;
    gap: 12px;
    padding: 12px;
  }
}
</style>
