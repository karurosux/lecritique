<script lang="ts">
  type InputType = 'text' | 'email' | 'password' | 'number' | 'tel' | 'url' | 'search';

  export let type: InputType = 'text';
  export let value = '';
  export let placeholder = '';
  export let disabled = false;
  export let required = false;
  export let error = '';
  export let label = '';
  export let id = '';

  const inputClasses = 'block w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-gray-900 placeholder-gray-500 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 disabled:cursor-not-allowed disabled:bg-gray-50 disabled:text-gray-500';
  const errorClasses = 'border-red-300 focus:border-red-500 focus:ring-red-500';
</script>

<div class="space-y-1">
  {#if label}
    <label for={id} class="block text-sm font-medium text-gray-700">
      {label}
      {#if required}
        <span class="text-red-500">*</span>
      {/if}
    </label>
  {/if}
  
  <input
    {id}
    {type}
    {placeholder}
    {disabled}
    {required}
    bind:value
    class="{inputClasses} {error ? errorClasses : ''}"
    on:input
    on:focus
    on:blur
    on:keydown
  />
  
  {#if error}
    <p class="text-sm text-red-600">{error}</p>
  {/if}
</div>