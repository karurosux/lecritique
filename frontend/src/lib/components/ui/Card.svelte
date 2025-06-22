<script lang="ts">
  export let padding = true;
  export let shadow = true;
  export let hover = false;
  
  const baseClasses = 'bg-white rounded-lg border border-gray-200';
  const paddingClasses = padding ? 'p-6' : '';
  const shadowClasses = shadow ? 'shadow-sm' : '';
  const hoverClasses = hover ? 'hover:shadow-md transition-shadow duration-200' : '';
</script>

<div class="{baseClasses} {paddingClasses} {shadowClasses} {hoverClasses}">
  <slot />
</div>