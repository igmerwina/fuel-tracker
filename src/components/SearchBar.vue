<script setup lang="ts">
defineProps<{
  query: string;
  resultCount: number;
}>();

const emit = defineEmits<{
  'update-query': [value: string];
}>();
</script>

<template>
  <section class="search-shell" aria-label="Pencarian SPBU dan BBM">
    <div class="search-icon" aria-hidden="true">
      <i class="fa-solid fa-magnifying-glass"></i>
    </div>
    <label class="search-field" for="mainSearch">
      <span>Cari BBM termurah di Jakarta</span>
      <input
        id="mainSearch"
        :value="query"
        type="search"
        placeholder="Ketik Menteng, Pertamax, Shell, RON 92, atau alamat SPBU"
        autocomplete="off"
        @input="emit('update-query', ($event.target as HTMLInputElement).value)"
      />
    </label>
    <div class="search-meta">
      <strong>{{ resultCount }}</strong>
      <span>hasil</span>
    </div>
  </section>
</template>

<style scoped>
.search-shell {
  display: grid;
  grid-template-columns: 48px minmax(0, 1fr) auto;
  align-items: center;
  gap: 14px;
  margin: 16px 18px 0;
  padding: 12px;
  border: 1px solid rgba(216, 224, 234, 0.88);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.92);
  box-shadow: var(--small-shadow);
  backdrop-filter: blur(10px);
}

.search-icon {
  display: grid;
  width: 48px;
  height: 48px;
  place-items: center;
  border-radius: 8px;
  color: #082f49;
  background: #ffdb4d;
}

.search-field {
  display: grid;
  gap: 2px;
  min-width: 0;
}

.search-field span {
  color: #475569;
  font-size: 12px;
  font-weight: 900;
  letter-spacing: 0.06em;
  text-transform: uppercase;
}

input {
  width: 100%;
  min-width: 0;
  border: 0;
  outline: 0;
  color: var(--ink);
  background: transparent;
  font-size: 20px;
  font-weight: 800;
}

input::placeholder {
  color: #94a3b8;
  font-weight: 650;
}

.search-meta {
  display: grid;
  min-width: 76px;
  min-height: 52px;
  place-items: center;
  border-radius: 8px;
  color: #ffffff;
  background: var(--brand);
  font-weight: 900;
}

.search-meta strong {
  line-height: 1;
}

.search-meta span {
  margin-top: -8px;
  font-size: 11px;
  text-transform: uppercase;
}

@media (max-width: 760px) {
  .search-shell {
    grid-template-columns: 40px minmax(0, 1fr);
    margin: 12px 12px 0;
  }

  .search-icon {
    width: 40px;
    height: 40px;
  }

  .search-meta {
    grid-column: 1 / -1;
    min-height: 38px;
  }

  input {
    font-size: 16px;
  }
}
</style>
