<template>
  <div class="relative h-full w-full overflow-hidden bg-trotoar">
    <div ref="mapContainer" class="h-full w-full"></div>
  </div>
</template>

<script setup lang="ts">
import { nextTick, onMounted, ref, watch } from 'vue';
import L from 'leaflet';
import type { FuelType, StationResult } from '../types/index';
import { JAKARTA_BOUNDS } from '../data/stations';

const props = defineProps<{
  stations: StationResult[];
  activeFuel: FuelType;
  selectedStationId: string;
}>();

const emit = defineEmits<{
  'select-station': [station: StationResult | null];
}>();

const mapContainer = ref<HTMLElement>();
let map: L.Map | null = null;
let markerLayer: L.LayerGroup | null = null;
const landingCenter: L.LatLngExpression = [-6.1856, 106.8272];
const landingZoom = 12;
const stationZoom = 15;

// Esri World Street Map: keyless, attribution required.
const TILES = 'https://server.arcgisonline.com/ArcGIS/rest/services/World_Street_Map/MapServer/tile/{z}/{y}/{x}';
const ATTRIBUTION = 'Tiles &copy; Esri, HERE, Garmin &middot; &copy; OpenStreetMap contributors';

const activePrice = (station: StationResult) =>
  station.prices.find((price) => price.type === props.activeFuel)?.price ?? station.selectedPrice;

const escapeHtml = (value: string) =>
  value
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&#39;');

const formatPrice = (price: number) => price.toLocaleString('id-ID');

const routeUrl = (station: StationResult) =>
  station.googleMapsUrl || `https://www.google.com/maps/search/?api=1&query=${station.latitude},${station.longitude}`;

const clusterStations = () => {
  if (!map || map.getZoom() >= 14) return props.stations.map((station) => [station]);

  const buckets = new globalThis.Map<string, StationResult[]>();
  props.stations.forEach((station) => {
    const key = `${Math.round(station.latitude / 0.035)}:${Math.round(station.longitude / 0.035)}`;
    buckets.set(key, [...(buckets.get(key) ?? []), station]);
  });
  return [...buckets.values()];
};

const createStationMarker = (station: StationResult) => {
  const tone = props.selectedStationId === station.id ? 'selected' : station.isCheapest ? 'cheapest' : 'default';
  const content = `
    <div class="pin pin--${tone}" role="button" tabindex="0" aria-label="${escapeHtml(station.name)}, Rp ${formatPrice(activePrice(station))} per liter">
      <span class="pin__dot pin__dot--${station.brand.toLowerCase()}"></span>
      <span class="pin__price">${formatPrice(activePrice(station))}</span>
      <span class="pin__card">
        <strong>${escapeHtml(station.name)}</strong>
        <span>${escapeHtml(station.brand)}, ${station.distanceKm.toFixed(1)} km</span>
        <a href="${escapeHtml(routeUrl(station))}" target="_blank" rel="noopener noreferrer">Buka rute</a>
      </span>
    </div>`;

  return L.marker([station.latitude, station.longitude], {
    icon: L.divIcon({ html: content, iconSize: [0, 0], iconAnchor: [0, 0], className: 'fuel-marker' }),
    zIndexOffset: tone === 'selected' ? 1000 : tone === 'cheapest' ? 500 : 0,
  }).on('click', () => {
    emit('select-station', props.selectedStationId === station.id ? null : station);
  });
};

const createClusterMarker = (stations: StationResult[]) => {
  const cheapest = [...stations].sort((a, b) => activePrice(a) - activePrice(b))[0];
  return L.marker([cheapest.latitude, cheapest.longitude], {
    icon: L.divIcon({
      html: `<button class="cluster" aria-label="${stations.length} SPBU, mulai Rp ${formatPrice(activePrice(cheapest))}">
          <strong>${stations.length}</strong><span>dari ${formatPrice(activePrice(cheapest))}</span>
        </button>`,
      iconSize: [0, 0],
      iconAnchor: [0, 0],
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
    markerLayer?.addLayer(group.length > 1 ? createClusterMarker(group) : createStationMarker(group[0]));
  });
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
    minZoom: landingZoom,
    maxZoom: 18,
    zoomControl: false,
  });

  L.control.zoom({ position: 'bottomright' }).addTo(map);
  L.tileLayer(TILES, { maxZoom: 18, noWrap: true, attribution: ATTRIBUTION }).addTo(map);

  markerLayer = L.layerGroup().addTo(map);
  map.on('zoomend', syncMarkers);
  map.setView(landingCenter, landingZoom);
  window.setTimeout(() => map?.invalidateSize(), 250);
  syncMarkers();
};

const flyToStation = (station: StationResult) => {
  map?.flyTo([station.latitude, station.longitude], stationZoom, { duration: 0.9 });
};

defineExpose({ flyToStation });

onMounted(async () => {
  await nextTick();
  initMap();
});

watch(() => [props.stations, props.activeFuel, props.selectedStationId], syncMarkers, { deep: true });
</script>

<style>
.fuel-marker {
  background: none !important;
  border: none !important;
}

.leaflet-container {
  font-family: var(--font-sans);
  background: var(--color-trotoar);
}

.leaflet-control-zoom {
  border: 1px solid var(--color-garis) !important;
  border-radius: 10px !important;
  overflow: hidden;
  box-shadow: none !important;
}

.leaflet-control-zoom a {
  color: var(--color-aspal) !important;
}

/* Mobile: station sheet covers the bottom 42vh; lift controls above it, pinch replaces zoom buttons. */
@media (max-width: 1023px) {
  .leaflet-bottom { bottom: 42vh; }
  .leaflet-control-zoom { display: none; }
}

.leaflet-control-attribution {
  font-size: 10px;
  background: rgb(255 255 255 / 0.8) !important;
}

/* Price pin: anchored by its bottom-center tip. */
.pin {
  position: absolute;
  bottom: 6px;
  left: 0;
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 3px 8px 3px 5px;
  border-radius: 999px;
  background: #fff;
  color: var(--color-petrol);
  border: 1px solid var(--color-petrol);
  box-shadow: 0 2px 6px rgb(30 35 40 / 0.18);
  transform: translateX(-50%);
  white-space: nowrap;
  cursor: pointer;
}

.pin::after {
  content: '';
  position: absolute;
  left: 50%;
  bottom: -5px;
  width: 8px;
  height: 8px;
  background: inherit;
  border-right: inherit;
  border-bottom: inherit;
  transform: translateX(-50%) rotate(45deg);
}

.pin__price {
  font-family: var(--font-board);
  font-weight: 800;
  font-size: 14px;
  line-height: 1;
  font-variant-numeric: tabular-nums;
}

.pin__dot {
  width: 10px;
  height: 10px;
  border-radius: 999px;
  box-shadow: 0 0 0 2px #fff;
}

.pin__dot--pertamina { background: #d7262e; }
.pin__dot--shell { background: #f2c500; }
.pin__dot--bp { background: #128a3e; }
.pin__dot--vivo { background: #1f5fd1; }

.pin--cheapest {
  background: var(--color-petrol-deep);
  border-color: var(--color-petrol-deep);
  color: var(--color-lampu);
}

.pin--selected {
  background: var(--color-teal);
  border-color: var(--color-teal);
  color: #fff;
  z-index: 10;
}

.pin__card {
  position: absolute;
  bottom: calc(100% + 10px);
  left: 50%;
  display: grid;
  gap: 2px;
  width: 220px;
  padding: 12px 14px;
  border-radius: 12px;
  background: #fff;
  color: var(--color-aspal);
  border: 1px solid var(--color-garis);
  box-shadow: 0 10px 30px rgb(30 35 40 / 0.18);
  font-size: 12px;
  white-space: normal;
  opacity: 0;
  pointer-events: none;
  transform: translate(-50%, 4px);
  transition: opacity 140ms ease, transform 140ms ease;
}

.pin__card strong {
  font-size: 14px;
  font-weight: 700;
}

.pin__card span {
  color: var(--color-redup);
}

.pin__card a {
  margin-top: 8px;
  color: var(--color-teal) !important;
  font-weight: 700;
  text-decoration: underline;
  text-underline-offset: 3px;
}

.pin:hover .pin__card,
.pin:focus-visible .pin__card,
.pin--selected .pin__card {
  opacity: 1;
  pointer-events: auto;
  transform: translate(-50%, 0);
}

.cluster {
  position: absolute;
  left: 0;
  top: 0;
  display: grid;
  place-items: center;
  width: max-content;
  min-width: 56px;
  padding: 6px 10px;
  white-space: nowrap;
  border: 2px solid #fff;
  border-radius: 999px;
  background: var(--color-petrol);
  color: #fff;
  box-shadow: 0 3px 10px rgb(30 35 40 / 0.25);
  transform: translate(-50%, -50%);
  font: inherit;
  cursor: pointer;
}

.cluster strong {
  font-family: var(--font-board);
  font-size: 18px;
  font-weight: 800;
  line-height: 1;
}

.cluster span {
  font-size: 10px;
  color: var(--color-lampu);
  line-height: 1.2;
}
</style>
