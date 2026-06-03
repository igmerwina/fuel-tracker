<template>
  <div class="relative h-full w-full overflow-hidden bg-slate-200">
    <div ref="mapContainer" class="h-full w-full"></div>

    <div class="absolute left-3 top-3 z-[500] rounded-[12px] bg-white/95 p-2 shadow-lg shadow-slate-950/10 backdrop-blur">
      <div class="flex flex-wrap items-center gap-2 text-xs font-black text-slate-600">
        <span class="inline-flex items-center gap-1"><span class="h-2.5 w-2.5 rounded-full bg-emerald-500"></span> Cheapest</span>
        <span class="inline-flex items-center gap-1"><span class="h-2.5 w-2.5 rounded-full bg-blue-600"></span> Average</span>
        <span class="inline-flex items-center gap-1"><span class="h-2.5 w-2.5 rounded-full bg-slate-400"></span> Higher</span>
      </div>
    </div>

    <article
      v-if="selectedStation"
      class="absolute bottom-4 left-4 right-4 z-[500] rounded-[12px] bg-white p-4 shadow-2xl shadow-slate-950/20 md:left-auto md:right-4 md:top-4 md:bottom-auto md:w-[380px]"
    >
      <div class="flex items-start justify-between gap-3">
        <div class="min-w-0">
          <p class="text-xs font-black uppercase tracking-wide text-blue-600">
            {{ selectedStation.brand }} · {{ selectedStation.distanceKm.toFixed(1) }} km
          </p>
          <h3 class="mt-1 truncate text-xl font-black text-slate-950">{{ selectedStation.name }}</h3>
        </div>
        <button
          class="grid h-10 w-10 place-items-center rounded-full bg-slate-100 text-slate-500 hover:bg-slate-200"
          type="button"
          title="Close details"
          @click="selectedStation = null"
        >
          <i class="fa-solid fa-xmark" aria-hidden="true"></i>
        </button>
      </div>

      <p class="mt-2 text-sm font-semibold text-slate-500">{{ selectedStation.address }}</p>
      <div class="mt-4 rounded-xl bg-emerald-50 p-3 text-emerald-900" v-if="selectedStation.savingsPerLiter > 0">
        <strong class="block text-sm font-black">
          You save Rp {{ selectedStation.savingsPerLiter.toLocaleString('id-ID') }}/L compared to average
        </strong>
      </div>

      <div class="mt-4 grid grid-cols-2 gap-2">
        <div v-for="price in selectedStation.prices" :key="price.type" class="rounded-xl bg-slate-100 p-3">
          <p class="text-xs font-black text-slate-500">{{ getFuelTypeLabel(price.type) }}</p>
          <p class="text-base font-black text-slate-950">Rp {{ price.price.toLocaleString('id-ID') }}</p>
        </div>
      </div>
    </article>
  </div>
</template>

<script setup lang="ts">
import { nextTick, onMounted, ref, watch } from 'vue';
import L from 'leaflet';
import type { FuelType, StationResult } from '../types/index';
import { fuelTypes, JAKARTA_BOUNDS, JAKARTA_CENTER } from '../data/stations';

const props = defineProps<{
  stations: StationResult[];
  activeFuel: FuelType;
  averagePrice: number;
}>();

const mapContainer = ref<HTMLElement>();
const selectedStation = ref<StationResult | null>(null);
let map: L.Map | null = null;
let markerLayer: L.LayerGroup | null = null;

const getFuelTypeLabel = (typeId: string): string => {
  const fuelType = fuelTypes.find((fuel) => fuel.id === typeId);
  return fuelType?.label || typeId;
};

const activePrice = (station: StationResult) =>
  station.prices.find((price) => price.type === props.activeFuel)?.price ?? station.selectedPrice;

const markerTone = (station: StationResult) => {
  const price = activePrice(station);
  if (station.isCheapest) return 'cheapest';
  if (price <= props.averagePrice) return 'average';
  return 'expensive';
};

const clusterStations = () => {
  if (!map || map.getZoom() >= 13) return props.stations.map((station) => [station] as StationResult[]);

  const buckets = new globalThis.Map<string, StationResult[]>();
  props.stations.forEach((station) => {
    const key = `${Math.round(station.latitude / 0.025)}:${Math.round(station.longitude / 0.025)}`;
    buckets.set(key, [...(buckets.get(key) ?? []), station]);
  });
  return [...buckets.values()];
};

const createStationMarker = (station: StationResult) => {
  const price = activePrice(station);
  const tone = markerTone(station);
  return L.marker([station.latitude, station.longitude], {
    icon: L.divIcon({
      html: `
        <button class="gm-price-marker gm-price-marker--${tone}" aria-label="${station.name} Rp ${price.toLocaleString('id-ID')}">
          <span>${station.isCheapest ? 'Best' : getFuelTypeLabel(props.activeFuel)}</span>
          <strong>Rp ${price.toLocaleString('id-ID')}</strong>
        </button>
      `,
      iconSize: [94, 42],
      iconAnchor: [47, 42],
      className: 'fuel-marker',
    }),
  }).on('click', () => {
    selectedStation.value = station;
  });
};

const createClusterMarker = (stations: StationResult[]) => {
  const cheapest = [...stations].sort((a, b) => activePrice(a) - activePrice(b))[0];
  return L.marker([cheapest.latitude, cheapest.longitude], {
    icon: L.divIcon({
      html: `
        <button class="gm-cluster" aria-label="${stations.length} nearby stations">
          <strong>${stations.length}</strong>
          <span>from Rp ${activePrice(cheapest).toLocaleString('id-ID')}</span>
        </button>
      `,
      iconSize: [88, 48],
      iconAnchor: [44, 44],
      className: 'fuel-marker',
    }),
  }).on('click', () => {
    map?.flyTo([cheapest.latitude, cheapest.longitude], 14, { duration: 0.8 });
  });
};

const syncMarkers = () => {
  if (!map || !markerLayer) return;
  markerLayer.clearLayers();

  clusterStations().forEach((group) => {
    const marker = group.length > 1 ? createClusterMarker(group) : createStationMarker(group[0]);
    markerLayer?.addLayer(marker);
  });

  if (
    selectedStation.value &&
    !props.stations.some((station) => station.id === selectedStation.value?.id)
  ) {
    selectedStation.value = null;
  }
};

const initMap = () => {
  if (!mapContainer.value) return;

  const jakartaBounds = L.latLngBounds(
    [JAKARTA_BOUNDS.southWest.lat, JAKARTA_BOUNDS.southWest.lng],
    [JAKARTA_BOUNDS.northEast.lat, JAKARTA_BOUNDS.northEast.lng],
  );

  map = L.map(mapContainer.value, {
    maxBounds: jakartaBounds,
    maxBoundsViscosity: 1,
    minZoom: 11,
    zoomControl: false,
    worldCopyJump: false,
  }).setView([JAKARTA_CENTER.lat, JAKARTA_CENTER.lng], 13);

  L.control.zoom({ position: 'bottomright' }).addTo(map);
  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '&copy; OpenStreetMap contributors',
    bounds: jakartaBounds,
    maxZoom: 19,
    minZoom: 11,
    noWrap: true,
  }).addTo(map);

  markerLayer = L.layerGroup().addTo(map);
  map.on('zoomend', syncMarkers);
  map.setMaxBounds(jakartaBounds);
  syncMarkers();
};

const flyToStation = (station: StationResult) => {
  if (!map) return;
  map.flyTo([station.latitude, station.longitude], 15, { duration: 0.9 });
  selectedStation.value = station;
};

defineExpose({ flyToStation });

onMounted(async () => {
  await nextTick();
  initMap();
});

watch(
  () => [props.stations, props.activeFuel, props.averagePrice],
  () => syncMarkers(),
  { deep: true },
);
</script>

<style>
.fuel-marker {
  background: none !important;
  border: none !important;
}

.gm-price-marker,
.gm-cluster {
  display: grid;
  min-width: 90px;
  padding: 6px 9px;
  border: 0;
  border-radius: 999px;
  color: #ffffff;
  box-shadow: 0 10px 24px rgba(15, 23, 42, 0.28);
  cursor: pointer;
  font-family: Inter, ui-sans-serif, system-ui, sans-serif;
  text-align: center;
}

.gm-price-marker span,
.gm-cluster span {
  font-size: 10px;
  font-weight: 900;
  line-height: 1;
}

.gm-price-marker strong,
.gm-cluster strong {
  margin-top: 2px;
  font-size: 13px;
  font-weight: 950;
  line-height: 1;
}

.gm-price-marker--cheapest {
  background: #10b981;
}

.gm-price-marker--average {
  background: #2563eb;
}

.gm-price-marker--expensive {
  background: #64748b;
}

.gm-cluster {
  background: #0f172a;
}

.leaflet-control-zoom a {
  color: #0f172a;
}
</style>
