<script setup lang="ts">
import { nextTick, ref, watch } from 'vue';
import type { FuelType, SortMode, StationResult } from '../types';

const props = defineProps<{
  stations: StationResult[];
  activeFuel: FuelType;
  sortMode: SortMode;
  selectedStationId: string;
}>();

const emit = defineEmits<{
  flyto: [station: StationResult];
  'toggle-favorite': [stationId: string];
  'start-navigation': [station: StationResult];
  'reset-filters': [];
}>();

const cardRefs = ref<Record<string, HTMLElement>>({});

const setCardRef = (stationId: string, element: Element | null | unknown) => {
  if (element instanceof HTMLElement) cardRefs.value[stationId] = element;
};

const price = (station: StationResult, type: FuelType) =>
  station.prices.find((item) => item.type === type)?.price.toLocaleString('id-ID') ?? '-';

const activePrice = (station: StationResult, type: FuelType) =>
  station.prices.find((item) => item.type === type)?.price.toLocaleString('id-ID') ?? station.selectedPrice.toLocaleString('id-ID');

watch(
  () => props.selectedStationId,
  async (stationId) => {
    await nextTick();
    cardRefs.value[stationId]?.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
  },
);
</script>

<template>
  <aside
    class="fixed bottom-0 left-2 right-2 z-[900] max-h-[42vh] overflow-hidden rounded-t-[24px] bg-white shadow-2xl shadow-slate-950/16 lg:static lg:max-h-none lg:rounded-2xl lg:shadow-sm lg:ring-1 lg:ring-slate-100"
    aria-label="Hasil SPBU"
  >
    <div class="mx-auto mt-2 h-1.5 w-12 rounded-full bg-slate-300 lg:hidden"></div>

    <div class="flex h-14 items-center justify-between px-5">
      <div>
        <p class="text-[11px] font-black uppercase tracking-wide text-slate-500">
          {{ sortMode === 'nearest' ? 'Jarak terdekat' : 'Harga terendah' }}
        </p>
        <h2 class="text-base font-black text-slate-950">{{ stations.length }} SPBU</h2>
      </div>
    </div>

    <div class="max-h-[calc(42vh-64px)] overflow-y-auto px-3 pb-4 lg:max-h-[calc(100vh-150px)]">
      <div v-if="!stations.length" class="rounded-2xl bg-slate-50 p-5 text-center ring-1 ring-slate-100">
        <p class="text-sm font-black text-slate-950">Tidak ada SPBU cocok</p>
        <p class="mt-1 text-xs font-bold text-slate-500">Coba ubah kata kunci, merek, wilayah, atau jenis BBM.</p>
        <button
          type="button"
          class="mt-4 inline-flex h-9 items-center justify-center rounded-full bg-blue-600 px-4 text-xs font-black text-white"
          @click="emit('reset-filters')"
        >
          Reset filter
        </button>
      </div>

      <article
        v-for="station in stations"
        :key="station.id"
        :ref="(element) => setCardRef(station.id, element)"
        class="mb-2.5 grid min-h-[104px] w-full cursor-pointer gap-2 rounded-2xl bg-white p-3 text-left shadow-sm ring-1 transition duration-200 hover:-translate-y-0.5 hover:shadow-lg hover:shadow-slate-200/80 focus:outline-none focus:ring-2 focus:ring-blue-600"
        :class="[
          selectedStationId === station.id ? 'ring-blue-500 shadow-blue-100 animate-selected-card' : station.isCheapest ? 'ring-emerald-300 bg-emerald-50/35' : 'ring-slate-100',
        ]"
        role="button"
        tabindex="0"
        @click="emit('flyto', station)"
        @keydown.enter="emit('flyto', station)"
      >
        <div class="flex items-start justify-between gap-2">
          <div class="min-w-0">
            <div class="flex items-center gap-2">
              <span
                class="grid h-8 w-8 shrink-0 place-items-center rounded-full text-xs font-black"
                :class="station.brand === 'Shell' ? 'bg-amber-100 text-red-800' : station.brand === 'Vivo' ? 'bg-blue-600 text-white' : station.brand === 'BP' ? 'bg-emerald-100 text-emerald-700' : 'bg-red-500 text-white'"
              >
                <i v-if="station.brand === 'Pertamina'" class="fa-solid fa-gas-pump" aria-hidden="true"></i>
                <i v-else-if="station.brand === 'Shell'" class="fa-solid fa-fan" aria-hidden="true"></i>
                <i v-else-if="station.brand === 'BP'" class="fa-solid fa-leaf" aria-hidden="true"></i>
                <span v-else>V</span>
              </span>
              <div class="min-w-0">
                <h3 class="truncate text-sm font-black leading-5 text-slate-950">{{ station.name }}</h3>
                <p class="truncate text-xs font-bold text-slate-500">{{ station.brand }} · {{ station.region }}</p>
              </div>
              <span v-if="station.isCheapest" class="shrink-0 rounded-full bg-emerald-500 px-2 py-0.5 text-[10px] font-black text-white">
                <i class="fa-solid fa-fire-flame-curved" aria-hidden="true"></i>
              </span>
              <span v-else-if="station.hasSamePrice" class="shrink-0 rounded-full bg-emerald-50 px-2 py-0.5 text-[10px] font-black text-emerald-700">
                sama
              </span>
            </div>
          </div>
          <button
            type="button"
            class="grid h-8 w-8 shrink-0 place-items-center rounded-full text-slate-400 transition hover:bg-rose-50 hover:text-rose-500"
            :aria-label="station.isFavorite ? 'Hapus favorit' : 'Simpan favorit'"
            @click.stop="emit('toggle-favorite', station.id)"
          >
            <i :class="station.isFavorite ? 'fa-solid fa-heart text-rose-500' : 'fa-regular fa-heart'" aria-hidden="true"></i>
          </button>
        </div>

        <div class="flex flex-wrap items-center gap-x-3 gap-y-1 text-xs font-bold text-slate-500">
          <span class="text-slate-950">Aktif Rp{{ activePrice(station, activeFuel) }}/L</span>
          <span>RON90 Rp{{ price(station, 'RON_90') }}</span>
          <span>RON92 Rp{{ price(station, 'RON_92') }}</span>
          <span>RON95 Rp{{ price(station, 'RON_95') }}</span>
        </div>

        <div class="flex items-center justify-between gap-2">
          <p class="min-w-0 truncate text-xs font-bold text-slate-500">
            <i class="fa-solid fa-location-dot text-slate-400" aria-hidden="true"></i>
            {{ station.distanceKm.toFixed(1) }} km · Hemat Rp{{ station.savingsPerLiter.toLocaleString('id-ID') }}/L
          </p>
          <button
            type="button"
            class="inline-flex h-8 shrink-0 items-center gap-1.5 rounded-full px-3 text-xs font-black transition"
            :class="selectedStationId === station.id || station.isCheapest ? 'bg-blue-600 text-white shadow-md shadow-blue-600/20 hover:bg-blue-700' : 'bg-slate-100 text-slate-700 hover:bg-slate-200'"
            @click.stop="emit('start-navigation', station)"
          >
            <i class="fa-solid fa-route" aria-hidden="true"></i>
            <span v-if="selectedStationId === station.id || station.isCheapest">Rute</span>
          </button>
        </div>
      </article>
    </div>
  </aside>
</template>

<style scoped>
@keyframes selectedCard {
  0%,
  100% {
    box-shadow: 0 10px 28px rgba(37, 99, 235, 0.1);
  }
  50% {
    box-shadow: 0 14px 36px rgba(37, 99, 235, 0.18);
  }
}

.animate-selected-card {
  animation: selectedCard 1.8s ease-in-out infinite;
}
</style>
