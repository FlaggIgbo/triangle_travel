<script lang="ts">
  import { error } from '@sveltejs/kit';

  // 404: unknown route – triggers +error.svelte
  export function load() {
    error(404, 'Page not found');
  }
</script>
