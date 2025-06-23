<script lang="ts">
  let {
    label = '',
    value = '',
    variant = 'blue',
    removable = true,
    onremove = (value: string) => {},
    ...restProps
  }: {
    label?: string;
    value?: string;
    variant?: 'blue' | 'green' | 'purple' | 'orange' | 'gray';
    removable?: boolean;
    onremove?: (value: string) => void;
    class?: string;
    [key: string]: any;
  } = $props();
  
  let className = restProps.class || '';

  const variantClasses = {
    blue: 'bg-blue-100 text-blue-800 hover:bg-blue-200',
    green: 'bg-green-100 text-green-800 hover:bg-green-200',
    purple: 'bg-purple-100 text-purple-800 hover:bg-purple-200',
    orange: 'bg-orange-100 text-orange-800 hover:bg-orange-200',
    gray: 'bg-gray-100 text-gray-800 hover:bg-gray-200'
  };

  const removeButtonClasses = {
    blue: 'hover:text-blue-600',
    green: 'hover:text-green-600',
    purple: 'hover:text-purple-600',
    orange: 'hover:text-orange-600',
    gray: 'hover:text-gray-600'
  };

  function handleRemove() {
    onremove(value);
  }
</script>

<span class="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium transition-colors duration-200 {variantClasses[variant]} {className}">
  {label}
  {#if removable}
    <button 
      class="ml-1 {removeButtonClasses} transition-colors duration-200" 
      onclick={handleRemove}
      aria-label="Remove filter"
    >
      ×
    </button>
  {/if}
</span>