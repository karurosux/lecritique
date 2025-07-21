<script lang="ts">
  let {
    checked = $bindable(false),
    disabled = false,
    id = '',
    class: className = '',
    ...rest
  }: {
    checked?: boolean;
    disabled?: boolean;
    id?: string;
    class?: string;
    [key: string]: any;
  } = $props();
</script>

<input
  type="checkbox"
  {id}
  bind:checked
  {disabled}
  class="h-4 w-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500 focus:ring-2 {disabled ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'} {className}"
  {...rest}
/>