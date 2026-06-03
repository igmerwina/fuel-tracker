<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import Header from './components/Header.vue';
import FilterBar from './components/FilterBar.vue';
import StationList from './components/StationList.vue';
import Map from './components/Map.vue';
import type { FilterState, FuelBrand, FuelStation, FuelType } from './types/index';
import { fuelStations } from './data/stations';
import { fetchRealtimePrices } from './services/prices';

const stations = ref<FuelStation[]>(fuelStations);

const filterState = ref<FilterState>({
  region: '',
  brand: '',
  fuelType: '',
});

const mapRef = ref<InstanceType<typeof Map>>();

const filteredStations = computed(() => {
  return stations.value
    .filter((station) => {
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

const handleFiltersUpdate = (newFilters: FilterState) => {
  filterState.value = { ...newFilters };
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
