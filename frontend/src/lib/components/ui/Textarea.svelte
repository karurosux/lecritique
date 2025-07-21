<script lang="ts">
  let {
    value = $bindable(''),
    placeholder = '',
    rows = 3,
    disabled = false,
    required = false,
    id = '',
    class: className = '',
    ...rest
  }: {
    value?: string;
    placeholder?: string;
    rows?: number;
    disabled?: boolean;
    required?: boolean;
    id?: string;
    class?: string;
    [key: string]: any;
  } = $props();
</script>

<textarea
  {id}
  bind:value
  {placeholder}
  {rows}
  {disabled}
  {required}
  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors resize-none {disabled ? 'bg-gray-100 text-gray-500' : 'bg-white'} {className}"
  {...rest}
></textarea>