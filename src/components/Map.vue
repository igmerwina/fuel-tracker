<template>
  <div class="relative h-full w-full overflow-hidden bg-slate-200">
    <div ref="mapContainer" class="h-full w-full"></div>

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
          title="Tutup detail"
          @click="selectedStation = null"
        >
          <i class="fa-solid fa-xmark" aria-hidden="true"></i>
        </button>
      </div>

      <p class="mt-2 text-sm font-semibold text-slate-500">{{ selectedStation.address }}</p>
      <div class="mt-4 rounded-xl bg-emerald-50 p-3 text-emerald-900" v-if="selectedStation.savingsPerLiter > 0">
        <strong class="block text-sm font-black">
          Hemat Rp {{ selectedStation.savingsPerLiter.toLocaleString('id-ID') }}/L dibanding rata-rata
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
  selectedStationId: string;
}>();

const emit = defineEmits<{
  'select-station': [station: StationResult];
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
  if (props.selectedStationId === station.id) return 'selected';
  const price = activePrice(station);
  if (station.isCheapest) return 'cheapest';
  if (price <= props.averagePrice) return 'average';
  return 'expensive';
};

const brandLogo = (brand: string) => {
  const logos: Record<string, string> = {
    BP: '<span class="brand-logo brand-logo--bp">bp</span>',
    Shell: '<span class="brand-logo brand-logo--shell">S</span>',
    Vivo: '<span class="brand-logo brand-logo--vivo">V</span>',
    Pertamina: '<span class="brand-logo brand-logo--pertamina">P</span>',
  };
  return logos[brand] ?? `<span class="brand-logo">${brand.charAt(0)}</span>`;
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
  const zoom = map?.getZoom() ?? 13;
  const mode = zoom < 13 ? 'brand' : zoom < 15 ? 'price' : 'selected';
  const label = station.isCheapest ? '🔥' : '';
  const content =
    mode === 'brand'
      ? `<button class="brand-pin brand-pin--${tone}" aria-label="${station.brand} ${station.name}">${brandLogo(station.brand)}</button>`
      : `<button class="gm-price-marker gm-price-marker--${tone} gm-price-marker--${mode}" aria-label="${station.name} Rp ${price.toLocaleString('id-ID')}">
          <strong>${label} Rp ${price.toLocaleString('id-ID')}</strong>
          <span class="gm-price-marker__body">${brandLogo(station.brand)}</span>
        </button>`;
  return L.marker([station.latitude, station.longitude], {
    icon: L.divIcon({
      html: content,
      iconSize: mode === 'brand' ? [44, 44] : props.selectedStationId === station.id ? [104, 112] : [82, 92],
      iconAnchor: mode === 'brand' ? [22, 22] : props.selectedStationId === station.id ? [52, 108] : [41, 88],
      className: 'fuel-marker',
    }),
  }).on('click mouseover', () => {
    selectedStation.value = station;
    emit('select-station', station);
  });
};

const createClusterMarker = (stations: StationResult[]) => {
  const cheapest = [...stations].sort((a, b) => activePrice(a) - activePrice(b))[0];
  return L.marker([cheapest.latitude, cheapest.longitude], {
    icon: L.divIcon({
      html: `
        <button class="gm-cluster" aria-label="${stations.length} SPBU terdekat">
          <strong>${stations.length}</strong>
          <span>mulai Rp ${activePrice(cheapest).toLocaleString('id-ID')}</span>
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
  L.tileLayer('https://{s}.tile-cyclosm.openstreetmap.fr/cyclosm/{z}/{x}/{y}.png', {
    attribution: '&copy; OpenStreetMap contributors, CyclOSM',
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
  () => [props.stations, props.activeFuel, props.averagePrice, props.selectedStationId],
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
  min-width: 82px;
  overflow: visible;
  padding: 0;
  border: 0;
  border-radius: 12px;
  color: #ffffff;
  box-shadow: 0 10px 24px rgba(15, 23, 42, 0.28);
  cursor: pointer;
  font-family: Inter, ui-sans-serif, system-ui, sans-serif;
  text-align: center;
}

.gm-price-marker {
  position: relative;
  grid-template-rows: 34px 50px;
  background: #ffffff;
  color: #0f172a;
}

.gm-price-marker::after {
  position: absolute;
  bottom: -13px;
  left: 50%;
  width: 20px;
  height: 20px;
  border-right: 1px solid #cbd5e1;
  border-bottom: 1px solid #cbd5e1;
  background: #ffffff;
  content: '';
  transform: translateX(-50%) rotate(45deg);
}

.gm-price-marker strong,
.gm-cluster strong {
  display: grid;
  place-items: center;
  border-radius: 12px 12px 0 0;
  color: #ffffff;
  font-size: 16px;
  font-weight: 950;
  line-height: 1;
}

.gm-price-marker__body {
  display: grid;
  place-items: center;
  border-right: 1px solid #cbd5e1;
  border-left: 1px solid #cbd5e1;
  border-radius: 0 0 12px 12px;
  background: #ffffff;
  color: #0f172a;
}

.gm-price-marker--cheapest {
  border: 1px solid #10b981;
}

.gm-price-marker--cheapest strong {
  background: #10b981;
}

.gm-price-marker--average {
  border: 1px solid #2563eb;
}

.gm-price-marker--average strong {
  background: #2563eb;
}

.gm-price-marker--expensive {
  border: 1px solid #64748b;
}

.gm-price-marker--expensive strong {
  background: #64748b;
}

.gm-price-marker--selected {
  transform: scale(1.14);
  z-index: 999;
  border: 2px solid #2563eb;
}

.gm-price-marker--selected strong {
  background: #2563eb;
}

.gm-price-marker--brand {
  display: none;
}

.gm-cluster {
  background: #0f172a;
}

.leaflet-control-zoom a {
  color: #0f172a;
}

.brand-logo {
  display: grid;
  width: 34px;
  height: 34px;
  place-items: center;
  border-radius: 999px;
  font-size: 15px;
  font-weight: 950;
}

.brand-logo--bp {
  color: #15803d;
  background:
    radial-gradient(circle at center, #facc15 0 18%, transparent 19%),
    conic-gradient(#16a34a 0 10%, #facc15 10% 16%, #16a34a 16% 26%, #facc15 26% 32%, #16a34a 32% 42%, #facc15 42% 48%, #16a34a 48% 58%, #facc15 58% 64%, #16a34a 64% 74%, #facc15 74% 80%, #16a34a 80% 90%, #facc15 90% 100%);
  font-size: 0;
}

.brand-logo--shell {
  color: #991b1b;
  background: #facc15;
}

.brand-logo--vivo {
  color: #ffffff;
  background: #2563eb;
}

.brand-logo--pertamina {
  color: #ffffff;
  background: #ef4444;
}

.brand-pin {
  display: grid;
  width: 42px;
  height: 42px;
  place-items: center;
  border: 3px solid #ffffff;
  border-radius: 999px;
  background: #ffffff;
  box-shadow: 0 8px 18px rgba(15, 23, 42, 0.22);
}

.brand-pin--cheapest {
  outline: 3px solid #10b981;
}

.brand-pin--selected {
  outline: 3px solid #2563eb;
}
</style>
