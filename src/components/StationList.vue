<script setup lang="ts">
import { nextTick, ref, watch } from 'vue';
import type { SortMode, StationResult } from '../types';

const props = defineProps<{
  stations: StationResult[];
  sortMode: SortMode;
  selectedStationId: string;
  locationNote: string;
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

const brandDot: Record<string, string> = {
  Pertamina: 'bg-[#d7262e]',
  Shell: 'bg-[#f2c500]',
  BP: 'bg-[#128a3e]',
  Vivo: 'bg-[#1f5fd1]',
};

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
    class="fixed inset-x-0 bottom-0 z-[900] flex max-h-[42vh] flex-col rounded-t-2xl border-t border-garis bg-white shadow-[0_-8px_30px_rgb(30_35_40/0.12)] lg:static lg:max-h-none lg:rounded-none lg:border-t-0 lg:border-l lg:shadow-none"
    aria-label="Daftar SPBU"
  >
    <div class="mx-auto mt-2 h-1 w-10 shrink-0 rounded-full bg-garis lg:hidden" aria-hidden="true"></div>

    <div class="flex shrink-0 items-baseline justify-between gap-3 px-4 pt-2 pb-3 lg:px-5 lg:pt-4">
      <h2 class="text-base font-bold">
        {{ stations.length }} SPBU
        <span class="font-normal text-redup">{{ sortMode === 'nearest' ? 'dari yang terdekat' : 'dari yang termurah' }}</span>
      </h2>
      <span v-if="locationNote" class="text-xs text-redup">{{ locationNote }}</span>
    </div>

    <div class="min-h-0 flex-1 overflow-y-auto border-t border-garis">
      <div v-if="!stations.length" class="px-5 py-10 text-center">
        <p class="font-semibold">Tidak ada SPBU yang cocok</p>
        <p class="mt-1 text-sm text-redup">Longgarkan pencarian, merek, atau wilayah.</p>
        <button
          type="button"
          class="mt-4 h-9 rounded-lg bg-teal px-4 text-sm font-semibold text-white hover:bg-teal-dark"
          @click="emit('reset-filters')"
        >
          Hapus semua filter
        </button>
      </div>

      <ul v-else>
        <li
          v-for="station in stations"
          :key="station.id"
          :ref="(element) => setCardRef(station.id, element)"
          class="group relative grid cursor-pointer grid-cols-[minmax(0,1fr)_auto] gap-x-3 border-b border-garis px-4 py-3 transition-colors hover:bg-trotoar/60 lg:px-5"
          :class="selectedStationId === station.id ? 'bg-teal/[0.06] shadow-[inset_3px_0_0_var(--color-teal)]' : ''"
          tabindex="0"
          :aria-current="selectedStationId === station.id ? 'true' : undefined"
          @click="emit('flyto', station)"
          @keydown.enter="emit('flyto', station)"
        >
          <div class="min-w-0">
            <div class="flex items-center gap-2">
              <span class="h-2.5 w-2.5 shrink-0 rounded-full" :class="brandDot[station.brand]" aria-hidden="true"></span>
              <h3 class="truncate text-[15px] font-semibold">{{ station.name }}</h3>
              <span
                v-if="station.isCheapest"
                class="shrink-0 rounded bg-lampu px-1.5 py-px text-[11px] font-bold text-petrol-deep"
              >Termurah</span>
            </div>
            <p class="mt-0.5 truncate text-[13px] text-redup">{{ station.brand }}, {{ station.region }}</p>
            <p class="mt-1 text-[13px] text-aspal-soft">
              {{ station.distanceKm.toFixed(1) }} km
              <template v-if="station.savingsPerLiter > 0">
                <span class="text-garis" aria-hidden="true">/</span>
                <span class="text-teal">hemat Rp{{ station.savingsPerLiter.toLocaleString('id-ID') }}/L</span>
              </template>
            </p>
          </div>

          <div class="flex flex-col items-end justify-between">
            <p class="price text-[26px] leading-none font-extrabold text-petrol">
              <span class="mr-0.5 text-sm font-semibold text-redup">Rp</span>{{ station.selectedPrice.toLocaleString('id-ID') }}
            </p>
            <div class="mt-2 flex items-center gap-1">
              <button
                type="button"
                class="grid h-8 w-8 place-items-center rounded-md text-redup transition hover:bg-white hover:text-[#c2334a]"
                :aria-label="station.isFavorite ? `Hapus ${station.name} dari favorit` : `Simpan ${station.name} ke favorit`"
                :aria-pressed="station.isFavorite"
                @click.stop="emit('toggle-favorite', station.id)"
              >
                <i :class="station.isFavorite ? 'fa-solid fa-heart text-[#c2334a]' : 'fa-regular fa-heart'" aria-hidden="true"></i>
              </button>
              <button
                type="button"
                class="inline-flex h-8 items-center gap-1.5 rounded-md border border-garis bg-white px-2.5 text-[13px] font-semibold transition hover:border-teal hover:text-teal"
                @click.stop="emit('start-navigation', station)"
              >
                <i class="fa-solid fa-diamond-turn-right text-xs" aria-hidden="true"></i>
                Rute
              </button>
            </div>
          </div>
        </li>
      </ul>
    </div>
  </aside>
</template>
