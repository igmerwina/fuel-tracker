<template>
  <div class="map-container">
    <div ref="mapContainer" class="map"></div>

    <div class="map-toolbar" aria-label="Legenda brand SPBU">
      <span v-for="item in brandLegend" :key="item.name" class="legend-item">
        <span class="legend-dot" :style="{ backgroundColor: item.color }"></span>
        {{ item.name }}
      </span>
    </div>

    <article v-if="selectedStation" class="detail-panel">
      <div class="popup-header">
        <div>
          <p class="brand-line" :style="{ color: getBrandColor(selectedStation.brand) }">
            {{ selectedStation.brand }} · {{ selectedStation.region }}
          </p>
          <h3>{{ selectedStation.name }}</h3>
        </div>
        <button class="close-btn" title="Tutup detail" type="button" @click="selectedStation = null">
          <i class="fa-solid fa-xmark" aria-hidden="true"></i>
        </button>
      </div>

      <p class="address">
        <i class="fa-solid fa-location-dot" aria-hidden="true"></i>
        {{ selectedStation.address }}
      </p>

      <a
        class="maps-link"
        :href="selectedStation.googleMapsUrl"
        target="_blank"
        rel="noopener noreferrer"
      >
        <i class="fa-solid fa-arrow-up-right-from-square" aria-hidden="true"></i>
        Buka Google Maps
      </a>

      <div class="prices-table">
        <div v-for="price in selectedStation.prices" :key="price.type" class="price-row">
          <span>{{ getFuelTypeLabel(price.type) }}</span>
          <strong>Rp {{ price.price.toLocaleString('id-ID') }}</strong>
        </div>
      </div>
    </article>
  </div>
</template>

<script setup lang="ts">
import { nextTick, onMounted, ref, watch } from 'vue';
import L from 'leaflet';
import type { FuelBrand, FuelStation } from '../types/index';
import { fuelTypes, JAKARTA_BOUNDS, JAKARTA_CENTER } from '../data/stations';

const props = defineProps<{
  stations: FuelStation[];
}>();

const mapContainer = ref<HTMLElement>();
const selectedStation = ref<FuelStation | null>(null);
const markers: Record<string, L.Marker> = {};
let map: L.Map | null = null;

const brandLegend: { name: FuelBrand; color: string }[] = [
  { name: 'Pertamina', color: '#d71e2b' },
  { name: 'Shell', color: '#f4c300' },
  { name: 'Vivo', color: '#1d4ed8' },
  { name: 'BP', color: '#159447' },
];

const getBrandColor = (brand: FuelBrand): string => {
  return brandLegend.find((item) => item.name === brand)?.color ?? '#0b4f7a';
};

const getFuelTypeLabel = (typeId: string): string => {
  const fuelType = fuelTypes.find((fuel) => fuel.id === typeId);
  return fuelType?.label || typeId;
};

const getRon92Price = (station: FuelStation) => {
  return station.prices.find((price) => price.type === 'RON_92')?.price ?? station.prices[0]?.price ?? 0;
};

const createMarker = (station: FuelStation) => {
  if (!map) return;

  const color = getBrandColor(station.brand);
  const ron92Price = getRon92Price(station);
  const marker = L.marker([station.latitude, station.longitude], {
    icon: L.divIcon({
      html: `
        <div class="price-marker" style="--pin-color: ${color}">
          <div class="price-marker__price">
            <span>RON 92</span>
            <strong>Rp ${ron92Price.toLocaleString('id-ID')}</strong>
          </div>
          <div class="price-marker__pin">
            <i class="fa-solid fa-gas-pump"></i>
          </div>
        </div>
      `,
      iconSize: [98, 58],
      iconAnchor: [49, 56],
      className: 'fuel-marker',
    }),
  });

  marker.on('click', () => {
    selectedStation.value = station;
  });

  marker.addTo(map);
  markers[station.id] = marker;
};

const syncMarkers = () => {
  if (!map) return;

  Object.values(markers).forEach((marker) => marker.remove());
  Object.keys(markers).forEach((key) => delete markers[key]);
  props.stations.forEach(createMarker);

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

  map.setMaxBounds(jakartaBounds);
  syncMarkers();
};

const flyToStation = (station: FuelStation) => {
  if (!map) return;

  map.flyTo([station.latitude, station.longitude], 15, {
    duration: 1.1,
  });
  selectedStation.value = station;
};

defineExpose({ flyToStation });

onMounted(async () => {
  await nextTick();
  initMap();
});

watch(
  () => props.stations,
  () => syncMarkers(),
  { deep: true },
);
</script>

<style scoped>
.map-container {
  position: relative;
  width: 100%;
  height: 100%;
  overflow: hidden;
  background: #dce7ef;
}

.map {
  width: 100%;
  height: 100%;
}

.map-toolbar {
  position: absolute;
  top: 14px;
  left: 14px;
  z-index: 500;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  max-width: min(520px, calc(100% - 28px));
  padding: 8px;
  border: 1px solid rgba(203, 213, 225, 0.88);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.94);
  box-shadow: var(--small-shadow);
}

.legend-item {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: #334155;
  font-size: 12px;
  font-weight: 800;
}

.legend-dot {
  width: 10px;
  height: 10px;
  border-radius: 999px;
}

.detail-panel {
  position: absolute;
  top: 14px;
  right: 14px;
  z-index: 500;
  width: min(360px, calc(100% - 28px));
  padding: 16px;
  border: 1px solid var(--line);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.98);
  box-shadow: var(--shadow);
}

.popup-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}

.brand-line {
  margin: 0 0 4px;
  font-size: 12px;
  font-weight: 900;
  text-transform: uppercase;
}

h3 {
  margin: 0;
  color: var(--ink);
  font-size: 19px;
  line-height: 1.2;
}

.close-btn {
  display: grid;
  width: 34px;
  height: 34px;
  flex: 0 0 auto;
  place-items: center;
  border-radius: 8px;
  color: #475569;
  background: #eef2f7;
  cursor: pointer;
}

.close-btn:hover {
  color: var(--brand-strong);
  background: #e2e8f0;
}

.address {
  display: flex;
  gap: 8px;
  margin: 14px 0;
  color: var(--muted);
  font-size: 13px;
}

.maps-link {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 14px;
  color: var(--brand);
  font-size: 13px;
  font-weight: 900;
  text-decoration: none;
}

.maps-link:hover {
  color: var(--brand-strong);
  text-decoration: underline;
}

.prices-table {
  display: grid;
  overflow: hidden;
  border: 1px solid var(--line);
  border-radius: 8px;
}

.price-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 10px 12px;
  background: #ffffff;
}

.price-row + .price-row {
  border-top: 1px solid var(--line);
}

.price-row span {
  color: #475569;
  font-size: 13px;
  font-weight: 800;
}

.price-row strong {
  color: var(--brand-strong);
  font-size: 15px;
}

@media (max-width: 760px) {
  .map-toolbar {
    top: 10px;
    left: 10px;
  }

  .detail-panel {
    top: auto;
    right: 10px;
    bottom: 10px;
    width: calc(100% - 20px);
  }
}
</style>

<style>
.price-marker {
  display: grid;
  justify-items: center;
  gap: 2px;
  width: 98px;
  pointer-events: auto;
}

.price-marker__price {
  display: grid;
  min-width: 86px;
  padding: 5px 7px;
  border: 2px solid var(--pin-color);
  border-radius: 8px;
  color: #0f172a;
  background: #ffffff;
  box-shadow: 0 8px 20px rgba(15, 23, 42, 0.22);
  text-align: center;
}

.price-marker__price span {
  color: #64748b;
  font-size: 9px;
  font-weight: 900;
  line-height: 1;
}

.price-marker__price strong {
  margin-top: 2px;
  color: #0f172a;
  font-size: 12px;
  line-height: 1;
  white-space: nowrap;
}

.price-marker__pin {
  position: relative;
  display: grid;
  width: 28px;
  height: 28px;
  place-items: center;
  border: 3px solid #ffffff;
  border-radius: 50% 50% 50% 4px;
  color: #ffffff;
  background: var(--pin-color);
  box-shadow: 0 8px 18px rgba(15, 23, 42, 0.34);
  transform: rotate(-45deg);
}

.price-marker__pin i {
  font-size: 11px;
  transform: rotate(45deg);
}
</style>
