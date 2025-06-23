<script lang="ts">
  type InputType = 'text' | 'email' | 'password' | 'number' | 'tel' | 'url' | 'search';
  type InputVariant = 'default' | 'filled' | 'glass';

  let {
    type = 'text',
    variant = 'default',
    value = $bindable(''),
    placeholder = '',
    disabled = false,
    required = false,
    error = '',
    label = '',
    id = '',
    icon = '',
    oninput = undefined,
    onkeydown = undefined,
    ...restProps
  }: {
    type?: InputType;
    variant?: InputVariant;
    value?: string;
    placeholder?: string;
    disabled?: boolean;
    required?: boolean;
    error?: string;
    label?: string;
    id?: string;
    icon?: string;
    oninput?: ((event: Event) => void) | undefined;
    onkeydown?: ((event: KeyboardEvent) => void) | undefined;
    class?: string;
    [key: string]: any;
  } = $props();
  
  let focused = $state(false);

  const variantClasses = {
    default: 'bg-white border-2 border-gray-200 focus:border-blue-500 focus:bg-white',
    filled: 'bg-gray-50 border-2 border-transparent focus:border-blue-500 focus:bg-white',
    glass: 'bg-white/50 backdrop-blur-xl border-2 border-white/20 focus:border-blue-500 focus:bg-white/80'
  };

  let inputClasses = $derived(`
    block w-full rounded-xl px-4 py-3 text-gray-900 placeholder-gray-500 
    transition-all duration-300 ease-out
    focus:outline-none focus:ring-4 focus:ring-blue-500/20 focus:scale-[1.02]
    disabled:cursor-not-allowed disabled:bg-gray-100 disabled:text-gray-500 disabled:opacity-50
    ${variantClasses[variant]}
    ${error ? 'border-red-400 focus:border-red-500 focus:ring-red-500/20' : ''}
    ${icon ? 'pl-12' : ''}
  `);

  function handleFocus() {
    focused = true;
  }

  function handleBlur() {
    focused = false;
  }
</script>

<div class="space-y-2">
  {#if label}
    <label for={id} class="block text-sm font-semibold text-gray-700 transition-colors duration-200">
      {label}
      {#if required}
        <span class="text-red-500 ml-1">*</span>
      {/if}
    </label>
  {/if}
  
  <div class="relative">
    {#if icon}
      <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
        <svg class="h-5 w-5 text-gray-400 transition-colors duration-200" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={icon}></path>
        </svg>
      </div>
    {/if}
    
    <input
      {id}
      {type}
      {placeholder}
      {disabled}
      {required}
      bind:value
      class={inputClasses}
      oninput={oninput}
      onfocus={handleFocus}
      onblur={handleBlur}
      onkeydown={onkeydown}
    />
    
    {#if focused && !error}
      <div class="absolute inset-0 rounded-xl ring-4 ring-blue-500/10 pointer-events-none transition-all duration-300"></div>
    {/if}
  </div>
  
  {#if error}
    <div class="flex items-center space-x-2">
      <svg class="h-4 w-4 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
      </svg>
      <p class="text-sm text-red-600 font-medium">{error}</p>
    </div>
  {/if}
</div>