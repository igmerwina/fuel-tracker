<template>
  <div class="relative h-full w-full overflow-hidden bg-white">
    <div
      v-if="!mapReady"
      class="absolute inset-0 z-[400] animate-pulse bg-[linear-gradient(110deg,#f8fafc_8%,#eef2f7_18%,#f8fafc_33%)] bg-[length:200%_100%]"
      aria-hidden="true"
    ></div>
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
  averagePrice: number;
  selectedStationId: string;
}>();

const emit = defineEmits<{
  'select-station': [station: StationResult | null];
}>();

const mapContainer = ref<HTMLElement>();
const selectedStation = ref<StationResult | null>(null);
const mapReady = ref(false);
let map: L.Map | null = null;
let markerLayer: L.LayerGroup | null = null;
const landingCenter: L.LatLngExpression = [-6.1856, 106.8272];
const landingZoom = 11;
const stationZoom = 15;

const activePrice = (station: StationResult) =>
  station.prices.find((price) => price.type === props.activeFuel)?.price ?? station.selectedPrice;

const ron92Price = (station: StationResult) =>
  station.prices.find((price) => price.type === 'RON_92')?.price ?? station.selectedPrice;

const markerTone = (station: StationResult) => {
  if (props.selectedStationId === station.id) return 'selected';
  if (station.isCheapest) return 'cheapest';
  return 'default';
};

const escapeHtml = (value: string) =>
  value
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&#39;');

const escapeAttribute = (value: string) => escapeHtml(value);

const formatPrice = (price: number) => price.toLocaleString('id-ID');

const routeUrl = (station: StationResult) =>
  station.googleMapsUrl || `https://www.google.com/maps/search/?api=1&query=${station.latitude},${station.longitude}`;

const brandLogo = (brand: string) => {
  const logos: Record<string, string> = {
    BP: '<span class="brand-logo brand-logo--bp"><i class="fa-solid fa-leaf" aria-hidden="true"></i></span>',
    Shell: '<span class="brand-logo brand-logo--shell"><i class="fa-solid fa-fan" aria-hidden="true"></i></span>',
    Vivo: '<span class="brand-logo brand-logo--vivo">V</span>',
    Pertamina: '<span class="brand-logo brand-logo--pertamina"><i class="fa-solid fa-gas-pump" aria-hidden="true"></i></span>',
  };
  return logos[brand] ?? `<span class="brand-logo">${brand.charAt(0)}</span>`;
};

const clusterStations = () => {
  if (!map || map.getZoom() >= 14) return props.stations.map((station) => [station] as StationResult[]);

  const buckets = new globalThis.Map<string, StationResult[]>();
  props.stations.forEach((station) => {
    const key = `${Math.round(station.latitude / 0.035)}:${Math.round(station.longitude / 0.035)}`;
    buckets.set(key, [...(buckets.get(key) ?? []), station]);
  });
  return [...buckets.values()];
};

const clusterTone = (count: number) => {
  if (count >= 10) return 'hot';
  if (count >= 6) return 'warm';
  if (count >= 3) return 'fresh';
  return 'soft';
};

const createStationMarker = (station: StationResult) => {
  const tone = markerTone(station);
  const price = ron92Price(station);
  const savings = Math.max(0, station.savingsPerLiter);
  const content = `
    <div class="brand-pin brand-pin--${tone}" role="button" tabindex="0" aria-label="${escapeAttribute(station.brand)} ${escapeAttribute(station.name)}">
      ${brandLogo(station.brand)}
      ${station.isCheapest ? '<span class="brand-pin__flame"><i class="fa-solid fa-fire-flame-curved" aria-hidden="true"></i></span>' : ''}
      <span class="price-tooltip">
        <span class="price-tooltip__meta">${escapeHtml(station.brand)} · ${station.distanceKm.toFixed(1)} km</span>
        <strong>${escapeHtml(station.name)}</strong>
        <span class="price-tooltip__price">RON92 Rp ${formatPrice(price)}/L</span>
        <span class="price-tooltip__save">Hemat Rp ${formatPrice(savings)}/L</span>
        <a class="price-tooltip__route" href="${escapeAttribute(routeUrl(station))}" target="_blank" rel="noopener noreferrer">
          <i class="fa-solid fa-route" aria-hidden="true"></i>
          Rute
        </a>
      </span>
    </div>`;

  return L.marker([station.latitude, station.longitude], {
    icon: L.divIcon({
      html: content,
      iconSize: [42, 42],
      iconAnchor: [21, 21],
      className: 'fuel-marker',
    }),
  }).on('click', () => {
    if (props.selectedStationId === station.id) {
      selectedStation.value = null;
      emit('select-station', null);
      return;
    }
    selectedStation.value = station;
    emit('select-station', station);
  });
};

const createClusterMarker = (stations: StationResult[]) => {
  const cheapest = [...stations].sort((a, b) => activePrice(a) - activePrice(b))[0];
  return L.marker([cheapest.latitude, cheapest.longitude], {
    icon: L.divIcon({
      html: `
        <button class="gm-cluster gm-cluster--${clusterTone(stations.length)}" aria-label="${stations.length} SPBU terdekat">
          <strong>${stations.length}</strong>
          <span>mulai Rp ${formatPrice(ron92Price(cheapest))}</span>
        </button>
      `,
      iconSize: [74, 44],
      iconAnchor: [37, 37],
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
    minZoom: landingZoom,
    zoomControl: false,
    worldCopyJump: false,
  });

  L.control.zoom({ position: 'bottomright' }).addTo(map);
  L.tileLayer('https://{s}.basemaps.cartocdn.com/rastertiles/voyager/{z}/{x}/{y}{r}.png', {
    attribution: '&copy; OpenStreetMap contributors &copy; CARTO',
    bounds: jakartaBounds,
    maxZoom: 19,
    minZoom: landingZoom,
    noWrap: true,
  }).addTo(map);

  markerLayer = L.layerGroup().addTo(map);
  map.on('zoomend', syncMarkers);
  map.setMaxBounds(jakartaBounds);
  map.setView(landingCenter, landingZoom);
  syncMarkers();
  mapReady.value = true;
};

const flyToStation = (station: StationResult) => {
  if (!map) return;
  map.flyTo([station.latitude, station.longitude], stationZoom, { duration: 0.9 });
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

.leaflet-control-zoom a {
  color: #0f172a;
}

.leaflet-tile-pane {
  filter: saturate(0.74) contrast(0.98) brightness(1.04);
}

.leaflet-control-zoom {
  overflow: hidden;
  border: 0 !important;
  border-radius: 14px !important;
  box-shadow: 0 14px 32px rgba(15, 23, 42, 0.16) !important;
}

.brand-logo {
  display: grid;
  width: 28px;
  height: 28px;
  place-items: center;
  border-radius: 999px;
  font-size: 13px;
  font-weight: 900;
}

.brand-logo--bp {
  color: #15803d;
  background: #dcfce7;
}

.brand-logo--shell {
  color: #991b1b;
  background: #fef3c7;
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
  position: relative;
  display: grid;
  width: 38px;
  height: 38px;
  place-items: center;
  border: 2px solid #ffffff;
  border-radius: 999px;
  background: #ffffff;
  box-shadow: 0 12px 26px rgba(15, 23, 42, 0.18);
  cursor: pointer;
  transform-origin: center;
  transition:
    box-shadow 180ms ease,
    outline-color 180ms ease,
    transform 180ms ease;
}

.brand-pin:hover {
  z-index: 800;
  box-shadow: 0 16px 32px rgba(15, 23, 42, 0.24);
  transform: scale(1.08);
}

.brand-pin--cheapest {
  outline: 3px solid rgba(16, 185, 129, 0.5);
  box-shadow: 0 14px 30px rgba(16, 185, 129, 0.28);
}

.brand-pin--selected {
  z-index: 900;
  outline: 3px solid rgba(37, 99, 235, 0.72);
  transform: scale(1.14);
  animation: selectedPulse 1.8s ease-in-out infinite;
}

.brand-pin__flame {
  position: absolute;
  top: -7px;
  right: -5px;
  display: grid;
  width: 18px;
  height: 18px;
  place-items: center;
  border: 2px solid #ffffff;
  border-radius: 999px;
  background: #10b981;
  color: #ffffff;
  font-size: 9px;
}

.price-tooltip {
  position: absolute;
  bottom: calc(100% + 12px);
  left: 50%;
  z-index: 1000;
  display: grid;
  width: 210px;
  gap: 4px;
  padding: 12px;
  border: 1px solid rgba(226, 232, 240, 0.9);
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.98);
  box-shadow: 0 24px 55px rgba(15, 23, 42, 0.2);
  color: #0f172a;
  font-family: Inter, ui-sans-serif, system-ui, sans-serif;
  opacity: 0;
  pointer-events: none;
  text-align: left;
  transform: translate(-50%, 8px);
  transition:
    opacity 160ms ease,
    transform 160ms ease;
}

.price-tooltip::after {
  position: absolute;
  bottom: -6px;
  left: 50%;
  width: 12px;
  height: 12px;
  border-right: 1px solid rgba(226, 232, 240, 0.9);
  border-bottom: 1px solid rgba(226, 232, 240, 0.9);
  background: #ffffff;
  content: '';
  transform: translateX(-50%) rotate(45deg);
}

.brand-pin:hover .price-tooltip,
.brand-pin--selected .price-tooltip {
  opacity: 1;
  pointer-events: auto;
  transform: translate(-50%, 0);
}

.price-tooltip__meta,
.price-tooltip__save {
  color: #64748b;
  font-size: 11px;
  font-weight: 800;
}

.price-tooltip strong {
  overflow: hidden;
  font-size: 13px;
  font-weight: 950;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.price-tooltip__price {
  color: #0f172a;
  font-size: 15px;
  font-weight: 950;
}

.price-tooltip__route {
  z-index: 1;
  display: inline-flex;
  width: max-content;
  align-items: center;
  gap: 6px;
  margin-top: 6px;
  padding: 8px 12px;
  border-radius: 999px;
  border: 1px solid rgba(29, 78, 216, 0.22);
  background: linear-gradient(135deg, #2563eb, #1d4ed8);
  color: #ffffff !important;
  font-size: 12px;
  font-weight: 950;
  text-decoration: none;
  box-shadow: 0 10px 22px rgba(37, 99, 235, 0.28);
}

.price-tooltip__route i {
  color: #ffffff;
}

.gm-cluster {
  display: grid;
  min-width: 74px;
  gap: 1px;
  padding: 8px 11px;
  border: 2px solid #ffffff;
  border-radius: 999px;
  background: linear-gradient(135deg, #38bdf8, #2563eb);
  color: #ffffff;
  box-shadow: 0 14px 30px rgba(15, 23, 42, 0.22);
  cursor: pointer;
  font-family: Inter, ui-sans-serif, system-ui, sans-serif;
  text-align: center;
  transition: transform 180ms ease;
}

.gm-cluster:hover {
  transform: scale(1.08);
}

.gm-cluster--soft {
  background: linear-gradient(135deg, #38bdf8, #2563eb);
}

.gm-cluster--fresh {
  background: linear-gradient(135deg, #22c55e, #14b8a6);
}

.gm-cluster--warm {
  background: linear-gradient(135deg, #f59e0b, #f97316);
}

.gm-cluster--hot {
  background: linear-gradient(135deg, #f472b6, #ef4444);
  box-shadow: 0 16px 34px rgba(239, 68, 68, 0.28);
}

.gm-cluster strong {
  font-size: 14px;
  font-weight: 950;
  line-height: 1;
}

.gm-cluster span {
  font-size: 9px;
  font-weight: 800;
  opacity: 0.82;
}

@keyframes selectedPulse {
  0%,
  100% {
    box-shadow: 0 0 0 0 rgba(37, 99, 235, 0.28), 0 14px 30px rgba(37, 99, 235, 0.22);
  }
  50% {
    box-shadow: 0 0 0 7px rgba(37, 99, 235, 0), 0 18px 34px rgba(37, 99, 235, 0.26);
  }
}
</style>
