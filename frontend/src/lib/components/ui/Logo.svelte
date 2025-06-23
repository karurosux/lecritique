<script lang="ts">
  import logoImage from '../../../assets/logo.png';
  
  export let size: 'sm' | 'md' | 'lg' | 'xl' = 'md';
  export let className = '';

  const sizes = {
    sm: 'h-6',
    md: 'h-8', 
    lg: 'h-10',
    xl: 'h-12'
  };

  $: currentSize = sizes[size];
</script>

<div class="flex items-center {className}">
  <img 
    src={logoImage} 
    alt="LeCritique Logo" 
    class="{currentSize} object-contain transition-transform duration-200 hover:scale-105"
  />
</div>